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
)

var (
	// @BotFather в телеграме даст вам это
	BotToken = "5827575728:AAGzyCtfF98NhB8cr700536evIF6rW27tyM"

	// урл выдаст вам нгрок или хероку
	WebhookURL = "https://8c13-178-217-27-193.eu.ngrok.io"
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
		command := update.Message.Text

		switch command {
		case "/tasks":
			if len(pull.Tasks) == 0 {
				bot.Send(tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Список задач пуст",
				))
			} else {
				err := showTasks(pull, update.Message.Chat.ID, bot)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(
						update.Message.Chat.ID,
						"Failed at showing tasks",
					))
				}
			}
		default:
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Иди поспи",
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
