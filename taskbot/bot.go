package main

// сюда писать код

import (
	"context"
	"fmt"
	tgbotapi "github.com/skinass/telegram-bot-api/v5"
	handlers "gitlab.com/mailru-go/lectures-2022-1/04_net2/99_hw/taskbot/tmpl"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	// @BotFather в телеграме даст вам это
	BotToken = "5827575728:AAGzyCtfF98NhB8cr700536evIF6rW27tyM"

	// урл выдаст вам нгрок или хероку
	WebhookURL = "https://f9e4-91-193-176-7.eu.ngrok.io"
)

func startTaskBot(ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatalf("NewBotAPI failed: %s", err)
	}

	bot.Debug = true

	wh, err := tgbotapi.NewWebhook(WebhookURL)
	if err != nil {
		log.Fatalf("NewWebhook failed: %s", err)
	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("SetWebhook failed: %s", err)
	}

	updates := bot.ListenForWebhook("/")

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all is working"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	go func() {
		log.Fatalln("http err:", http.ListenAndServe(":"+port, nil))
	}()
	fmt.Println("start listen :" + port)

	pool := new(handlers.Pool)

	for update := range updates {
		log.Printf("upd: %#v\n", update)
		text := strings.Split(update.Message.Text, " ")
		command := strings.Split(text[0], "_")

		switch command[0] {
		case "/tasks":
			err := handlers.ShowTasks(pool, &update, bot)
			if err != nil {
				log.Printf("Error at showing tasks: %v", err)
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Не удалось показать задачи",
				))
			}
		case "/new":
			err := handlers.CreateNewTask(pool, &update, bot, strings.Join(text[1:], " "))
			if err != nil {
				log.Printf("Error at creating tasks: %v", err)
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Не удалось создать задачу",
				))
			}
		case "/assign":
			err := handlers.AssignTask(pool, &update, bot, command)
			if err != nil {
				log.Printf("Error at assigning tasks: %v", err)
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Не удалось назначить задачу",
				))
			}
		case "/unassign":
			err := handlers.UnassignTask(pool, &update, bot, command)
			if err != nil {
				log.Printf("Error at unassigning task: %v", err)
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Не удалось снять задачу",
				))
			}
		case "/resolve":
			err := handlers.ResolveTask(pool, &update, bot, command)
			if err != nil {
				log.Printf("Error at resolving task: %v", err)
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Не удалось выполнить задачу",
				))
			}

		case "/help":
			fHelp, err := os.ReadFile("./tmpl/Help.txt")
			if err != nil {
				log.Printf("Error at opening file: %v", err)
			}
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				string(fHelp),
			))
		case "/my":
			err := handlers.MyTasks(pool, &update, bot)
			if err != nil {
				log.Printf("Error at showing tasks: %v", err)
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Не удалось показать задачи",
				))
			}
		case "/owner":
			err := handlers.OwnTasks(pool, &update, bot)
			if err != nil {
				log.Printf("Error at showing users's tasks: %v", err)
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Не удалось показать задачи",
				))
			}

		default:
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Неопознанная команда",
			))
		}
	}

	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
