package tmpl

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
	"strconv"
	"text/template"
)

func ShowTasks(pool *Pool, update *tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	var (
		chatId   = update.Message.Chat.ID
		sender   = update.Message.From.UserName
		tpl      bytes.Buffer
		TempShow = template.Must(template.New("ShowTasks.txt").Funcs(template.FuncMap{"inc": inc, "deref": deref, "isMe": isMe(sender), "isActive": isActive}).ParseFiles("./tmpl/ShowTasks.txt"))
	)

	if len(pool.Tasks) == 0 {
		_, err := bot.Send(tgbotapi.NewMessage(
			chatId,
			"Список задач пуст",
		))
		if err != nil {
			return err
		}
		return nil
	}

	err := TempShow.Execute(&tpl, pool.Tasks)
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

func CreateNewTask(pool *Pool, update *tgbotapi.Update, bot *tgbotapi.BotAPI, task string) error {
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
	pool.Tasks = append(pool.Tasks, NewTask)

	_, err := bot.Send(tgbotapi.NewMessage(
		chatId,
		fmt.Sprintf("Задача \"%v\" создана, id=%v", task, len(pool.Tasks)),
	))
	if err != nil {
		return err
	}
	return nil
}

func AssignTask(pool *Pool, update *tgbotapi.Update, bot *tgbotapi.BotAPI, command []string) error {
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

	if len(pool.Tasks) < id {
		_, err := bot.Send(tgbotapi.NewMessage(
			chatId,
			"Вы ввели слишком большой номер. Такой задачи не существует",
		))
		if err != nil {
			return err
		}
		return nil
	}

	if pool.Tasks[id-1].Executor != nil && pool.Tasks[id-1].Executor.ChatID == chatId {
		_, err = bot.Send(tgbotapi.NewMessage(
			pool.Tasks[id-1].Executor.ChatID,
			"Эта задача уже назначена на вас",
		))
		if err != nil {
			return err
		}

		return nil
	}

	if pool.Tasks[id-1].Executor != nil && pool.Tasks[id-1].Executor.ChatID == pool.Tasks[id-1].Author.ChatID {
		_, err = bot.Send(tgbotapi.NewMessage(
			pool.Tasks[id-1].Executor.ChatID,
			fmt.Sprintf("Задача \"%v\" назначена на @%v", pool.Tasks[id-1].Content, sender),
		))
		if err != nil {
			return err
		}
	} else if chatId != pool.Tasks[id-1].Author.ChatID {
		_, err = bot.Send(tgbotapi.NewMessage(
			pool.Tasks[id-1].Author.ChatID,
			fmt.Sprintf("Задача \"%v\" назначена на @%v", pool.Tasks[id-1].Content, sender),
		))
		if err != nil {
			return err
		}
	}

	user := &User{ID: sender, ChatID: chatId}
	pool.Tasks[id-1].Executor = user

	_, err = bot.Send(tgbotapi.NewMessage(
		chatId,
		fmt.Sprintf("Задача \"%v\" назначена на вас", pool.Tasks[id-1].Content),
	))
	if err != nil {
		return err
	}

	return nil
}

func UnassignTask(pool *Pool, update *tgbotapi.Update, bot *tgbotapi.BotAPI, command []string) error {
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

	if len(pool.Tasks) < id {
		_, err := bot.Send(tgbotapi.NewMessage(
			chatId,
			"Вы ввели слишком большой номер. Такой задачи не существует",
		))
		if err != nil {
			return err
		}
		return nil
	}

	if pool.Tasks[id-1].Executor.ID != sender {
		_, err = bot.Send(tgbotapi.NewMessage(
			chatId,
			"Вы не являетесь исполнителем этой задачи",
		))
		if err != nil {
			return err
		}

		return nil
	}

	pool.Tasks[id-1].Executor = nil

	_, err = bot.Send(tgbotapi.NewMessage(
		chatId,
		"Принято",
	))
	if err != nil {
		return err
	}

	if chatId != pool.Tasks[id-1].Author.ChatID {
		_, err = bot.Send(tgbotapi.NewMessage(
			pool.Tasks[id-1].Author.ChatID,
			fmt.Sprintf("Задача \"%v\" осталась без исполнителя", pool.Tasks[id-1].Content),
		))
		if err != nil {
			return err
		}
	}

	return nil

}

func ResolveTask(pool *Pool, update *tgbotapi.Update, bot *tgbotapi.BotAPI, command []string) error {
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

	if len(pool.Tasks) < id {
		_, err := bot.Send(tgbotapi.NewMessage(
			chatId,
			"Вы ввели слишком большой номер. Такой задачи не существует",
		))
		if err != nil {
			return err
		}
		return nil
	}

	if pool.Tasks[id-1].Executor.ID != sender {
		_, err = bot.Send(tgbotapi.NewMessage(
			chatId,
			"Вы не являетесь исполнителем этой задачи",
		))
		if err != nil {
			return err
		}

		return nil
	}

	pool.Tasks[id-1].Executor = nil

	_, err = bot.Send(tgbotapi.NewMessage(
		chatId,
		fmt.Sprintf("Задача \"%v\" выполнена", pool.Tasks[id-1].Content),
	))
	if err != nil {
		return err
	}

	if chatId != pool.Tasks[id-1].Author.ChatID {
		_, err = bot.Send(tgbotapi.NewMessage(
			pool.Tasks[id-1].Author.ChatID,
			fmt.Sprintf("Задача \"%v\" выполнена @%v", pool.Tasks[id-1].Content, sender),
		))
		if err != nil {
			return err
		}
	}

	pool.Tasks = append(pool.Tasks[:id-1], pool.Tasks[id:]...)

	return nil

}

func MyTasks(pool *Pool, update *tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	var (
		chatId   = update.Message.Chat.ID
		sender   = update.Message.From.UserName
		tpl      bytes.Buffer
		TempShow = template.Must(template.New("MyTasks.txt").Funcs(template.FuncMap{"inc": inc, "deref": deref, "isMe": isMe(sender)}).ParseFiles("./tmpl/MyTasks.txt"))
	)

	err := TempShow.Execute(&tpl, pool.Tasks)
	if err != nil {
		return err
	}

	if len(tpl.String()) == 0 {
		_, err := bot.Send(tgbotapi.NewMessage(
			chatId,
			"Список задач пуст",
		))
		if err != nil {
			return err
		}
		return nil
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

func OwnTasks(pool *Pool, update *tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	var (
		chatId   = update.Message.Chat.ID
		sender   = update.Message.From.UserName
		tpl      bytes.Buffer
		TempShow = template.Must(template.New("OwnTasks.txt").Funcs(template.FuncMap{"inc": inc, "deref": deref, "isMe": isMe(sender), "isActive": isActive}).ParseFiles("./tmpl/OwnTasks.txt"))
	)

	err := TempShow.Execute(&tpl, pool.Tasks)
	if err != nil {
		return err
	}

	if len(tpl.String()) == 0 {
		_, err := bot.Send(tgbotapi.NewMessage(
			chatId,
			"Список задач, созданных вами, пуст",
		))
		if err != nil {
			return err
		}
		return nil
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
