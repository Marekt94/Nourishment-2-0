package modules

import (
	"nourishment_20/internal/api"
	"nourishment_20/internal/auth"
	meal "nourishment_20/internal/mealDomain"
	"nourishment_20/kernel"

	"github.com/gin-gonic/gin"
)

type ModuleMealsInDay struct {
	Repo          meal.MealsInDayRepoIntf
	Engine        *gin.Engine
	AuthValidator *api.AuthMiddleware
	MethodExposer *api.MealsInDayAPI
	PermRepo      auth.PermissionsIntf
}

func (m *ModuleMealsInDay) ExposeMethods() {
	m.Engine.GET(api.PATH_MEALSINDAY, m.AuthValidator.Middleware, m.MethodExposer.GetMealsInDay)
	m.Engine.GET(api.PATH_MEALSINDAY_WITH_ID, m.AuthValidator.Middleware, m.MethodExposer.GetMealInDay)
	m.Engine.POST(api.PATH_MEALSINDAY, m.AuthValidator.Middleware, m.MethodExposer.CreateMealInDay)
	m.Engine.PUT(api.PATH_MEALSINDAY, m.AuthValidator.Middleware, m.MethodExposer.UpdateMealInDay)
	m.Engine.DELETE(api.PATH_MEALSINDAY_WITH_ID, m.AuthValidator.Middleware, m.MethodExposer.DeleteMealInDay)
}

func (m *ModuleMealsInDay) RegisterPermissions() {
	m.PermRepo.RegisterPermissions(api.RESOURCE_MEALSINDAY, []string{"read", "write"})
	m.PermRepo.RegisterUserPermission(kernel.ADMIN_USER_NAME, api.RESOURCE_MEALSINDAY, "read")
	m.PermRepo.RegisterUserPermission(kernel.ADMIN_USER_NAME, api.RESOURCE_MEALSINDAY, "write")
	m.PermRepo.RegisterUserPermission(kernel.READER_USER_NAME, api.RESOURCE_MEALSINDAY, "read")
}

func (m *ModuleMealsInDay) GetName() string {
	return "MealsInDay"
}
