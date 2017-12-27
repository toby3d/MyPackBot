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

	setUploaded = "uploaded"

	patternUsers    = "users"
	patternUserSets = "user_sets"
)

var db *buntdb.DB

func dbInit() {
	log.Ln("Open database file...")
	var err error
	db, err = buntdb.Open("stickers.db")
	errCheck(err)

	log.Ln("Creating users index...")
	err = db.CreateIndex(
		patternUsers,       // name
		"user:*",           // pattern
		buntdb.IndexString, // options
	)
	errCheck(err)

	log.Ln("Creating user_sets index...")
	err = db.CreateIndex(
		patternUserSets,    // name
		"user:*:set:*",     // pattern
		buntdb.IndexString, // options
	)
	errCheck(err)

	select {}
}

func dbGetUsers() ([]int, error) {
	var users []int
	err := db.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys(
			"user:*:state",
			func(key, val string) bool {
				log.Ln(key, "=", val)
				subKeys := strings.Split(key, ":")
				id, err := strconv.Atoi(subKeys[1])
				if err == nil {
					users = append(users, id)
				}

				return true
			},
		)
	})
	if err == buntdb.ErrNotFound {
		return nil, nil
	}

	return users, err
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

func dbAddSticker(userID int, setName, fileID, emoji string) (bool, error) {
	log.Ln("Trying to add", fileID, "sticker from", userID, "user")
	if setName == "" {
		setName = setUploaded
	}

	var exists bool
	err := db.Update(func(tx *buntdb.Tx) error {
		var err error
		_, exists, err = tx.Set(
			fmt.Sprint("user:", userID, ":set:", setName, ":sticker:", fileID), // key
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

func dbDeleteSticker(userID int, setName, fileID string) (bool, error) {
	log.Ln("Trying to remove", fileID, "sticker from", userID, "user")
	if setName == "" {
		setName = setUploaded
	}

	err := db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(fmt.Sprint("user:", userID, ":set:", setName, ":sticker:", fileID))
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

		err := tx.AscendKeys(
			patternUserSets, // index
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

func dbGetUserStickers(userID, offset int, query string) ([]string, int, error) {
	log.Ln("Trying to get", userID, "stickers")
	var total, count int
	var stickers []string
	offset = offset * 50

	err := db.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys(
			fmt.Sprint("user:", userID, ":set:*"), // index
			func(key, val string) bool { // iterator
				log.Ln(key, "=", val)
				subKeys := strings.Split(key, ":")
				if subKeys[1] != strconv.Itoa(userID) {
					return true
				}

				total++

				if count >= 51 {
					return true
				}

				if total < offset {
					return true
				}

				if query != "" &&
					query != val {
					return true
				}

				count++
				stickers = append(stickers, subKeys[5])
				return true
			},
		)
	})

	switch {
	case err == buntdb.ErrNotFound:
		log.Ln("Not found stickers")
		return nil, total, nil
	case err != nil:
		return nil, total, err
	}

	return stickers, total, nil
}
