package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	staticFileDir := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDir))
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")
	r.HandleFunc("/hello", handler).Methods("GET")
	r.HandleFunc("/birds", getBirdHandler).Methods("GET")
	r.HandleFunc("/birds", createBirdHandler).Methods("POST")
	return r
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func main() {
	r := newRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}
