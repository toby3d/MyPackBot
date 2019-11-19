package users

import "gitlab.com/toby3d/mypackbot/internal/model"

type Manager interface {
	Create(*model.User) error
	Get(int) *model.User
	GetOrCreate(*model.User) (*model.User, error)
	Remove(int) error
	Update(*model.User) error
}
