package main

import (
	"sync"
	"time"

	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

var waitForwards = new(sync.WaitGroup)

func updateChannelPost(post *tg.Message) {
	if post.Chat.ID != channelID {
		log.Ln(post.Chat.ID, "!=", channelID)
		return
	}

	users, err := dbGetUsers()
	errCheck(err)

	for i := range users {
		waitForwards.Add(1)
		if _, err = bot.ForwardMessage(
			tg.NewForwardMessage(post.Chat.ID, int64(users[i]), post.ID),
		); err != nil {
			log.Ln(err.Error())
		}
		waitForwards.Done()

		time.Sleep(time.Second / 10) // For avoid spamming
	}
}
