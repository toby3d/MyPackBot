package sticker

import "gitlab.com/toby3d/mypackbot/internal/models"

type Store interface {
	GetByID(string) (*models.Sticker, error)
	Create(*models.Sticker) error
	Update(*models.Sticker) error
	Delete(*models.Sticker) error
	List(int, int) ([]models.Sticker, int, error)
	ListByEmoji(string, int, int) ([]models.Sticker, int, error)
	GetSet(string, int, int) ([]models.Sticker, int, error)
}
