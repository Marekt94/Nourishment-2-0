package modules

import (
	"nourishment_20/internal/api"
	"nourishment_20/internal/auth"
	meal "nourishment_20/internal/mealDomain"
	"nourishment_20/kernel"

	"github.com/gin-gonic/gin"
)

type ModuleMeals struct {
	Repo          meal.MealsRepoIntf
	Engine        *gin.Engine
	AuthValidator *api.AuthMiddleware
	MethodExposer *api.MealsAPI
	PermRepo      auth.PermissionsIntf
}

func (m *ModuleMeals) ExposeMethods() {
	m.Engine.GET(api.PATH_MEALS, m.AuthValidator.Middleware, m.MethodExposer.GetMeals)
	m.Engine.GET(api.PATH_MEALS_WITH_ID, m.AuthValidator.Middleware, m.MethodExposer.GetMeal)
	m.Engine.POST(api.PATH_MEALS, m.AuthValidator.Middleware, m.MethodExposer.CreateMeal)
	m.Engine.PUT(api.PATH_MEALS, m.AuthValidator.Middleware, m.MethodExposer.UpdateMeal)
	m.Engine.DELETE(api.PATH_MEALS_WITH_ID, m.AuthValidator.Middleware, m.MethodExposer.DeleteMeal)
}

func (m *ModuleMeals) RegisterPermissions() {
	m.PermRepo.RegisterPermissions(api.RESOURCE_MEALS, []string{"read", "write"})
	m.PermRepo.RegisterUserPermission(kernel.ADMIN_USER_NAME, api.RESOURCE_MEALS, "read")
	m.PermRepo.RegisterUserPermission(kernel.ADMIN_USER_NAME, api.RESOURCE_MEALS, "write")
	m.PermRepo.RegisterUserPermission(kernel.READER_USER_NAME, api.RESOURCE_MEALS, "read")
}

func (m *ModuleMeals) GetName() string {
	return "Meals"
}
