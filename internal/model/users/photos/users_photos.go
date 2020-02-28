package photos

import "gitlab.com/toby3d/mypackbot/internal/model"

type (
	ReadWriter interface {
		Reader
		Writer
	}

	Reader interface {
		Get(up *model.UserPhoto) *model.Photo
		GetList(offset int, limit int, filter *model.UserPhoto) (model.Photos, int, error)
	}

	Writer interface {
		Add(up *model.UserPhoto) error
		Update(up *model.UserPhoto) error
		Remove(up *model.UserPhoto) error
	}
)
