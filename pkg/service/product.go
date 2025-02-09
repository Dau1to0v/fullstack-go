package service

import (
	"github.com/Dau1to0v/fullstack-go/models"
	"github.com/Dau1to0v/fullstack-go/pkg/repository"
)

type ProductService struct {
	repo          repository.Product
	warehouseRepo repository.Warehouse
}

func NewProductService(repo repository.Product, warehouseRepo repository.Warehouse) *ProductService {
	return &ProductService{repo: repo, warehouseRepo: warehouseRepo}
}

func (s *ProductService) Create(userId, warehouseId int, product models.Product) (int, error) {
	_, err := s.warehouseRepo.GetById(userId, warehouseId)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(userId, warehouseId, product)
}

func (s *ProductService) GetAll(userId, warehouseId int) ([]models.Product, error) {
	return s.repo.GetAll(userId, warehouseId)
}

func (s *ProductService) Delete(userId, productId int) error {
	return s.repo.Delete(userId, productId)
}
