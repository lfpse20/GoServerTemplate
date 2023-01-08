package api

import (
	"GoServerTemplate/api/endpoints"
	"GoServerTemplate/api/endpoints/fooEndpoints"
	"GoServerTemplate/business/services"
	"fmt"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints map[string]endpoint.Endpoint

func InitEndpoints(services services.ServiceList) Endpoints {
	eps := Endpoints{
		endpoints.GetFooEndpointName: fooEndpoints.GetFoo(services.FooService),
	}

	return eps
}

func MustGetEndpoint(endpoints Endpoints, endpointName string) endpoint.Endpoint {
	ep, ok := endpoints[endpointName]
	if !ok {
		panic(fmt.Sprintf("endpoint with name %s", endpointName))
	}

	return ep
}
