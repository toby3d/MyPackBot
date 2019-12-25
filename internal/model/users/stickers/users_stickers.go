package stickers

import "gitlab.com/toby3d/mypackbot/internal/model"

type Manager interface {
	Add(*model.UserSticker) error
	AddSet(uint64, string) error
	Get(*model.UserSticker) *model.UserSticker
	GetList(uint64, int, int, string) (model.Stickers, int)
	GetSet(uint64, int, int, string) (model.Stickers, int)
	Remove(*model.UserSticker) error
	RemoveSet(uint64, string) error
	Update(*model.UserSticker) error
}
