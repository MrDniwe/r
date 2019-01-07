package handlers

import (
	"net/http"

	"github.com/mrdniwe/r/pkg/templator"
)

func Home(t *templator.Pages) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t.Items["home.html"].Execute(w, r)
	}
}
