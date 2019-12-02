package users

import "gitlab.com/toby3d/mypackbot/internal/model"

type Manager interface {
	Create(*model.User) error
	Get(uint64) *model.User
	GetByUserID(int64) *model.User
	GetOrCreate(*model.User) (*model.User, error)
	Update(*model.User) error
}
