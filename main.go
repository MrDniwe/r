package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Default format of %v", "string")
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprintf(w, "Looking for post id %v", id)
}

func suchMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)
	r.HandleFunc("/post/{id}", postHandler)
	r.Use(suchMiddleware)
	http.Handle("/", r)

	fmt.Println("Server is running on :3000")
	http.ListenAndServe(":3000", nil)
}
