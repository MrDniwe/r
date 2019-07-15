package errors

import (
	"errors"
	ers "github.com/pkg/errors"
	"net/http"
)

//TODO унифицировать типы ошибок
const (
	ServerError     string = "Server error"
	BadRequestError string = "Bad request"
)

var (
	BadRequestErr = errors.New(BadRequestError)
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
