package modules

import (
	"nourishment_20/internal/api"
	"nourishment_20/internal/auth"
	"nourishment_20/kernel"

	"github.com/gin-gonic/gin"
)

type ModuleOptimizeMeal struct {
	Engine        *gin.Engine
	AuthValidator *api.AuthMiddleware
	MethodExposer *api.AIOptimizerAPI
	PermRepo      auth.PermissionsIntf
}

func (m *ModuleOptimizeMeal) ExposeMethods() {
	m.Engine.POST(api.PATH_OPTIMIZEMEAL, m.AuthValidator.Middleware, m.MethodExposer.OptimizeMeal)
	m.Engine.POST(api.PATH_OPTIMIZEMEAL_WITH_ID, m.AuthValidator.Middleware, m.MethodExposer.OptimizeMealFromRepo)
}

func (m *ModuleOptimizeMeal) RegisterPermissions() {
	m.PermRepo.RegisterPermissions(api.RESOURCE_OPTIMIZEMEAL, []string{"write"})
	m.PermRepo.RegisterUserPermission(kernel.ADMIN_USER_NAME, api.RESOURCE_OPTIMIZEMEAL, "write")
}

func (m *ModuleOptimizeMeal) GetName() string {
	return "OptimizeMeal"
}
