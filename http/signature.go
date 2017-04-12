package http

import (
	"net/http"
	"github.com/davepgreene/turnstile/store"
	"github.com/davepgreene/turnstile/errors"
	"strconv"
	log "github.com/Sirupsen/logrus"
	"crypto"
	"crypto/hmac"
	"strings"
	"fmt"
	"encoding/base64"
)

func signature(db store.Store, algorithm crypto.Hash) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		identifier := GetCorrelationId(rw)
		metadata := map[string]interface{}{
			"identifier": identifier,
		}
		identity := r.Header.Get("identity")

		keys, err := db.Lookup(identity, identifier)
		if err != nil {
			errors.ErrorWriter(err, rw)
			return
		}

		requestSignature := r.Header.Get("signature")

		log.WithFields(log.Fields{
			"identifier": identifier,
		}).Debugf("Request signature: %s", requestSignature)

		for _, key := range keys {
			// Create and validate the signature
			signature, err := generateSignature(algorithm, key, r, identifier)

			if err != nil {
				errors.ErrorWriter(err, rw)
				return
			}

			// return early and invoke the next middleware if the key is valid
			if hmac.Equal([]byte(signature), []byte(requestSignature)) == true {
				log.WithFields(log.Fields{
					"identifier": identifier,
				}).Debug("Authenticated. Forwarding request")
				next(rw, r)
				return
			}
		}

		errors.ErrorWriter(errors.NewAuthorizationError("Invalid authentication factors", metadata), rw)
		return
	}
}

func generateSignature(algorithm crypto.Hash, secret string, r *http.Request, identifier string) (string, errors.HTTPWrappedError) {
	fields := log.Fields{
		"identifier": identifier,
	}

	// Method
	method := strings.ToUpper(r.Method)
	log.WithFields(fields).Debugf("Using method: %s", method)

	// URI
	uri := r.RequestURI
	log.WithFields(fields).Debugf("Using URI: %s", uri)

	// Host
	host := r.Host
	log.WithFields(fields).Debugf("Using host: %s", host)

	// Date
	dateMsInt, _ := strconv.ParseInt(r.Header.Get("date"), 10, 64)
	log.WithFields(fields).Debugf("Using date: %d", dateMsInt)


	// Identity
	identity := r.Header.Get("identity")
	log.WithFields(fields).Debugf("Using identity: %s", identity)

	// Digest
	digest := r.Header.Get("digest")
	log.WithFields(fields).Debugf("Using digest: %s", digest)

	mac := hmac.New(algorithm.New, []byte(secret))

	mac.Write([]byte(fmt.Sprintf("%s %s\n", method, uri)))
	mac.Write([]byte(fmt.Sprintf("%s\n", host)))
	mac.Write([]byte(fmt.Sprintf("%d\n", dateMsInt)))
	mac.Write([]byte(fmt.Sprintf("%s\n", identity)))
	mac.Write([]byte(fmt.Sprintf("%s\n", digest)))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}
