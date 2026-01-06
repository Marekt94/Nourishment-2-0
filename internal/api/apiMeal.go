package api

import (
	"net/http"
	meal "nourishment_20/internal/mealDomain"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MealsAPI struct {
	Repo meal.MealsRepoIntf
}

// GetMeal godoc
// @Security BearerAuth
// @Summary      Get meal by id
// @Description  Get a single meal by its ID
// @Tags         meals
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Meal ID"
// @Success      200  {object}  meal.Meal
// @Failure      400  {object}  Error
// @Router       /meals/{id} [get]
func (ms *MealsAPI) GetMeal(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err == nil {
		c.IndentedJSON(http.StatusOK, ms.Repo.GetMeal(id))
	} else {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
	}
}

// GetMeals godoc
// @Security BearerAuth
// @Summary      List meals
// @Description  Get list of all meals
// @Tags         meals
// @Accept       json
// @Produce      json
// @Success      200  {array}   meal.Meal
// @Router       /meals [get]
func (ms *MealsAPI) GetMeals(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, ms.Repo.GetMeals())
}

// CreateMeal godoc
// @Security BearerAuth
// @Summary      Create a new meal
// @Description  Create a new meal from JSON body
// @Tags         meals
// @Accept       json
// @Produce      json
// @Param        meal  body      meal.Meal  true  "Meal payload"
// @Success      200   {object}  map[string]int64
// @Failure      400   {object}  Error
// @Router       /meals [post]
func (ms *MealsAPI) CreateMeal(c *gin.Context) {
	var m meal.Meal // [AI REFACTOR]
	if err := c.ShouldBindJSON(&m); err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	id := ms.Repo.CreateMeal(&m)
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

// UpdateMeal godoc
// @Security BearerAuth
// @Summary      Update an existing meal
// @Description  Update meal by JSON body (must contain ID)
// @Tags         meals
// @Accept       json
// @Produce      json
// @Param        meal  body      meal.Meal  true  "Meal payload"
// @Success      200   {object}  nil
// @Failure      400   {object}  Error
// @Router       /meals [put]
func (ms *MealsAPI) UpdateMeal(c *gin.Context) {
	var m meal.Meal // [AI REFACTOR]
	if err := c.ShouldBindJSON(&m); err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	ms.Repo.UpdateMeal(&m)
	c.Status(http.StatusOK)
}

// DeleteMeal godoc
// @Security BearerAuth
// @Summary      Delete a meal
// @Description  Delete meal by id
// @Tags         meals
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Meal ID"
// @Success      200  {object}  nil
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Router       /meals/{id} [delete]
func (ms *MealsAPI) DeleteMeal(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	ok := ms.Repo.DeleteMeal(id)
	if ok {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}
