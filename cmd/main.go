package main

import (
	"log"

	"github.com/OzkrOssa/mkgram-go/internal/commands"
	"github.com/OzkrOssa/mkgram-go/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	bot, err := tgbotapi.NewBotAPI(config.BotToken)

	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)

	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	commandHandler := commands.CommandHandler{}

	for update := range updates {
		if update.Message != nil {
			err := commandHandler.HandlerCommands(bot, &update)
			if err != nil {
				log.Println(err)
			}
		}
	}

}
