package model

import (
	"gorm.io/gorm"
)
 
type Chef struct {
	gorm.Model
	Name     string   `json:"name"     gorm:"not null"`
	Email    string   `json:"email"    gorm:"unique;not null"`
	Password string   `json:"-"`                            
	Recipes  []Recipe `json:"recipes,omitempty"             gorm:"foreignKey:ChefID"`
}
 
type RegisterInput struct {
	Name     string `json:"name"     binding:"required,min=2"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
 
type LoginInput struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
 
type AuthResponse struct {
	Token string `json:"token"`
	Chef  Chef   `json:"chef"`
}