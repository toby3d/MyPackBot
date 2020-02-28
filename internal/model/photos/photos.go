package photos

import "gitlab.com/toby3d/mypackbot/internal/model"

type (
	Manager interface {
		Reader
		Writer
		ReadWriter
	}

	ReadWriter interface {
		GetOrCreate(*model.Photo) (*model.Photo, error)
	}

	Reader interface {
		Get(string) *model.Photo
		GetList(int, int, *model.Photo) (model.Photos, int, error)
	}

	Writer interface {
		Create(*model.Photo) error
		Update(*model.Photo) error
		Remove(string) error
	}
)
