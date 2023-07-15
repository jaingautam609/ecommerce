package routes

import (
	"ecommerce/database/middleware"
	"ecommerce/handler"
	"github.com/gin-gonic/gin"
)

func ServerRoutes(r1 *gin.Engine) {

	r1.POST("/login", handler.Login)
	r1.POST("/register", handler.Register)
	userRouter := r1.Group("/admin")
	userRouter.Use(middleware.AdminMiddleware())
	{
		userRouter.POST("/type", handler.ItemsType)
		userRouter.POST("/items", handler.AddItems)
		userRouter.DELETE("/item/:id", handler.DeleteItem)
		userRouter.GET("users", handler.Users)
		userRouter.POST("/item-image/:id", handler.Upload)
	}
	userRouter = r1.Group("/customer")
	userRouter.Use(middleware.AccessMiddleware())
	{
		userRouter.DELETE("/account", handler.DeleteAccount)
		userRouter.GET("/products", handler.Product)
		userRouter.GET("/products-type/:id", handler.ProductByType)
		userRouter.GET("/products-id/:id", handler.ProductById)
		userRouter.POST("/add-to-cart/:id", handler.AddToCart)
		userRouter.DELETE("/delete-from-cart/:id", handler.DeleteFromCart)
		userRouter.GET("/show-cart", handler.ShowCart)
		userRouter.DELETE("/payment", handler.Payment)
	}
}
