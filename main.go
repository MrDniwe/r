package main

import (
	"fmt"
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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)
	r.HandleFunc("/post/{id}", postHandler)
	http.Handle("/", r)

	fmt.Println("Server is running on :3000")
	http.ListenAndServe(":3000", nil)
}
