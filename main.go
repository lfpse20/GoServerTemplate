package main

import (
	"GoServerTemplate/api"
	"GoServerTemplate/api/endpoints"
	"GoServerTemplate/business"
	"context"
	transport "github.com/go-kit/kit/transport/http"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	services := business.InitServices()
	eps := api.InitEndpoints(services)

	r := mux.NewRouter()
	setUpHandlers(r, eps)

	log.Println("listening on port 8080")

	http.ListenAndServe(":8080", r)
}

func setUpHandlers(r *mux.Router, eps api.Endpoints) {
	for _, definition := range endpoints.Definitions {
		r.Methods(definition.Method).
			Path(definition.Path).
			Handler(
				transport.NewServer(
					api.MustGetEndpoint(eps, definition.EndpointName),
					definition.Decoder,
					encoder(),
				))
	}
}

func encoder() transport.EncodeResponseFunc {
	return func(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
		return nil
	}
}
