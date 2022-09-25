package middleware

import (
	"auth-and-db-service/dotEnv"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt"

	"fmt"
	"log"
)

type TokenBody struct {
	Token string
}

type Authentication struct {
	EMail string
	Type  string
}

func UserAuth(c *gin.Context) (*Authentication, error) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("Recovered in f", r)
	// 		fmt.Println("User not logged in")
	// 		c.JSON(http.StatusOK, gin.H{"message": "loginNeeded"})
	// 	}
	// }()

	var tokenBody TokenBody
	if err := c.ShouldBindBodyWith(&tokenBody, binding.JSON); err != nil {
		log.Printf("%+v", err)
	}

	if tokenBody.Token == "" {
		return nil, errors.New("empty token")
	}

	clientToken := tokenBody.Token
	secretToken := dotEnv.GoDotEnvVariable("TOKEN")
	hmacSampleSecret := []byte(secretToken)

	token, err := jwt.Parse(clientToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var auth = Authentication{
			EMail: claims["mail"].(string),
			Type:  claims["type"].(string),
		}
		return &auth, nil

	} else {
		return nil, err
	}
}
