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
	WebhookURL = "https://fc79-195-19-61-105.eu.ngrok.io"
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

	pull := new(handlers.Pull)

	for update := range updates {
		log.Printf("upd: %#v\n", update)
		text := strings.Split(update.Message.Text, " ")
		command := strings.Split(text[0], "_")

		switch command[0] {
		case "/tasks":
			err := handlers.ShowTasks(pull, &update, bot)
			if err != nil {
				log.Printf("Error at showing tasks: %v", err)
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Не удалось показать задачи",
				))
			}
		case "/new":
			err := handlers.CreateNewTask(pull, &update, bot, strings.Join(text[1:], " "))
			if err != nil {
				log.Printf("Error at creating tasks: %v", err)
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Не удалось создать задачу",
				))
			}
		case "/assign":
			err := handlers.AssignTask(pull, &update, bot, command)
			if err != nil {
				log.Printf("Error at assigning tasks: %v", err)
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Не удалось назначить задачу",
				))
			}
		case "/unassign":
			err := handlers.UnassignTask(pull, &update, bot, command)
			if err != nil {
				log.Printf("Error at unassigning task: %v", err)
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Не удалось снять задачу",
				))
			}
		case "/resolve":
			{
				err := handlers.ResolveTask(pull, &update, bot, command)
				if err != nil {
					log.Printf("Error at resolving task: %v", err)
					bot.Send(tgbotapi.NewMessage(
						update.Message.Chat.ID,
						"Не удалось выполнить задачу",
					))
				}
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
