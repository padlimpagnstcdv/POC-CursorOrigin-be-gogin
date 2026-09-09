package repository

import (
	"sort"
	"sync"
	"time"

	"be-gin-go/internal/core/domain"
	"be-gin-go/internal/core/port"
)

// productMemoryRepository menyimpan data di memory (data dummy).
type productMemoryRepository struct {
	mu     sync.RWMutex
	data   map[int64]domain.Product
	lastID int64
}

// NewProductMemoryRepository membuat repository beserta data dummy awal.
func NewProductMemoryRepository() port.ProductRepository {
	repo := &productMemoryRepository{
		data: make(map[int64]domain.Product),
	}

	now := time.Now()
	seeds := []domain.Product{
		{Name: "Kopi Arabica 250gr", Price: 85000, Stock: 25},
		{Name: "Teh Hijau Premium", Price: 45000, Stock: 40},
		{Name: "Gula Aren Cair", Price: 30000, Stock: 15},
	}

	for _, seed := range seeds {
		repo.lastID++
		seed.ID = repo.lastID
		seed.CreatedAt = now
		seed.UpdatedAt = now
		repo.data[seed.ID] = seed
	}

	return repo
}

func (r *productMemoryRepository) FindAll() ([]domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]domain.Product, 0, len(r.data))
	for _, product := range r.data {
		products = append(products, product)
	}

	sort.Slice(products, func(i, j int) bool { return products[i].ID < products[j].ID })
	return products, nil
}

func (r *productMemoryRepository) FindByID(id int64) (domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, ok := r.data[id]
	if !ok {
		return domain.Product{}, domain.ErrProductNotFound
	}
	return product, nil
}

func (r *productMemoryRepository) Create(product domain.Product) (domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.lastID++
	product.ID = r.lastID
	r.data[product.ID] = product

	return product, nil
}

func (r *productMemoryRepository) Update(product domain.Product) (domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[product.ID]; !ok {
		return domain.Product{}, domain.ErrProductNotFound
	}
	r.data[product.ID] = product

	return product, nil
}

func (r *productMemoryRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return domain.ErrProductNotFound
	}
	delete(r.data, id)

	return nil
}
