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

var METHOD_MAP_RIGHTS = map[string]string{
	"GET":    "read",
	"POST":   "write",
	"PUT":    "write",
	"DELETE": "write",
}

type AuthMiddleware struct {
	JwtGenerator auth.JWTGenerator
}

// Mapowanie URL path na nazwę zasobu
func (m *AuthMiddleware) getResourceFromPath(path string) string {
	// Usuń parametry z URL (np. /meals/123 -> /meals)
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return ""
	}

	// Pierwsza część URL to resource
	return parts[0]
}

// Sprawdź czy użytkownik ma uprawnienia do zasobu
func (m *AuthMiddleware) hasPermission(scope map[string][]string, resource, method string) bool {
	requiredRight, exists := METHOD_MAP_RIGHTS[method]
	if !exists {
		return false
	}

	// Sprawdź czy użytkownik ma wymagane uprawnienie
	rights, hasResource := scope[resource]
	if !hasResource {
		return false
	}

	for _, right := range rights {
		if right == requiredRight {
			return true
		}
	}
	return false
}

// AuthMiddleware: weryfikacja JWT i wstrzyknięcie claims do kontekstu
func (m *AuthMiddleware) Middleware(c *gin.Context) {
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
	token, err := jwt.ParseWithClaims(raw, &auth.InternalClaims{}, m.JwtGenerator.Validate,
		jwt.WithAudience(audience), jwt.WithIssuer(issuer),
		jwt.WithLeeway(time.Duration(clockSkew)*time.Second))
	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// 3) Zweryfikuj registered claims (exp/iat/nbf sprawdzane przez opcje powyżej)
	claims, ok := token.Claims.(*auth.InternalClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		return
	}

	// 4) Wyciągnij tożsamość i kontekst
	sub, _ := claims.GetSubject()

	// Scopes mogą być w dwóch formatach: "scopes": []string lub "scope": "a b c"
	scope := m.JwtGenerator.StringSliceToScope(claims.Scope)

	// 5) Sprawdź uprawnienia do danego endpointa
	resource := m.getResourceFromPath(c.Request.URL.Path)
	method := c.Request.Method

	if !m.hasPermission(scope, resource, method) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error":    "insufficient permissions",
			"resource": resource,
			"method":   method,
		})
		return
	}

	// 6) Umieść minimalny kontekst do Gin
	c.Set("sub", sub)
	c.Set("scope", scope)

	c.Next()
}
