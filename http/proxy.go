package http

import (
	"net/http"
	"github.com/vulcand/oxy/stream"
	prox "github.com/davepgreene/turnstile/proxy"
)

func proxy(url string) http.Handler {
	proxy := prox.New(url)

	stream, err := stream.New(proxy)
	if err != nil {
		panic(err)
	}

	return stream
}
