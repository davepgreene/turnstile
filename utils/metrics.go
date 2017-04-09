package utils

import (
	"sync"

	"fmt"

	"github.com/DataDog/datadog-go/statsd"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

type singleton struct {
}

var instance *statsd.Client
var once sync.Once

// Metrics creates a singleton client for submitting metrics to DogStatsD
func Metrics() *statsd.Client {
	once.Do(func() {
		conn := fmt.Sprintf("%s:%d", viper.GetString("metrics.client.host"), viper.GetInt("metrics.client.port"))
		instance, err := statsd.New(conn)
		if err != nil {
			log.Info("Unable to initialize statsd client.")
		}
		instance.Namespace = viper.GetString("metrics.client.prefix")

	})
	return instance
}
