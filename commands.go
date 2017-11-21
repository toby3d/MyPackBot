package main

import (
	"strings"

	"github.com/toby3d/go-telegram" // My Telegram bindings
)

func commands(msg *telegram.Message) {
	switch strings.ToLower(msg.Command()) {
	case "start":
		commandStart(msg)
	case "help":
		commandHelp(msg)
	case "add":
		commandAdd(msg)
	case "remove":
		commandRemove(msg)
	case "reset":
		commandReset(msg)
	case "cancel":
		commandCancel(msg)
	}
}
