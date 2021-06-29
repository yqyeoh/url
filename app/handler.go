package app

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/yqyeoh/url/config"
	"github.com/yqyeoh/url/platform/handler"
	"go.uber.org/zap"
)

type Handler struct {
	handler.Handler
	service Service
	conf    config.Config
}

func NewHandler(logger *zap.SugaredLogger, service Service, conf config.Config) Handler {
	return Handler{
		Handler: handler.New(logger),
		service: service,
		conf:    conf,
	}
}

func (h Handler) AddRoutes(router *mux.Router) {
	router.HandleFunc("/url", h.getShortenedURL).Methods(http.MethodPost)
	router.HandleFunc("/url", h.getOriginalURL).Methods(http.MethodGet)
}

func (h Handler) getShortenedURL(w http.ResponseWriter, req *http.Request) {
	var (
		ctx           = req.Context()
		logger        = h.Logger()
		shortenURLReq ShortenURLReq
	)

	if err := h.FromJSON(&shortenURLReq, req); err != nil {
		h.ReplyError(w, http.StatusBadRequest, "Bad JSON Request")
		return
	}

	originalURL := shortenURLReq.URL

	if isValidURL := govalidator.IsURL(originalURL); !isValidURL {
		h.ReplyError(w, http.StatusBadRequest, "Invalid URL")
		return
	}

	if !strings.Contains(originalURL, "://") {
		originalURL = "http://" + originalURL
	}

	code, err := h.service.Create(ctx, originalURL)
	if err != nil {
		logger.Errorf("failed to get URL Code: %v", err)
		h.ReplyError(w, http.StatusInternalServerError, "Unexpected error occurred")
		return
	}

	shortURL := fmt.Sprintf("%s/%s", h.conf.FrontEnd.BaseURL, code)

	response := URLRes{
		URL: shortURL,
	}
	h.ReplyJSON(w, http.StatusOK, response)
}

func (h Handler) getOriginalURL(w http.ResponseWriter, req *http.Request) {
	var (
		ctx               = req.Context()
		logger            = h.Logger()
		getOriginalURLReq GetOriginalURLReq
	)

	if err := h.FromQuery(&getOriginalURLReq, req); err != nil {
		h.ReplyError(w, http.StatusBadRequest, "Bad JSON Request")
		return
	}

	shortURL := getOriginalURLReq.URL
	code := ""

	for _, prefix := range h.conf.FrontEnd.URLPrefixesToMatch {
		isValidPrefix := strings.HasPrefix(strings.ToLower(shortURL), prefix) && len(shortURL) > len(prefix)
		if isValidPrefix {
			code = after(shortURL, prefix)
			break
		}
	}

	url, err := h.service.FindURLByCode(ctx, code)
	switch {
	case errors.Is(err, ErrInvalidCode):
		logger.Errorf("Invalid Code: %s", code)
		h.ReplyError(w, http.StatusBadRequest, "Invalid URL")
	case err != nil:
		logger.Errorf("failed to get original URL: %v", err)
		h.ReplyError(w, http.StatusInternalServerError, "Unexpected error occurred")
	default:
		response := URLRes{
			URL: url,
		}
		h.ReplyJSON(w, http.StatusOK, response)
	}
}

func after(wholestring string, substring string) string {
	pos := strings.LastIndex(wholestring, substring)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(substring)
	if adjustedPos >= len(wholestring) {
		return ""
	}
	return wholestring[adjustedPos:len(wholestring)]
}
