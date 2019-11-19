package model

import (
	"fmt"

	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/xerrors"
)

type (
	User struct {
		ID        int   `json:"id"`
		CreatedAt int64 `json:"created_at"`
		UpdatedAt int64 `json:"updated_at"`

		LanguageCode string `json:"language_code"`
		LastSeen     int64  `json:"last_seen"`
	}

	Users []*User

	Sticker struct {
		ID        string `json:"id"`
		CreatedAt int64  `json:"created_at"`
		UpdatedAt int64  `json:"updated_at"`

		Width      int    `json:"width"`
		Height     int    `json:"height"`
		IsAnimated bool   `json:"is_animated"`
		SetName    string `json:"set_name"`
		Emoji      string `json:"emoji"`
	}

	Stickers []*Sticker

	UserSticker struct {
		StickerID string `json:"sticker_id"`
		UserID    int    `json:"user_id"`
		CreatedAt int64  `json:"created_at"`

		SetName string `json:"set_name"`
		Emojis  string `json:"emojis"`
	}

	UserStickers []*UserSticker

	UpdateFunc func(*Context) error

	Context struct {
		*tg.Bot
		*tg.Update
		printer *message.Printer
		Sticker *Sticker
		User    *User
	}

	Error struct {
		Message string
		frame   xerrors.Frame
	}
)

func (err Error) FormatError(p xerrors.Printer) error {
	p.Printf("🐛 %s", err.Message)
	err.frame.Format(p)
	return nil
}

func (err Error) Format(f fmt.State, c rune) {
	xerrors.FormatError(err, f, c)
}

func (err Error) Error() string {
	return fmt.Sprint(err)
}

func (users Users) GetByID(id int) *User {
	var u *User

	for i := range users {
		if users[i].ID != id {
			continue
		}

		u = users[i]

		break
	}

	return u
}

func (stickers Stickers) GetByID(id string) *Sticker {
	var s *Sticker

	for i := range stickers {
		if stickers[i].ID != id {
			continue
		}

		s = stickers[i]
	}

	return s
}

func (stickers Stickers) GetSet(name string) (Stickers, int) {
	set := make(Stickers, 0)

	for i := range stickers {
		if stickers[i].SetName != name {
			continue
		}

		set = append(set, stickers[i])
	}

	return set, len(set)
}

func (userStickers UserStickers) GetByID(uid int, sid string) *UserSticker {
	var us *UserSticker

	for i := range userStickers {
		if userStickers[i].UserID != uid || userStickers[i].StickerID != sid {
			continue
		}

		us = userStickers[i]
	}

	return us
}

func (ctx *Context) T() *message.Printer {
	if ctx.printer != nil {
		return ctx.printer
	}

	code := language.English
	if ctx.User != nil {
		code = language.Make(ctx.User.LanguageCode)
	}

	tag, _, _ := message.DefaultCatalog.Matcher().Match(code)
	ctx.printer = message.NewPrinter(tag)

	return ctx.T()
}

func (ctx *Context) Error(err error) error {
	switch {
	case ctx.IsCallbackQuery():
		answer := tg.NewAnswerCallbackQuery(ctx.CallbackQuery.ID)
		answer.Text = "🐞 " + err.Error()

		if _, sendErr := ctx.AnswerCallbackQuery(answer); sendErr != nil {
			err = sendErr
		}
	case ctx.IsInlineQuery():
		answer := tg.NewAnswerInlineQuery(ctx.InlineQuery.ID)
		answer.IsPersonal = true
		answer.SwitchPrivateMessageParameter = "error"
		answer.SwitchPrivateMessageText = "🐞 " + err.Error()

		if _, sendErr := ctx.AnswerInlineQuery(answer); sendErr != nil {
			err = sendErr
		}
	}

	return err
}
