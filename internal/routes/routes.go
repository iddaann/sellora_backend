package routes

import (
	"github.com/gin-gonic/gin"

	"sellora_backend/internal/handlers"
)

func Setup(
	r *gin.Engine,
	categoryHandler *handlers.CategoryHandler,
	productHandler *handlers.ProductHandler,
	transactionHandler *handlers.TransactionHandler,
) {
	api := r.Group("/api")
	{
		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAll)
			categories.POST("", categoryHandler.Create)
			categories.DELETE("/:id", categoryHandler.Delete)
		}

		products := api.Group("/products")
		{
			products.GET("", productHandler.GetAll)
			products.POST("", productHandler.Create)
			products.PUT("/:id", productHandler.Update)
			products.DELETE("/:id", productHandler.Delete)
		}

		transactions := api.Group("/transactions")
		{
			transactions.GET("", transactionHandler.GetAll)
			transactions.POST("", transactionHandler.Create)
		}
	}
}