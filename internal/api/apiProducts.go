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
func (ms *ProductsAPI) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product := ms.Repo.GetProduct(id)
	if product.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, product)
}

func (ms *ProductsAPI) GetProducts(c *gin.Context) {
	products := ms.Repo.GetProducts()
	c.IndentedJSON(http.StatusOK, products)
}

func (ms *ProductsAPI) CreateProduct(c *gin.Context) {
	var p meal.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := ms.Repo.CreateProduct(&p)
	if id <= 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "CreateProduct failed"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func (ms *ProductsAPI) UpdateProduct(c *gin.Context) {
	var p meal.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ms.Repo.UpdateProduct(&p)
	c.Status(http.StatusOK)
}

func (ms *ProductsAPI) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok := ms.Repo.DeleteProduct(id)
	if ok {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}
