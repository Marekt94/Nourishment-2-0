package modules

import (
	"nourishment_20/internal/api"
	"nourishment_20/internal/auth"
	meal "nourishment_20/internal/mealDomain"
	"nourishment_20/kernel"

	"github.com/gin-gonic/gin"
)

type ModuleShoppingList struct {
	Repo          meal.ShoppingListRepoIntf
	Engine        *gin.Engine
	AuthValidator *api.AuthMiddleware
	ListExposer   *api.ShoppingListAPI
	ProdExposer   *api.ProductsInShoppingListAPI
	PermRepo      auth.PermissionsIntf
}

func (m *ModuleShoppingList) ExposeMethods() {
	// Shopping Lists
	m.Engine.GET(api.PATH_SHOPPINGLISTS, m.AuthValidator.Middleware, m.ListExposer.GetShoppingLists)
	m.Engine.GET(api.PATH_SHOPPINGLISTS_WITH_ID, m.AuthValidator.Middleware, m.ListExposer.GetShoppingList)
	m.Engine.POST(api.PATH_SHOPPINGLISTS, m.AuthValidator.Middleware, m.ListExposer.CreateShoppingList)
	m.Engine.POST(api.PATH_SHOPPINGLISTS_GENERATE, m.AuthValidator.Middleware, m.ListExposer.GenerateShoppingList)
	m.Engine.PUT(api.PATH_SHOPPINGLISTS, m.AuthValidator.Middleware, m.ListExposer.UpdateShoppingList)
	m.Engine.DELETE(api.PATH_SHOPPINGLISTS_WITH_ID, m.AuthValidator.Middleware, m.ListExposer.DeleteShoppingList)

	// Products in Shopping Lists
	m.Engine.POST(api.PATH_SHOPPINGLIST_PRODS, m.AuthValidator.Middleware, m.ProdExposer.AddProduct)
	m.Engine.PUT(api.PATH_SHOPPINGLIST_PRODS, m.AuthValidator.Middleware, m.ProdExposer.UpdateProduct)
	m.Engine.DELETE(api.PATH_SHOPPINGLIST_PRODS_WITH_ID, m.AuthValidator.Middleware, m.ProdExposer.DeleteProduct)
}

func (m *ModuleShoppingList) RegisterPermissions() {
	m.PermRepo.RegisterPermissions(api.RESOURCE_SHOPPINGLISTS, []string{"read", "write"})
	m.PermRepo.RegisterUserPermission(kernel.ADMIN_USER_NAME, api.RESOURCE_SHOPPINGLISTS, "read")
	m.PermRepo.RegisterUserPermission(kernel.ADMIN_USER_NAME, api.RESOURCE_SHOPPINGLISTS, "write")

	m.PermRepo.RegisterPermissions(api.RESOURCE_SHOPPINGLIST_PRODS, []string{"read", "write"})
	m.PermRepo.RegisterUserPermission(kernel.ADMIN_USER_NAME, api.RESOURCE_SHOPPINGLIST_PRODS, "read")
	m.PermRepo.RegisterUserPermission(kernel.ADMIN_USER_NAME, api.RESOURCE_SHOPPINGLIST_PRODS, "write")
}

func (m *ModuleShoppingList) GetName() string {
	return "ShoppingList"
}
