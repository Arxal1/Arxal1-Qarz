package database

import (
	"fmt"
	"log"
	"qarzi/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(cfg config.DatabaseConfig) {

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	m, err := migrate.New("file://migration", dbURL)
	if err != nil {
		log.Fatal("❌ Ошибка инициализации миграций: ", err)
	}

	err = m.Up()
	if err != nil {

		if err == migrate.ErrNoChange {
			log.Println("ℹ️ База данных в актуальном состоянии, новые миграции не требуются.")
			return
		}
		log.Fatal("❌ Ошибка при применении миграций: ", err)
	}

	log.Println("✅ Миграции успешно применены! Таблицы созданы.")
}
