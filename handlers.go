package main

import (
	"log"
	"net/http"
)

type handlers map[string]func(http.ResponseWriter, *http.Request)

var h handlers

func init() {
	h = make(handlers)

	// Emty handler
	h["dummy"] = func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "oummy.html", r)
		if err != nil {
			log.Println(err)
		}
	}

}
