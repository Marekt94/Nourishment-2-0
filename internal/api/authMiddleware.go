package api

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"nourishment_20/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	jwtGenerator auth.JWTGenerator
}

// AuthMiddleware: weryfikacja JWT i wstrzyknięcie claims do kontekstu
func (m *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1) Pobierz i sprawdź nagłówek Authorization: Bearer <token>
		authz := c.GetHeader("Authorization")
		if authz == "" || !strings.HasPrefix(authz, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
			return
		}
		raw := strings.TrimSpace(strings.TrimPrefix(authz, "Bearer "))

		audience := os.Getenv("JWT_AUDIENCE")
		issuer := os.Getenv("JWT_ISSUER")
		clockSkew, err := strconv.Atoi(os.Getenv("JWT_CLOCK_SKEW"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error parsing clock skew"})
		}

		// 2) Parsuj i zweryfikuj sygnaturę z białą listą algorytmów
		token, err := jwt.Parse(raw, m.jwtGenerator.Validate, jwt.WithAudience(audience),
			jwt.WithIssuer(issuer), jwt.WithLeeway(time.Duration(clockSkew)*time.Second))
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// 3) Zweryfikuj registered claims (exp/iat/nbf sprawdzane przez opcje powyżej)
		claims, ok := token.Claims.(auth.InternalClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		// 4) Wyciągnij tożsamość i kontekst
		sub, _ := claims.GetSubject()

		// Scopes mogą być w dwóch formatach: "scopes": []string lub "scope": "a b c"
		scope := m.jwtGenerator.StringSliceToScope(claims.Scope)

		// 5) Umieść minimalny kontekst do Gin
		c.Set("sub", sub)
		c.Set("scopes", scope)

		c.Next()
	}
}
