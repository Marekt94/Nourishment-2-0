package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DefaultReq(c *gin.Context){
	str := "{\"status\": \"ok\"}";
	c.IndentedJSON(http.StatusOK, str)
}