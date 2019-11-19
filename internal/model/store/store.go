package store

import model "gitlab.com/toby3d/mypackbot/internal/model"

type Manager interface {
	AddSticker(*model.User, *model.Sticker) error
	AddStickersSet(*model.User, string) error
	GetSticker(*model.User, *model.Sticker) (*model.UserSticker, error)
	GetStickersList(*model.User, int, int, string) (model.Stickers, int)
	GetStickersSet(*model.User, int, int, string) (model.Stickers, int)
	RemoveSticker(*model.User, *model.Sticker) error
	RemoveStickersSet(*model.User, string) error
}
