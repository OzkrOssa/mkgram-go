package cronjobs

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/OzkrOssa/mkgram-go/internal/config"
	"github.com/OzkrOssa/mkgram-go/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
)

type ResourcesMonitor struct {
	bot *tgbotapi.BotAPI
}

func NewResourcesMonitor(bot *tgbotapi.BotAPI) *ResourcesMonitor {
	return &ResourcesMonitor{bot}
}

func (tm *ResourcesMonitor) CheckResources() {
	providerConfig, err := config.LoadProviderConfig()
	resultsChanResources := make(chan repository.Resources, len(providerConfig.Providers))
	if err != nil {
		log.Println(err)
	}

	for _, provider := range providerConfig.Providers {

		go func(p config.ProviderData) {

			mk, err := repository.New(p.LocalAddress, "telegram-api", "1017230619", "8728")

			if err != nil {
				log.Println(err)
			}

			resources, err := mk.GetResources()
			if err != nil {
				log.Println(err)
			}

			resultsChanResources <- resources

			cpu, err := strconv.Atoi(resources.Cpu)

			if err != nil {
				log.Println(err)
			}

			if cpu > 70 {
				log.Printf("Current CPU load: %d", cpu)
				textMessage := fmt.Sprintf("⚡ La CPU en <b><i>%s</i></b> supero el <b><i>%d</i></b> ⚡", p.Name, cpu)
				message := tgbotapi.NewMessage(config.GroupChatID, textMessage)
				message.ParseMode = "Html"
				tm.bot.Send(message)
			}
		}(provider)
	}
	go func() {
		time.Sleep(time.Second)
		close(resultsChanResources)
	}()
}

func StartResourcesMonitorJob(bot *tgbotapi.BotAPI) {
	log.Println("Resources Monitor Job Started")
	cron := cron.New()
	monitor := NewResourcesMonitor(bot)
	cron.AddFunc("* * * * *", monitor.CheckResources)
	cron.Start()

}
