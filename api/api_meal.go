package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMeal(c *gin.Context){
	id, err := strconv.Atoi(c.Params.ByName("id"));
	if err == nil {
		
		c.IndentedJSON(http.StatusOK, meal.GetMeal(id))
	} else {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
}

func GetMeals(c *gin.Context){
	c.IndentedJSON(http.StatusOK, meal.GetMeals())

}

func CreateOrUpdateMeal(c *gin.Context) {
}

func DeleteMeal(c *gin.Context){
	
}