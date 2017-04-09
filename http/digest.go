package http

import (
	"net/http"
	"github.com/davepgreene/turnstile/errors"
	"crypto"
	"io/ioutil"
	"bytes"
	"encoding/base64"
	"crypto/subtle"
)

func digest(algorithm crypto.Hash) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		metadata := map[string]interface{}{
			"identifier": GetCorrelationId(rw),
		}

		// Save a copy of this request so we can put it back in the request once we consume it.
		buf, _ := ioutil.ReadAll(r.Body)
		body := bytes.NewBuffer(buf)
		deferredBody := ioutil.NopCloser(bytes.NewBuffer(buf))

		h := algorithm.New()
		h.Write(body.Bytes())
		signature := []byte(base64.URLEncoding.EncodeToString(h.Sum(nil)))

		digestHeader := []byte(r.Header.Get("digest"))

		if val := subtle.ConstantTimeCompare(signature, digestHeader); val == 0 {
			errors.ErrorWriter(errors.NewAuthorizationError("Digest header does not match request body", metadata), rw)
			return
		}

		r.Body = deferredBody
		next(rw, r)
	}
}