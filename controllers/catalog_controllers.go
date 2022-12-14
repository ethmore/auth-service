package controllers

import (
	"auth-and-db-service/repositories/postgresql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ProductResponse struct {
	Token string
	Id    string
}

func GetProduct(productRepo postgresql.IProductRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestBody ProductResponse
		if bodyErr := ctx.ShouldBindBodyWith(&requestBody, binding.JSON); bodyErr != nil {
			fmt.Println("body: ", bodyErr)
			ctx.Status(http.StatusBadRequest)
			return
		}

		product, pqErr := productRepo.GetProduct(requestBody.Id)
		if pqErr != nil {
			fmt.Println("postgresql (get)", pqErr, "req id: ", requestBody.Id)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		fmt.Println(product)
		ctx.JSON(http.StatusOK, gin.H{"product": product})
	}
}

func GetAllProducts(productRepo postgresql.IProductRepo) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		products, err := productRepo.GetAllProducts()
		if err != nil {
			fmt.Println(err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"products": products})
	}
}
