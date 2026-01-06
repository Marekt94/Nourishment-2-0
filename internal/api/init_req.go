package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Error string `json:"error"`
}

func DefaultReq(c *gin.Context) {
	str := "{\"status\": \"ok\"}" // [AI REFACTOR]
	c.IndentedJSON(http.StatusOK, str)
}
