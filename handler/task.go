package handler

import (
	"ecommerce/database"
	"ecommerce/database/dbHelper"
	"ecommerce/models"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func ItemsType(c *gin.Context) {
	adminIdInterface, flag := c.Get("adminId")
	if !flag {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error to get id",
		})
		return
	}

	adminId, ok := adminIdInterface.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "adminId is not of type int",
		})
		return
	}
	var itemTypes models.ItemType
	if err := c.BindJSON(&itemTypes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err := dbHelper.AddItemType(database.Todo, adminId, itemTypes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Item added successfully",
	})
	return
}
func AddItems(c *gin.Context) {
	adminIdInterface, flag := c.Get("adminId")
	if !flag {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error to get id",
		})
		return
	}

	adminId, ok := adminIdInterface.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "adminId is not of type int",
		})
		return
	}
	var ItemInfo models.Item
	if err := c.BindJSON(&ItemInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := dbHelper.AddItem(database.Todo, ItemInfo, adminId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Item added successfully",
	})
	return
}
func DeleteItem(c *gin.Context) {
	id := c.Param("id")
	itemId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Id is of type string",
		})
		return
	}
	err = dbHelper.DeleteItem(database.Todo, itemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Item deleted successfully",
	})
	return
}
func Users(c *gin.Context) {
	var pageInfo models.PageInfo
	if err := c.BindJSON(&pageInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	userCount, err := dbHelper.CountUsers(database.Todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	var allInfo []models.Register
	offset := pageInfo.PageNo * pageInfo.Limit
	if offset < userCount {
		allInfo, err = dbHelper.AllUsers(database.Todo, pageInfo.Limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No more Users",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"count": userCount,
		"data":  allInfo,
	})
	return
}
func Product(c *gin.Context) {
	var pageInfo models.PageInfo
	if err := c.BindJSON(&pageInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	productCount, err := dbHelper.CountProduct(database.Todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	offset := pageInfo.PageNo * pageInfo.Limit
	var item []models.Item
	if offset < productCount {
		item, err = dbHelper.AllProducts(database.Todo, pageInfo.Limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"ItemCount": productCount,
		"data":      item,
	})
}
func ProductById(c *gin.Context) {
	id := c.Param("id")
	itemId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Id is of type string",
		})
		return
	}
	item, err := dbHelper.ProductById(database.Todo, itemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}
func ProductByType(c *gin.Context) {
	id := c.Param("id")
	typeId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Id is of type string",
		})
		return
	}
	var pageInfo models.PageInfo
	if err := c.BindJSON(&pageInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	productCount, err := dbHelper.CountProductByType(database.Todo, typeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	offset := pageInfo.PageNo * pageInfo.Limit
	var item []models.Item
	if offset < productCount {
		item, err = dbHelper.ProductByType(database.Todo, typeId, pageInfo.Limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"count": productCount,
		"data":  item,
	})
}

func AddToCart(c *gin.Context) {
	id := c.Param("id")
	itemId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	var itemQuantity models.Pieces
	if err := c.BindJSON(&itemQuantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
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
	getItem, err := dbHelper.ProductById(database.Todo, itemId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	itemType, err := dbHelper.GetType(database.Todo, getItem.TypeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	cartId, err := dbHelper.GetCartId(database.Todo, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	quantity, err := dbHelper.GetQuantity(database.Todo, itemId, cartId)
	if quantity == 0 {
		err = dbHelper.AddToCart(database.Todo, getItem, itemQuantity.Quantity, cartId, itemType, itemId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		err = dbHelper.IncreaseInCart(database.Todo, cartId, itemId, itemQuantity.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Item added to cart",
	})
}
func DeleteFromCart(c *gin.Context) {
	id := c.Param("id")
	itemId, err := strconv.Atoi(id) //type conversion
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
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
	cartId, err := dbHelper.GetCartId(database.Todo, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	quantity, err := dbHelper.GetQuantity(database.Todo, itemId, cartId)
	if quantity == 1 {
		err = dbHelper.DeleteFromCart(database.Todo, itemId, cartId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else if quantity > 1 {
		err = dbHelper.DecreaseFromCart(database.Todo, itemId, cartId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "already empty",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Removed from cart",
	})
}
func ShowCart(c *gin.Context) {
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
	cartId, err := dbHelper.GetCartId(database.Todo, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	cartItems, err := dbHelper.ShowCart(database.Todo, cartId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": cartItems,
	})
	return
}
func Payment(c *gin.Context) {
	var checkout models.Checkout
	if err := c.BindJSON(&checkout); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
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
	cartId, err := dbHelper.GetCartId(database.Todo, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = dbHelper.Checkout(database.Todo, cartId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"error": "payment completed",
	})
}
func CreateSession(awsConfig models.AWSConfig) *session.Session {
	sess := session.Must(session.NewSession(
		&aws.Config{
			Region: aws.String(awsConfig.Region),
			Credentials: credentials.NewStaticCredentials(
				awsConfig.AccessKeyID,
				awsConfig.AccessKeySecret,
				"",
			),
		},
	))
	return sess
}
func Upload(c *gin.Context) {
	id := c.Param("id")
	itemId, err := strconv.Atoi(id) //type conversion
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = c.Request.ParseMultipartForm(32 << 20) // 32 MB is the maximum file size
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer file.Close()

	bucket := "gautam001"
	fileName := handler.Filename + time.Now().String()
	var awsConfig models.AWSConfig
	awsConfig.Region = os.Getenv("region")
	awsConfig.AccessKeySecret = os.Getenv("secretKey")
	awsConfig.AccessKeyID = os.Getenv("accessKey")
	awsConfig.BucketName = os.Getenv("bucketName")
	sess := CreateSession(awsConfig)
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Printf("Successfully uploaded %q to %q\n", fileName, bucket)
	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
	urlStr, err := req.Presign(240 * time.Hour)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var upload models.Uploads
	upload.Url = urlStr
	upload.Name = handler.Filename
	upload.Path = "./uploads/" + upload.Name
	f, err := os.OpenFile(upload.Path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer f.Close()

	// Copy the contents of the file to the new file
	_, err = io.Copy(f, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
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
	uploadId, err := dbHelper.Upload(tx, upload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		_ = tx.Rollback()
		return
	}
	err = dbHelper.ItemImage(tx, itemId, uploadId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		_ = tx.Rollback()
		return
	}
	err = tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Uploaded successfully",
	})
}
