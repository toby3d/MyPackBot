package store

import model "gitlab.com/toby3d/mypackbot/internal/model"

type (
	Manager interface {
		AddSticker(*model.User, *model.Sticker, string) error
		GetSticker(*model.User, *model.Sticker) (*model.UserSticker, error)
		GetStickersList(*model.User, int, int, string) (model.Stickers, int)
		GetStickersSet(*model.User, int, int, string) (model.Stickers, int)
		HitSticker(*model.User, *model.Sticker) error
		RemoveSticker(*model.User, *model.Sticker) error
		Stickers() StickersManager
		Users() UsersManager
	}

	UsersManager interface {
		Create(*model.User) error
		Get(int) *model.User
		Update(*model.User) error
		Remove(int) error
		GetOrCreate(*model.User) (*model.User, error)
	}

	StickersManager interface {
		Create(*model.Sticker) error
		Get(string) *model.Sticker
		GetOrCreate(*model.Sticker) (*model.Sticker, error)
		Hit(string) error
		Remove(string) error
		Update(*model.Sticker) error
	}
)
