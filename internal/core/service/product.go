package service

import (
	"strings"
	"time"

	"be-gin-go/internal/core/domain"
	"be-gin-go/internal/core/port"
)

type productService struct {
	repo port.ProductRepository
}

// NewProductService membuat service dengan dependency repository (interface).
func NewProductService(repo port.ProductRepository) port.ProductService {
	return &productService{repo: repo}
}

func (s *productService) List() ([]domain.Product, error) {
	return s.repo.FindAll()
}

func (s *productService) GetByID(id int64) (domain.Product, error) {
	return s.repo.FindByID(id)
}

func (s *productService) Create(name string, price float64, stock int) (domain.Product, error) {
	if err := validate(name, price, stock); err != nil {
		return domain.Product{}, err
	}

	now := time.Now()
	return s.repo.Create(domain.Product{
		Name:      strings.TrimSpace(name),
		Price:     price,
		Stock:     stock,
		CreatedAt: now,
		UpdatedAt: now,
	})
}

func (s *productService) Update(id int64, name string, price float64, stock int) (domain.Product, error) {
	if err := validate(name, price, stock); err != nil {
		return domain.Product{}, err
	}

	existing, err := s.repo.FindByID(id)
	if err != nil {
		return domain.Product{}, err
	}

	existing.Name = strings.TrimSpace(name)
	existing.Price = price
	existing.Stock = stock
	existing.UpdatedAt = time.Now()

	return s.repo.Update(existing)
}

func (s *productService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func validate(name string, price float64, stock int) error {
	if strings.TrimSpace(name) == "" || price < 0 || stock < 0 {
		return domain.ErrInvalidProduct
	}
	return nil
}
