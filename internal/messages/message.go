package messages

import (
	"strings"

	"github.com/toby3d/MyPackBot/internal/actions"
	"github.com/toby3d/MyPackBot/internal/commands"
	"github.com/toby3d/MyPackBot/internal/errors"
	"github.com/toby3d/MyPackBot/internal/i18n"
	tg "github.com/toby3d/telegram"
)

// Message checks user message on response, stickers, reset key phrase, else do
// Actions
func Message(msg *tg.Message) {
	T, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	switch {
	case strings.EqualFold(msg.Text, T("button_add_sticker")):
		commands.Add(msg, false)
	case strings.EqualFold(msg.Text, T("button_add_pack")):
		commands.Add(msg, true)
	case strings.EqualFold(msg.Text, T("button_del_sticker")):
		commands.Delete(msg, false)
	case strings.EqualFold(msg.Text, T("button_del_pack")):
		commands.Delete(msg, true)
	case strings.EqualFold(msg.Text, T("button_reset")):
		commands.Reset(msg)
	case strings.EqualFold(msg.Text, T("button_cancel")):
		commands.Cancel(msg)
	default:
		actions.Action(msg)
	}
}
