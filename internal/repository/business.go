package repository

import (
	"database/sql"
	"log"
	"qarzi/internal/model"
)

type BusinessRepo struct {
	DB *sql.DB
}

func NewBusinessRepo(db *sql.DB) *BusinessRepo {
	return &BusinessRepo{DB: db}
}

func (r *BusinessRepo) CreateBusiness(b *model.Business) error {
	query := `
				INSERT INTO businesses (owner_id, name, phone)
				VALUES ($1, $2, $3)
				RETURNING id, created_at
		`
	err := r.DB.QueryRow(query, b.OwnerID, b.Name, b.Phone).Scan(&b.ID, &b.CreatedAt)
	if err != nil {
		log.Println("❌ Ошибка при создании бизнеса:", err)
		return err
	}

	return nil
}
