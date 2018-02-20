package updates

import (
	"time"

	log "github.com/kirillDanshin/dlog"
	"github.com/toby3d/MyPackBot/internal/bot"
	"github.com/toby3d/MyPackBot/internal/config"
	"github.com/toby3d/MyPackBot/internal/db"
	"github.com/toby3d/MyPackBot/internal/errors"
	tg "github.com/toby3d/telegram"
)

// ChannelPost checks ChannelPost update for forwarding content to bot users
func ChannelPost(post *tg.Message) {
	if post.Chat.ID != config.ChannelID {
		log.Ln(post.Chat.ID, "!=", config.ChannelID)
		return
	}

	users, err := db.Users()
	errors.Check(err)

	for i := range users {
		errors.WaitForwards.Add(1)
		if _, err = bot.Bot.ForwardMessage(
			tg.NewForwardMessage(post.Chat.ID, int64(users[i]), post.ID),
		); err != nil {
			log.Ln(err.Error())
		}
		errors.WaitForwards.Done()

		time.Sleep(time.Second / 10) // For avoid spamming
	}
}
