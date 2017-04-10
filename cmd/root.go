package cmd

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/davepgreene/turnstile/http"
	"github.com/davepgreene/turnstile/errors"
	"github.com/davepgreene/turnstile/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/davepgreene/turnstile/config"
)

var cfgFile string
var verbose bool

// TurnstileCmd represents the base command when called without any subcommands
var TurnstileCmd = &cobra.Command{
	Use:   "turnstile",
	Short: "An HTTP proxy service to add role-based access control to any REST-like API",
	Long: `Turnstile is a minimalistic middleware server that
	applies consistent access control policies to any HTTP
	service with path-based routing. The goal of the
	project is to support pluggable providers for authentication,
	authorization-policy, rate-limiting, and logging/reporting.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		err := initializeConfig()
		initializeLog()
		if err != nil {
			return err
		}

		// Prevent boot if we aren't using a supported algorithm
		confAlg := viper.GetString("local.algorithm")
		if val, ok := config.SUPPORTED_ALGORITHMS[strings.ToUpper(confAlg)]; ok {
			// Set the actual crypto algorithm instead of a string representation
			viper.Set("local.algorithm", val)
			return boot()
		}

		supported := strings.Join(utils.MapKeys(config.SUPPORTED_ALGORITHMS), ", ")
		message := fmt.Sprintf("Turnstile currently supports the following encryption algorithms: %s. You specified %s.", supported, confAlg)
		panic(errors.NewUnsupportedAlgorithmError(message, strings.ToUpper(confAlg)))
	},
}

func boot() error {
	router := http.Handler()
	log.Error(router)
	return router
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the TokendCmd.
func Execute() {
	if err := TurnstileCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	TurnstileCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	TurnstileCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose level logging")
	validConfigFilenames := []string{"json"}
	TurnstileCmd.PersistentFlags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)
}

func initializeLog() {
	log.RegisterExitHandler(func() {
		log.Info("Shutting down")
	})

	// Set logging options based on config
	if lvl, err := log.ParseLevel(viper.GetString("log.level")); err == nil {
		log.SetLevel(lvl)
	} else {
		log.Info("Unable to parse log level in settings. Defaulting to INFO")
	}

	// If using verbose mode, log at debug level
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	if viper.GetBool("log.json") {
		log.SetFormatter(&log.JSONFormatter{})
	}

	if cfgFile != "" {
		log.WithFields(log.Fields{
			"file": viper.ConfigFileUsed(),
		}).Info("Loaded config file")
	}

}

func initializeConfig(subCmdVs ...*cobra.Command) error {
	config.Defaults()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	viper.AutomaticEnv() // read in environment variables that match
	log.Info(viper.Get("local.db.path"))

	return nil
}
