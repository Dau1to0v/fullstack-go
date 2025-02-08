package repository

import (
	"errors"
	"fmt"
	"github.com/Dau1to0v/fullstack-go/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
)

type WarehousePostgres struct {
	db *sqlx.DB
}

func NewWarehousePostgres(db *sqlx.DB) *WarehousePostgres {
	return &WarehousePostgres{db: db}
}

func (r *WarehousePostgres) Create(userId int, warehouse models.Warehouse) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO warehouses (name, location, user_id) VALUES ($1, $2, $3) RETURNING id")

	row := r.db.QueryRow(query, warehouse.Name, warehouse.Location, userId)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (r *WarehousePostgres) GetAll(userId int) ([]models.Warehouse, error) {
	var warehouse []models.Warehouse

	query := "SELECT id, name, location, user_id FROM warehouses WHERE user_id = $1"
	err := r.db.Select(&warehouse, query, userId)

	if err != nil {
		logrus.Printf("Ошибка при выполнении запроса: %v", err)
		return nil, err
	}

	return warehouse, err
}

func (r *WarehousePostgres) GetById(userId, warehouseId int) (models.Warehouse, error) {
	var warehouse models.Warehouse

	query := "SELECT id, name, location, user_id, created_at FROM warehouses WHERE id = $1 AND user_id = $2"

	// Логируем запрос перед выполнением
	logrus.Debugf("Executing query: %s with warehouseId=%d, userId=%d", query, warehouseId, userId)

	err := r.db.Get(&warehouse, query, warehouseId, userId)
	if err != nil {
		logrus.Errorf("Error fetching updated warehouse: %v", err)
		return warehouse, err
	}

	return warehouse, nil
}

func (r *WarehousePostgres) Delete(userId, warehouseId int) error {
	query := "DELETE FROM warehouses WHERE id = $1 AND user_id = $2"

	result, err := r.db.Exec(query, warehouseId, userId)
	if err != nil {
		return err
	}

	// Проверяем, был ли удалён хотя бы 1 склад
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("warehouse not found or access denied")
	}

	return nil
}

func (r *WarehousePostgres) Update(userId, warehouseId int, input models.UpdateWarehouseInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name = $%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Location != nil {
		setValues = append(setValues, fmt.Sprintf("location = $%d", argId))
		args = append(args, *input.Location)
		argId++
	}

	if len(setValues) == 0 {
		return errors.New("empty update fields")
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE warehouses SET %s WHERE id = $%d AND user_id = $%d", setQuery, argId, argId+1)
	args = append(args, warehouseId, userId) // ✅ Добавляем warehouseId и userId в аргументы

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %v", args) // ✅ Исправлено, теперь правильно печатает аргументы

	_, err := r.db.Exec(query, args...)
	return err
}
