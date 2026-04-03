package middleware

import (
	"net/http"
	"strings"
 
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)
 
// AuthRequired valida o JWT no header Authorization e injeta o chef_id no contexto.
func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token não informado"})
			return
		}
 
		// Espera o formato: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido"})
			return
		}
 
		tokenString := parts[1]
 
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})
 
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido ou expirado"})
			return
		}
 
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token mal formado"})
			return
		}
 
		// Injeta o chef_id no contexto para uso nos handlers
		chefID := uint(claims["chef_id"].(float64))
		c.Set("chef_id", chefID)
 
		c.Next()
	}
}