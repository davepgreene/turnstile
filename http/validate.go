package http

import (
	"github.com/davepgreene/turnstile/utils"
	"github.com/davepgreene/turnstile/errors"
	"net/http"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"math"
	"time"
	"github.com/spf13/viper"
)

var REQUIRED_HEADERS = [3]string{"authorization", "date", "digest"}

func validate(skew int) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		metadata := map[string]interface{}{
			"identifier": GetCorrelationId(rw),
		}

		for _, header := range REQUIRED_HEADERS {
			if r.Header.Get(header) == "" {
				errors.ErrorWriter(errors.NewRequestError(fmt.Sprintf("Missing header: %s", header), metadata), rw)
				return
			}
		}

		// Go promotes the `Host` header to a field on the request so we
		// have to validate it there
		if r.Host != fmt.Sprintf("%s:%d", viper.GetString("listen.bind"), viper.GetInt("listen.port")) {
			errors.ErrorWriter(errors.NewRequestError("Host header mismatch", metadata), rw)
			return
		}

		date, err := utils.MsToTime(r.Header.Get("date"))
		if err != nil {
			errors.ErrorWriter(errors.NewRequestError("Invalid date header", metadata), rw)
			return
		}

		// Verify that the request date is close to $NOW
		now := time.Now().Unix()
		drift := now - date.Unix()

		log.Debugf("Date skew %d ms", drift)

		if math.Abs(float64(now - drift)) > float64(skew) {
			errors.ErrorWriter(errors.NewRequestError("Request date skew is too large", metadata), rw)
			return
		}

		next(rw, r)
	}
}
