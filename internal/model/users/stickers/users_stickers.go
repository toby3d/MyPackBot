package stickers

import "gitlab.com/toby3d/mypackbot/internal/model"

type (
	ReadWriter interface {
		Reader
		Writer
	}

	Reader interface {
		Get(up *model.UserSticker) *model.Sticker
		GetList(offset int, limit int, filter *model.UserSticker) (model.Stickers, int, error)
	}

	Writer interface {
		Add(up *model.UserSticker) error
		AddSet(uid int, setName string) error
		Update(up *model.UserSticker) error
		Remove(up *model.UserSticker) error
		RemoveSet(uid int, setName string) error
	}
)
