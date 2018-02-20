package db

import (
	"fmt"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
	"github.com/toby3d/MyPackBot/internal/models"
)

// AddSticker add sticker FileID, Emoji and SetName meta for UserID
func AddSticker(userID int, setName, fileID, emoji string) (bool, error) {
	log.Ln("Trying to add", fileID, "sticker from", userID, "user")
	if setName == "" {
		setName = models.SetUploaded
	}

	var exists bool
	err := DB.Update(func(tx *buntdb.Tx) error {
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
