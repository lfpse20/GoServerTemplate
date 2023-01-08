package decoders

import (
	"context"
	transport "github.com/go-kit/kit/transport/http"
	"net/http"
)

func GetFooDecoder() transport.DecodeRequestFunc {
	return func(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
		return nil, nil
	}
}
