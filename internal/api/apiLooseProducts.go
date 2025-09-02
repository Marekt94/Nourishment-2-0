package api

import (
	"net/http"
	meal "nourishment_20/internal/mealDomain"
	"nourishment_20/internal/mealOptimizer"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LooseProductsInDayAPI struct {
	Repo     meal.LooseProductsInDayRepoIntf
	AIClient mealOptimizer.Optimizer
}

// CRUD dla LooseProductInDay
func (ms *LooseProductsInDayAPI) GetLooseProductInDay(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product := ms.Repo.GetLooseProductInDay(id)
	if product.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, product)
}

func (ms *LooseProductsInDayAPI) GetLooseProductsInDay(c *gin.Context) {
	dayIdStr := c.Query("dayId")
	if dayIdStr == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "missing query parameter: dayId"})
		return
	}
	dayId, err := strconv.Atoi(dayIdStr)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	products := ms.Repo.GetLooseProductsInDay(dayId)
	c.IndentedJSON(http.StatusOK, products)
}

func (ms *LooseProductsInDayAPI) CreateLooseProductInDay(c *gin.Context) {
	var p meal.LooseProductInDay
	if err := c.ShouldBindJSON(&p); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := ms.Repo.CreateLooseProductInDay(&p)
	if id <= 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "CreateLooseProductInDay failed"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func (ms *LooseProductsInDayAPI) UpdateLooseProductInDay(c *gin.Context) {
	var p meal.LooseProductInDay
	if err := c.ShouldBindJSON(&p); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ms.Repo.UpdateLooseProductInDay(&p)
	c.Status(http.StatusOK)
}

func (ms *LooseProductsInDayAPI) DeleteLooseProductInDay(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok := ms.Repo.DeleteLooseProductInDay(id)
	if ok {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}
