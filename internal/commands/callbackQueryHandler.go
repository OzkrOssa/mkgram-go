package commands

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/OzkrOssa/mkgram-go/internal/config"
	"github.com/OzkrOssa/mkgram-go/internal/repository"
	"github.com/OzkrOssa/mkgram-go/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackQueryHandler struct{}

type btsResult struct {
	Name      string
	Resources repository.Resources
	Traffic   repository.Traffic
}

func (ch *CallbackQueryHandler) HandlerCallbackQuery(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	var wg sync.WaitGroup
	var message tgbotapi.MessageConfig
	btsConfig, err := config.LoadBtsConfig()
	if err != nil {
		return err
	}

	resultChan := make(chan btsResult, len(btsConfig.Bts))

	for _, bts := range btsConfig.Bts {
		wg.Add(1)
		go func(b config.BtsData) {
			defer wg.Done()
			mk, err := repository.New(b.LocalAddress, os.Getenv("TELEGRAM_USER"), os.Getenv("TELEGRAM_PASSWORD"), "8728")
			if err != nil {
				log.Println("Login", err)
			}
			identity, err := mk.GetIndentity()
			if err != nil {
				log.Println("Identity", err)
			}

			resources, err := mk.GetResources()
			if err != nil {
				log.Println("Resouces", err)
			}

			traffic, err := mk.GetTraffic(b.WAN)

			if err != nil {
				log.Println("Traffic", err)
			}

			resultChan <- btsResult{
				Name:      identity,
				Resources: resources,
				Traffic:   traffic,
			}
		}(bts)
	}
	wg.Wait()
	close(resultChan)

	for ch := range resultChan {
		if update.CallbackQuery.Data == ch.Name {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				log.Println("Request", err)
			}
			tx, _ := strconv.Atoi(ch.Traffic.Tx)
			rx, _ := strconv.Atoi(ch.Traffic.Rx)
			textMessage := fmt.Sprintf("<b>%s</b>\n<b><i>Cpu:</i></b> %s <b><i>Uptime:</i></b> %s\n<b><i>Rx:</i></b> %s <b><i>Tx:</i></b> %s", ch.Name, ch.Resources.Cpu, ch.Resources.Uptime, utils.FormatSize(int64(rx)), utils.FormatSize(int64(tx)))
			message = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, textMessage)
			message.ParseMode = "Html"
			_, err = bot.Send(message)
			if err != nil {
				return err
			}
		}

	}

	return err
}
