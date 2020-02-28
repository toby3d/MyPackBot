package users

import "gitlab.com/toby3d/mypackbot/internal/model"

type (
	Manager interface {
		Reader
		Writer
		ReadWriter
	}

	ReadWriter interface {
		GetOrCreate(u *model.User) (*model.User, error)
	}

	Reader interface {
		Get(id int) *model.User
		GetList(offset, limit int, filter *model.User) (model.Users, int, error)
	}

	Writer interface {
		Create(u *model.User) error
		Update(u *model.User) error
	}
)
