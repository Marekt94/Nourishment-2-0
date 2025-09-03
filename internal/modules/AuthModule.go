package modules

import (
	"nourishment_20/internal/api"
	"nourishment_20/internal/auth"

	"github.com/gin-gonic/gin"
)

type ModuleAuth struct {
	Engine     *gin.Engine
	AuthServer *api.AuthServer
	PermRepo   auth.PermissionsIntf
}

func (m ModuleAuth) ExposeMethods() {
	m.Engine.POST(api.PATH_LOGIN, m.AuthServer.GenerateToken)
}

func (m ModuleAuth) RegisterPermissions() {

}

func (m ModuleAuth) GetName() string {
	return "Auth"
}
