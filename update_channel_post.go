package main

import (
	"time"

	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	tg "github.com/toby3d/telegram"     // My Telegram bindings
)

func updateChannelPost(post *tg.Message) {
	if post.Chat.ID != channelID {
		log.Ln(post.Chat.ID, "!=", channelID)
		return
	}

	users, err := dbGetUsers()
	errCheck(err)

	for i := range users {
		if _, err = bot.ForwardMessage(
			tg.NewForwardMessage(post.Chat.ID, int64(users[i]), post.ID),
		); err != nil {
			log.Ln(err.Error())
		}

		time.Sleep(time.Second / 10) // For avoid spamming
	}
}
