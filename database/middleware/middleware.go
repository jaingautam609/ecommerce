package middleware

import (
	"ecommerce/database"
	"ecommerce/database/authentication"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		adminId, err := ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		var flag bool
		flag, err = authentication.ValidateAdmin(database.Todo, adminId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"restaurants": err.Error(),
			})
			return
		}
		if flag == false {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "No access",
			})
			return
		}
		c.Set("adminId", adminId)
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
