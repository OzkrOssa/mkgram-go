package commands

import (
	"github.com/OzkrOssa/mkgram-go/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartCommand struct{}

func (sc StartCommand) Execute(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	//TODO: add custom message to show all bot commands
	startText := "Welcome to Red Planet"
	message := tgbotapi.NewMessage(config.GroupChatID, startText)
	_, err := bot.Send(message)

	return err
}
