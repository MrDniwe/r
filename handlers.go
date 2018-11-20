package main

import (
	"fmt"
	"log"
	"net/http"
)

func emptyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside empty handler")
	err := tmpl.ExecuteTemplate(w, "post.html", fmt.Sprintf("Request URI: %v", r.RequestURI))
	if err != nil {
		log.Println(err)
	}
}
