package fooEndpoints

import (
	"GoServerTemplate/business/services"
	"context"
	"github.com/go-kit/kit/endpoint"
)

func GetFoo(foodService services.FooServicer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return "Hello, World!", nil
	}
}
