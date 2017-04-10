package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"crypto"
)

// Defaults generates a set of default configuration options
func Defaults() {
	viper.SetDefault("listen", map[string]interface{}{
		"bind":		"0.0.0.0",
		"port":  9300,
		"limit": "10mb",
	})

	viper.SetDefault("log", map[string]interface{}{
		"level": log.InfoLevel,
		"json":     	true,
		"requests": 	true,
	})

	viper.SetDefault("correlation", map[string]interface{}{
		"enable": true,
		"header": "X-Request-Identifier",
	})

	viper.SetDefault("local", map[string]interface{}{
		"algorithm": "SHA256",
		"skew":      1000,
	})

	viper.SetDefault("local.db", map[string]interface{}{
		"path":  "http://localhost:9100/v1/properties",
		"signal": "SIGHUP",
		"propsd": true,
		"prefix": "",
	})

	viper.SetDefault("service", map[string]interface{}{
		"port":     9301,
		"hostname": "127.0.0.1",
		"limit":    "10mb",
		"protocol": "http://",
	})

	viper.SetDefault("metrics", map[string]interface{}{
		"enabled": true,
	})

	viper.SetDefault("metrics.client", map[string]interface{}{
		"host":   "localhost",
		"port":   8125,
		"prefix": "turnstile.",
	})
}

// Local fixes an issue where config files make nested values in the same
// map disappear.
func Local() map[string]interface{} {
	conf := make(map[string]interface{})
	conf["algorithm"] = viper.GetString("local.algorithm")
	conf["skew"] = viper.GetInt("local.skew")

	db := make(map[string]interface{})
	db["path"] = viper.GetString("local.db.path")
	db["signal"] = viper.GetString("local.db.signal")
	db["propsd"] = viper.GetBool("local.db.propsd")
	db["prefix"] = viper.GetString("local.db.prefix")
	conf["db"] = db

	return conf
}

// Metrics fixes an issue where config files make nested values in the same
// map disappear.
func Metrics() map[string]interface{} {
	conf := make(map[string]interface{})
	conf["enabled"] = viper.GetBool("metrics.enabled")

	client := make(map[string]interface{})
	client["host"] = viper.GetString("metrics.client.host")
	client["port"] = viper.GetInt("metrics.client.port")
	client["prefix"] = viper.GetString("metrics.client.prefix")
	conf["client"] = client

	return conf
}

var SUPPORTED_ALGORITHMS = map[string]interface{}{
	"SHA256": crypto.SHA256,
	"SHA512": crypto.SHA512,
}

var SUPPORT_ALGORITHMS_LOOKUP = map[crypto.Hash]string {
	crypto.SHA256: "SHA256",
	crypto.SHA512: "SHA512",
}