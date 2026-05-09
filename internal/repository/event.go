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

// CreateShipment создает событие "отгрузка" и СРАЗУ увеличивает долг клиента
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

	// Выполняем именно INSERT запрос
	err = tx.QueryRow(insertEventQuery, e.ClientID, e.Amount, e.Description, e.CreatedBy).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
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

func (r *EventRepo) GetUpcomingPayments() ([]model.InstallmentWithClient, error) {
	query := `
		SELECT i.id, i.amount, i.due_date, c.phone, c.name, u.telegram_id
		FROM installments i
		JOIN events e ON i.event_id = e.id
		JOIN clients c ON e.client_id = c.id
		LEFT JOIN users u ON c.user_id = u.id
		WHERE i.due_date = CURRENT_DATE + INTERVAL '1 day'
		  AND i.status = 'pending'
		  AND i.reminder_sent = FALSE
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []model.InstallmentWithClient

	for rows.Next() {
		var p model.InstallmentWithClient
		err := rows.Scan(&p.ID, &p.Amount, &p.DueDate, &p.ClientPhone, &p.ClientName, &p.TelegramID)
		if err != nil {
			log.Println("❌ Ошибка при чтении данных рассрочки:", err)
			continue
		}
		payments = append(payments, p)
	}

	return payments, nil
}

func (r *EventRepo) MarkReminderAsSent(installmentID int) error {
	query := `UPDATE installments SET reminder_sent = TRUE WHERE id = $1`
	_, err := r.DB.Exec(query, installmentID)
	return err
}

func (r *EventRepo) GetClientEvents(clientID int) ([]model.Event, error) {
	query := `
		SELECT id, client_id, event_type, amount, status, description, created_by, created_at, updated_at
		FROM events
		WHERE client_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var e model.Event
		err := rows.Scan(&e.ID, &e.ClientID, &e.EventType, &e.Amount, &e.Status, &e.Description, &e.CreatedBy, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			log.Println("❌ Ошибка при чтении события:", err)
			continue
		}
		events = append(events, e)
	}

	if events == nil {
		events = []model.Event{}
	}

	return events, nil
}
