package service

import (
	"github.com/Dau1to0v/fullstack-go/models"
	"github.com/Dau1to0v/fullstack-go/pkg/repository"
)

type WarehouseService struct {
	repo repository.Warehouse
}

func NewWarehouseService(repo repository.Warehouse) *WarehouseService {
	return &WarehouseService{repo: repo}
}

func (s *WarehouseService) Create(userId int, warehouse models.Warehouse) (int, error) {
	return s.repo.Create(userId, warehouse)
}

func (s *WarehouseService) GetAll(userId int) ([]models.Warehouse, error) {
	return s.repo.GetAll(userId)
}

func (s *WarehouseService) Delete(userId, warehouseId int) error {
	return s.repo.Delete(userId, warehouseId)
}

func (s *WarehouseService) Update(userId, warehouseId int, input models.UpdateWarehouseInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, warehouseId, input)
}

func (s *WarehouseService) GetById(userId, warehouseId int) (models.Warehouse, error) {
	return s.repo.GetById(userId, warehouseId)
}

func (s *WarehouseService) CalculateWarehousesValue(userId int) ([]models.WarehouseNetWorth, error) {
	return s.repo.CalculateWarehousesValue(userId)
}
