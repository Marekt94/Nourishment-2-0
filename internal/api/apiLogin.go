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

type GenerateTokenRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GenerateTokenResponse struct {
	Token string `json:"token"`
}

// @Summary      Generate Auth Token
// @Description  Generate JWT token for user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body api.GenerateTokenRequest true "User login and password"
// @Success      200 {object} api.GenerateTokenResponse
// @Failure      401 {object} Error
// @Failure      403 {object} Error
// @Failure      500 {object} Error
// @Router       /login [post]
func (a *AuthServer) GenerateToken(c *gin.Context) {
	var loginRequest GenerateTokenRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(400, Error{Error: "Invalid request format"})
		return
	}

	// Currently using: a.Repo.UserExists(login, password)
	// Check if this method exists in PermissionsIntf
	exists := a.UserRepo.IsUserExists(loginRequest.Login, loginRequest.Password)
	if exists == auth.NO_USER_ID {
		c.JSON(401, Error{Error: "Invalid credentials"})
		return
	}

	// Currently using: a.JWTGenerator.GenerateToken(login)
	// Check if this method exists in JWTGenerator
	tokenRaw, err := a.JWTGenerator.GetJWT(loginRequest.Login)
	if err != nil {
		c.JSON(500, Error{Error: "Failed to generate token"})
		return
	}
	token, err := a.JWTGenerator.JWTToString(tokenRaw)
	if err != nil {
		c.JSON(500, Error{Error: "Failed to generate token"})
		return
	}

	resp := GenerateTokenResponse{
		Token: *token,
	}

	c.JSON(200, resp)
}
