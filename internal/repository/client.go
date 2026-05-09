package repository

import (
	"database/sql"
	"log"
	"qarzi/internal/model"
)

type ClientRepo struct {
	DB *sql.DB
}

func NewClientRepo(db *sql.DB) *ClientRepo {
	return &ClientRepo{DB: db}
}

func (r *ClientRepo) CreateClient(c *model.Client) error {
	query := `
		INSERT INTO clients (business_id, name, phone, debt_limit)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	err := r.DB.QueryRow(query, c.BusinessID, c.Name, c.Phone, c.DebtLimit).Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		log.Println("❌ Ошибка при создании клиента:", err)
		return err
	}
	return nil
}
