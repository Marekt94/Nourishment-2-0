package api

import (
	"net/http"
	meal "nourishment_20/internal/mealDomain"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MealsInDayAPI struct {
	Repo meal.MealsInDayRepoIntf
}

// CRUD dla MealsInDay

// GetMealInDay godoc
// @Security BearerAuth
// @Summary      Get meal in day by id
// @Description  Get a single meal in day by its ID
// @Tags         mealsinday
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "MealInDay ID"
// @Success      200  {object}  meal.MealInDay
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Router       /mealsinday/{id} [get]
func (ms *MealsInDayAPI) GetMealInDay(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	meal := ms.Repo.GetMealsInDay(id)
	if meal.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, meal)
}

// GetMealsInDay godoc
// @Security BearerAuth
// @Summary      List meals in day
// @Description  Get list of all meals in day
// @Tags         mealsinday
// @Accept       json
// @Produce      json
// @Success      200  {array}   meal.MealInDay
// @Router       /mealsinday [get]
func (ms *MealsInDayAPI) GetMealsInDay(c *gin.Context) {
	repo := ms.Repo.(meal.MealsInDayRepoIntf)
	meals := repo.GetMealsInDays()
	c.IndentedJSON(http.StatusOK, meals)
}

// CreateMealInDay godoc
// @Security BearerAuth
// @Summary      Create a new meal in day
// @Description  Create a new meal in day from JSON body
// @Tags         mealsinday
// @Accept       json
// @Produce      json
// @Param        mealinday  body      meal.MealInDay  true  "MealInDay payload"
// @Success      200        {object}  map[string]int64
// @Failure      400        {object}  Error
// @Failure      500        {object}  Error
// @Router       /mealsinday [post]
func (ms *MealsInDayAPI) CreateMealInDay(c *gin.Context) {
	var m meal.MealInDay
	if err := c.ShouldBindJSON(&m); err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	id := ms.Repo.CreateMealsInDay(&m)
	if id <= 0 {
		c.IndentedJSON(http.StatusInternalServerError, Error{Error: "CreateMealInDay failed"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

// UpdateMealInDay godoc
// @Security BearerAuth
// @Summary      Update an existing meal in day
// @Description  Update meal in day by JSON body (must contain ID)
// @Tags         mealsinday
// @Accept       json
// @Produce      json
// @Param        mealinday  body      meal.MealInDay  true  "MealInDay payload"
// @Success      200        {object}  nil
// @Failure      400        {object}  Error
// @Router       /mealsinday [put]
func (ms *MealsInDayAPI) UpdateMealInDay(c *gin.Context) {
	var m meal.MealInDay
	if err := c.ShouldBindJSON(&m); err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	ms.Repo.UpdateMealsInDay(&m)
	c.Status(http.StatusOK)
}

// DeleteMealInDay godoc
// @Security BearerAuth
// @Summary      Delete a meal in day
// @Description  Delete meal in day by id
// @Tags         mealsinday
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "MealInDay ID"
// @Success      200  {object}  nil
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Router       /mealsinday/{id} [delete]
func (ms *MealsInDayAPI) DeleteMealInDay(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	ok := ms.Repo.DeleteMealsInDay(id)
	if ok {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}
