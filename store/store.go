package store

import (
	"github.com/davepgreene/turnstile/errors"
)

type Store interface {
	Load()
	Path() string
	Lookup(identity string, identifier string) ([]string, errors.HTTPWrappedError)
}

type AbstractStore struct {
	keys map[string][]string
	path string
}

func newAbstractStore(path string) AbstractStore {
	return AbstractStore{
		keys: make(map[string][]string),
		path: path,
	}
}

// Lookup does a basic
func (s *AbstractStore) Lookup(identity string, identifier string) ([]string, errors.HTTPWrappedError) {
	if val, ok := s.keys[identity]; ok {
		return val, nil
	}
	metadata := map[string]interface{}{
		"identifier": identifier,
	}

	return nil, errors.NewAuthorizationError("Invalid authentication factors", metadata)
}
