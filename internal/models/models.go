package models

import (
	tg "gitlab.com/toby3d/telegram"
)

// Commands... represents available and supported bot commands
const (
	CommandAddPack       = "addPack"
	CommandAddSticker    = "addSticker"
	CommandCancel        = "cancel"
	CommandDeleteSticker = "delSticker"
	CommandDeletePack    = "delPack"
	CommandReset         = "reset"
)

// SetUploaded is a mimic set name of uploaded stickers without any parent set
const (
	SetFavorite = "!"
	SetUploaded = "?"
)

// AllowedUpdates is a filter list of updates from Telegram
var AllowedUpdates = []string{
	tg.UpdateInlineQuery, // For searching and sending stickers
	tg.UpdateMessage,     // For get commands and messages
	tg.UpdateChannelPost, // For forwarding announcements
}
