package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BtsCommand struct{}

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Moron", "BTS_Moron"),
		tgbotapi.NewInlineKeyboardButtonData("Calera", "BTS_Calera"),
		tgbotapi.NewInlineKeyboardButtonData("Tabuyo", "BTS_Tabuyo"),
		tgbotapi.NewInlineKeyboardButtonData("Blandon", "BTS_Blandon"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Guamal", "BTS_Guamal"),
		tgbotapi.NewInlineKeyboardButtonData("Cabuyal", "BTS_Cabuyal"),
		tgbotapi.NewInlineKeyboardButtonData("Iberia", "BTS_Iberia"),
		tgbotapi.NewInlineKeyboardButtonData("Irra", "BTS_Irra"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Cruces", "BTS_Cruces"),
		tgbotapi.NewInlineKeyboardButtonData("Pueblo Viejo", "BTS_Pueblo_Viejo"),
		tgbotapi.NewInlineKeyboardButtonData("Clavijo", "BTS_Clavijo"),
		tgbotapi.NewInlineKeyboardButtonData("San lorenzo", "BTS_SL_200"),
		tgbotapi.NewInlineKeyboardButtonData("Irra Fibra", "BTS_Irra_200"),
	),
)

func (bc *BtsCommand) Execute(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

	msg.ReplyMarkup = numericKeyboard

	_, err := bot.Send(msg)

	return err
}
