package store

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"fmt"
	"errors"
	err "github.com/davepgreene/turnstile/errors"
	//log "github.com/Sirupsen/logrus"
)

type FileStore struct {
	store AbstractStore
}

func NewFileStore(conf map[string]interface{}) (Store, error) {
	path, ok := conf["path"].(string)

	if !ok {
		return nil, errors.New(fmt.Sprintf("%s is required for the file store", "path"))
	}

	store := &FileStore{
		store: newAbstractStore(path),
	}
	store.Load()

	return store, nil
}

func (s *FileStore) Path() string {
	return s.store.path
}

func (s *FileStore) Load() {
	path, err := filepath.Abs(s.store.path)
	if err != nil {
		panic(err)
	}
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var objmap map[string][]string
	err = json.Unmarshal(dat, &objmap)
	if err != nil {
		panic(err)
	}

	s.store.keys = objmap
}

func (s *FileStore) Lookup(identity string, identifier string) ([]string, err.HTTPWrappedError) {
	return s.store.Lookup(identity, identifier)
}
