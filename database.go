package main

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/tidwall/buntdb"         // Redis-like database
)

const (
	stateNone       = "none"
	stateAddSticker = "addSticker"
	stateAddPack    = "addPack"
	stateDelete     = "del"
	stateReset      = "reset"
)

var db *buntdb.DB

func dbInit() {
	log.Ln("Open database file...")
	var err error
	db, err = buntdb.Open("stickers.db")
	errCheck(err)

	log.Ln("Creating user_stickers index...")
	err = db.CreateIndex(
		"user_stickers",    // name
		"user:*:sticker:*", // pattern
		buntdb.IndexString, // options
	)
	errCheck(err)

	select {}
}

func dbChangeUserState(userID int, state string) error {
	log.Ln("Trying to change", userID, "state to", state)
	return db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(
			fmt.Sprint("user:", userID, ":state"), // key
			state, // val
			nil,   // options
		)
		return err
	})
}

func dbGetUserState(userID int) (string, error) {
	log.Ln("Trying to get", userID, "state")
	var state string
	err := db.View(func(tx *buntdb.Tx) error {
		var err error
		state, err = tx.Get(fmt.Sprint("user:", userID, ":state"))
		return err
	})

	switch err {
	case buntdb.ErrNotFound:
		log.Ln(userID, "not found, create new one")
		if err := dbChangeUserState(userID, stateNone); err != nil {
			return state, err
		}
	}

	return state, err
}

func dbAddSticker(userID int, fileID, emoji string) (bool, error) {
	log.Ln("Trying to add", fileID, "sticker from", userID, "user")
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

func dbDeleteSticker(userID int, fileID string) (bool, error) {
	log.Ln("Trying to remove", fileID, "sticker from", userID, "user")
	err := db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(fmt.Sprint("user:", userID, ":sticker:", fileID))
		return err
	})

	switch err {
	case buntdb.ErrNotFound:
		log.Ln(userID, "not found, create new one")
		return true, nil
	}

	return false, err
}

func dbResetUserStickers(userID int) error {
	log.Ln("Trying reset all stickers of", userID, "user")
	return db.Update(func(tx *buntdb.Tx) error {
		var keys []string

		err := tx.Ascend(
			"user_stickers", // index
			func(key, val string) bool { // iterator
				subKeys := strings.Split(key, ":")
				if subKeys[1] == strconv.Itoa(userID) {
					keys = append(keys, key)
				}
				return true
			},
		)
		if err != nil {
			return err
		}

		for i := range keys {
			_, err = tx.Delete(keys[i])
			if err != nil {
				break
			}
		}

		return err
	})
}

func dbGetUserStickers(userID, offset int, emoji string) ([]string, error) {
	log.Ln("Trying to get", offset, "page of", userID, "stickers")
	count := 0
	var stickers []string
	err := db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend(
			"user_stickers", // index
			func(key, val string) bool { // iterator
				subKeys := strings.Split(key, ":")
				if subKeys[1] != strconv.Itoa(userID) {
					return true
				}

				count++

				if len(stickers) >= 50 {
					return false
				}

				if emoji != "" && val != emoji {
					return true

				}

				if offset >= 1 && count <= ((offset-1)*50) {
					return true
				}

				stickers = append(stickers, subKeys[3])
				return true
			},
		)
	})

	switch err {
	case buntdb.ErrNotFound:
		log.Ln("Not found stickers")
		return nil, nil
	}

	return stickers, err
}
