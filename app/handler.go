package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"github.com/asaskevich/govalidator"
	"github.com/yqyeoh/url/platform/handler"
)

type Handler struct {
	handler.Handler
	service Service
	conf config.Config
}

func NewHandler(logger *zap.SugaredLogger, service Service, conf config.Config) Handler {
	return Handler{
		Handler: handler.New(logger),
		service: service,
		conf: config.Config,
	}
}

func (h Handler) AddRoutes(router *mux.Router) {
	router.HandlerFunc("/url", h.getOriginalURL).Methods(http.MethodGet)
	router.HandlerFunc("/url", h.getShortenedURL).Methods(http.MethodPost)
}

func (h Handler) getShortenedURL((w http.ResponseWriter, req *http.Request){
	var (
		ctx = req.Contextogger()
		logger = h.Logger()
		shortenURLReq ShortenURLReq
	)
	
	if err := h.FromJSON(&shortenURLReq, req); err!=nil{
		h.ReplyError(w, req, http.StatusBadRequest, "Bad JSON Request")
		return
	}

	url := shortenURLReq.URL

	if isValidURL := govalidator.IsURL(url); !isValidURL{
		h.ReplyError(w, req, http.StatusBadRequest, "Invalid URL")
	}

	code, err := h.service.FindOrCreateCode(ctx, url)
	if err!=nil{
		logger.Errorf("failed to get URL Code: %v", err)
		h.ReplyError(w, req, http.StatusInternalServerError, "Unexpected error occurred")
	}

	shortenedURL := fmt.Sprintf("%s/%s", h.conf.FrontEnd.BaseURL, code)

	response := ShortenURLRespons{
		original : url,
		shortened: shortenedURL,
	}
	h.ReplyJSON(w, req, http.StatusOK, response)
}
