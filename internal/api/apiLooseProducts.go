package api

import (
	"net/http"
	meal "nourishment_20/internal/mealDomain"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LooseProductsInDayAPI struct {
	Repo meal.LooseProductsInDayRepoIntf
}

// CRUD dla LooseProductInDay

// GetLooseProductInDay godoc
// @Security BearerAuth
// @Summary      Get loose product in day by id
// @Description  Get a single loose product in day by its ID
// @Tags         looseproductsinday
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "LooseProductInDay ID"
// @Success      200  {object}  meal.LooseProductInDay
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Router       /looseproductsinday/{id} [get]
func (ms *LooseProductsInDayAPI) GetLooseProductInDay(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	product := ms.Repo.GetLooseProductInDay(id)
	if product.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, product)
}

// GetLooseProductsInDay godoc
// @Security BearerAuth
// @Summary      List loose products in day
// @Description  Get list of loose products in day by day ID
// @Tags         looseproductsinday
// @Accept       json
// @Produce      json
// @Param        dayId  query     int  true  "Day ID"
// @Success      200    {array}   meal.LooseProductInDay
// @Failure      400    {object}  Error
// @Router       /looseproductsinday [get]
func (ms *LooseProductsInDayAPI) GetLooseProductsInDay(c *gin.Context) {
	dayIdStr := c.Query("dayId")
	if dayIdStr == "" {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: "missing query parameter: dayId"})
		return
	}
	dayId, err := strconv.Atoi(dayIdStr)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	products := ms.Repo.GetLooseProductsInDay(dayId)
	c.IndentedJSON(http.StatusOK, products)
}

// CreateLooseProductInDay godoc
// @Security BearerAuth
// @Summary      Create a new loose product in day
// @Description  Create a new loose product in day from JSON body
// @Tags         looseproductsinday
// @Accept       json
// @Produce      json
// @Param        looseproduct  body      meal.LooseProductInDay  true  "LooseProductInDay payload"
// @Success      200           {object}  map[string]int64
// @Failure      400           {object}  Error
// @Failure      500           {object}  Error
// @Router       /looseproductsinday [post]
func (ms *LooseProductsInDayAPI) CreateLooseProductInDay(c *gin.Context) {
	var p meal.LooseProductInDay
	if err := c.ShouldBindJSON(&p); err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	id := ms.Repo.CreateLooseProductInDay(&p)
	if id <= 0 {
		c.IndentedJSON(http.StatusInternalServerError, Error{Error: "CreateLooseProductInDay failed"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

// UpdateLooseProductInDay godoc
// @Security BearerAuth
// @Summary      Update an existing loose product in day
// @Description  Update loose product in day by JSON body (must contain ID)
// @Tags         looseproductsinday
// @Accept       json
// @Produce      json
// @Param        looseproduct  body      meal.LooseProductInDay  true  "LooseProductInDay payload"
// @Success      200           {object}  nil
// @Failure      400           {object}  Error
// @Router       /looseproductsinday [put]
func (ms *LooseProductsInDayAPI) UpdateLooseProductInDay(c *gin.Context) {
	var p meal.LooseProductInDay
	if err := c.ShouldBindJSON(&p); err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	ms.Repo.UpdateLooseProductInDay(&p)
	c.Status(http.StatusOK)
}

// DeleteLooseProductInDay godoc
// @Security BearerAuth
// @Summary      Delete a loose product in day
// @Description  Delete loose product in day by id
// @Tags         looseproductsinday
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "LooseProductInDay ID"
// @Success      200  {object}  nil
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Router       /looseproductsinday/{id} [delete]
func (ms *LooseProductsInDayAPI) DeleteLooseProductInDay(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	ok := ms.Repo.DeleteLooseProductInDay(id)
	if ok {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}
