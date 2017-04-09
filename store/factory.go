package store

import (
	"fmt"
	"strings"
	log "github.com/Sirupsen/logrus"
)

type StoreFactory func(conf map[string]interface{}) (Store, error)

var storeFactories = make(map[string]StoreFactory)

func register(name string, factory StoreFactory) {
	if factory == nil {
		log.Panicf("Store factory %s does not exist.", name)
	}
	_, registered := storeFactories[name]
	if registered {
		log.Errorf("Store factory %s already registered. Ignoring.", name)
	}
	storeFactories[name] = factory
}

func CreateStore(name string, conf map[string]interface{}) (Store, error) {
	storeFactory, ok := storeFactories[strings.ToLower(name)]

	if !ok {
		// Factory has not been registered.
		// Make a list of all available datastore factories for logging.
		availableStores := make([]string, len(storeFactories))
		for k := range storeFactories {
			availableStores = append(availableStores, k)
		}
		return nil, fmt.Errorf("Invalid Datastore name. Must be one of: %s", strings.Join(availableStores, ", "))
	}

	// Run the factory with the configuration.
	return storeFactory(conf)
}

func init() {
	register("propsd", NewPropsdStore)
	register("file", NewFileStore)
}
