package photos

import "gitlab.com/toby3d/mypackbot/internal/model"

type Manager interface {
	Add(*model.UserPhoto) error
	Get(*model.UserPhoto) *model.UserPhoto
	GetList(uint64, int, int, string) (model.Photos, int)
	Remove(*model.UserPhoto) error
	Update(*model.UserPhoto) error
}
