package router

import (
	"myapp/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", controllers.HomeHandler)
	r.POST("/products", controllers.CreateProduct)
	r.PUT("/products/:id", controllers.UpdateProduct)
	r.DELETE("/products/:id", controllers.DeleteProduct)
	r.GET("/products/:id", controllers.GetProductByID)
	r.GET("/products", controllers.GetAllProducts)

	return r
}
