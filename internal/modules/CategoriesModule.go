package modules

import (
	"nourishment_20/internal/api"
	"nourishment_20/internal/auth"
	meal "nourishment_20/internal/mealDomain"
	"nourishment_20/kernel"

	"github.com/gin-gonic/gin"
)

type ModuleCategories struct {
	Repo          meal.CategoriesRepoIntf
	Engine        *gin.Engine
	AuthValidator *api.AuthMiddleware
	MethodExposer *api.CategoriesAPI
	PermRepo      auth.PermissionsIntf
}

func (m *ModuleCategories) ExposeMethods() {
	m.Engine.GET(api.PATH_CATEGORIES, m.AuthValidator.Middleware, m.MethodExposer.GetCategories)
	m.Engine.GET(api.PATH_CATEGORIES_WITH_ID, m.AuthValidator.Middleware, m.MethodExposer.GetCategory)
	m.Engine.POST(api.PATH_CATEGORIES, m.AuthValidator.Middleware, m.MethodExposer.CreateCategory)
	m.Engine.PUT(api.PATH_CATEGORIES, m.AuthValidator.Middleware, m.MethodExposer.UpdateCategory)
	m.Engine.DELETE(api.PATH_CATEGORIES_WITH_ID, m.AuthValidator.Middleware, m.MethodExposer.DeleteCategory)
}

func (m *ModuleCategories) RegisterPermissions() {
	m.PermRepo.RegisterPermissions(api.RESOURCE_CATEGORIES, []string{"read", "write"})
	m.PermRepo.RegisterUserPermission(kernel.ADMIN_USER_NAME, api.RESOURCE_CATEGORIES, "read")
	m.PermRepo.RegisterUserPermission(kernel.ADMIN_USER_NAME, api.RESOURCE_CATEGORIES, "write")
	m.PermRepo.RegisterUserPermission(kernel.READER_USER_NAME, api.RESOURCE_CATEGORIES, "read")
}

func (m *ModuleCategories) GetName() string {
	return "Categories"
}
