package api

import (
	"net/http"
	meal "nourishment_20/internal/mealDomain"
	"nourishment_20/internal/mealOptimizer"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoriesAPI struct {
	Repo     meal.CategoriesRepoIntf
	AIClient mealOptimizer.Optimizer
}

// CRUD dla Categories
func (ms *CategoriesAPI) GetCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat := ms.Repo.GetCategory(id)
	if cat.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, cat)
}

func (ms *CategoriesAPI) GetCategories(c *gin.Context) {
	cats := ms.Repo.GetCategories()
	c.IndentedJSON(http.StatusOK, cats)
}

func (ms *CategoriesAPI) CreateCategory(c *gin.Context) {
	var cat meal.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := ms.Repo.CreateCategory(&cat)
	if id <= 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "CreateCategory failed"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func (ms *CategoriesAPI) UpdateCategory(c *gin.Context) {
	var cat meal.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ms.Repo.UpdateCategory(&cat)
	c.Status(http.StatusOK)
}

func (ms *CategoriesAPI) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok := ms.Repo.DeleteCategory(id)
	if ok {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}
