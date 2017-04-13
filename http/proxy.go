package http

import (
	"net/http"
	"github.com/vulcand/oxy/stream"
	"github.com/davepgreene/turnstile/proxy"
)

func proxy(url string) http.Handler {
	prox := proxy.New(url)

	stream, err := stream.New(prox)
	if err != nil {
		panic(err)
	}

	return stream
}
