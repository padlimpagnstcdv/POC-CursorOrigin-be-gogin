package port

import "be-gin-go/internal/core/domain"

// ProductRepository adalah kontrak yang harus dipenuhi oleh data layer.
type ProductRepository interface {
	FindAll() ([]domain.Product, error)
	FindByID(id int64) (domain.Product, error)
	Create(product domain.Product) (domain.Product, error)
	Update(product domain.Product) (domain.Product, error)
	Delete(id int64) error
}

// ProductService adalah kontrak yang dipakai oleh presentation layer.
type ProductService interface {
	List() ([]domain.Product, error)
	GetByID(id int64) (domain.Product, error)
	Create(name string, price float64, stock int) (domain.Product, error)
	Update(id int64, name string, price float64, stock int) (domain.Product, error)
	Delete(id int64) error
}
