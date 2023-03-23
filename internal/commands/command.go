package commands

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Command interface {
	Execute(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error
}
