package cronjobs

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/OzkrOssa/mkgram-go/internal/config"
	"github.com/OzkrOssa/mkgram-go/internal/repository"
	"github.com/OzkrOssa/mkgram-go/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

type TrafficMonitor struct {
	bot *tgbotapi.BotAPI
}

func NewTrafficMonitor(bot *tgbotapi.BotAPI) *TrafficMonitor {
	return &TrafficMonitor{bot}
}

func (tm *TrafficMonitor) CheckTraffic() {
	var wg sync.WaitGroup
	providerConfig, err := config.LoadProviderConfig()
	resultsChanTraffic := make(chan repository.ChanTraffic, len(providerConfig.Providers))
	if err != nil {
		log.Println(err)
	}

	for _, provider := range providerConfig.Providers {
		wg.Add(1)
		go func(p config.ProviderData) {
			defer wg.Done()
			mk, err := repository.New(p.LocalAddress, "telegram-api", "1017230619", "8728")

			if err != nil {
				log.Println(err)
			}

			traffic, err := mk.GetTraffic(p.WAN)
			if err != nil {
				log.Println(err)
			}

			result := repository.ChanTraffic{
				Name:       p.Name,
				Saturation: p.Saturation,
				Traffic: repository.Traffic{
					Rx: traffic.Rx,
					Tx: traffic.Tx,
				},
			}
			resultsChanTraffic <- result

			Rx, err := strconv.Atoi(traffic.Rx)

			if err != nil {
				log.Println(err)
			}
			log.Println("Current RX: ", utils.FormatSize(int64(Rx)))
		}(provider)
	}
	wg.Wait()
	close(resultsChanTraffic)
	for ch := range resultsChanTraffic {
		Rx, err := strconv.Atoi(ch.Rx)

		if err != nil {
			log.Println(err)
		}
		switch {
		case int64(Rx) > ch.Saturation:
			textMessage := fmt.Sprintf("⚠️ <b><i>%s</i></b> supero el umbral de trafico de <b><i>%s</i></b> ⚠️", ch.Name, utils.FormatSize(ch.Saturation))
			message := tgbotapi.NewMessage(config.GroupChatID, textMessage)
			message.ParseMode = "Html"
			tm.bot.Send(message)
		case int64(Rx) < 100000000:
			textMessage := fmt.Sprintf("❌ El Trafico cayo a <b><i>%s</i></b> en <b><i>%s</i></b> ❌", utils.FormatSize(int64(Rx)), ch.Name)
			message := tgbotapi.NewMessage(config.GroupChatID, textMessage)
			message.ParseMode = "Html"
			tm.bot.Send(message)
		}
	}
}

func StartTrafficMonitorJob(bot *tgbotapi.BotAPI) {
	log.Println("Traffic Monitor Job Started")
	cron := cron.New()
	monitor := NewTrafficMonitor(bot)
	cron.AddFunc("* 6-23 * * *", monitor.CheckTraffic)
	cron.Start()

}
