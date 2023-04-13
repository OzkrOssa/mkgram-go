package commands

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/OzkrOssa/mkgram-go/internal/utils"

	"github.com/OzkrOssa/mkgram-go/internal/config"
	"github.com/OzkrOssa/mkgram-go/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type internetProviderCommand struct{}

type ProviderResult struct {
	Name      string
	Resources repository.Resources
	Traffic   repository.Traffic
}

func (ip *internetProviderCommand) Execute(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {

	providerConfig, err := config.LoadProviderConfig()

	resultsChan := make(chan ProviderResult, len(providerConfig.Providers))

	if err != nil {
		return err
	}

	for _, provider := range providerConfig.Providers {
		go func(p config.ProviderData) {
			mk, err := repository.New(p.LocalAddress, "telegram-api", "1017230619", "8728")
			if err != nil {
				log.Fatalln(err)
			}
			identity, err := mk.GetIndentity()
			if err != nil {
				log.Fatalln(err)
			}

			resources, err := mk.GetResources()
			if err != nil {
				log.Fatalln(err)
			}

			traffic, err := mk.GetTraffic(p.WAN)

			if err != nil {
				log.Fatalln(err)
			}

			resultsChan <- ProviderResult{
				Name:      identity,
				Resources: resources,
				Traffic:   traffic,
			}

		}(provider)
	}

	go func() {
		time.Sleep(time.Second) // esperar un poco para asegurarnos de que todos los resultados se hayan enviado
		close(resultsChan)
	}()

	var message tgbotapi.MessageConfig
	for ch := range resultsChan {

		tx, _ := strconv.Atoi(ch.Traffic.Tx)
		rx, _ := strconv.Atoi(ch.Traffic.Rx)

		log.Println(ch.Name, ch.Resources, ch.Traffic)
		textMessage := fmt.Sprintf("<b>%s</b>\n<b><i>Cpu:</i></b> %s <b><i>Uptime:</i></b> %s\n<b><i>Rx:</i></b> %s <b><i>Tx:</i></b> %s", ch.Name, ch.Resources.Cpu, ch.Resources.Uptime, utils.FormatSize(int64(rx)), utils.FormatSize(int64(tx)))
		message = tgbotapi.NewMessage(update.Message.Chat.ID, textMessage)
		message.ParseMode = "Html"
		_, err = bot.Send(message)
	}

	return err
}
