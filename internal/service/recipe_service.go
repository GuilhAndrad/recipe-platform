package service

import (
	"errors"
	"math"
	"mime/multipart"

	"github.com/GuilhAndrad/recipe-platform/internal/model"
	"github.com/GuilhAndrad/recipe-platform/internal/repository"
)
 
type RecipeService struct {
	repo         *repository.RecipeRepository
	imageService *ImageService
}
 
func NewRecipeService(repo *repository.RecipeRepository, imageService *ImageService) *RecipeService {
	return &RecipeService{repo: repo, imageService: imageService}
}
 
// Create valida, processa a imagem e persiste a receita.
func (s *RecipeService) Create(
	input model.CreateRecipeInput,
	imageFile *multipart.FileHeader,
	chefID uint,
) (*model.Recipe, error) {
	var imageURL string
 
	if imageFile != nil {
		path, err := s.imageService.Process(imageFile)
		if err != nil {
			return nil, err
		}
		imageURL = path
	}
 
	recipe := &model.Recipe{
		Title:       input.Title,
		Description: input.Description,
		Labels:      input.Labels,
		ImageURL:    imageURL,
		ChefID:      chefID,
	}
 
	if err := s.repo.Create(recipe); err != nil {
		// Se falhou ao salvar, remove a imagem que já foi processada
		s.imageService.Delete(imageURL)
		return nil, errors.New("erro ao salvar receita")
	}
 
	return recipe, nil
}
 
// List normaliza os filtros e retorna a página com metadados de paginação.
func (s *RecipeService) List(filter model.RecipeFilter) (*model.PaginatedRecipes, error) {
	// Valores padrão de paginação
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 50 {
		filter.Limit = 10
	}
 
	recipes, total, err := s.repo.List(filter)
	if err != nil {
		return nil, errors.New("erro ao buscar receitas")
	}
 
	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))
 
	return &model.PaginatedRecipes{
		Data:       recipes,
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
 
//busca uma receita pública pelo ID.
func (s *RecipeService) GetByID(id uint) (*model.Recipe, error) {
	recipe, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("receita não encontrada")
	}
	return recipe, nil
}
 
//garante que apenas o chef dono da receita pode removê-la.
func (s *RecipeService) Delete(recipeID uint, chefID uint) error {
	// Busca primeiro para pegar o imageURL antes de deletar
	recipe, err := s.repo.FindByID(recipeID)
	if err != nil {
		return errors.New("receita não encontrada")
	}
 
	if recipe.ChefID != chefID {
		return errors.New("você não tem permissão para deletar esta receita")
	}
 
	if err := s.repo.Delete(recipeID, chefID); err != nil {
		return errors.New("erro ao deletar receita")
	}
 
	// Remove a imagem do disco após deletar do banco
	s.imageService.Delete(recipe.ImageURL)
 
	return nil
}