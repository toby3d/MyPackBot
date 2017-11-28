package main

import (
	"strings"

	log "github.com/kirillDanshin/dlog" // Insert logs only in debug builds
	"github.com/toby3d/botan"
	"github.com/toby3d/go-telegram" // My Telegram bindings
)

const (
	cmdAddPack    = "addPack"
	cmdAddSticker = "addSticker"
	cmdCancel     = "cancel"
	cmdHelp       = "help"
	cmdDelete     = "del"
	cmdReset      = "reset"
	cmdStart      = "start"
)

func commands(msg *telegram.Message) {
	log.Ln("Received a", msg.Command(), "command")
	switch strings.ToLower(msg.Command()) {
	case strings.ToLower(cmdStart):
		appMetrika.TrackAsync(
			"Start", msg.From.ID, *msg,
			func(answer *botan.Answer, err error) {
				log.Ln("Asynchonous:", answer.Status)
				metrika <- true
			},
		)

		commandStart(msg)

		<-metrika
	case strings.ToLower(cmdHelp):
		appMetrika.TrackAsync(
			"Help", msg.From.ID, *msg,
			func(answer *botan.Answer, err error) {
				log.Ln("Asynchonous:", answer.Status)
				metrika <- true
			},
		)

		commandHelp(msg)

		<-metrika
	case strings.ToLower(cmdAddSticker):
		appMetrika.TrackAsync(
			"Add single sticker", msg.From.ID, *msg,
			func(answer *botan.Answer, err error) {
				log.Ln("Asynchonous:", answer.Status)
				metrika <- true
			},
		)

		commandAdd(msg, false)

		<-metrika
	case strings.ToLower(cmdAddPack):
		appMetrika.TrackAsync(
			"Add pack", msg.From.ID, *msg,
			func(answer *botan.Answer, err error) {
				log.Ln("Asynchonous:", answer.Status)
				metrika <- true
			},
		)

		commandAdd(msg, true)

		<-metrika
	case strings.ToLower(cmdDelete):
		appMetrika.TrackAsync(
			"Delete single sticker", msg.From.ID, *msg,
			func(answer *botan.Answer, err error) {
				log.Ln("Asynchonous:", answer.Status)
				metrika <- true
			},
		)

		commandDelete(msg)

		<-metrika
	case strings.ToLower(cmdReset):
		appMetrika.TrackAsync(
			"Reset", msg.From.ID, *msg,
			func(answer *botan.Answer, err error) {
				log.Ln("Asynchonous:", answer.Status)
				metrika <- true
			},
		)

		commandReset(msg)

		<-metrika
	case strings.ToLower(cmdCancel):
		appMetrika.TrackAsync(
			"Cancel", msg.From.ID, *msg,
			func(answer *botan.Answer, err error) {
				log.Ln("Asynchonous:", answer.Status)
				metrika <- true
			},
		)

		commandCancel(msg)

		<-metrika
	default:
		appMetrika.Track("Command", msg.From.ID, *msg)
	}
}
