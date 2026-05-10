package repository

import (
	"database/sql"
	"fmt"
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

func (r *EventRepo) InitiatePayment(e *model.Event) error {
	query := `
		INSERT INTO events (client_id, event_type, amount, status, description, created_by)
		VALUES ($1, 'payment_cash', $2, 'pending', $3, $4)
		RETURNING id, created_at, updated_at
	`
	err := r.DB.QueryRow(query, e.ClientID, e.Amount, e.Description, e.CreatedBy).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		log.Println("❌ Ошибка при инициации оплаты:", err)
		return err
	}
	return nil
}

func (r *EventRepo) ConfirmPayment(eventID int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var amount int64
	var clientID int
	var currentStatus string

	err = tx.QueryRow(`SELECT amount, client_id, status FROM events WHERE id = $1 FOR UPDATE`, eventID).Scan(&amount, &clientID, &currentStatus)
	if err != nil {
		return err
	}

	if currentStatus != "pending" {
		return fmt.Errorf("оплата уже обработана или отменена")
	}

	_, err = tx.Exec(`UPDATE events SET status = 'confirmed', updated_at = CURRENT_TIMESTAMP WHERE id = $1`, eventID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE clients SET debt_limit = debt_limit - $1 WHERE id = $2`, amount, clientID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
func (r *EventRepo) DisputeEvent(eventID int, clientID int) error {

	query := `
		UPDATE events 
		SET status = 'disputed', updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1 AND client_id = $2
	`
	res, err := r.DB.Exec(query, eventID, clientID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("событие не найдено или не принадлежит этому клиенту")
	}

	return nil
}

func (r *EventRepo) AdjustEvent(originalEventID int, newAmount int64, userID int, comment string) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var oldAmount int64
	var clientID int
	err = tx.QueryRow(`SELECT amount, client_id FROM events WHERE id = $1 FOR UPDATE`, originalEventID).Scan(&oldAmount, &clientID)
	if err != nil {
		return fmt.Errorf("оригинальное событие не найдено: %v", err)
	}

	difference := newAmount - oldAmount
	if difference == 0 {
		return fmt.Errorf("новая сумма совпадает со старой")
	}

	var adjustmentEventID int
	err = tx.QueryRow(`
		INSERT INTO events (client_id, event_type, amount, status, description, created_by)
		VALUES ($1, 'adjustment', $2, 'confirmed', $3, $4)
		RETURNING id
	`, clientID, difference, comment, userID).Scan(&adjustmentEventID)
	if err != nil {
		return fmt.Errorf("ошибка при создании корректировки: %v", err)
	}

	_, err = tx.Exec(`UPDATE clients SET debt_limit = debt_limit + $1 WHERE id = $2`, difference, clientID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении баланса: %v", err)
	}

	oldVals := fmt.Sprintf(`{"amount": %d}`, oldAmount)
	newVals := fmt.Sprintf(`{"amount": %d, "adjustment_event_id": %d}`, newAmount, adjustmentEventID)

	_, err = tx.Exec(`
		INSERT INTO audit_logs (user_id, action, entity_type, entity_id, old_values, new_values)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, userID, "adjustment", "events", originalEventID, oldVals, newVals)
	if err != nil {
		return fmt.Errorf("ошибка при записи в лог аудита: %v", err)
	}

	return tx.Commit()
}
