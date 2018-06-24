package db

import (
	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
)

type DataBase struct{ *buntdb.DB }

// DB is a main object of current database connection
var DB *DataBase

// Open just open connection to database for work
func Open(path string) (*DataBase, error) {
	log.Ln("Open database file...")
	db, err := buntdb.Open(path)
	return &DataBase{db}, err
}
