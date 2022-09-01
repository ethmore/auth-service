package controllers

import (
	"fmt"
	"net/http"
	"time"

	"e-comm/authService/bcrypt"
	"e-comm/authService/dotEnv"

	"e-comm/authService/middleware"

	"e-comm/authService/repositories/mongodb"

	"github.com/golang-jwt/jwt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserRegisterBody struct {
	Name          string
	Surname       string
	Email         string
	Password      string
	PasswordAgain string
}

type LoginBody struct {
	Email    string
	Password string
	Type     string
}

func UserRegisterPostHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userBody UserRegisterBody
		if bodyErr := ctx.ShouldBindBodyWith(&userBody, binding.JSON); bodyErr != nil {
			fmt.Println("body: ", bodyErr)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		if userBody.Password != userBody.PasswordAgain {
			fmt.Println("passwords does not match")
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "passwords does not match"})
			return
		}

		user, mongoErr := mongodb.FindOneUser(userBody.Email)
		if mongoErr != nil {
			fmt.Println("mongodb (findOne): ", mongoErr)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		if user != nil {
			fmt.Println("email already registered")
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "email already registered"})
			return
		}

		salt := dotEnv.GoDotEnvVariable("SALT")
		saltedPassword := userBody.Password + salt
		hash, _ := bcrypt.HashPassword(saltedPassword)

		insertErr := mongodb.InsertOneUser(userBody.Name, userBody.Surname, userBody.Email, hash)
		if insertErr != nil {
			fmt.Println("mongodb (insert): ", insertErr)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		fmt.Println("User registered: ", userBody.Email)
		ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
	}
}

func UserLoginPostHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userBody LoginBody
		if bodyErr := ctx.ShouldBindBodyWith(&userBody, binding.JSON); bodyErr != nil {
			fmt.Println("body: ", bodyErr)
			ctx.Status(http.StatusInternalServerError)
			return
		}

		user, mongoErr := mongodb.FindOneUser(userBody.Email)
		if mongoErr != nil {
			fmt.Println("mongodb (findOne): ", mongoErr)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		if user == nil {
			fmt.Println("email not registered")
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "wrong credentials"})
			return
		}

		salt := dotEnv.GoDotEnvVariable("SALT")
		saltedPassword := userBody.Password + salt
		match := bcrypt.CheckPasswordHash(saltedPassword, user.Password)

		if !match {
			fmt.Println("wrong password")
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "wrong credentials"})
		}

		secretToken := dotEnv.GoDotEnvVariable("TOKEN")
		hmacSampleSecret := []byte(secretToken)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"mail": userBody.Email,
			"type": userBody.Type,
			"nbf":  time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		})
		tokenString, tokenErr := token.SignedString(hmacSampleSecret)
		if tokenErr != nil {
			fmt.Println("token: ", tokenErr)
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad token"})
			return
		}

		fmt.Println("User logged in: ", userBody.Email)
		ctx.JSON(http.StatusOK, gin.H{"message": "OK", "token": tokenString})

	}
}

func UserProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth, err := middleware.UserAuth(ctx)
		if err != nil {
			fmt.Println("authentication: ", err)
			ctx.JSON(http.StatusOK, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "OK", "mail": auth.EMail, "type": auth.Type})
	}
}