package errors

import (
	"errors"
	ers "github.com/pkg/errors"
	"net/http"
)

const (
	PostgresError   string = "Postgres request error"
	ValidationError string = "Validation error"
	UnmarshalError  string = "JSON unmarshalling failed"
	MailError       string = "Mail request error"
)

var (
	BadRequestErr = errors.New("Bad request")
	NotFoundErr   = errors.New("Not found")
	ServerErr     = errors.New("Server error")
	ForbiddenErr  = errors.New("Forbidden")
	UnknownErr    = errors.New("Unknown error")
)

type StackTracer interface {
	StackTrace() ers.StackTrace
}

func HandleError(err error, w http.ResponseWriter, r *http.Request) {
	switch err {
	case BadRequestErr:
		http.Redirect(w, r, "/errors/badrequest", http.StatusMovedPermanently)
	case NotFoundErr:
		http.Redirect(w, r, "/errors/notfound", http.StatusMovedPermanently)
	case ServerErr:
		http.Redirect(w, r, "/errors/server", http.StatusMovedPermanently)
	case ForbiddenErr:
		http.Redirect(w, r, "/errors/forbidden", http.StatusMovedPermanently)
	case UnknownErr:
		http.Redirect(w, r, "/errors/unknown", http.StatusMovedPermanently)
	default:
		http.Redirect(w, r, "/errors/unknown", http.StatusMovedPermanently)
	}
}
