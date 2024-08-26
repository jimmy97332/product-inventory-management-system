package router

import (
	"myapp/controllers"
	"myapp/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/login/:user", controllers.Login)
	// the APIs protect by using JWT
	authorized := r.Group("/protected")
	authorized.Use(middlewares.JWTAuthMiddleware())
	{
		authorized.GET("/", controllers.HomeHandler)
		authorized.POST("/products", controllers.CreateProduct)
		authorized.PUT("/products/:id", controllers.UpdateProduct)
		authorized.DELETE("/products/:id", controllers.DeleteProduct)
		authorized.GET("/products/:id", controllers.GetProductByID)
		authorized.GET("/products", controllers.GetAllProducts)
	}
	return r
}
