package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"be-gin-go/internal/core/domain"
	"be-gin-go/internal/core/port"
	"be-gin-go/internal/presentation/http/dto"
)

// ProductHandler menerjemahkan HTTP request ke core layer.
type ProductHandler struct {
	service port.ProductService
}

func NewProductHandler(service port.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// List GET /api/products
func (h *ProductHandler) List(c *gin.Context) {
	products, err := h.service.List()
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "success", Data: products})
}

// Detail GET /api/products/:id
func (h *ProductHandler) Detail(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Message: "invalid id"})
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "success", Data: product})
}

// Create POST /api/products
func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Message: err.Error()})
		return
	}

	product, err := h.service.Create(req.Name, req.Price, req.Stock)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.Response{Message: "product created", Data: product})
}

// Update PUT /api/products/:id
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Message: "invalid id"})
		return
	}

	var req dto.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Message: err.Error()})
		return
	}

	product, err := h.service.Update(id, req.Name, req.Price, req.Stock)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "product updated", Data: product})
}

// Delete DELETE /api/products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Message: "invalid id"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "product deleted"})
}

func parseID(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.Param("id"), 10, 64)
}

// writeError memetakan error domain ke status HTTP.
func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		c.JSON(http.StatusNotFound, dto.Response{Message: err.Error()})
	case errors.Is(err, domain.ErrInvalidProduct):
		c.JSON(http.StatusBadRequest, dto.Response{Message: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, dto.Response{Message: "internal server error"})
	}
}
