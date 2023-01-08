package endpoints

import (
	"GoServerTemplate/api/paths"
	"GoServerTemplate/api/resources"
	transport "github.com/go-kit/kit/transport/http"
	"net/http"
)

type Definition struct {
	Name         string                      // documentation
	Method       string                      // used to init handler
	Path         string                      // used to init handler
	EndpointName string                      // used to init handler
	Decoder      transport.DecodeRequestFunc // used to init handler
	Resource     any                         // documentation
}

var Definitions = []Definition{
	{
		Name:         "Get Foo",
		Method:       http.MethodGet,
		Path:         paths.Foo,
		EndpointName: GetFooEndpointName,
		//Decoder:  GetFooDecoder(),
		Resource: resources.Foo{},
	},
}
