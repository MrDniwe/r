package handlers

import (
	"net/http"

	"github.com/mrdniwe/r/pkg/templator"
)

func Dummy(t *templator.Pages) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t.Items["dummy.html"].Execute(w, r)
	}
}
