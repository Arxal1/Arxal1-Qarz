package bot

import (
	"context"
	"log"

	"gopkg.in/telebot.v3"
)

func (b *Bot) AuthMiddleware() telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			tgID := c.Sender().ID

			user, err := b.UserRepo.GetUserByTelegramID(context.Background(), tgID)
			if err != nil {
				log.Printf("❌ Ошибка при получении пользователя (Telegram ID: %d): %v\n", tgID, err)
				return c.Send("Произошла ошибка сервера. Пожалуйста, попробуйте позже.")
			}

			c.Set("user", user)

			return next(c)
		}
	}
}
