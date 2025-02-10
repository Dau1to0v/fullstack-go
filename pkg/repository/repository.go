package repository

import (
	"github.com/Dau1to0v/fullstack-go/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
	GetUserById(id int) (models.User, error)
	UpdateUser(userId int, input models.UpdateUserInput) error
}

type Warehouse interface {
	Create(userId int, warehouse models.Warehouse) (int, error)
	GetAll(userId int) ([]models.Warehouse, error)
	GetById(userId, warehouseId int) (models.Warehouse, error)
	Delete(userId, warehouseId int) error
	Update(userId int, warehouseId int, input models.UpdateWarehouseInput) error
	CalculateWarehousesValue() ([]models.WarehouseNetWorth, error)
}

type Product interface {
	Create(userId, warehouseId int, product models.Product) (int, error)
	GetAll(userId, warehouseId int) ([]models.Product, error)
	GetById(userId, productId int) (models.Product, error)
	Delete(userId, productId int) error
	Update(userId, productId int, input models.UpdateProductInput) error
}

type Repository struct {
	Authorization
	Warehouse
	Product
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Warehouse:     NewWarehousePostgres(db),
		Product:       NewProductPostgres(db),
	}
}
