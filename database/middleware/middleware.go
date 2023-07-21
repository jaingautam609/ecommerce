package middleware

import (
	"ecommerce/database"
	"ecommerce/database/authentication"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdInterface, flag := c.Get("userId")
		if !flag {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "error to get id",
			})
			return
		}

		userId, ok := userIdInterface.(int)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "userId is not of type int",
			})
			return
		}
		var isTrue bool
		isTrue, err := authentication.ValidateAdmin(database.Todo, userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"restaurants": err.Error(),
			})
			return
		}
		if isTrue == false {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "No access",
			})
			return
		}
		c.Set("adminId", userId)
		c.Next()
	}
}
func AccessMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		userId, err := ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.Set("userId", userId)
		c.Next()
	}
}
