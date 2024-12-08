package handlers

import (
	"ecommerce/models"
	"net/http"
	"ecommerce/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateOrder(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var request struct {
            UserID uint `json:"user_id"`
            CartID uint `json:"cart_id"`
        }

        if err := c.ShouldBindJSON(&request); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
            return
        }

        var order models.Order
        if err := db.Preload("Cart").Preload("Cart.Item").Where("cart_id = ?", request.CartID).First(&order).Error; err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Order creation failed"})
            return
        }

        // Additional logic here for saving the order, checking availability, etc.

        c.JSON(http.StatusOK, gin.H{
            "message": "Order placed successfully",
            "order": order,
        })
    }
}


func ListOrders(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the token from the Authorization header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			return
		}

		// Validate the JWT token and retrieve the username
		username, err := utils.VerifyJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Find the user by username
		var user models.User
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		// Fetch the user's orders, including the associated carts and items
		var orders []models.Order
		if err := db.Preload("Cart").Preload("Cart.Item").Where("user_id = ?", user.ID).Find(&orders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch orders"})
			return
		}

		// Return the orders for the authenticated user
		c.JSON(http.StatusOK, orders)
	}
}
