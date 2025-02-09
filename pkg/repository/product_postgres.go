package repository

import (
	"errors"
	"github.com/Dau1to0v/fullstack-go/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ProductPostgres struct {
	db *sqlx.DB
}

func NewProductPostgres(db *sqlx.DB) *ProductPostgres {
	return &ProductPostgres{db: db}
}

func (r *ProductPostgres) Create(userId, warehouseId int, product models.Product) (int, error) {
	var id int

	query := "INSERT INTO products (name, quantity, price, category, description, image, user_id, warehouse_id ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"

	row := r.db.QueryRow(query, product.Name, product.Quantity, product.Price, product.Category, product.Description, product.Image, userId, warehouseId)
	err := row.Scan(&id)
	if err != nil {
		logrus.WithError(err).Error("Ошибка при добавлении продукта в базу данных")
		return 0, err
	}

	return id, nil
}

func (r *ProductPostgres) GetAll(userId, warehouseId int) ([]models.Product, error) {
	var products []models.Product

	query := "SELECT id, name, quantity, price, category, description, image, user_id, warehouse_id FROM products WHERE user_id = $1 AND warehouse_id = $2"
	err := r.db.Select(&products, query, userId, warehouseId)

	if err != nil {
		logrus.Printf("Ошибка при выполнении запроса: %v", err)
		return nil, err
	}

	return products, nil
}

func (r *ProductPostgres) Delete(userId, productId int) error {
	query := "DELETE FROM products WHERE id = $1 AND user_id = $2"

	result, err := r.db.Exec(query, productId, userId)
	if err != nil {
		return err
	}

	// Проверяем, был ли удалён хотя бы 1 товар
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product not found or access denied")
	}

	return nil
}
