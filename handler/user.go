package handler

import (
	"ecommerce/database"
	"ecommerce/database/authentication"
	"ecommerce/database/dbHelper"
	"ecommerce/database/middleware"
	"ecommerce/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	var user models.Users
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	uId, err := authentication.Login(database.Todo, user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	signedToken, err := middleware.GenerateToken(uId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": signedToken,
	})
	return
}

func Register(c *gin.Context) {
	var info models.RegisterUser
	if err := c.BindJSON(&info); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	tx, err := database.Todo.Beginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	userId, err := authentication.Create(tx, info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		_ = tx.Rollback()
		return
	}
	//role := "customer"
	err = authentication.AddRole(tx, userId, "customer")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		_ = tx.Rollback()
		return
	}
	err = dbHelper.AssignCart(tx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		_ = tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Created Successful",
	})
	return
}
func DeleteAccount(c *gin.Context) {
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
			"message": "adminId is not of type int",
		})
		return
	}
	err := authentication.Delete(database.Todo, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "deleted Successful",
	})
}
