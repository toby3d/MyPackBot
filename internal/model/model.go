package model

import (
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/message"
)

type (
	User struct {
		ID        int   `json:"id"`
		CreatedAt int64 `json:"created_at"`
		UpdatedAt int64 `json:"updated_at"`

		LanguageCode string `json:"language_code"`
		LastSeen     int64  `json:"last_seen"`
	}

	Users []*User

	Sticker struct {
		ID        string `json:"id"`
		CreatedAt int64  `json:"created_at"`
		UpdatedAt int64  `json:"updated_at"`

		Width      int    `json:"width"`
		Height     int    `json:"height"`
		IsAnimated bool   `json:"is_animated"`
		SetName    string `json:"set_name"`
		Emoji      string `json:"emoji"`
	}

	Stickers []*Sticker

	UserSticker struct {
		StickerID string `json:"sticker_id"`
		UserID    int    `json:"user_id"`
		CreatedAt int64  `json:"created_at"`

		SetName string `json:"set_name"`
		Emojis  string `json:"emojis"`
	}

	UserStickers []*UserSticker

	UpdateFunc func(*Context) error

	Context struct {
		User    *User
		Sticker *Sticker
		Printer *message.Printer
		*tg.Update
	}
)

func (users Users) GetByID(id int) *User {
	var u *User

	for i := range users {
		if users[i].ID != id {
			continue
		}

		u = users[i]

		break
	}

	return u
}

func (stickers Stickers) GetByID(id string) *Sticker {
	var s *Sticker

	for i := range stickers {
		if stickers[i].ID != id {
			continue
		}

		s = stickers[i]
	}

	return s
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
	var us *UserSticker

	for i := range userStickers {
		if userStickers[i].UserID != uid || userStickers[i].StickerID != sid {
			continue
		}

		us = userStickers[i]
	}

	return us
}
