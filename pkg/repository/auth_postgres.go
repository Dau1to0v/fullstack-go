package repository

import (
	"errors"
	"fmt"
	"github.com/Dau1to0v/fullstack-go/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO users (username, email, password_hash, created_at) VALUES ($1, $2, $3, $4) RETURNING id")

	row := r.db.QueryRow(query, user.Username, user.Email, user.Password, user.CreatedAt)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (models.User, error) {
	var user models.User

	query := fmt.Sprintf("SELECT id FROM users WHERE username = $1 AND password_hash = $2")
	err := r.db.Get(&user, query, username, password)

	return user, err
}

func (r *AuthPostgres) GetUserById(userId int) (models.User, error) {
	var user models.User
	query := "SELECT id, username, email, password_hash FROM users WHERE id = $1"

	err := r.db.Get(&user, query, userId)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *AuthPostgres) UpdateUser(userId int, input models.UpdateUserInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Username != nil {
		setValues = append(setValues, fmt.Sprintf("username = $%d", argId)) // Исправлено имя поля
		args = append(args, *input.Username)
		argId++
	}

	if input.Email != nil {
		setValues = append(setValues, fmt.Sprintf("email = $%d", argId))
		args = append(args, *input.Email)
		argId++
	}

	if len(setValues) == 0 {
		return errors.New("empty update fields")
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", setQuery, argId) // Исправлен WHERE

	args = append(args, userId) // Добавляем userId в аргументы

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *AuthPostgres) UpdatePassword(userId int, newPassword string) error {
	query := `UPDATE users SET password_hash = $1 WHERE id = $2`

	_, err := r.db.Exec(query, newPassword, userId)
	return err
}
