package tmpl

import (
	"bytes"
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
	"text/template"
)

func ShowTasks(pull *Pull, chatId int64, sender string, bot *tgbotapi.BotAPI) error {
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

	var (
		tpl      bytes.Buffer
		TempShow = template.Must(template.New("templates.txt").Funcs(template.FuncMap{"inc": inc, "deref": deref, "isMe": isMe(sender), "isActive": isActive}).ParseFiles("./tmpl/templates.txt"))
	)
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
