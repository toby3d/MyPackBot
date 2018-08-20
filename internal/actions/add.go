package actions

import (
	log "github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/models"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

// Add action add sticker or set to user's pack
func Add(msg *tg.Message, pack bool) {
	if !msg.IsSticker() {
		return
	}

	p := utils.NewPrinter(msg.From.LanguageCode)

	_, err := bot.Bot.SendChatAction(msg.Chat.ID, tg.ActionTyping)
	errors.Check(err)

	reply := tg.NewMessage(msg.Chat.ID, p.Sprintf("The sticker was successfully added to your set!"))
	reply.ParseMode = tg.StyleMarkdown

	if !pack {
		var exist bool
		sticker := msg.Sticker
		exist, err = db.DB.AddSticker(msg.From.ID, sticker)
		errors.Check(err)

		if exist {
			reply.Text = p.Sprintf("This sticker is already in your collection.")
		}

		reply.ReplyMarkup = utils.CancelButton(p)
		_, err = bot.Bot.SendMessage(reply)
		errors.Check(err)
		return
	}

	reply.Text = p.Sprintf(
		"It seems you're trying to add your own sticker. Use the /%s command for this.",
		models.CommandAddSticker,
	)

	if msg.Sticker.SetName != "" {
		var set *tg.StickerSet
		set, err = bot.Bot.GetStickerSet(msg.Sticker.SetName)
		errors.Check(err)

		log.Ln("SetTitle:", set.Title)
		reply.Text = p.Sprintf(
			"The set *%s* was successfully added to yours!",
			set.Title,
		)

		allExists := true
		for i := range set.Stickers {
			var exist bool
			exist, err = db.DB.AddSticker(msg.From.ID, &set.Stickers[i])
			errors.Check(err)

			if !exist {
				allExists = false
			}
		}

		log.Ln("All exists?", allExists)
		if allExists {
			reply.Text = p.Sprintf(
				"All of the *%s* stickers are already in your collection.",
				set.Title,
			)
		}
	}

	reply.ReplyMarkup = utils.CancelButton(p)
	_, err = bot.Bot.SendMessage(reply)
	errors.Check(err)
}
