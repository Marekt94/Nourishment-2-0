package api

import (
	"net/http"
	meal "nourishment_20/internal/mealDomain"
	"nourishment_20/internal/mealOptimizer"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AIOptimizerAPI struct {
	Repo     meal.MealsRepoIntf
	AIClient mealOptimizer.Optimizer
}

func (ms *AIOptimizerAPI) OptimizeMeal(c *gin.Context) {
	kcalStr := c.Query("kcal")
	kcal, err := strconv.ParseFloat(kcalStr, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var meal meal.Meal
	if err := c.ShouldBindJSON(&meal); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := ms.AIClient.OptimizeMeal(&meal, kcal)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, res)
}

func (ms *AIOptimizerAPI) OptimizeMealFromRepo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	kcal, err := strconv.ParseFloat(c.Query("kcal"), 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	meal := ms.Repo.GetMeal(id)
	if meal.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	res, err := ms.AIClient.OptimizeMeal(&meal, kcal)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, res)
}
