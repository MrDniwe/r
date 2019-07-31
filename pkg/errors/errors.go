package errors

import (
	"errors"
	"net/http"

	ers "github.com/pkg/errors"
)

const (
	PostgresError   string = "Postgres request error"
	ValidationError string = "Validation error"
	UnmarshalError  string = "JSON unmarshalling failed"
	MailError       string = "Mail request error"
	ToSoonCode      string = "codegen_request_to_soon"
	NotFoundCode    string = "user_not_found_by_email"
	WrongPassword   string = "wrong_password"
)

var (
	BadRequestErr    = errors.New("Bad request")
	NotFoundErr      = errors.New("Not found")
	ServerErr        = errors.New("Server error")
	ForbiddenErr     = errors.New("Forbidden")
	UnknownErr       = errors.New("Unknown error")
	DelayErr         = errors.New("Delay error")
	WrongPasswordErr = errors.New("Wrong password")
	InvalidTokenErr  = errors.New("Invalid token")
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
