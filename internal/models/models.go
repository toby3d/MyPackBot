package models

import tg "gitlab.com/toby3d/telegram"

// Commands... represents available and supported bot commands
const (
	CommandAddPack       = "addPack"
	CommandAddSticker    = "addSticker"
	CommandCancel        = "cancel"
	CommandDeleteSticker = "delSticker"
	CommandDeletePack    = "delPack"
	CommandReset         = "reset"
)

// State... represents current state of user action
const (
	StateNone          = "none"
	StateAddSticker    = CommandAddSticker
	StateAddPack       = CommandAddPack
	StateDeleteSticker = CommandDeleteSticker
	StateDeletePack    = CommandDeletePack
	StateReset         = CommandReset
)

// SetUploaded is a mimic set name of uploaded stickers without any parent set
const SetUploaded = "?"

// AllowedUpdates is
var AllowedUpdates = []string{
	tg.UpdateInlineQuery, // For searching and sending stickers
	tg.UpdateMessage,     // For get commands and messages
	tg.UpdateChannelPost, // For forwarding announcements
}
