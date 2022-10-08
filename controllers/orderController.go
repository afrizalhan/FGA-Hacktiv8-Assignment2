package controllers

import (
	"rest-api-go/database"
	"rest-api-go/models"
	"strconv"
	"net/http"
	"time"
	"fmt"

	// "gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

func CreateOrder(ctx *gin.Context) {
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
		if check {
			break
		}
	}

	if check {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "item code already exist"})
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("item code already exist"))
		return
	}

	err := db.Create(&newOrder).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed storing item"))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "order has been recorded"})
}

func GetOrder(ctx *gin.Context) {
	db := database.GetDB()

	var results []models.Order

	db.Find(&results)

	if len(results) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No record found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": results})
}

func UpdateOrder(ctx *gin.Context) {

	orderId, _ := strconv.Atoi(ctx.Param("orderId"))

	db := database.GetDB()

	var body models.Body
	var order models.Order
	var newOrder models.Order

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	newOrder.CustomerName = body.CustomerName
	newOrder.OrderedAt = time.Now()

	err := db.Model(&order).Where("order_id = ?", orderId).Updates(newOrder).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if len(body.Items) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "order has been updated"})
		return
	}

	var flag bool

	for _, item := range body.Items {
		var itemMod models.Item
		var newItem models.Item

		newItem.Description = item.Description
		newItem.Quantity = item.Quantity

		err := db.Model(&itemMod).Where("item_id = ? AND order_id = ?", item.ItemID, orderId).Updates(newItem).Error
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			flag = true
		}
	}

	if flag {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed updating an item"})
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed updating an item"))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "order has been updated"})
}

func DeleteOrder(ctx *gin.Context) {
	orderId, _ := strconv.Atoi(ctx.Param("orderId"))

	db := database.GetDB()

	var order models.Order
	var item models.Item

	err := db.Where("order_id = ?", orderId).Delete(item).Error

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed deleting item"))
		return
	}

	err = db.Where("order_id = ?", orderId).Delete(order).Error

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed deleting order"))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "order and items related to it has beed deleted"})
}

func checkCodeExist(id string) bool {
	db := database.GetDB()

	var item models.Item

	r := db.Where("item_code = ?", id).First(&item)

	return r.RowsAffected > 0
}
