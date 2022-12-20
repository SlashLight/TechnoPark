package tmpl

import (
	"bytes"
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

func showTasks(pull Pull, chatId int64, bot *tgbotapi.BotAPI) error {
	var tpl bytes.Buffer
	err := TempShow.Execute(&tpl, pull.Tasks)
	if err != nil {
		return err
	}

	bot.Send(tgbotapi.NewMessage(
		chatId,
		tpl.String(),
	))
	return nil
}
