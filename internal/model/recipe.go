package model

import (
	"time"
 
	"gorm.io/gorm"
)
 
type Recipe struct {
	gorm.Model
	Title       string    `json:"title"       gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	ImageURL    string    `json:"image_url"`
	Labels      string    `json:"labels"`
	PublishedAt time.Time `json:"published_at" gorm:"autoCreateTime"`
	ChefID      uint      `json:"chef_id"`
	Chef        Chef      `json:"chef,omitempty" gorm:"foreignKey:ChefID"`
}
 
type CreateRecipeInput struct {
	Title       string `form:"title"       binding:"required,min=3"`
	Description string `form:"description" binding:"required,min=10"`
	Labels      string `form:"labels"`
}
 
type RecipeFilter struct {
	Keyword        string `form:"keyword"`
	ChefID         uint   `form:"chef_id"`
	Label          string `form:"label"`
	PublishedAfter string `form:"published_after"`
	Page           int    `form:"page"`
	Limit          int    `form:"limit"`
}
 
type PaginatedRecipes struct {
	Data       []Recipe `json:"data"`
	Page       int      `json:"page"`
	Limit      int      `json:"limit"`
	Total      int64    `json:"total"`
	TotalPages int      `json:"total_pages"`
}