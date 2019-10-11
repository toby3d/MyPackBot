package model

type (
	// User represent a simple bot user
	User struct {
		ID           int    `json:"id"`
		CreatedAt    int64  `json:"created_at"`
		UpdatedAt    int64  `json:"updated_at"`
		LastSeen     int64  `json:"last_seen"`
		LanguageCode string `json:"language_code"`
	}

	Users []*User

	Sticker struct {
		ID         string `json:"id"`
		SetName    string `json:"set_name"`
		CreatedAt  int64  `json:"created_at"`
		IsAnimated bool   `json:"is_animated"`
		Hits       int    `json:"hits"`
	}

	Stickers []*Sticker

	UserSticker struct {
		UserID    int    `json:"user_id"`
		Hits      int    `json:"hits"`
		CreatedAt int64  `json:"created_at"`
		StickerID string `json:"sticker_id"`
		Emoji     string `json:"emoji"`
	}

	UserStickers []*UserSticker
)

func (users Users) GetByID(id int) *User {
	for i := range users {
		if users[i].ID != id {
			continue
		}
		return users[i]
	}
	return nil
}

func (stickers Stickers) GetByID(id string) *Sticker {
	for i := range stickers {
		if stickers[i].ID != id {
			continue
		}
		return stickers[i]
	}
	return nil
}

func (userStickers UserStickers) GetByID(uid int, sid string) *UserSticker {
	for i := range userStickers {
		if userStickers[i].UserID != uid || userStickers[i].StickerID != sid {
			continue
		}
		return userStickers[i]
	}
	return nil
}
