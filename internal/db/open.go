package db

import (
	"path/filepath"

	bolt "github.com/etcd-io/bbolt"
)

// Open just open connection to database for work
func Open(path string) (*DB, error) {
	var (
		db  DB
		err error
	)
	if db.path, err = filepath.Abs(path); err != nil {
		return nil, err
	}

	if db.db, err = bolt.Open(db.path, 0600, nil); err != nil {
		return nil, err
	}

	return &db, nil
}
