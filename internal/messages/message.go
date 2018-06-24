package messages

import (
	"strings"

	"gitlab.com/toby3d/mypackbot/internal/actions"
	"gitlab.com/toby3d/mypackbot/internal/commands"
	"gitlab.com/toby3d/mypackbot/internal/errors"
	"gitlab.com/toby3d/mypackbot/internal/i18n"
	tg "gitlab.com/toby3d/telegram"
)

// Message checks user message on response, stickers, reset key phrase, else do
// Actions
func Message(msg *tg.Message) {
	t, err := i18n.SwitchTo(msg.From.LanguageCode)
	errors.Check(err)

	switch {
	case strings.EqualFold(msg.Text, t("button_add_sticker")):
		commands.Add(msg, false)
	case strings.EqualFold(msg.Text, t("button_add_pack")):
		commands.Add(msg, true)
	case strings.EqualFold(msg.Text, t("button_del_sticker")):
		commands.Delete(msg, false)
	case strings.EqualFold(msg.Text, t("button_del_pack")):
		commands.Delete(msg, true)
	case strings.EqualFold(msg.Text, t("button_reset")):
		commands.Reset(msg)
	case strings.EqualFold(msg.Text, t("button_cancel")):
		commands.Cancel(msg)
	default:
		actions.Action(msg)
	}
}
