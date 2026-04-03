package handler

import (
	"net/http"
	"strconv"

	"github.com/GuilhAndrad/recipe-platform/internal/model"
	"github.com/GuilhAndrad/recipe-platform/internal/service"
	"github.com/gin-gonic/gin"
)
 
type RecipeHandler struct {
	service *service.RecipeService
}
 
func NewRecipeHandler(service *service.RecipeService) *RecipeHandler {
	return &RecipeHandler{service: service}
}
 
// ListRecipes godoc
// GET /recipes?page=1&limit=10&keyword=massa&chef_id=2&label=vegan&published_after=2024-01-01
// Rota pública — sem autenticação.
func (h *RecipeHandler) ListRecipes(c *gin.Context) {
	var filter model.RecipeFilter
 
	// ShouldBindQuery preenche a struct com os query params automaticamente
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
 
	result, err := h.service.List(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
 
	c.JSON(http.StatusOK, result)
}
 
// GetRecipe godoc
// GET /recipes/:id — pública
func (h *RecipeHandler) GetRecipe(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
 
	recipe, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
 
	c.JSON(http.StatusOK, recipe)
}
 
// CreateRecipe godoc
// POST /recipes  (multipart/form-data) — protegida por JWT
// Campos: title, description, labels, image (arquivo)
func (h *RecipeHandler) CreateRecipe(c *gin.Context) {
	var input model.CreateRecipeInput
 
	// ShouldBind funciona com form-data (necessário para upload de arquivo)
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
 
	// Imagem é opcional — sem imagem a receita ainda é criada
	imageFile, _ := c.FormFile("image")
 
	chefID, _ := c.Get("chef_id")
 
	recipe, err := h.service.Create(input, imageFile, chefID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
 
	c.JSON(http.StatusCreated, recipe)
}
 
// DeleteRecipe godoc
// DELETE /recipes/:id — protegida por JWT
func (h *RecipeHandler) DeleteRecipe(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}
 
	chefID, _ := c.Get("chef_id")
 
	if err := h.service.Delete(uint(id), chefID.(uint)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
 
	c.JSON(http.StatusOK, gin.H{"message": "receita deletada com sucesso"})
}