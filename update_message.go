package main

import (
	log "github.com/kirillDanshin/dlog"
	tg "github.com/toby3d/telegram"
)

func updateMessage(msg *tg.Message) {
	if bot.IsMessageFromMe(msg) || bot.IsForwardFromMe(msg) {
		log.Ln("Ignore message update")
		return
	}

	switch {
	case bot.IsCommandToMe(msg):
		commands(msg)
	case msg.IsText():
		messages(msg)
	default:
		actions(msg)
	}
}
