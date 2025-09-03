package modules

import (
	"nourishment_20/internal/api"
	"nourishment_20/internal/auth"
	meal "nourishment_20/internal/mealDomain"
	"nourishment_20/kernel"

	"github.com/gin-gonic/gin"
)

type ModuleProducts struct {
	Repo          meal.ProductsRepoIntf
	Engine        *gin.Engine
	AuthValidator *api.AuthMiddleware
	MethodExposer *api.ProductsAPI
	PermRepo      auth.PermissionsIntf
}

func (m *ModuleProducts) ExposeMethods() {
	m.Engine.GET(api.PATH_PRODUCTS, m.AuthValidator.Middleware, m.MethodExposer.GetProducts)
	m.Engine.GET(api.PATH_PRODUCTS_WITH_ID, m.AuthValidator.Middleware, m.MethodExposer.GetProduct)
	m.Engine.POST(api.PATH_PRODUCTS, m.AuthValidator.Middleware, m.MethodExposer.CreateProduct)
	m.Engine.PUT(api.PATH_PRODUCTS, m.AuthValidator.Middleware, m.MethodExposer.UpdateProduct)
	m.Engine.DELETE(api.PATH_PRODUCTS_WITH_ID, m.AuthValidator.Middleware, m.MethodExposer.DeleteProduct)
}

func (m *ModuleProducts) RegisterPermissions() {
	m.PermRepo.RegisterPermissions(api.RESOURCE_PRODUCTS, []string{"read", "write"})
	m.PermRepo.RegisterUserPermission(kernel.ADMIN_USER_NAME, api.RESOURCE_PRODUCTS, "read")
	m.PermRepo.RegisterUserPermission(kernel.ADMIN_USER_NAME, api.RESOURCE_PRODUCTS, "write")
	m.PermRepo.RegisterUserPermission(kernel.READER_USER_NAME, api.RESOURCE_PRODUCTS, "read")
}

func (m *ModuleProducts) GetName() string {
	return "Products"
}
