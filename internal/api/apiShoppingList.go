package api

import (
	"net/http"
	meal "nourishment_20/internal/mealDomain"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShoppingListAPI struct {
	Repo meal.ShoppingListRepoIntf
}

// GetShoppingLists godoc
// @Security BearerAuth
// @Summary      List shopping lists
// @Description  Get all shopping lists with their products
// @Tags         shopping-lists
// @Accept       json
// @Produce      json
// @Success      200  {array}   meal.ShoppingList
// @Router       /shopping-lists [get]
func (api *ShoppingListAPI) GetShoppingLists(c *gin.Context) {
	lists := api.Repo.GetShoppingLists()
	c.IndentedJSON(http.StatusOK, lists)
}

// GetShoppingList godoc
// @Security BearerAuth
// @Summary      Get shopping list by id
// @Description  Get a single shopping list with its products
// @Tags         shopping-lists
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Shopping List ID"
// @Success      200  {object}  meal.ShoppingList
// @Failure      404  {object}  Error
// @Router       /shopping-lists/{id} [get]
func (api *ShoppingListAPI) GetShoppingList(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	list := api.Repo.GetShoppingList(id)
	if list.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, list)
}

// CreateShoppingList godoc
// @Security BearerAuth
// @Summary      Create a shopping list
// @Description  Create a shopping list, optionally with products
// @Tags         shopping-lists
// @Accept       json
// @Produce      json
// @Param        list  body      meal.ShoppingList  true  "Shopping List payload"
// @Success      200   {object}  map[string]int64
// @Failure      400   {object}  Error
// @Router       /shopping-lists [post]
func (api *ShoppingListAPI) CreateShoppingList(c *gin.Context) {
	var list meal.ShoppingList
	if err := c.ShouldBindJSON(&list); err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	id := api.Repo.CreateShoppingList(&list)
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

// UpdateShoppingList godoc
// @Security BearerAuth
// @Summary      Update a shopping list
// @Description  Update shopping list metadata (name)
// @Tags         shopping-lists
// @Accept       json
// @Produce      json
// @Param        list  body      meal.ShoppingList  true  "Shopping List payload"
// @Success      200   {object}  nil
// @Failure      400   {object}  Error
// @Router       /shopping-lists [put]
func (api *ShoppingListAPI) UpdateShoppingList(c *gin.Context) {
	var list meal.ShoppingList
	if err := c.ShouldBindJSON(&list); err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	api.Repo.UpdateShoppingList(&list)
	c.Status(http.StatusOK)
}

// DeleteShoppingList godoc
// @Security BearerAuth
// @Summary      Delete a shopping list
// @Description  Delete shopping list and all its products
// @Tags         shopping-lists
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Shopping List ID"
// @Success      200  {object}  nil
// @Router       /shopping-lists/{id} [delete]
func (api *ShoppingListAPI) DeleteShoppingList(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	api.Repo.DeleteShoppingList(id)
	c.Status(http.StatusOK)
}

type ProductsInShoppingListAPI struct {
	Repo meal.ShoppingListRepoIntf
}

// AddProduct godoc
// @Security BearerAuth
// @Summary      Add product to list
// @Description  Add a product to an existing shopping list
// @Tags         shopping-list-products
// @Accept       json
// @Produce      json
// @Param        product  body      meal.ProductInShoppingList  true  "Product payload"
// @Success      200      {object}  map[string]int64
// @Failure      400      {object}  Error
// @Router       /shopping-list-products [post]
func (api *ProductsInShoppingListAPI) AddProduct(c *gin.Context) {
	var p meal.ProductInShoppingList
	if err := c.ShouldBindJSON(&p); err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	id := api.Repo.AddProductToShoppingList(&p)
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

// UpdateProduct godoc
// @Security BearerAuth
// @Summary      Update product on list
// @Description  Update weight or bought status of a product on a list
// @Tags         shopping-list-products
// @Accept       json
// @Produce      json
// @Param        product  body      meal.ProductInShoppingList  true  "Product payload"
// @Success      200      {object}  nil
// @Failure      400      {object}  Error
// @Router       /shopping-list-products [put]
func (api *ProductsInShoppingListAPI) UpdateProduct(c *gin.Context) {
	var p meal.ProductInShoppingList
	if err := c.ShouldBindJSON(&p); err != nil {
		c.IndentedJSON(http.StatusBadRequest, Error{Error: err.Error()})
		return
	}
	api.Repo.UpdateProductInShoppingList(&p)
	c.Status(http.StatusOK)
}

// DeleteProduct godoc
// @Security BearerAuth
// @Summary      Remove product from list
// @Description  Remove a product entry from a shopping list
// @Tags         shopping-list-products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product Entry ID"
// @Success      200  {object}  nil
// @Router       /shopping-list-products/{id} [delete]
func (api *ProductsInShoppingListAPI) DeleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	api.Repo.DeleteProductFromShoppingList(id)
	c.Status(http.StatusOK)
}
