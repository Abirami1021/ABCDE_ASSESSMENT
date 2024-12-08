package main

import (
	"ecommerce/handlers"
	"ecommerce/middleware"
	"ecommerce/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	db, _ := gorm.Open(sqlite.Open("ecommerce.db"), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.Item{}, &models.Cart{}, &models.Order{})

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the E-Commerce API"})
	})
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.JSON(204, nil)
	})

	r.POST("/users", handlers.CreateUser(db)) 
	r.GET("/users", handlers.GetUsers(db)) 
	r.POST("/users/login", handlers.LoginUser(db))

	r.POST("/items", handlers.CreateItem(db))
	r.GET("/items", handlers.ListItems(db)) 

	auth := r.Group("/", middleware.Authenticate())
	{
		auth.POST("/carts", handlers.AddToCart(db))
		auth.GET("/carts", handlers.ListCarts(db))
		auth.POST("/orders", handlers.CreateOrder(db))
		auth.GET("/orders", handlers.ListOrders(db))
	}

	r.Run(":8080")
}
