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

// OptimizeMeal godoc
// @Security BearerAuth
// @Summary      Optimize meal by AI
// @Description  Optimize meal ingredients to match target kcal using AI
// @Tags         ai-optimizer
// @Accept       json
// @Produce      json
// @Param        kcal  query     float64    true  "Target kcal"
// @Param        meal  body      meal.Meal  true  "Meal to optimize"
// @Success      200   {object}  meal.Meal
// @Failure      400   {object}  Error
// @Failure      500   {object}  Error
// @Router       /optimizemeal [post]
func (ms *AIOptimizerAPI) OptimizeMeal(c *gin.Context) {
	kcalStr := c.Query("kcal")
	kcal, err := strconv.ParseFloat(kcalStr, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	var meal meal.Meal
	if err := c.ShouldBindJSON(&meal); err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	res, err := ms.AIClient.OptimizeMeal(&meal, kcal)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, Error{Error: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, res)
}

// OptimizeMealFromRepo godoc
// @Security BearerAuth
// @Summary      Optimize meal from repository by AI
// @Description  Get meal from repository and optimize its ingredients to match target kcal using AI
// @Tags         ai-optimizer
// @Accept       json
// @Produce      json
// @Param        id    path      int      true   "Meal ID"
// @Param        kcal  query     float64  true   "Target kcal"
// @Success      200   {object}  meal.Meal
// @Failure      400   {object}  Error
// @Failure      404   {object}  Error
// @Failure      500   {object}  Error
// @Router       /optimizemeal/{id} [post]
func (ms *AIOptimizerAPI) OptimizeMealFromRepo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	kcal, err := strconv.ParseFloat(c.Query("kcal"), 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	meal := ms.Repo.GetMeal(id)
	if meal.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	res, err := ms.AIClient.OptimizeMeal(&meal, kcal)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, Error{Error: err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, res)
}
