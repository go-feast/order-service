package httpstatus

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/go-multierror"
	"net/http"
)

var Marshaler = json.Marshal

type apierror struct {
	Status string   `json:"status"`
	Errors []detail `json:"errors"`
}

type detail struct {
	Message string `json:"message"`
}

func formatErrorResponse(w http.ResponseWriter, err error, status int) {
	if err == nil {
		panic("http api: error shouldn`t be nil")
	}

	details := make([]detail, 0)

	var merror *multierror.Error

	switch {
	case errors.As(err, &merror):
		for _, e := range merror.Errors {
			details = append(details, detail{Message: e.Error()})
		}
	default:
		details = append(details, detail{Message: err.Error()})
	}

	apierr := apierror{
		Status: http.StatusText(status),
		Errors: details,
	}

	bytes, _ := Marshaler(apierr)

	w.Header().Set("Content-Type", "application/json")
	http.Error(w, string(bytes), status)
}

func formatSuccessfulResponse(w http.ResponseWriter, i interface{}, status int) {
	bytes, err := Marshaler(i)
	if err != nil {
		return
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bytes)
}

/////////// 200 ///////////

func Ok(w http.ResponseWriter, i interface{}) {
	formatSuccessfulResponse(w, i, http.StatusOK)
}

func Created(w http.ResponseWriter, i interface{}) {
	formatSuccessfulResponse(w, i, http.StatusCreated)
}

/////////// 400 ///////////

func BadRequest(w http.ResponseWriter, err error) {
	formatErrorResponse(w, err, http.StatusBadRequest)
}

/////////// 500 ///////////

func InternalServerError(w http.ResponseWriter, err error) {
	formatErrorResponse(w, err, http.StatusInternalServerError)
}
