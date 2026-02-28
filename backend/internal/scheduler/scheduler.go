package scheduler

import (
	"log"

	"github.com/caner/home-gider/internal/services"
	"github.com/robfig/cron/v3"
)

func Start(recurringService *services.RecurringService) *cron.Cron {
	c := cron.New()

	// Her gün saat 00:05'te taksit ve sabit giderleri kontrol et
	c.AddFunc("5 0 * * *", func() {
		log.Println("Taksit/sabit gider kontrolü başlatıldı...")
		if err := recurringService.ProcessRecurring(); err != nil {
			log.Printf("Taksit/sabit gider işleme hatası: %v", err)
		} else {
			log.Println("Taksit/sabit gider kontrolü tamamlandı")
		}
	})

	c.Start()
	log.Println("Zamanlayıcı başlatıldı")
	return c
}
