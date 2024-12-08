package handlers

import (
	"ecommerce/models"
	"ecommerce/utils"
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddToCart(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			return
		}
		var request struct {
			ItemID   uint `json:"item_id"`
			Quantity int  `json:"quantity"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		if request.ItemID == 0 || request.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID or quantity"})
			return
		}

		var item models.Item
		if err := db.First(&item, request.ItemID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Item not found"})
			return
		}

		username, err := utils.VerifyJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		var user models.User
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		cartItem := models.Cart{
			UserID:   user.ID,
			ItemID:   request.ItemID,
			Quantity: request.Quantity,
		}

		if err := db.Create(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Item added to cart successfully"})
	}
}

func ListCarts(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
            return
        }

        username, err := utils.VerifyJWT(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        var user models.User
        if err := db.Where("username = ?", username).First(&user).Error; err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
            return
        }

        var carts []models.Cart
        if err := db.Preload("Item").Where("user_id = ?", user.ID).Find(&carts).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch carts"})
            return
        }

        c.JSON(http.StatusOK, carts)
    }
}


