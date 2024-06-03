package httpstatus

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/hashicorp/go-multierror"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

var Marshaler = json.Marshal

type apierror struct {
	Status string   `json:"status"`
	Errors []detail `json:"errors"`
}

func (a apierror) Error() string {
	bytes, _ := Marshaler(a)
	return string(bytes)
}

type detail struct {
	Message string `json:"message"`
}

func formatErrorResponse(ctx context.Context, w http.ResponseWriter, err error, status int) {
	if err == nil {
		panic("http api: error shouldn`t be nil")
	}

	span := trace.SpanFromContext(ctx)

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

	err = apierror{
		Status: http.StatusText(status),
		Errors: details,
	}

	span.RecordError(err)

	span.SetStatus(codes.Error, err.Error())

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	http.Error(w, err.Error(), status)
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

func BadRequest(ctx context.Context, w http.ResponseWriter, err error) {
	formatErrorResponse(ctx, w, err, http.StatusBadRequest)
}

/////////// 500 ///////////

func InternalServerError(ctx context.Context, w http.ResponseWriter, err error) {
	formatErrorResponse(ctx, w, err, http.StatusInternalServerError)
}
