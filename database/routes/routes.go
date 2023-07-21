package routes

import (
	"ecommerce/database/middleware"
	"ecommerce/handler"
	"github.com/gin-gonic/gin"
)

func ServerRoutes(r1 *gin.Engine) {

	r1.POST("/send-email", handler.SendOtpByEmail)
	r1.POST("/verify-email", handler.VerifyEmail)

	r1.POST("/send-sms", handler.SendSms)
	r1.POST("/verify-sms", handler.VerifyNumber)

	r1.POST("/user-details", handler.UserDetails)

	r1.POST("/login", handler.Login)

	userRouter := r1.Group("/user")
	userRouter.Use(middleware.AccessMiddleware())
	{
		userRouter.DELETE("/account", handler.DeleteAccount)
		userRouter.GET("/products", handler.Product)
		userRouter.GET("/products-type/:id", handler.ProductByType)
		userRouter.GET("/products-id/:id", handler.ProductById)
		userRouter.POST("/add-to-cart/:id", handler.AddToCart)
		userRouter.DELETE("/delete-from-cart/:id", handler.DeleteFromCart)
		userRouter.GET("/show-cart", handler.ShowCart)
		userRouter.POST("/payment", handler.Payment)

		adminRouter := userRouter.Group("/admin")
		adminRouter.Use(middleware.AdminMiddleware())
		{
			adminRouter.POST("/type", handler.ItemsType)
			adminRouter.POST("/items", handler.AddItems)
			adminRouter.DELETE("/item/:id", handler.DeleteItem)
			adminRouter.GET("/users", handler.Users)
			adminRouter.POST("/item-image/:id", handler.Upload)
		}
	}
}
