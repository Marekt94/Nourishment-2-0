package api

import (
	"nourishment_20/internal/auth"

	"github.com/gin-gonic/gin"
)

type AuthServer struct {
	UserRepo     auth.UserRepoIntf
	PermRepo     auth.PermissionsIntf
	JWTGenerator *auth.JWTGenerator
}

func (a *AuthServer) GenerateToken(c *gin.Context) {
	var loginRequest struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	// Currently using: a.Repo.UserExists(login, password)
	// Check if this method exists in PermissionsIntf
	exists := a.UserRepo.IsUserExists(loginRequest.Login, loginRequest.Password)
	if exists == auth.NO_USER_ID {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// Currently using: a.JWTGenerator.GenerateToken(login)
	// Check if this method exists in JWTGenerator
	tokenRaw, err := a.JWTGenerator.GetJWT(loginRequest.Login)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}
	token, err := a.JWTGenerator.JWTToString(tokenRaw)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}
