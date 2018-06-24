package updates

import (
	"time"

	log "github.com/kirillDanshin/dlog"
	"gitlab.com/toby3d/mypackbot/internal/bot"
	"gitlab.com/toby3d/mypackbot/internal/config"
	"gitlab.com/toby3d/mypackbot/internal/db"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	tg "gitlab.com/toby3d/telegram"
)

// ChannelPost checks ChannelPost update for forwarding content to bot users
func ChannelPost(post *tg.Message) {
	channelID := config.Config.GetInt64("telegram.channel")
	if post.Chat.ID != channelID {
		log.Ln(post.Chat.ID, "!=", channelID)
		return
	}

	users, err := db.DB.GetUsers()
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
