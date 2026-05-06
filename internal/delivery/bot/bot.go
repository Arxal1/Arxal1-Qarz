package bot

import (
	"log"
	"qarzi/internal/repository"
	"time"

	"gopkg.in/telebot.v3"
)

type Bot struct {
	Bot      *telebot.Bot
	UserRepo *repository.UserRepo
	// потом тут будет бизнес репо и ивентрепо
}

func NewBot(token string, userRepo *repository.UserRepo) (*Bot, error) {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Bot:      b,
		UserRepo: userRepo,
	}, nil
}

func (b *Bot) Start() {
	b.setupRoutes()
	log.Println("Telegram бот успешно запущен")
	b.Bot.Start()
}

func (b *Bot) setupRoutes() {
	b.Bot.Use(b.AuthMiddleware())
	b.Bot.Handle("/start", func(c telebot.Context) error {
		user := c.Get("user")
		menu := &telebot.ReplyMarkup{}
		btnWebApp := menu.WebApp("🚀 Открыть Qarzi", &telebot.WebApp{URL: "https://qarzi.uz"})
		menu.Inline(menu.Row(btnWebApp))

		if user == nil {
			return c.Send("Привет! Я бот для управления твоими долгами. Пожалуйста, зарегистрируйся, чтобы начать использовать бота.")
		}
		return c.Send("С возвращением! Управляй своими клиентами прямо в приложении:", menu)
	})
}
