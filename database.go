package main

import (
	"fmt"
	"strings"

	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/tidwall/buntdb"         // Redis-like database
)

var db *buntdb.DB

func dbInit() {
	log.Ln("[dbInit] Open database file...")
	var err error
	db, err = buntdb.Open("bot.db")
	errCheck(err)

	select {}
}

func dbChangeUserState(userID int, state string) (string, bool, error) {
	var prevState string
	var changed bool
	err := db.Update(func(tx *buntdb.Tx) error {
		var err error
		prevState, changed, err = tx.Set(
			fmt.Sprint("user:", userID, ":state"), // key
			state, // val
			nil,   // options
		)
		return err
	})
	return prevState, changed, err
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
	return state, err
}

func dbAddSticker(userID int, fileID, emoji string) error {
	if err := db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(
			fmt.Sprint("user:", userID, ":sticker:", fileID), // key
			emoji, // value
			nil,   // options
		)
		return err
	}); err != nil {
		return err
	}

	err := dbUpdateUserStickersIndex(userID)
	return err
}

func dbDeleteSticker(userID int, fileID string) error {
	if err := db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(
			fmt.Sprint("user:", userID, ":sticker:", fileID), // key
		)
		return err
	}); err != nil {
		return err
	}

	err := dbUpdateUserStickersIndex(userID)
	return err
}

func dbUpdateUserStickersIndex(userID int) error {
	return db.CreateIndex(
		fmt.Sprint("stickers", userID),            // name
		fmt.Sprint("user:", userID, ":sticker:*"), // pattern
		buntdb.IndexString,                        // options
	)
}

func dbGetUserStickers(userID int, emoji string) ([]string, error) {
	var stickers []string
	err := db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend(
			fmt.Sprint("stickers", userID), // index
			func(key, val string) bool { // iterator
				fileID := strings.TrimPrefix(
					key, // source
					fmt.Sprint("user:", userID, ":sticker:"), // prefix
				)

				if emoji != "" {
					if val != emoji {
						return true
					}
				}

				stickers = append(stickers, fileID)
				return true
			},
		)
	})

	return stickers, err
}
