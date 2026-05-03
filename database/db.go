package database

import (
	"database/sql"
	"fmt"
	"log"
	"qarzi/config"

	_ "github.com/lib/pq"
)

func Connect(cfg config.DatabaseConfig) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Ошибка инициализации базы данных:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("База данных не отвечает:", err)
	}

	log.Println("✅ Успешное подключение к PostgreSQL!")
	return db
}
