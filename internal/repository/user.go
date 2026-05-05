package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"qarzi/internal/model"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) CreateUser(u *model.User) error {

	query := `
		INSERT INTO users (telegram_id, full_name)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	err := r.DB.QueryRow(query, u.TelegramID, u.FullName).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		log.Println("❌ Ошибка при сохранении пользователя в БД:", err)
		return fmt.Errorf("❌ не удалось создать пользователя: %w", err)
	}

	return nil

}

func (r *UserRepo) GetUserByTelegramID(ctx context.Context, telegramID int64) (*model.User, error) {
	query := `
		SELECT id, telegram_id, phone, full_name, role, is_identified, pinfl, created_at
		FROM users
		WHERE telegram_id = $1
	`

	var user model.User

	err := r.DB.QueryRowContext(ctx, query, telegramID).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Phone,
		&user.FullName,
		&user.Role,
		&user.IsIdentified,
		&user.PINFL,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		log.Printf("❌ Ошибка при поиске пользователя (ID: %d): %v\n", telegramID, err)
		return nil, err
	}

	return &user, nil
}
