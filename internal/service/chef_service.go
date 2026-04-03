package service

import (
	"errors"
	"time"
	"github.com/GuilhAndrad/recipe-platform/internal/model"
	"github.com/GuilhAndrad/recipe-platform/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)
 
type ChefService struct {
	repo               *repository.ChefRepository
	jwtSecret          string
	jwtExpirationHours int
}
 
func NewChefService(repo *repository.ChefRepository, jwtSecret string, jwtExpirationHours int) *ChefService {
	return &ChefService{
		repo:               repo,
		jwtSecret:          jwtSecret,
		jwtExpirationHours: jwtExpirationHours,
	}
}
 
//cria um novo chef após validar e fazer hash da senha.
func (s *ChefService) Register(input model.RegisterInput) (*model.AuthResponse, error) {
	// Verifica se o e-mail já está em uso
	existing, _ := s.repo.FindByEmail(input.Email)
	if existing != nil {
		return nil, errors.New("e-mail já cadastrado")
	}
 
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("erro ao processar senha")
	}
 
	chef := &model.Chef{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}
 
	if err := s.repo.Create(chef); err != nil {
		return nil, errors.New("erro ao criar chef")
	}
 
	token, err := s.generateToken(chef)
	if err != nil {
		return nil, err
	}
 
	return &model.AuthResponse{Token: token, Chef: *chef}, nil
}
 
//valida as credenciais e retorna um token JWT.
func (s *ChefService) Login(input model.LoginInput) (*model.AuthResponse, error) {
	chef, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		// Mensagem genérica para não revelar se o e-mail existe
		return nil, errors.New("credenciais inválidas")
	}
 
	if err := bcrypt.CompareHashAndPassword([]byte(chef.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("credenciais inválidas")
	}
 
	token, err := s.generateToken(chef)
	if err != nil {
		return nil, err
	}
 
	return &model.AuthResponse{Token: token, Chef: *chef}, nil
}
 
//busca o chef pelo ID extraído do token.
func (s *ChefService) GetProfile(chefID uint) (*model.Chef, error) {
	return s.repo.FindByID(chefID)
}
 
//cria um JWT assinado com as claims do chef.
func (s *ChefService) generateToken(chef *model.Chef) (string, error) {
	claims := jwt.MapClaims{
		"chef_id": chef.ID,
		"email":   chef.Email,
		"exp":     time.Now().Add(time.Duration(s.jwtExpirationHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
 
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}