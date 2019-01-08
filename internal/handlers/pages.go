package handlers

import (
	"net/http"

	"github.com/mrdniwe/r/pkg/templator"
)

func Home(t *templator.Pages) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t.Items["mainpage"].Execute(w, r)
	}
}

func Post(t *templator.Pages) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t.Items["post"].Execute(w, r)
	}
}

func Static(t *templator.Pages) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t.Items["static"].Execute(w, r)
	}
}

func Dummy(t *templator.Pages) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t.Items["dummy"].Execute(w, r)
	}
}
