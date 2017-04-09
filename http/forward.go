package http

import (
	"encoding/json"
	"net/http"

	"github.com/davepgreene/turnstile/errors"
)

func forward(rw http.ResponseWriter, r *http.Request) {
	err := errors.NewHTTPError(http.StatusBadRequest, "WTF", make(map[string]interface{}))
	if err != nil {
		errors.ErrorWriter(err, rw)
		return
	}

	rw.WriteHeader(http.StatusNotImplemented)
	m := make(map[string]interface{})
	b, _ := json.Marshal(m)

	rw.Write(b)
}
