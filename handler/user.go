package handler

import (
	"ecommerce/database"
	"ecommerce/database/authentication"
	"ecommerce/database/dbHelper"
	"ecommerce/database/middleware"
	"ecommerce/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/twilio/twilio-go"
	"github.com/twilio/twilio-go/rest/api/v2010"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const otpLength = 6

func generateOTP() string {
	// Initialize the random number generator with a seed based on the current time
	rand.Seed(time.Now().UnixNano())

	// Define the characters to be used in the OTP
	// You can customize this string according to your requirements
	characters := "0123456789"

	otp := make([]byte, otpLength)
	for i := 0; i < otpLength; i++ {
		otp[i] = characters[rand.Intn(len(characters))]
	}

	return string(otp)
}

func SendOtpByEmail(c *gin.Context) {
	var email models.EmailInfo
	if err := c.BindJSON(&email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	otp := generateOTP()

	Sender := "jaingautam601@gmail.com"
	Recipient := email.Email

	Subject := "Amazon SES Test (AWS SDK for Go)"
	htmlBody := "<h1>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
		"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
		"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.Your otp for verification is:- " + otp + "</p>"

	TextBody := "This email was sent with Amazon SES using the AWS SDK for Go."
	CharSet := "UTF-8"
	var awsConfig models.AWSConfig
	awsConfig.Region = os.Getenv("region")
	awsConfig.AccessKeySecret = os.Getenv("secretKey")
	awsConfig.AccessKeyID = os.Getenv("accessKey")
	awsConfig.BucketName = os.Getenv("bucketName")
	sess := CreateSession(awsConfig)

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(htmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(Sender),
	}

	result, err := svc.SendEmail(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = authentication.StoreEmailOtp(database.Todo, otp, email.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    result,
		"message": "OTP Sent",
	})
}

func VerifyEmail(c *gin.Context) {
	var otp models.EmailVerify
	if err := c.BindJSON(&otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	email, err := authentication.VerifyEmailOtp(database.Todo, otp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = dbHelper.EnterEmail(database.Todo, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Verified!! ",
	})
}

func SendSms(c *gin.Context) {
	var number models.NumberInfo
	if err := c.BindJSON(&number); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error to get input",
			"error":   err.Error(),
		})
		return
	}

	otp := generateOTP()

	client := twilio.NewRestClient()

	params := &openapi.CreateMessageParams{}
	params.SetTo("+91" + number.Number)
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetBody(otp)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Email sent!!",
		})
	}

	err = authentication.StoreNumberOtp(database.Todo, otp, number.Number)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "OTP Sent",
	})
}
func VerifyNumber(c *gin.Context) {
	var otp models.NumberVerify
	if err := c.BindJSON(&otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	number, err := authentication.VerifyNumberOtp(database.Todo, otp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = dbHelper.EnterNumber(database.Todo, number)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Verified!! ",
	})
}

func UserDetails(c *gin.Context) {
	var info models.RegisterUser
	if err := c.BindJSON(&info); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Validator are applied here for name and password only
	validate := validator.New()
	err := validate.Struct(info)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user not validated",
		})
		return
	}
	// checking if user is verified or not
	flag, err := authentication.CheckVerified(database.Todo, info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	var userId int
	//Transaction Here
	tx, err := database.Todo.Beginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Checking if user is verified by email or phone number
	if info.PhoneNumber == "" && flag == true { // if user is registered then number or email must be provided by front-end . so anyhow , every request will fall under one of these.
		userId, err = authentication.CreateUserByEmail(tx, info)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			tx.Rollback()
			return
		}
	} else if flag == true {
		userId, err = authentication.CreateUserByNumber(tx, info)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			tx.Rollback()
			return
		}
	}
	// Normally adding their role and assigning cart
	err = authentication.AddRole(tx, userId, "customer")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		tx.Rollback()
		return
	}
	err = dbHelper.AssignCart(tx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Created Successful",
	})
	return
}

// Login function here
func Login(c *gin.Context) {
	var user models.Users
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please register yourself",
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
			"message": "userId is not of type int",
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
