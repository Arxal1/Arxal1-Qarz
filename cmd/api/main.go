package main

import (
	"fmt"
	"log"
	"net/http"
	"qarzi/config"
	"qarzi/database"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg.Database)
	defer db.Close()

	database.RunMigrations(cfg.Database)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) //ошибка 500

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status": "ok", "project": "Qarzi"}`)
	})

	addr := ":" + cfg.App.Port
	log.Println("Сервер запущен на http://localhost:8080", addr, cfg.App.Env)
	log.Fatal(http.ListenAndServe(addr, r))

}
