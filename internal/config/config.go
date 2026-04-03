package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
 
	"github.com/joho/godotenv"
)
 
type Config struct {
	Port               string
	DatabaseURL        string
	JWTSecret          string
	JWTExpirationHours int
}
 
// Load lê as variáveis do arquivo .env e retorna a config da aplicação.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}
 
	jwtExpiration, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		log.Fatal("JWT_EXPIRATION_HOURS deve ser um número inteiro")
	}
 
	return &Config{
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        buildDatabaseURL(),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		JWTExpirationHours: jwtExpiration,
	}
}
 
func buildDatabaseURL() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "recipe_platform")
 
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
}
 
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}