package model

import (
	"context"

	tg "gitlab.com/toby3d/telegram"
)

type (
	User struct {
		ID           int    `json:"id"`
		LanguageCode string `json:"language_code"`
		CreatedAt    int64  `json:"created_at"`
		UpdatedAt    int64  `json:"updated_at"`
		LastSeen     int64  `json:"last_seen"`
	}

	Users []*User

	Sticker struct {
		ID         string `json:"id"`
		Emoji      string `json:"emoji"`
		SetName    string `json:"set_name"`
		IsAnimated bool   `json:"is_animated"`
		CreatedAt  int64  `json:"created_at"`
	}

	Stickers []*Sticker

	UserSticker struct {
		CreatedAt int64  `json:"created_at"`
		Emojis    string `json:"emojis"`
		SetName   string `json:"set_name"`
		StickerID string `json:"sticker_id"`
		UserID    int    `json:"user_id"`
	}

	UserStickers []*UserSticker

	UpdateFunc func(context.Context, *tg.Update) error
)

func (users Users) GetByID(id int) *User {
	for i := range users {
		if users[i].ID != id {
			continue
		}
		return users[i]
	}
	return nil
}

func (stickers Stickers) GetByID(id string) *Sticker {
	for i := range stickers {
		if stickers[i].ID != id {
			continue
		}
		return stickers[i]
	}
	return nil
}

func (stickers Stickers) GetSet(name string) (Stickers, int) {
	set := make(Stickers, 0)
	for i := range stickers {
		if stickers[i].SetName != name {
			continue
		}
		set = append(set, stickers[i])
	}
	return set, len(set)
}

func (userStickers UserStickers) GetByID(uid int, sid string) *UserSticker {
	for i := range userStickers {
		if userStickers[i].UserID != uid || userStickers[i].StickerID != sid {
			continue
		}
		return userStickers[i]
	}
	return nil
}
