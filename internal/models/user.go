//go:generate ffjson $GOFILE
package models

type (
	User struct {
		ID       int    `json:"id"`
		Language string `json:"language"`
	}

	UsersStickers struct {
		UserID    int    `json:"user_id"`
		StickerID string `json:"sticker_id"`
	}
)
