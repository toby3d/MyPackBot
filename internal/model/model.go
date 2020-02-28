package model

import (
	"strings"

	"gitlab.com/toby3d/mypackbot/internal/common"
	tg "gitlab.com/toby3d/telegram"
)

type (
	User struct {
		ID           int `boltholdKey:"ID"`
		CreatedAt    int64
		UpdatedAt    int64
		LanguageCode string
		LastSeen     int64
	}

	Users []*User

	Sticker struct {
		ID         string `boltholdKey:"ID"`
		CreatedAt  int64
		UpdatedAt  int64
		FileID     string
		Width      int
		Height     int
		IsAnimated bool
		SetName    string
		Emoji      string
	}

	Stickers []*Sticker

	Photo struct {
		ID        string `boltholdKey:"ID"`
		CreatedAt int64
		UpdatedAt int64
		FileID    string
		Width     int
		Height    int
	}

	Photos []*Photo

	UserSticker struct {
		ID        uint64 `boltholdKey:"ID"`
		CreatedAt int64
		UpdatedAt int64
		UserID    int
		StickerID string
		Query     string
	}

	UserStickers []*UserSticker

	UserPhoto struct {
		ID        uint64 `boltholdKey:"ID"`
		CreatedAt int64
		UpdatedAt int64
		UserID    int
		PhotoID   string
		Query     string
	}

	UserPhotos []*UserPhoto

	InlineResult interface {
		GetType() string
		GetID() string
		GetFileID() string
		GetUpdatedAt() int64
	}
)

func (s *Sticker) InSet() bool {
	return s.SetName != "" && !strings.EqualFold(s.SetName, common.SetNameUploaded)
}

func (Sticker) GetType() string { return tg.TypeSticker }

func (s Sticker) GetID() string { return s.ID }

func (s Sticker) GetFileID() string { return s.FileID }

func (s Sticker) GetUpdatedAt() int64 { return s.UpdatedAt }

func (Photo) GetType() string { return tg.TypePhoto }

func (p Photo) GetID() string { return p.ID }

func (p Photo) GetFileID() string { return p.FileID }

func (p Photo) GetUpdatedAt() int64 { return p.UpdatedAt }
