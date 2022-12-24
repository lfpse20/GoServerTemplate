package main

import (
	"GoServerTemplate/business/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// todo services will get repos and gateways past in

	_ = initServices() // todo inject service list in endpoints

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

func initServices() services.ServiceList {
	serviceList := services.ServiceList{}

	serviceList.FooService = services.NewFooService()
	serviceList.BarService = services.NewBarService()

	return serviceList
}
