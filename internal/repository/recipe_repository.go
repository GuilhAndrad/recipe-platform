package repository

import (
	"github.com/GuilhAndrad/recipe-platform/internal/model"
	"gorm.io/gorm"
)
 
type RecipeRepository struct {
	db *gorm.DB
}
 
func NewRecipeRepository(db *gorm.DB) *RecipeRepository {
	return &RecipeRepository{db: db}
}
 
func (r *RecipeRepository) Create(recipe *model.Recipe) error {
	return r.db.Create(recipe).Error
}
 
func (r *RecipeRepository) FindByID(id uint) (*model.Recipe, error) {
	var recipe model.Recipe
	err := r.db.Preload("Chef").First(&recipe, id).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}
 
// List aplica os filtros opcionais e retorna a página solicitada.
func (r *RecipeRepository) List(filter model.RecipeFilter) ([]model.Recipe, int64, error) {
	var recipes []model.Recipe
	var total int64
 
	query := r.db.Model(&model.Recipe{}).Preload("Chef")
 
	// Filtro por palavra-chave no título ou descrição
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", like, like)
	}
 
	// Filtro por chef
	if filter.ChefID > 0 {
		query = query.Where("chef_id = ?", filter.ChefID)
	}
 
	// Filtro por label (busca dentro da string separada por vírgulas)
	if filter.Label != "" {
		query = query.Where("labels ILIKE ?", "%"+filter.Label+"%")
	}
 
	// Filtro por data de publicação (ex: "2024-01-01")
	if filter.PublishedAfter != "" {
		query = query.Where("published_at >= ?", filter.PublishedAfter)
	}
 
	// Conta o total antes de aplicar paginação (para o response)
	query.Count(&total)
 
	// Paginação
	offset := (filter.Page - 1) * filter.Limit
	err := query.
		Order("published_at DESC").
		Limit(filter.Limit).
		Offset(offset).
		Find(&recipes).Error
 
	return recipes, total, err
}
 
func (r *RecipeRepository) Delete(id uint, chefID uint) error {
	result := r.db.Where("id = ? AND chef_id = ?", id, chefID).Delete(&model.Recipe{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}