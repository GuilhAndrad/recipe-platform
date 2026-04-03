package repository

import (
	"github.com/GuilhAndrad/recipe-platform/internal/model"
	"gorm.io/gorm"
)
 
type ChefRepository struct {
	db *gorm.DB
}
 
func NewChefRepository(db *gorm.DB) *ChefRepository {
	return &ChefRepository{db: db}
}
 
func (r *ChefRepository) Create(chef *model.Chef) error {
	return r.db.Create(chef).Error
}
 
func (r *ChefRepository) FindByEmail(email string) (*model.Chef, error) {
	var chef model.Chef
	err := r.db.Where("email = ?", email).First(&chef).Error
	if err != nil {
		return nil, err
	}
	return &chef, nil
}
 
func (r *ChefRepository) FindByID(id uint) (*model.Chef, error) {
	var chef model.Chef
	err := r.db.Preload("Recipes").First(&chef, id).Error
	if err != nil {
		return nil, err
	}
	return &chef, nil
}
