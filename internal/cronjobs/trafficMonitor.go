package cronjobs

import (
	"fmt"
	"log"
	"strconv"
	"time"

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
	providerConfig, err := config.LoadConfig()
	resultsChan := make(chan repository.Traffic, len(providerConfig.Providers))
	if err != nil {
		log.Println(err)
	}

	for _, provider := range providerConfig.Providers {

		go func(p config.ProviderData) {

			mk, err := repository.New(p.LocalAddress, "telegram-api", "1017230619", "8728")

			if err != nil {
				log.Println(err)
			}

			traffic, err := mk.GetTraffic(p.WAN)
			if err != nil {
				log.Println(err)
			}

			resultsChan <- traffic

			Rx, err := strconv.Atoi(traffic.Rx)

			if err != nil {
				log.Println(err)
			}

			switch {

			case int64(Rx) > p.Saturation:
				textMessage := fmt.Sprintf("⚠️ <b><i>%s</i></b> supero el umbral de trafico de <b><i>%s</i></b> ⚠️", p.Name, utils.FormatSize(p.Saturation))
				message := tgbotapi.NewMessage(config.GroupChatID, textMessage)
				message.ParseMode = "Html"
				tm.bot.Send(message)
			case int64(Rx) < 100000000:
				textMessage := fmt.Sprintf("❌ Cayo el trafico a menos de <b><i>%s</i></b> en <b><i>%s</i></b> ❌", utils.FormatSize(int64(Rx)), p.Name)
				message := tgbotapi.NewMessage(config.GroupChatID, textMessage)
				message.ParseMode = "Html"
				tm.bot.Send(message)
			}
		}(provider)
	}
	go func() {
		time.Sleep(time.Second)
		close(resultsChan)
	}()
}

func StartTrafficMonitorJob(bot *tgbotapi.BotAPI) {
	log.Println("Traffic Monitor Job Started")
	cron := cron.New()
	monitor := NewTrafficMonitor(bot)
	cron.AddFunc("* * * * *", monitor.CheckTraffic)
	cron.Start()

}
