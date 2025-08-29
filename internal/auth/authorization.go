package auth

import (
	"fmt"
	"nourishment_20/internal/logging"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const PERMS_DIVIDOR = ":"

type JWTGenerator struct {
	Repo PermissionsIntf
}

type InternalClaims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles,omitempty"`
	Scope []string `json:"scope,omitempty"`
}

func scopeToStringSlice(s []Permission) []string {
	var res []string
	for _, perm := range s {
		for _, right := range perm.Rights {
			res = append(res, fmt.Sprintf("%s"+PERMS_DIVIDOR+"%s", perm.Resource, right))
		}
	}
	return res
}

func (j *JWTGenerator) GetJWT(subject string) (*jwt.Token, error) {
	claims := InternalClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("JWT_ISSUER"),
			Subject:   subject,
			Audience:  jwt.ClaimStrings{os.Getenv("JWT_AUDIENCE")},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	scope := j.Repo.GetPermissions(subject)
	logging.Global.Debugf("scope: %v", scope)
	scopeSlice := scopeToStringSlice(scope)
	logging.Global.Debugf("scope slice: %v", scopeSlice)

	claims.Scope = scopeSlice

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token, nil
}

func (j *JWTGenerator) JWTToString(t *jwt.Token) (*string, error) {
	tokenString, err := t.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		logging.Global.Panicf("Failed to sign JWT token: %v", err)
		return nil, err
	}
	logging.Global.Tracef(`token: %v`, tokenString)
	return &tokenString, nil
}

func validate(token *jwt.Token) (any, error) {
	// Sprawdź metodę podpisywania
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(os.Getenv("JWT_SECRET")), nil
}

func stringSliceToScope(s []string) map[string][]string {
	scope := map[string][]string{}
	for _, scopeTmp := range s {
		scopeSlice := strings.Split(scopeTmp, PERMS_DIVIDOR)
		key := scopeSlice[0]
		val := scopeSlice[1]
		if valOld, e := scope[key]; e {
			valNew := append(valOld, val)
			scope[key] = valNew
		} else {
			scope[key] = []string{val}
		}
	}
	return scope
}

func (j *JWTGenerator) GetScope(t string) map[string][]string {
	var claims InternalClaims
	_, error := jwt.ParseWithClaims(t, &claims, validate)
	if error != nil {
		logging.Global.Warnf("Failed to parse JWT token: %v", error)
		return nil
	}
	return stringSliceToScope(claims.Scope)
}
