package service

import (
	"github.com/Dau1to0v/fullstack-go/models"
	"github.com/Dau1to0v/fullstack-go/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Warehouse interface {
	Create(userId int, warehouse models.Warehouse) (int, error)
	GetAll(userId int) ([]models.Warehouse, error)
	GetById(userId, warehouseId int) (models.Warehouse, error)
	Delete(userId, warehouseId int) error
	Update(userId, warehouseId int, input models.UpdateWarehouseInput) error
}

type Product interface {
	Create(userId, warehouseId int, product models.Product) (int, error)
	GetAll(userId, warehouseId int) ([]models.Product, error)
	Delete(userId, productId int) error
}
type Service struct {
	Authorization
	Warehouse
	Product
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Warehouse:     NewWarehouseService(repos.Warehouse),
		Product:       NewProductService(repos.Product, repos.Warehouse),
	}
}
