package db

import (
	"fmt"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
	"github.com/toby3d/MyPackBot/internal/models"
)

// DeleteSticker just remove specified sticker key from database.
func DeleteSticker(userID int, setName, fileID string) (bool, error) {
	log.Ln("Trying to remove", fileID, "sticker from", userID, "user")
	if setName == "" {
		setName = models.SetUploaded
	}

	err := DB.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(
			fmt.Sprint("user:", userID, ":set:", setName, ":sticker:", fileID),
		)
		return err
	})

	if err == buntdb.ErrNotFound {
		log.Ln(userID, "not found, create new one")
		return true, nil
	}

	return false, err
}
