package tmpl

import (
	"bytes"
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

func showTasks(pull Pull, chatId int64, bot *tgbotapi.BotAPI) error {
	if len(pull.Tasks) == 0 {
		_, err := bot.Send(tgbotapi.NewMessage(
			chatId,
			"Список задач пуст",
		))
		if err != nil {
			return err
		}

		return nil
	}

	var tpl bytes.Buffer
	err := TempShow.Execute(&tpl, pull.Tasks)
	if err != nil {
		return err
	}

	_, err = bot.Send(tgbotapi.NewMessage(
		chatId,
		tpl.String(),
	))
	if err != nil {
		return err
	}
	return nil
}
