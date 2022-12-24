package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", helloWorldHandler())
	log.Println("listening on port 8080")
	http.ListenAndServe(":8080", r)
}

func helloWorldHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}
}
