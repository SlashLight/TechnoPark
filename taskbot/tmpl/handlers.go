package tmpl

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
	"strconv"
	"text/template"
)

func ShowTasks(pull *Pull, update *tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	var (
		chatId   = update.Message.Chat.ID
		sender   = update.Message.From.UserName
		tpl      bytes.Buffer
		TempShow = template.Must(template.New("templates.txt").Funcs(template.FuncMap{"inc": inc, "deref": deref, "isMe": isMe(sender), "isActive": isActive}).ParseFiles("./tmpl/templates.txt"))
	)

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

func CreateNewTask(pull *Pull, update *tgbotapi.Update, bot *tgbotapi.BotAPI, task string) error {
	var (
		chatId = update.Message.Chat.ID
		sender = update.Message.From.UserName
	)

	if len(task) == 0 {
		_, err := bot.Send(tgbotapi.NewMessage(
			chatId,
			"Напишите саму задачу после команды",
		))
		if err != nil {
			return err
		}
		return nil
	}

	NewTask := Task{
		Content:  task,
		Author:   User{ID: sender, ChatID: chatId},
		Executor: nil,
	}
	pull.Tasks = append(pull.Tasks, NewTask)

	_, err := bot.Send(tgbotapi.NewMessage(
		chatId,
		fmt.Sprintf("Задача \"%v\" создана, id=%v", task, len(pull.Tasks)),
	))
	if err != nil {
		return err
	}
	return nil
}

func AssignTask(pull *Pull, update *tgbotapi.Update, bot *tgbotapi.BotAPI, command []string) error {
	var (
		chatId = update.Message.Chat.ID
		sender = update.Message.From.UserName
	)

	if len(command) == 1 {
		_, err := bot.Send(tgbotapi.NewMessage(
			chatId,
			"Вы не ввели номер задачи",
		))
		if err != nil {
			return err
		}
		return nil
	}

	id, err := strconv.Atoi(command[1])
	if err != nil {
		_, err := bot.Send(tgbotapi.NewMessage(
			chatId,
			"Введите существующий номер",
		))
		if err != nil {
			return err
		}
		return nil
	}

	if len(pull.Tasks) < id {
		_, err := bot.Send(tgbotapi.NewMessage(
			chatId,
			"Вы ввели слишком большой номер. Такой задачи не существует",
		))
		if err != nil {
			return err
		}
		return nil
	}

	if pull.Tasks[id-1].Executor != nil && pull.Tasks[id-1].Executor.ChatID == chatId {
		_, err = bot.Send(tgbotapi.NewMessage(
			pull.Tasks[id-1].Executor.ChatID,
			"Эта задача уже назначена на вас",
		))
		if err != nil {
			return err
		}

		return nil
	}

	if pull.Tasks[id-1].Executor != nil && pull.Tasks[id-1].Executor.ChatID == pull.Tasks[id-1].Author.ChatID {
		_, err = bot.Send(tgbotapi.NewMessage(
			pull.Tasks[id-1].Executor.ChatID,
			fmt.Sprintf("Задача \"%v\" назначена на %v", pull.Tasks[id-1].Content, sender),
		))
		if err != nil {
			return err
		}
	} else {
		_, err = bot.Send(tgbotapi.NewMessage(
			pull.Tasks[id-1].Author.ChatID,
			fmt.Sprintf("Задача \"%v\" назначена на %v", pull.Tasks[id-1].Content, sender),
		))
		if err != nil {
			return err
		}
	}

	user := &User{ID: sender, ChatID: chatId}
	pull.Tasks[id-1].Executor = user

	_, err = bot.Send(tgbotapi.NewMessage(
		chatId,
		fmt.Sprintf("Задача \"%v\" назначена на вас", pull.Tasks[id-1].Content),
	))
	if err != nil {
		return err
	}

	return nil
}

func UnassignTask(pull *Pull, update *tgbotapi.Update, bot *tgbotapi.BotAPI, command []string) error {
	var (
		chatId = update.Message.Chat.ID
		sender = update.Message.From.UserName
	)

	if len(command) == 1 {
		_, err := bot.Send(tgbotapi.NewMessage(
			chatId,
			"Вы не ввели номер задачи",
		))
		if err != nil {
			return err
		}
		return nil
	}

	id, err := strconv.Atoi(command[1])
	if err != nil {
		_, err := bot.Send(tgbotapi.NewMessage(
			chatId,
			"Введите существующий номер",
		))
		if err != nil {
			return err
		}
		return nil
	}

	if len(pull.Tasks) < id {
		_, err := bot.Send(tgbotapi.NewMessage(
			chatId,
			"Вы ввели слишком большой номер. Такой задачи не существует",
		))
		if err != nil {
			return err
		}
		return nil
	}

	if pull.Tasks[id-1].Executor.ID != sender {
		_, err = bot.Send(tgbotapi.NewMessage(
			chatId,
			"Вы не являетесь исполнителем этой задачи",
		))
		if err != nil {
			return err
		}

		return nil
	}

	pull.Tasks[id-1].Executor = nil

	_, err = bot.Send(tgbotapi.NewMessage(
		chatId,
		"Принято",
	))
	if err != nil {
		return err
	}

	_, err = bot.Send(tgbotapi.NewMessage(
		pull.Tasks[id-1].Author.ChatID,
		fmt.Sprintf("Задача \"%v\" осталась без исполнителя", pull.Tasks[id-1].Content),
	))
	if err != nil {
		return err
	}

	return nil

}

func ResolveTask(pull *Pull, update *tgbotapi.Update, bot *tgbotapi.BotAPI, command []string) error {
	var (
		chatId = update.Message.Chat.ID
		sender = update.Message.From.UserName
	)

	if len(command) == 1 {
		_, err := bot.Send(tgbotapi.NewMessage(
			chatId,
			"Вы не ввели номер задачи",
		))
		if err != nil {
			return err
		}
		return nil
	}

	id, err := strconv.Atoi(command[1])
	if err != nil {
		_, err := bot.Send(tgbotapi.NewMessage(
			chatId,
			"Введите существующий номер",
		))
		if err != nil {
			return err
		}
		return nil
	}

	if len(pull.Tasks) < id {
		_, err := bot.Send(tgbotapi.NewMessage(
			chatId,
			"Вы ввели слишком большой номер. Такой задачи не существует",
		))
		if err != nil {
			return err
		}
		return nil
	}

	if pull.Tasks[id-1].Executor.ID != sender {
		_, err = bot.Send(tgbotapi.NewMessage(
			chatId,
			"Вы не являетесь исполнителем этой задачи",
		))
		if err != nil {
			return err
		}

		return nil
	}

	pull.Tasks[id-1].Executor = nil

	_, err = bot.Send(tgbotapi.NewMessage(
		chatId,
		fmt.Sprintf("Задача \"%v\" выполнена", pull.Tasks[id-1].Content),
	))
	if err != nil {
		return err
	}

	_, err = bot.Send(tgbotapi.NewMessage(
		pull.Tasks[id-1].Author.ChatID,
		fmt.Sprintf("Задача \"%v\" выполнена %v", pull.Tasks[id-1].Content, sender),
	))
	if err != nil {
		return err
	}

	pull.Tasks = append(pull.Tasks[:id-1], pull.Tasks[id:]...)

	return nil

}
