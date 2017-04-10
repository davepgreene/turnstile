package http

import (
	"net/http"
	"fmt"
	"time"
	log "github.com/Sirupsen/logrus"
	"github.com/davepgreene/turnstile/store"
	"github.com/davepgreene/turnstile/utils"
	"github.com/gorilla/mux"
	"github.com/meatballhat/negroni-logrus"
	"github.com/spf13/viper"
	"github.com/thoas/stats"
	"github.com/urfave/negroni"
	"gopkg.in/tylerb/graceful.v1"
	"crypto"
	"github.com/davepgreene/turnstile/config"
	"github.com/davepgreene/turnstile/proxy"
)

const (
	// MaxRequestSize is the maximum accepted request size. This is to prevent
	// a denial of service attack where no Content-Length is provided and the server
	// is fed ever more data until it exhausts memory.
	MaxRequestSize = 32 * 1024 * 1024
)

// Handler returns an http.Handler for the API.
func Handler() error {
	db := db()
	r := mux.NewRouter()
	statsMiddleware := stats.New()
	r.HandleFunc("/stats", newAdminHandler(statsMiddleware).ServeHTTP)

	// Add middleware handlers
	n := negroni.New()

	// Add recovery handler that logs
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	n.Use(recovery)

	if viper.GetBool("log.requests") {
		n.Use(negronilogrus.NewCustomMiddleware(utils.GetLogLevel(), utils.GetLogFormatter(), "requests"))
	}

	n.Use(statsMiddleware)

	// Collect some metrics about incoming and active requests
	n.Use(negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		utils.Metrics().Incr("request.incoming", nil, 1)
		utils.Metrics().Gauge("request.active", 1, nil, 1)
		next(rw, r)
	}))

	// Correlation id middleware
	if viper.GetBool("correlation.enable") {
		n.Use(negroni.HandlerFunc(correlationMiddleware(viper.GetStringMap("correlation"))))
	}

	// Validate headers
	n.Use(negroni.HandlerFunc(validate(viper.GetInt("local.skew"))))

	// Validate the digest header
	algorithm := viper.Get("local.algorithm")
	n.Use(negroni.HandlerFunc(digest(algorithm.(crypto.Hash))))
	n.Use(negroni.HandlerFunc(authorization))
	n.Use(negroni.HandlerFunc(signature(db, algorithm.(crypto.Hash))))

	// All checks passed, forward the request
	forwardConn := fmt.Sprintf("%s%s:%d", viper.GetString("service.protocol"), viper.GetString("service.hostname"), viper.GetInt("service.port"))
	p := proxy.New(forwardConn)
	r.PathPrefix("/").HandlerFunc(p.Handle)

	n.UseHandler(r)

	// Set up connection
	conn := fmt.Sprintf("%s:%d", viper.GetString("listen.bind"), viper.GetInt("listen.port"))
	log.Infof("Listening on %s", conn)

	// Bombs away!
	return server(conn, n).ListenAndServe()
}

func server(conn string, handler http.Handler) *graceful.Server {
	return &graceful.Server{
		Timeout: 10 * time.Second,
		Server: &http.Server{
			Addr:    conn,
			Handler: handler,
		},
	}
}

func db() store.Store {
	local := config.Local()["db"].(map[string]interface{})
	// Create DB
	storeType := "file"
	if local["propsd"].(bool) == true {
		storeType = "propsd"
	}

	db, err := store.CreateStore(storeType, local)
	if err != nil {
		// If we have an invalid store, dump out before we start the server
		panic(err)
	}
	log.Infof("Using authentication controller with database: %s", db.Path())
	return db
}
