package db

import (
	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
	"github.com/toby3d/MyPackBot/internal/errors"
)

var DB *buntdb.DB

// Open just open connection to database for work
func Open(path string) {
	log.Ln("Open database file...")
	go func() {
		var err error
		DB, err = buntdb.Open(path)
		errors.Check(err)
	}()
}
