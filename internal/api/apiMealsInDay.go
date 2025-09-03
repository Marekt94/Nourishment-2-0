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
func (ms *MealsInDayAPI) GetMealInDay(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	meal := ms.Repo.GetMealsInDay(id)
	if meal.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, meal)
}

func (ms *MealsInDayAPI) GetMealsInDay(c *gin.Context) {
	repo := ms.Repo.(meal.MealsInDayRepoIntf)
	meals := repo.GetMealsInDays()
	c.IndentedJSON(http.StatusOK, meals)
}

func (ms *MealsInDayAPI) CreateMealInDay(c *gin.Context) {
	var m meal.MealInDay
	if err := c.ShouldBindJSON(&m); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := ms.Repo.CreateMealsInDay(&m)
	if id <= 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "CreateMealInDay failed"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func (ms *MealsInDayAPI) UpdateMealInDay(c *gin.Context) {
	var m meal.MealInDay
	if err := c.ShouldBindJSON(&m); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ms.Repo.UpdateMealsInDay(&m)
	c.Status(http.StatusOK)
}

func (ms *MealsInDayAPI) DeleteMealInDay(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok := ms.Repo.DeleteMealsInDay(id)
	if ok {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}
