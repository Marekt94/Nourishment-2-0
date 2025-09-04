package modules

import (
	"database/sql"
	"nourishment_20/internal/auth"
	db "nourishment_20/internal/database"
	log "nourishment_20/internal/logging"
	"nourishment_20/kernel"
	"os"
)

// TODO dodać logowanie
type MealKernel struct {
	kernel.Kernel
	dbAccess        *sql.DB
	permissionsRepo auth.PermissionsIntf
}

func NewMealKernel() kernel.KernelIntf {
	k := MealKernel{}
	k.Kernel = kernel.Kernel{}
	return &k
}

func (k *MealKernel) initDB() {
	DbEngine := db.FBDBEngine{BaseEngineIntf: &db.BaseEngine{}}

	conf := db.DBConf{
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		Address:    os.Getenv("DB_ADDRESS"),
		PathOrName: os.Getenv("DB_NAME"),
	}
	k.dbAccess = DbEngine.Connect(&conf)
}

func (k *MealKernel) initLogger() {
	log.SetGlobalLogger(log.NewZerologLogger())
}

func (k *MealKernel) initPermissionsRepo() {
	k.permissionsRepo = &auth.PermissionsRepo{Db: k.dbAccess}
}

func (k *MealKernel) Init() {
	k.initLogger()
	k.initDB()

	k.initPermissionsRepo()
}

func (k *MealKernel) Run() {
	k.Kernel.Run()
}
