package handler

import (
	"net/http"

	"github.com/GuilhAndrad/recipe-platform/internal/model"
	"github.com/GuilhAndrad/recipe-platform/internal/service"
	"github.com/gin-gonic/gin"
)
 
type ChefHandler struct {
	service *service.ChefService
}
 
func NewChefHandler(service *service.ChefService) *ChefHandler {
	return &ChefHandler{service: service}
}
 
// Register godoc
// POST /auth/register
func (h *ChefHandler) Register(c *gin.Context) {
	var input model.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
 
	response, err := h.service.Register(input)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
 
	c.JSON(http.StatusCreated, response)
}
 
// Login godoc
// POST /auth/login
func (h *ChefHandler) Login(c *gin.Context) {
	var input model.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
 
	response, err := h.service.Login(input)
	if err != nil {
		// 401 para credenciais inválidas
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
 
	c.JSON(http.StatusOK, response)
}
 
// GetProfile godoc
// GET /me  (rota protegida)
func (h *ChefHandler) GetProfile(c *gin.Context) {
	// O middleware já validou o token e injetou o chef_id
	chefID, _ := c.Get("chef_id")
 
	chef, err := h.service.GetProfile(chefID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "chef não encontrado"})
		return
	}
 
	c.JSON(http.StatusOK, chef)
}