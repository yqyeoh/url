package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/form"
	"go.uber.org/zap"
)

type Handler struct {
	logger      *zap.SugaredLogger
	formDecoder *form.Decoder
}

func New(logger *zap.SugaredLogger) Handler {
	formDecoder := form.NewDecoder()
	return Handler{
		logger,
		formDecoder,
	}
}

func (h Handler) Logger() *zap.SugaredLogger {
	return h.logger
}

// FromQuery unmarshals request query string to a struct
func (h Handler) FromQuery(dest interface{}, req *http.Request) error {
	if err := h.formDecoder.Decode(dest, req.URL.Query()); err != nil {
		return err
	}
	return nil
}

// FromJSON unmarshals request json body to a struct
func (h Handler) FromJSON(dest interface{}, req *http.Request) error {
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(dest); err != nil {
		return err
	}
	return nil
}

// ReplyError writes error to ResponseWriter
func (h Handler) ReplyError(w http.ResponseWriter, statusCode int, message string) {
	ReplyJSON(w, code, map[string]string{"error": message})
}

// ReplyJSON marshals a struct to json
func (h Handler) ReplyJSON(w http.ResponseWriter, statusCode int, payload interface{}) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		h.logger.Errorf("writing JSON to ResponseWriter failed on %T - error: %v", payload, err)
		fmt.Fprintln(w, e.Error())
	}
}
