package model

import (
	"strings"

	"gitlab.com/toby3d/mypackbot/internal/common"
)

type (
	Model struct {
		ID        uint64 `json:"id"`
		CreatedAt int64  `json:"created_at"`
		UpdatedAt int64  `json:"updated_at"`
	}

	User struct {
		Model
		UserID       int64  `json:"user_id"`
		LanguageCode string `json:"language_code"`
		LastSeen     int64  `json:"last_seen"`
	}

	Users []*User

	Sticker struct {
		Model
		FileID     string `json:"file_id"`
		Width      int    `json:"width"`
		Height     int    `json:"height"`
		IsAnimated bool   `json:"is_animated"`
		SetName    string `json:"set_name"`
		Emoji      string `json:"emoji"`
	}

	Stickers []*Sticker

	UserSticker struct {
		Model
		StickerID uint64 `json:"sticker_id"`
		UserID    uint64 `json:"user_id"`
		Query     string `json:"query"`
	}

	UserStickers []*UserSticker

	Photo struct {
		Model
		FileID string `json:"file_id"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	}

	Photos []*Photo

	UserPhoto struct {
		Model
		PhotoID uint64 `json:"photo_id"`
		UserID  uint64 `json:"user_id"`
		Query   string `json:"query"`
	}

	UserPhotos []*UserPhoto
)

func (s *Sticker) InSet() bool {
	return s.SetName != "" && !strings.EqualFold(s.SetName, common.SetNameUploaded)
}
