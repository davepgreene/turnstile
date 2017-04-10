package store

import (
	"fmt"
	"errors"
	err "github.com/davepgreene/turnstile/errors"
	"github.com/magiconair/properties"
	"strings"
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

	if s.prefix != "" {
		// Filter prefix then set keys
		p = p.FilterStripPrefix(s.prefix)
	}
	m := p.Map()
	// Iterate the properties and coerce each value to an array of strings.
	for k, v := range m {
		s.store.keys[k] = strings.Split(v, ",")
	}
}

func (s *PropsdStore) Lookup(identity string, identifier string) ([]string, err.HTTPWrappedError) {
	return s.store.Lookup(identity, identifier)
}
