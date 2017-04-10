package proxy

import (
	"net/url"
	"net/http/httputil"
	//"github.com/DataDog/datadog-go/statsd"
	"sync"
	"net/http"
)

type Proxy struct {
	target *url.URL
	proxy *httputil.ReverseProxy
}

var instance *Proxy
var once sync.Once

func New(target string) *Proxy {
	once.Do(func() {
		url, err := url.Parse(target)
		if err != nil {
			panic(err)
		}
		instance = &Proxy{
			target: url,
			proxy: httputil.NewSingleHostReverseProxy(url),
		}
	})
	return instance
}

func (p *Proxy) Handle(rw http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(rw, r)
}
