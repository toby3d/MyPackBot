package user

import "gitlab.com/toby3d/mypackbot/internal/models"

type Store interface {
	GetByID(int) (*models.User, error)
	Create(*models.User) error
	Update(*models.User) error
	AddSticker(int, string) error
	DeleteSticker(int, string) error
}
