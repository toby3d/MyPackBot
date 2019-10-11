package middleware

import (
	"gitlab.com/toby3d/mypackbot/internal/model"
	"gitlab.com/toby3d/mypackbot/internal/model/store"
	tg "gitlab.com/toby3d/telegram"
)

func GetUser(s store.Manager, u *tg.User, date int64) (*model.User, error) {
	return s.Users().GetOrCreate(&model.User{
		ID:           u.ID,
		LanguageCode: u.LanguageCode,
		CreatedAt:    date,
		UpdatedAt:    date,
		LastSeen:     date,
	})
}
