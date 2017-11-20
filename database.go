package main

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/tidwall/buntdb"         // Redis-like database
)

const (
	stateNone     = "none"
	stateAdding   = "add"
	stateDeleting = "del"
)

var db *buntdb.DB

func dbInit() {
	log.Ln("[dbInit] Open database file...")
	var err error
	db, err = buntdb.Open("bot.db")
	errCheck(err)

	err = db.CreateIndex(
		"user_stickers",    // name
		"user:*:sticker:*", // pattern
		buntdb.IndexString, // options
	)
	errCheck(err)

	select {}
}

func dbChangeUserState(userID int, state string) (string, error) {
	var prevState string
	err := db.Update(func(tx *buntdb.Tx) error {
		var err error
		prevState, _, err = tx.Set(
			fmt.Sprint("user:", userID, ":state"), // key
			state, // val
			nil,   // options
		)
		return err
	})

	return prevState, err
}

func dbGetUserState(userID int) (string, error) {
	var state string
	err := db.View(func(tx *buntdb.Tx) error {
		var err error
		state, err = tx.Get(
			fmt.Sprint("user:", userID, ":state"), // key
		)
		return err
	})

	switch err {
	case buntdb.ErrNotFound:
		state, err = dbChangeUserState(userID, stateNone)
	}

	return state, err
}

func dbAddSticker(userID int, fileID, emoji string) (bool, error) {
	var exists bool
	err := db.Update(func(tx *buntdb.Tx) error {
		var err error
		_, exists, err = tx.Set(
			fmt.Sprint("user:", userID, ":sticker:", fileID), // key
			emoji, // value
			nil,   // options
		)

		if err == buntdb.ErrIndexExists {
			exists = true
			return nil
		}

		return err
	})

	return exists, err
}

func dbDeleteSticker(userID int, fileID string) error {
	return db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(
			fmt.Sprint("user:", userID, ":sticker:", fileID), // key
		)
		return err
	})
}

func dbGetUserStickers(userID int, emoji string) ([]string, error) {
	var stickers []string
	err := db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend(
			"user_stickers", // index
			func(key, val string) bool { // iterator
				log.Ln(key, "=", val)

				subKeys := strings.Split(key, ":")
				if subKeys[1] != strconv.Itoa(userID) {
					return true
				}

				if emoji != "" {
					log.Ln("["+emoji+"]", "?=", "["+val+"]")
					if val != emoji {
						return true
					}
				}

				stickers = append(stickers, subKeys[3])
				return true
			},
		)
	})

	switch err {
	case buntdb.ErrNotFound:
		return nil, nil
	}

	return stickers, err
}
