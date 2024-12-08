package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Token    string
}

type Item struct {
	gorm.Model
	Name  string
	Price float64
}

type Cart struct {
	gorm.Model
	UserID   uint `json:"user_id"`
	ItemID   uint `json:"item_id"`
	Quantity int  `json:"quantity"`
	Item     Item `gorm:"foreignkey:ItemID"`
}

type Order struct {
	gorm.Model
    UserID uint
    CartID uint 
    Cart   Cart  `gorm:"foreignkey:CartID"`
}
