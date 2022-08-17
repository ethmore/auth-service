package controllers

import (
	"fmt"
	"log"
	"time"

	"e-comm/authService/bcrypt"
	"e-comm/authService/dotEnv"

	"e-comm/authService/middleware"

	"e-comm/authService/postgresql"

	"github.com/golang-jwt/jwt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type SellerRegisterBody struct {
	CompanyName   string
	Email         string
	Password      string
	PasswordAgain string
	Address       string
	PhoneNumber   string
}

type Product struct {
	Token       string
	Title       string
	Description string
	Price       string
	Stock       string
	Photo       string
}

func SellerRegisterPostHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestBody SellerRegisterBody

		if err := ctx.ShouldBindBodyWith(&requestBody, binding.JSON); err != nil {
			log.Printf("%+v", err)
		}

		companyName := requestBody.CompanyName
		email := requestBody.Email
		password := requestBody.Password
		passwordAgain := requestBody.PasswordAgain
		address := requestBody.Address
		phonenumber := requestBody.PhoneNumber

		salt := dotEnv.GoDotEnvVariable("SALT")

		if password == passwordAgain {
			_, checkedMail, _ := postgresql.GetSeller(email)
			if checkedMail == email {
				fmt.Println("email already registered")
				ctx.JSON(400, gin.H{"message": "email already registered"})
			} else {
				saltedPassword := password + salt
				hash, _ := bcrypt.HashPassword(saltedPassword)

				res := postgresql.Insert(companyName, email, hash, address, phonenumber)
				if res == 200 {
					ctx.JSON(200, gin.H{"message": "OK"})
				}
			}

		} else {
			fmt.Println("passwords does not match")
			ctx.JSON(400, gin.H{"message": "passwords does not match"})
		}
	}
}

func SellerLoginPostHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestBody LoginBody
		if err := ctx.ShouldBindBodyWith(&requestBody, binding.JSON); err != nil {
			log.Printf("%+v", err)
		}

		email := requestBody.Email
		password := requestBody.Password
		salt := dotEnv.GoDotEnvVariable("SALT")

		_, checkedMail, checkedPassword := postgresql.GetSeller(email)

		if checkedMail == email {
			saltedPassword := password + salt
			match := bcrypt.CheckPasswordHash(saltedPassword, checkedPassword)

			if match {
				secretToken := dotEnv.GoDotEnvVariable("TOKEN")
				hmacSampleSecret := []byte(secretToken)

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"mail": requestBody.Email,
					"nbf":  time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
				})
				tokenString, err := token.SignedString(hmacSampleSecret)
				if err != nil {
					fmt.Println(err)
				}

				fmt.Println("OK")
				ctx.JSON(200, gin.H{"message": "OK", "token": tokenString})
			} else {
				fmt.Println("wrong password")
				ctx.JSON(400, gin.H{"message": "wrong password"})
			}
		} else {
			fmt.Println("email not registered")
			ctx.JSON(400, gin.H{"message": "email not registered"})
		}
	}
}

func SellerDashboard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var mailAuth = middleware.UserAuth(ctx)
		if mailAuth != "" {
			ctx.JSON(200, gin.H{"message": "OK", "mail": mailAuth})
		}
	}
}

func AddProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var mailAuth = middleware.UserAuth(ctx)
		if mailAuth != "" {
			var requestBody Product
			if err := ctx.ShouldBindBodyWith(&requestBody, binding.JSON); err != nil {
				log.Printf("%+v", err)
			}

			fmt.Println(requestBody.Title, requestBody.Description, requestBody.Price, requestBody.Stock)
			postgresql.InsertProduct(mailAuth, requestBody.Title, requestBody.Price, requestBody.Description, requestBody.Stock, requestBody.Stock)
			ctx.JSON(200, gin.H{"message": "OK", "mail": mailAuth})
		}

	}
}