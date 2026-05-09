package worker

import (
	"fmt"
	"log"
	"time"

	"qarzi/internal/repository"

	"gopkg.in/telebot.v3"
)

func StartReminderWorker(bot *telebot.Bot, repo *repository.EventRepo) {

	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for range ticker.C {
			now := time.Now()

			if now.Hour() != 10 {
				continue
			}

			log.Println("⏰ Воркер: начинаю ежедневную проверку рассрочек...")

			payments, err := repo.GetUpcomingPayments()
			if err != nil {
				log.Println("❌ Ошибка при получении списка для напоминаний:", err)
				continue
			}

			for _, p := range payments {

				if p.TelegramID != nil {

					msg := fmt.Sprintf(
						"🔔 Напоминание о рассрочке!\n\nУважаемый(ая) %s, завтра наступает срок платежа.\nСумма: %d сум.",
						p.ClientName, p.Amount,
					)

					_, err := bot.Send(&telebot.User{ID: *p.TelegramID}, msg)
					if err != nil {
						log.Printf("⚠️ Не удалось отправить сообщение пользователю %d: %v\n", *p.TelegramID, err)
						continue
					}

					err = repo.MarkReminderAsSent(p.ID)
					if err != nil {
						log.Printf("❌ Ошибка при обновлении статуса напоминания ID %d: %v\n", p.ID, err)
					}
				}
			}
		}
	}()
}
