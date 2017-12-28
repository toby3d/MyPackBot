package main

import (
	"time"

	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	tg "github.com/toby3d/telegram"     // My Telegram bindings
)

func channelPost(post *tg.Message) {
	if post.Chat.ID != channelID {
		log.Ln(post.Chat.ID, "!=", channelID)
		return
	}

	users, err := dbGetUsers()
	errCheck(err)

	for i := range users {
		bot.ForwardMessage(
			tg.NewForwardMessage(post.Chat.ID, int64(users[i]), post.ID),
		)

		time.Sleep(time.Second / 10) // For avoid spamming
	}
}
