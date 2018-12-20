package main

import (
	"net/http"
)

type handlers map[string]func(http.ResponseWriter, *http.Request)

var h handlers

func init() {
	h = make(handlers)

	// Empty handler
	h["dummy"] = func(w http.ResponseWriter, r *http.Request) {
		tmpl["dummy.html"].execute(w, r)
	}

	// Home handler
	h["home"] = func(w http.ResponseWriter, r *http.Request) {
		tmpl["home.html"].execute(w, r)
	}

	// Post handler
	h["post"] = func(w http.ResponseWriter, r *http.Request) {
		tmpl["post.html"].execute(w, r)
	}

	// Static pages handler
	h["info"] = func(w http.ResponseWriter, r *http.Request) {
		tmpl["static.html"].execute(w, r)
	}

}
