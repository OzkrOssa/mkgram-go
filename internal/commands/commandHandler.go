package commands

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type CommandHandler struct {
}

func (ch CommandHandler) HandlerCommands(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	if update.Message.IsCommand() {
		commands := map[string]Command{
			"start":      &startCommand{},
			"provedores": &internetProviderCommand{},
			"bts":        &BtsCommand{},
		}

		command := update.Message.Command()

		if cmd, ok := commands[command]; ok {
			err := cmd.Execute(bot, update)

			if err != nil {
				return err
			}
		}

	}
	return nil
}
