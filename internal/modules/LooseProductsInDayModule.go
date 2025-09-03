package modules

import (
	"nourishment_20/internal/api"
	"nourishment_20/internal/auth"
	meal "nourishment_20/internal/mealDomain"
	"nourishment_20/kernel"

	"github.com/gin-gonic/gin"
)

type ModuleLooseProductsInDay struct {
	Repo          meal.LooseProductsInDayRepoIntf
	Engine        *gin.Engine
	AuthValidator *api.AuthMiddleware
	MethodExposer *api.LooseProductsInDayAPI
	PermRepo      auth.PermissionsIntf
}

func (m *ModuleLooseProductsInDay) ExposeMethods() {
	m.Engine.GET(api.PATH_LOOSEPRODUCTSINDAY, m.AuthValidator.Middleware, m.MethodExposer.GetLooseProductsInDay)
	m.Engine.GET(api.PATH_LOOSEPRODUCTSINDAY_WITH_ID, m.AuthValidator.Middleware, m.MethodExposer.GetLooseProductInDay)
	m.Engine.POST(api.PATH_LOOSEPRODUCTSINDAY, m.AuthValidator.Middleware, m.MethodExposer.CreateLooseProductInDay)
	m.Engine.PUT(api.PATH_LOOSEPRODUCTSINDAY, m.AuthValidator.Middleware, m.MethodExposer.UpdateLooseProductInDay)
	m.Engine.DELETE(api.PATH_LOOSEPRODUCTSINDAY_WITH_ID, m.AuthValidator.Middleware, m.MethodExposer.DeleteLooseProductInDay)
}

func (m *ModuleLooseProductsInDay) RegisterPermissions() {
	m.PermRepo.RegisterPermissions(api.RESOURCE_LOOSEPRODUCTSINDAY, []string{"read", "write"})
	m.PermRepo.RegisterUserPermission(kernel.ADMIN_USER_NAME, api.RESOURCE_LOOSEPRODUCTSINDAY, "read")
	m.PermRepo.RegisterUserPermission(kernel.ADMIN_USER_NAME, api.RESOURCE_LOOSEPRODUCTSINDAY, "write")
	m.PermRepo.RegisterUserPermission(kernel.READER_USER_NAME, api.RESOURCE_LOOSEPRODUCTSINDAY, "read")
}

func (m *ModuleLooseProductsInDay) GetName() string {
	return "LooseProductsInDay"
}
