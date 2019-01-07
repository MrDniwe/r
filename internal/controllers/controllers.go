package controllers

import (
	"github.com/gorilla/mux"
	"github.com/mrdniwe/r/internal/handlers"
	"github.com/mrdniwe/r/pkg/templator"
)

func Site(router *mux.Router, t *templator.Pages) {
	router.HandleFunc("/", handlers.Home(t)).Methods("GET")
	router.HandleFunc("/post/{id}", handlers.Post(t)).Methods("GET")
	router.HandleFunc("/info/{page}", handlers.Static(t)).Methods("GET")
}

// TODO: json api without templator
func Api(router *mux.Router, t *templator.Pages) {
	router.HandleFunc("/test", handlers.Home(t)).Methods("GET", "POST")
}
