package utilb

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/mahmudindes/orenocomic-bagicore/internal/model"
	"github.com/mahmudindes/orenocomic-bagicore/internal/utila"
)

func ResponseErr404(w http.ResponseWriter) {
	http.Error(w, "Not found.", http.StatusNotFound)
}

func ResponseErr500(w http.ResponseWriter) {
	http.Error(w, "Internal server error.", http.StatusInternalServerError)
}

type (
	ResError struct {
		Error ResErrorObject `json:"error"`
	}

	ResErrorObject struct {
		Status string `json:"status,omitempty"`
		Detail string `json:"detail,omitempty"`
	}
)

func ResponseJSON(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	data, _ := json.Marshal(v)
	w.Write(data)
}

func ResponseJSONErr(w http.ResponseWriter, err string, code int) {
	ResponseJSON(w, ResError{Error: ResErrorObject{
		Status: strconv.Itoa(code),
		Detail: err,
	}}, code)
}

func ResponseJSONErr404(w http.ResponseWriter) {
	ResponseJSONErr(w, "Not found.", http.StatusNotFound)
}

func ResponseJSONErr500(w http.ResponseWriter) {
	ResponseJSONErr(w, "Internal server error.", http.StatusInternalServerError)
}

func ResponseJSONAuthErr(w http.ResponseWriter, err error) {
	switch {
	case errors.As(err, &model.ErrGeneric):
		w.Header().Set("WWW-Authenticate", `Bearer error="invalid_token"`)
		ResponseJSONErr(w, utila.CapitalPeriod(err.Error()), http.StatusUnauthorized)
	default:
		ResponseJSONErr500(w)
	}
}

func ResponseJSONServiceErr(w http.ResponseWriter, err error) {
	switch {
	case errors.As(err, &model.ErrNotFound):
		ResponseJSONErr404(w)
	case errors.As(err, &model.ErrGeneric):
		ResponseJSONErr(w, utila.CapitalPeriod(err.Error()), http.StatusBadRequest)
	case errors.As(err, &model.ErrDatabase):
		ResponseJSONErr(w, "Database has encountered a problem.", http.StatusInternalServerError)
	default:
		ResponseJSONErr500(w)
	}
}
