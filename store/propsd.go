package store

import (
	"fmt"
	"errors"
	err "github.com/davepgreene/turnstile/errors"
	"github.com/magiconair/properties"
)

type PropsdStore struct {
	store AbstractStore
	prefix string
}

func NewPropsdStore(conf map[string]interface{}) (Store, error) {
	path, ok := conf["path"].(string)

	if !ok {
		return nil, errors.New(fmt.Sprintf("%s is required for the propsd store", "path"))
	}

	prefix, ok := conf["prefix"].(string)
	if !ok {
		prefix = ""
	}

	store := &PropsdStore{
		store: newAbstractStore(path),
		prefix: prefix,
	}
	store.Load()

	return store, nil
}

func (s *PropsdStore) Path() string {
	return s.store.path
}

func (s *PropsdStore) Load() {
	p := properties.MustLoadURL(s.Path())

	if s.prefix == "" {
		// Load all properties as key: [value]/value.([]string)
	}

	// Otherwise filter prefix then set keys
}

func (s *PropsdStore) Lookup(identity string, identifier string) ([]string, err.HTTPWrappedError) {
	return s.store.Lookup(identity, identifier)
}
