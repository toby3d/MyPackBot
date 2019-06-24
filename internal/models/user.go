//go:generate ffjson $GOFILE
package models

type (
	User struct {
		ID           int    `json:"id"`
		LanguageCode string `json:"language_code"`
		AutoSaving   bool   `json:"auto_saving"`
		StartedAt    int64  `json:"started_at"`
	}

	UsersStickers struct {
		UserID    int    `json:"user_id"`
		StickerID string `json:"sticker_id"`
	}
)
