package stickers

import "gitlab.com/toby3d/mypackbot/internal/model"

type (
	Manager interface {
		Reader
		Writer
		ReadWriter
	}

	ReadWriter interface {
		GetOrCreate(s *model.Sticker) (*model.Sticker, error)
	}

	Reader interface {
		Get(id string) *model.Sticker
		GetSet(name string) model.Stickers
		GetList(offset int, limit int, filter *model.Sticker) (model.Stickers, int, error)
	}

	Writer interface {
		Create(s *model.Sticker) error
		Remove(id string) error
		Update(s *model.Sticker) error
	}
)
