package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Dau1to0v/fullstack-go/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"log"
	"strings"
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

func (r *ProductPostgres) GetById(userId, productId int) (models.Product, error) {
	var product models.Product

	query := "SELECT id, name, quantity, price, category, description, image, user_id, warehouse_id, created_at FROM products WHERE id = $1 AND user_id = $2"

	err := r.db.Get(&product, query, productId, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return product, errors.New("product not found or access denied")
		}
		return product, err
	}
	return product, nil
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

func (r *ProductPostgres) Update(userId, productId int, input models.UpdateProductInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name = $%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Quantity != nil {
		setValues = append(setValues, fmt.Sprintf("quantity = $%d", argId))
		args = append(args, *input.Quantity)
		argId++
	}

	if input.Price != nil {
		setValues = append(setValues, fmt.Sprintf("price = $%d", argId))
		args = append(args, *input.Price)
		argId++
	}

	if input.Category != nil {
		setValues = append(setValues, fmt.Sprintf("category = $%d", argId))
		args = append(args, *input.Category)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Image != nil {
		setValues = append(setValues, fmt.Sprintf("image = $%d", argId))
		args = append(args, *input.Image)
		argId++
	}

	if len(setValues) == 0 {
		return errors.New("empty update fields")
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE products SET %s WHERE id = $%d AND user_id = $%d", setQuery, argId, argId+1)
	args = append(args, productId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %v", args)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product not found or access denied")
	}

	return nil
}

func (r *ProductPostgres) Search(warehouseId int, text, searchType string) ([]models.Product, error) {
	var products []models.Product

	var query string
	if searchType == "category" {
		query = `
			SELECT id, name, quantity, price, category, user_id, description, warehouse_id, image
			FROM products
			WHERE warehouse_id = $1
			AND (COALESCE($2, '') = '' OR category = $2);
		`
	} else {
		query = `
			SELECT id, name, quantity, price, category, user_id, description, warehouse_id, image
			FROM products
			WHERE warehouse_id = $1
			AND (COALESCE($2, '') = '' OR name ILIKE '%' || $2 || '%' OR description ILIKE '%' || $2 || '%');
		`
	}

	err := r.db.Select(&products, query, warehouseId, text)
	if err != nil {
		log.Printf("Ошибка выполнения SQL запроса: %v", err)
		return nil, err
	}

	return products, nil
}
