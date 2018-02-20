package main

import (
	"flag"

	log "github.com/kirillDanshin/dlog"
	_ "github.com/toby3d/MyPackBot/init"
	"github.com/toby3d/MyPackBot/internal/updates"
)

var flagWebhook = flag.Bool(
	"webhook", false,
	"enable work via webhooks (required valid certificates)",
)

// main function is a general function for work of this bot
func main() {
	flag.Parse() // Parse flagWebhook

	for update := range updates.Channel(*flagWebhook) {
		log.D(update)
		switch {
		case update.IsInlineQuery():
			updates.InlineQuery(update.InlineQuery)
		case update.IsMessage():
			updates.Message(update.Message)
		case update.IsChannelPost():
			updates.ChannelPost(update.ChannelPost)
		}
	}
}
