package controllers

import (
	"github.com/gorilla/mux"
	"github.com/mrdniwe/r/internal/handlers"
	"github.com/mrdniwe/r/pkg/templator"
)

func Site(router *mux.Router, t *templator.Pages) {
	router.HandleFunc("/", handlers.Dummy(t)).Methods("GET")
}
