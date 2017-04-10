package utils

import (
	"sync"
	"fmt"
	"github.com/davepgreene/turnstile/config"
	"github.com/DataDog/datadog-go/statsd"
	log "github.com/Sirupsen/logrus"
)

var instance *statsd.Client
var once sync.Once

// Metrics creates a singleton client for submitting metrics to DogStatsD
func Metrics() *statsd.Client {
	once.Do(func() {
		metricsClientConf := config.Metrics()["client"].(map[string]interface{})
		conn := fmt.Sprintf("%s:%d", metricsClientConf["host"], metricsClientConf["port"])
		instance, err := statsd.New(conn)
		if err != nil {
			log.Info("Unable to initialize statsd client.")
			panic(err)
		}
		instance.Namespace = metricsClientConf["prefix"].(string)

	})
	return instance
}
