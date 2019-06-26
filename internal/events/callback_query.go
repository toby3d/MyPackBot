package events

import (
	"strings"
	"time"

	"gitlab.com/toby3d/mypackbot/internal/models"
	tg "gitlab.com/toby3d/telegram"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (e *Events) CallbackQuery(b *tg.Bot, cc *tg.CallbackQuery) error {
	p := message.NewPrinter(language.English)
	answer := tg.NewAnswerCallbackQuery(cc.ID)
	answer.ShowAlert = false
	answer.Text = p.Sprintf("🤷🏻‍♂️ Unknown action")

	user, err := e.store.GetOrCreateUser(&models.User{
		ID:           cc.From.ID,
		LanguageCode: cc.From.LanguageCode,
		StartedAt:    time.Now().UTC().Unix(),
		AutoSaving:   true,
	})
	if err := e.store.CreateUser(user); err != nil {
		_, err = b.AnswerCallbackQuery(answer)
		return err
	}

	p = message.NewPrinter(language.Make(user.LanguageCode))

	parts := strings.Split(cc.Data, ":")
	switch parts[0] {
	case "set":
		switch parts[1] {
		case "import":
			set := e.store.GetSet(parts[2])
			if set == nil {
				tgSet, err := b.GetStickerSet(parts[2])
				if err != nil {
					_, err = b.AnswerCallbackQuery(answer)
					return err
				}

				if set, err = e.store.GetOrCreateSet(&models.Set{
					Name:  tgSet.Name,
					Title: tgSet.Title,
				}); err != nil {
					_, err = b.AnswerCallbackQuery(answer)
					return err
				}

				go func() {
					for _, tgSticker := range tgSet.Stickers {
						sticker, err := e.store.GetOrCreateSticker(&models.Sticker{
							Model: models.Model{
								ID:      tgSticker.FileID,
								SavedAt: time.Now().UTC().Unix(),
							},
							Emoji:   tgSticker.Emoji,
							SetName: set.Name,
						})
						if err != nil {
							continue
						}

						if err = e.store.AddSticker(user.ID, sticker.ID); err != nil {
							continue
						}
					}
				}()
			}
			answer.Text = p.Sprintf("📥 All stickers of %s set has been added!", set.Title)
		}
	}

	_, err = b.AnswerCallbackQuery(answer)
	return err
}
