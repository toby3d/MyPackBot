package stickers

import "gitlab.com/toby3d/mypackbot/internal/model"

type Manager interface {
	Create(*model.Sticker) error
	Get(uint64) *model.Sticker
	GetByFileID(string) *model.Sticker
	GetList(int, int, string) (model.Stickers, int)
	GetOrCreate(*model.Sticker) (*model.Sticker, error)
	GetSet(string) (model.Stickers, int)
	Remove(uint64) error
	Update(*model.Sticker) error
}
