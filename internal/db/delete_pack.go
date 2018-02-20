package db

import (
	"fmt"
	"strings"

	log "github.com/kirillDanshin/dlog"
	"github.com/tidwall/buntdb"
	"github.com/toby3d/MyPackBot/internal/models"
)

// DeletePack remove all keys for UserID which contains input SetName
func DeletePack(userID int, setName string) (bool, error) {
	log.Ln("Trying to remove all", setName, "sticker from", userID, "user")
	if setName == "" {
		setName = models.SetUploaded
	}

	var fileIDs []string
	err := DB.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys(
			fmt.Sprint("user:", userID, ":set:", setName, ":*"),
			func(key, val string) bool {
				keys := strings.Split(key, ":")
				fileIDs = append(fileIDs, keys[5])
				return true
			},
		)
	})

	if len(fileIDs) == 0 {
		return true, nil
	}

	for _, fileID := range fileIDs {
		var notExist bool
		notExist, err = DeleteSticker(userID, setName, fileID)
		if err != nil {
			return notExist, err
		}
	}

	switch err {
	case buntdb.ErrNotFound:
		log.Ln(userID, "not found")
		return true, nil
	}

	return false, err
}
