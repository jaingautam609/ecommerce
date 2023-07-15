package middleware

import (
	"fmt"
	"github.com/form3tech-oss/jwt-go"
	"os"
	"time"
)

func GenerateToken(userId int) (string, error) {
	secretKey := os.Getenv("key")
	//secretKey := "my_secret_key" // ***********do not hardcode the secret key here ,can set in the environment
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["userid"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}
func ValidateToken(tokenString string) (int, error) {
	//secretKey := "my_secret_key"
	secretKey := os.Getenv("key")
	var userId int
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return userId, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId = int(claims["userid"].(float64))

	} else {
		fmt.Println("Invalid token")
	}
	return userId, nil
}
