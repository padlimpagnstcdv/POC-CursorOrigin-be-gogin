package dto

// ProductRequest adalah payload untuk create & update product.
type ProductRequest struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"gte=0"`
	Stock int     `json:"stock" binding:"gte=0"`
}

// Response adalah bentuk response standar API.
type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
