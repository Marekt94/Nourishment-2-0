package api

import (
	"net/http"
	utils "nourishment_20/internal"
	meal "nourishment_20/internal/mealDomain"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoriesAPI struct {
	Repo meal.CategoriesRepoIntf
}

// CRUD dla Categories

// GetCategory godoc
// @Security BearerAuth
// @Summary      Get category by id
// @Description  Get a single category by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  meal.Category
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Router       /categories/{id} [get]
func (ms *CategoriesAPI) GetCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, utils.Error{Error: err.Error()})
		return
	}
	cat := ms.Repo.GetCategory(id)
	if cat.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, cat)
}

// GetCategories godoc
// @Security BearerAuth
// @Summary      List categories
// @Description  Get list of all categories
// @Tags         categories
// @Accept       json
// @Produce      json
// @Success      200  {array}   meal.Category
// @Router       /categories [get]
func (ms *CategoriesAPI) GetCategories(c *gin.Context) {
	cats := ms.Repo.GetCategories()
	c.IndentedJSON(http.StatusOK, cats)
}

// CreateCategory godoc
// @Security BearerAuth
// @Summary      Create a new category
// @Description  Create a new category from JSON body
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        category  body      meal.Category  true  "Category payload"
// @Success      200       {object}  map[string]int64
// @Failure      400       {object}  utils.Error
// @Failure      500       {object}  utils.Error
// @Router       /categories [post]
func (ms *CategoriesAPI) CreateCategory(c *gin.Context) {
	var cat meal.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.IndentedJSON(http.StatusBadRequest, utils.Error{Error: err.Error()})
		return
	}
	id := ms.Repo.CreateCategory(&cat)
	if id <= 0 {
		c.IndentedJSON(http.StatusInternalServerError, utils.Error{Error: "CreateCategory failed"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

// UpdateCategory godoc
// @Security BearerAuth
// @Summary      Update an existing category
// @Description  Update category by JSON body (must contain ID)
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        category  body      meal.Category  true  "Category payload"
// @Success      200       {object}  nil
// @Failure      400       {object}  utils.Error
// @Router       /categories [put]
func (ms *CategoriesAPI) UpdateCategory(c *gin.Context) {
	var cat meal.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.IndentedJSON(http.StatusBadRequest, utils.Error{Error: err.Error()})
		return
	}
	ms.Repo.UpdateCategory(&cat)
	c.Status(http.StatusOK)
}

// DeleteCategory godoc
// @Security BearerAuth
// @Summary      Delete a category
// @Description  Delete category by id
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  nil
// @Failure      400  {object}  utils.Error
// @Failure      404  {object}  utils.Error
// @Router       /categories/{id} [delete]
func (ms *CategoriesAPI) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, utils.Error{Error: err.Error()})
		return
	}
	ok := ms.Repo.DeleteCategory(id)
	if ok {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}
