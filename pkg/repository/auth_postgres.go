package repository

import (
	"fmt"
	"github.com/Dau1to0v/fullstack-go/models"
	"github.com/jmoiron/sqlx"
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
	query := "SELECT id, username, email FROM users WHERE id = $1"

	err := r.db.Get(&user, query, userId)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
