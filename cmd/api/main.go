package main

import (
	"fmt"
	"log"
	"net/http"
	"qarzi/config"
	"qarzi/database"
	"qarzi/internal/delivery/bot"
	"qarzi/internal/handler"
	"qarzi/internal/repository"
	"qarzi/internal/worker"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg.Database)
	defer db.Close()

	database.RunMigrations(cfg.Database)

	userRepo := repository.NewUserRepo(db)
	userHandler := handler.NewUserHandler(userRepo)

	//бизнес
	businessRepo := repository.NewBusinessRepo(db)
	businessHandler := handler.NewBusinessHandler(businessRepo)
	//клиент
	clientRepo := repository.NewClientRepo(db)
	clientHandler := handler.NewClientHandler(clientRepo)
	//событие
	eventRepo := repository.NewEventRepo(db)
	eventHandler := handler.NewEventHandler(eventRepo)

	tgBot, err := bot.NewBot(cfg.Telegram.BotToken, userRepo)

	if err != nil {
		log.Fatal("❌ Ошибка при инициализации Telegram бота: ", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) //ошибка 500

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status": "ok", "project": "Qarzi"}`)
	})

	r.Post("/api/users/register", userHandler.Register)
	r.Post("/api/businesses/register", businessHandler.RegisterBusiness)
	r.Post("/api/clients/register", clientHandler.CreateClient)
	r.Post("/api/events/shipment", eventHandler.CreateShipment)

	addr := ":" + cfg.App.Port

	go func() {
		log.Printf("HTTP-сервер запущен на http://localhost%s (режим: %s)", addr, cfg.App.Env)
		if err := http.ListenAndServe(addr, r); err != nil {
			log.Fatal("❌ Ошибка при запуске HTTP-сервера: ", err)
		}
	}()
	worker.StartReminderWorker(tgBot.Bot, eventRepo)
	tgBot.Start()
}
