package router

import (
	"github.com/gin-gonic/gin"

	"be-gin-go/internal/presentation/http/handler"
)

// New mendaftarkan semua route aplikasi.
func New(productHandler *handler.ProductHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	products := r.Group("/api/products")
	{
		products.GET("", productHandler.List)
		products.GET("/:id", productHandler.Detail)
		products.POST("", productHandler.Create)
		products.PUT("/:id", productHandler.Update)
		products.DELETE("/:id", productHandler.Delete)
	}

	return r
}
