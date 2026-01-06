package api

import (
	"net/http"
	meal "nourishment_20/internal/mealDomain"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductsAPI struct {
	Repo meal.ProductsRepoIntf
}

// CRUD dla Products

// GetProduct godoc
// @Security BearerAuth
// @Summary      Get product by id
// @Description  Get a single product by its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  meal.Product
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Router       /products/{id} [get]
func (ms *ProductsAPI) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	product := ms.Repo.GetProduct(id)
	if product.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, product)
}

// GetProducts godoc
// @Security BearerAuth
// @Summary      List products
// @Description  Get list of all products
// @Tags         products
// @Accept       json
// @Produce      json
// @Success      200  {array}   meal.Product
// @Router       /products [get]
func (ms *ProductsAPI) GetProducts(c *gin.Context) {
	products := ms.Repo.GetProducts()
	c.IndentedJSON(http.StatusOK, products)
}

// CreateProduct godoc
// @Security BearerAuth
// @Summary      Create a new product
// @Description  Create a new product from JSON body
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      meal.Product  true  "Product payload"
// @Success      200      {object}  map[string]int64
// @Failure      400      {object}  Error
// @Failure      500      {object}  Error
// @Router       /products [post]
func (ms *ProductsAPI) CreateProduct(c *gin.Context) {
	var p meal.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	id := ms.Repo.CreateProduct(&p)
	if id <= 0 {
		c.IndentedJSON(http.StatusInternalServerError, Error{Error: "CreateProduct failed"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

// UpdateProduct godoc
// @Security BearerAuth
// @Summary      Update an existing product
// @Description  Update product by JSON body (must contain ID)
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      meal.Product  true  "Product payload"
// @Success      200      {object}  nil
// @Failure      400      {object}  Error
// @Router       /products [put]
func (ms *ProductsAPI) UpdateProduct(c *gin.Context) {
	var p meal.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	ms.Repo.UpdateProduct(&p)
	c.Status(http.StatusOK)
}

// DeleteProduct godoc
// @Security BearerAuth
// @Summary      Delete a product
// @Description  Delete product by id
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  nil
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Router       /products/{id} [delete]
func (ms *ProductsAPI) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	ok := ms.Repo.DeleteProduct(id)
	if ok {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}
