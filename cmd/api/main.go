package main

import (
	"log"

	"github.com/GuilhAndrad/recipe-platform/internal/config"
	"github.com/GuilhAndrad/recipe-platform/internal/database"
	"github.com/GuilhAndrad/recipe-platform/internal/handler"
	"github.com/GuilhAndrad/recipe-platform/internal/middleware"
	"github.com/GuilhAndrad/recipe-platform/internal/model"
	"github.com/GuilhAndrad/recipe-platform/internal/repository"
	"github.com/GuilhAndrad/recipe-platform/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Carrega configurações do .env
	cfg := config.Load()

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET não pode ser vazio")
	}

	// 2. Conecta ao banco e executa migrations
	db := database.Connect(cfg.DatabaseURL,
		&model.Chef{},
		&model.Recipe{},
	)

	// 3. Instancia as camadas (repository → service → handler)
	chefRepo := repository.NewChefRepository(db)
	chefSvc := service.NewChefService(chefRepo, cfg.JWTSecret, cfg.JWTExpirationHours)
	chefHandler := handler.NewChefHandler(chefSvc)

	recipeRepo := repository.NewRecipeRepository(db)
	imageSvc := service.NewImageService()
	recipeSvc := service.NewRecipeService(recipeRepo, imageSvc)
	recipeHandler := handler.NewRecipeHandler(recipeSvc)

	// 4. Configura o Gin e as rotas
	r := gin.Default()

	// Serve imagens salvas em disco como arquivos estáticos
	// Ex: GET /uploads/1714000000000.jpg
	r.Static("/uploads", "./uploads")

	// — Rotas públicas —
	r.GET("/recipes", recipeHandler.ListRecipes)
	r.GET("/recipes/:id", recipeHandler.GetRecipe)

	auth := r.Group("/auth")
	{
		auth.POST("/register", chefHandler.Register)
		auth.POST("/login", chefHandler.Login)
	}

	// — Rotas protegidas por JWT —
	protected := r.Group("/")
	protected.Use(middleware.AuthRequired(cfg.JWTSecret))
	{
		protected.GET("/me", chefHandler.GetProfile)
		protected.POST("/recipes", recipeHandler.CreateRecipe)
		protected.DELETE("/recipes/:id", recipeHandler.DeleteRecipe)
	}

	// 5. Sobe o servidor
	log.Printf("Servidor rodando em http://localhost:%s\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}