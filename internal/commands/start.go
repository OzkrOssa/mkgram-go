package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type startCommand struct{}

func (sc startCommand) Execute(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	//TODO: add custom message to show all bot commands
	startText := "Welcome to Red Planet"
	message := tgbotapi.NewMessage(update.Message.Chat.ID, startText)
	_, err := bot.Send(message)

	return err
}
