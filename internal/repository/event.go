package repository

import (
	"database/sql"
	"log"
	"qarzi/internal/model"
)

type EventRepo struct {
	DB *sql.DB
}

func NewEventRepo(db *sql.DB) *EventRepo {
	return &EventRepo{DB: db}
}

func (r *EventRepo) CreateShipment(e *model.Event) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	insertEventQuery := `
			INSERT INTO events (client_id, event_type, amount, status, description, created_by)
			VALUES ($1, 'shipment', $2, 'confirmed', $3, $4)
			RETURNING id, created_at, updated_at
	`
	var currentDebt int64

	err = tx.QueryRow(`SELECT debt_limit FROM clients WHERE id = $1`, insertEventQuery, e.ClientID, e.Amount, e.Description, e.CreatedBy).Scan(&e.ID, &currentDebt, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		log.Println("❌ Ошибка при записи события отгрузки:", err)
		return err
	}

	updateClientQuery := `
		UPDATE clients
		SET debt_limit = debt_limit + $1
		WHERE id = $2
	`
	_, err = tx.Exec(updateClientQuery, e.Amount, e.ClientID)
	if err != nil {
		log.Println("❌ Ошибка при обновлении баланса клиента:", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil

}
