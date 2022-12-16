package main

// сюда писать код

import (
	"context"
)

var (
	// @BotFather в телеграме даст вам это
	BotToken = "5827575728:AAGzyCtfF98NhB8cr700536evIF6rW27tyM"

	// урл выдаст вам нгрок или хероку
	WebhookURL = "https://525f2cb5.ngrok.io"
)

func startTaskBot(ctx context.Context) error {
	// сюда пишите ваш код
	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
