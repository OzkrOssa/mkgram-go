package cronjobs

import (
	"fmt"
	"log"
	"strconv"
	"sync"

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
	var wg sync.WaitGroup
	providerConfig, err := config.LoadProviderConfig()
	resultsChanResources := make(chan repository.Resources, len(providerConfig.Providers))
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
	wg.Wait()
	close(resultsChanResources)
}

func StartResourcesMonitorJob(bot *tgbotapi.BotAPI) {
	log.Println("Resources Monitor Job Started")
	cron := cron.New()
	monitor := NewResourcesMonitor(bot)
	cron.AddFunc("* * * * *", monitor.CheckResources)
	cron.Start()

}
