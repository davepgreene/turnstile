package http

import (
	"github.com/davepgreene/turnstile/utils"
	"net/http"
	log "github.com/Sirupsen/logrus"
)

func GetCorrelationId(rw http.ResponseWriter) string {
	return rw.Header().Get("identifier")
}

func correlationMiddleware(correlation map[string]interface{}) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	header, ok := correlation["header"]
	if ok == false {
		panic("Missing required parameter `header`!")
	}

	log.Infof("Using Correlation-Identifier %s", header)

	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		id := r.Header.Get(header.(string))
		// If no header is set, create a correlation id
		if id == "" {
			id = utils.CreateCorrelationID()
			rw.Header().Add(header.(string), id)
			rw.Header().Add("identifier", id)

			log.WithFields(log.Fields{
				"identifier": id,
			}).Debugf("Setting Correlation-Identifier header %s:%s", header, id)

			next(rw, r)
		} else {
			rw.Header().Add("identifier", id)

			log.WithFields(log.Fields{
				"identifier": id,
			}).Debugf("Found Correlation-Identifier header %s:%s", header, id)

			next(rw, r)
		}
	}
}
