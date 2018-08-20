package messages

import (
	"fmt"
	"strings"

	"gitlab.com/toby3d/mypackbot/internal/actions"
	"gitlab.com/toby3d/mypackbot/internal/commands"
	"gitlab.com/toby3d/mypackbot/internal/utils"
	tg "gitlab.com/toby3d/telegram"
)

// Message checks user message on response, stickers, reset key phrase, else do
// Actions
func Message(msg *tg.Message) {
	p := utils.NewPrinter(msg.From.LanguageCode)

	switch {
	case strings.EqualFold(msg.Text, fmt.Sprintf("➕ %s", p.Sprintf("add a sticker"))):
		commands.Add(msg, false)
	case strings.EqualFold(msg.Text, fmt.Sprintf("📦 %s", p.Sprintf("add set"))):
		commands.Add(msg, true)
	case strings.EqualFold(msg.Text, fmt.Sprintf("🗑 %s", p.Sprintf("remove sticker"))):
		commands.Delete(msg, false)
	case strings.EqualFold(msg.Text, fmt.Sprintf("🗑 %s", p.Sprintf("delete set"))):
		commands.Delete(msg, true)
	case strings.EqualFold(msg.Text, fmt.Sprintf("🔥 %s", p.Sprintf("reset set"))):
		commands.Reset(msg)
	case strings.EqualFold(msg.Text, fmt.Sprintf("❌ %s", p.Sprintf("cancel"))):
		commands.Cancel(msg)
	default:
		actions.Action(msg)
	}
}
