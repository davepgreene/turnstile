package proxy

import (
	"net/url"
	//"github.com/DataDog/datadog-go/statsd"
	"net/http"
	//log "github.com/Sirupsen/logrus"
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/testutils"
)

type Proxy struct {
	url *url.URL
	fwd *forward.Forwarder
}

func New(url string) *Proxy {
	fwd, _ := forward.New(forward.RoundTripper(&transport{http.DefaultTransport}))

	return &Proxy{
		url: testutils.ParseURI(url),
		fwd: fwd,
	}
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// let us forward this request to another server
	r.URL = p.url
	p.fwd.ServeHTTP(rw, r)
}
