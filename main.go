package main

import (
	"fmt"
	"time"
	// "errors"
	// "fmt"
	"net/http"
	"rest-api-go/database"
	"rest-api-go/models"

	// "gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

func main(){
	port := ":8080"

	database.StartDB()

	router := gin.Default()

	router.POST("/orders", createOrder)

	fmt.Println("server is running on port", port)

	router.Run(port)
}

func createOrder(ctx *gin.Context){
	db := database.GetDB()

	var body models.Body
	var newOrder models.Order
	var items []models.Item

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	newOrder.CustomerName = body.CustomerName
	newOrder.OrderedAt = time.Now()

	var check bool

	for _, item := range body.Items {
		check = checkCodeExist(item.ItemCode)
	}

	if check {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("item code already exist"))
		return
	}

	err := db.Create(&newOrder).Error

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	for _, item := range body.Items {
		var newItem models.Item

		newItem.ItemCode = item.ItemCode
		newItem.Description = item.Description
		newItem.Quantity = item.Quantity
		newItem.OrderID = newOrder.OrderID

		items = append(items, newItem)
	}

	err = db.Create(&items).Error

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed storing item"))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "order has been recorded"})
}

func checkCodeExist(id string) bool{
	db := database.GetDB()

	var item models.Item

	r := db.Where("item_code = ?", id).First(&item)

	if r.RowsAffected > 0{
		return true
	} 

	return false
}