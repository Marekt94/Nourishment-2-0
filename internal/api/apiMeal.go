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

func (ms *MealsAPI) GetMeal(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err == nil {
		c.IndentedJSON(http.StatusOK, ms.Repo.GetMeal(id))
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (ms *MealsAPI) GetMeals(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, ms.Repo.GetMeals())
}

func (ms *MealsAPI) CreateMeal(c *gin.Context) {
	var m meal.Meal // [AI REFACTOR]
	if err := c.ShouldBindJSON(&m); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := ms.Repo.CreateMeal(&m)
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func (ms *MealsAPI) UpdateMeal(c *gin.Context) {
	var m meal.Meal // [AI REFACTOR]
	if err := c.ShouldBindJSON(&m); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ms.Repo.UpdateMeal(&m)
	c.Status(http.StatusOK)
}

func (ms *MealsAPI) DeleteMeal(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok := ms.Repo.DeleteMeal(id)
	if ok {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}
