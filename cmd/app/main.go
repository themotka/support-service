package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log/slog"
	"os"
	"support-bot/internal/config"
	"support-bot/internal/storage/postgres"
)

const managerID = 1

func main() {
	token := mustToken()
	cfg := config.MustLoad()

	_ = cfg
	storage, err := postgres.NewStorage("user=postgres password=123 dbname=postgres sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic("No storage found")
	}
	go func() {
		err = startTelegramBot(storage, token)
		if err != nil {
			fmt.Println(err)
			panic("Telegram bot failed")
		}
	}()
	select {}
}

func mustToken() string {
	token := os.Getenv("TOKEN")
	if token == "" {
		panic("No TOKEN environment variable")
	}
	return token
}

func startTelegramBot(db *postgres.Storage, token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		slog.Info("Received message from Telegram:", update.Message.Text, update.Message.Chat.ID)

		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text))

		err = handleMessage(db, int(update.Message.Chat.ID), "telegram", update.Message.Text)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleMessage(db *postgres.Storage, userID int, platform string, content string) error {
	const op = "app.handleMessage"

	var conversationID int
	_, err := db.UserExists(userID)
	if err != nil {
		platformID, err := db.PlatformID(platform)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		err = db.CreateUser(userID, platformID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		conversationID, err = db.CreateConversation(userID, managerID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	} else {
		conversationID, err = db.ConversationID(userID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	_, err = db.CreateMessage(content, conversationID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
