package models

import tg "github.com/toby3d/telegram"

const (
	CommandAddPack       = "addPack"
	CommandAddSticker    = "addSticker"
	CommandCancel        = "cancel"
	CommandHelp          = "help"
	CommandDeleteSticker = "delSticker"
	CommandDeletePack    = "delPack"
	CommandReset         = "reset"
	CommandStart         = "start"

	StateNone          = "none"
	StateAddSticker    = CommandAddSticker
	StateAddPack       = CommandAddPack
	StateDeleteSticker = CommandDeleteSticker
	StateDeletePack    = CommandDeletePack
	StateReset         = CommandReset

	SetUploaded = "?"

	LanguageFallback = "en"
)

var AllowedUpdates = []string{
	tg.UpdateInlineQuery, // For searching and sending stickers
	tg.UpdateMessage,     // For get commands and messages
	tg.UpdateChannelPost, // For forwarding announcements
}
