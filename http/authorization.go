package http

import (
	"net/http"
	"strings"
	"github.com/davepgreene/turnstile/errors"
	"encoding/base64"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

const AUTHN_PROTOCOL = "Rapid7-HMAC-V1-SHA256"

func authorization(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	identifier := GetCorrelationId(rw)
	metadata := map[string]interface{}{
		"identifier": identifier,
	}

	parts := strings.Split(r.Header.Get("authorization"), " ")

	if len(parts) != 2 {
		errors.ErrorWriter(errors.NewAuthorizationError("Invalid Authorization header", metadata), rw)
		return
	}

	if parts[0] != AUTHN_PROTOCOL {
		errors.ErrorWriter(errors.NewAuthorizationError(fmt.Sprintf("Invalid authentication protocol %s", parts[0]), metadata), rw)
		return
	}

	log.WithFields(log.Fields{
		"identifier": identifier,
	}).Debugf("Using Authorization Scheme: %s", parts[0])

	authParamBuf, _ := base64.StdEncoding.DecodeString(parts[1])
	parameters := strings.Split(string(authParamBuf), ":")

	if len(parameters) != 2 {
		errors.ErrorWriter(errors.NewAuthorizationError("Invalid authentication parameters", metadata), rw)
		return
	}

	r.Header.Set("identity", parameters[0])
	r.Header.Set("signature", parameters[1])

	next(rw, r)
}