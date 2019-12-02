package photos

import "gitlab.com/toby3d/mypackbot/internal/model"

type Manager interface {
	Create(*model.Photo) error
	Get(uint64) *model.Photo
	GetByFileID(string) *model.Photo
	GetList(int, int) (model.Photos, int)
	Update(*model.Photo) error
	Remove(uint64) error
	GetOrCreate(*model.Photo) (*model.Photo, error)
}
