package main

import (
	"log"
	"net/http"
)

type handlers map[string]func(http.ResponseWriter, *http.Request)

var h handlers

func init() {
	h = make(handlers)

	// Empty handler
	h["dummy"] = func(w http.ResponseWriter, r *http.Request) {
		err := tmpl["dummy"].ExecuteTemplate(w, "dummy.html", r)
		if err != nil {
			log.Println(err)
		}
	}

	// Home handler
	h["home"] = func(w http.ResponseWriter, r *http.Request) {
		err := tmpl["home"].ExecuteTemplate(w, "home.html", r)
		if err != nil {
			log.Println(err)
		}
	}

	// Post handler
	h["post"] = func(w http.ResponseWriter, r *http.Request) {
		err := tmpl["post"].ExecuteTemplate(w, "post.html", r)
		if err != nil {
			log.Println(err)
		}
	}

}
