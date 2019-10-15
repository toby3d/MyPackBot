package store

import model "gitlab.com/toby3d/mypackbot/internal/model"

type (
	Manager interface {
		AddSticker(*model.User, *model.Sticker) error
		AddStickersSet(*model.User, string) error
		GetSticker(*model.User, *model.Sticker) (*model.UserSticker, error)
		GetStickersList(*model.User, int, int, string) (model.Stickers, int)
		GetStickersSet(*model.User, int, int, string) (model.Stickers, int)
		RemoveSticker(*model.User, *model.Sticker) error
		RemoveStickersSet(*model.User, string) error
		Stickers() StickersManager
		Users() UsersManager
	}

	UsersManager interface {
		Create(*model.User) error
		Get(int) *model.User
		GetOrCreate(*model.User) (*model.User, error)
		Remove(int) error
		Update(*model.User) error
	}

	StickersManager interface {
		Create(*model.Sticker) error
		Get(string) *model.Sticker
		GetList(int, int, string) (model.Stickers, int)
		GetOrCreate(*model.Sticker) (*model.Sticker, error)
		GetSet(string) (model.Stickers, int)
		Remove(string) error
		Update(*model.Sticker) error
	}
)
