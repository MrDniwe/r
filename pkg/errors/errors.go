package errors

import (
	"errors"
	"net/http"
)

var (
	BadRequestErr  = errors.New("Bad request")
	NotFoundErr    = errors.New("Not found")
	ServerError    = errors.New("Server error")
	ForbiddenError = errors.New("Forbidden")
	UnknownError   = errors.New("Unknown error")
)

func HandleError(err error, w http.ResponseWriter, r *http.Request) {
	switch err {
	case BadRequestErr:
		http.Redirect(w, r, "/errors/badrequest", http.StatusMovedPermanently)
	case NotFoundErr:
		http.Redirect(w, r, "/errors/notfound", http.StatusMovedPermanently)
	case ServerError:
		http.Redirect(w, r, "/errors/server", http.StatusMovedPermanently)
	case ForbiddenError:
		http.Redirect(w, r, "/errors/forbidden", http.StatusMovedPermanently)
	case UnknownError:
		http.Redirect(w, r, "/errors/unknown", http.StatusMovedPermanently)
	default:
		http.Redirect(w, r, "/errors/unknown", http.StatusMovedPermanently)
	}
}
