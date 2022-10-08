package main

import (
	"fmt"
	"rest-api-go/database"
	"rest-api-go/controllers"

	// "gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

func main() {
	port := ":8080"

	database.StartDB()

	router := gin.Default()

	router.POST("/orders", controllers.CreateOrder)
	router.GET("/orders", controllers.GetOrder)
	router.PUT("/orders/:orderId", controllers.UpdateOrder)
	router.DELETE("/orders/:orderId", controllers.DeleteOrder)

	fmt.Println("server is running on port", port)

	router.Run(port)
}

