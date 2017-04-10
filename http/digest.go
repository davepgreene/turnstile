package http

import (
	"net/http"
	"github.com/davepgreene/turnstile/errors"
	"github.com/davepgreene/turnstile/config"
	"crypto"
	"io/ioutil"
	"bytes"
	"encoding/base64"
	"strings"
	"fmt"
	"crypto/hmac"
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
		digestHeader := r.Header.Get("digest")

		h := algorithm.New()
		h.Write(body.Bytes())

		strAlg := strings.ToUpper(config.SUPPORT_ALGORITHMS_LOOKUP[algorithm])
		signature := fmt.Sprintf("%s=%s", strAlg, base64.StdEncoding.EncodeToString(h.Sum(nil)))

		if hmac.Equal([]byte(signature), []byte(digestHeader)) == false {
			errors.ErrorWriter(errors.NewAuthorizationError("Digest header does not match request body", metadata), rw)
			return
		}

		r.Body = deferredBody
		next(rw, r)
	}
}