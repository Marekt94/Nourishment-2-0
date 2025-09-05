package modules

import (
	"database/sql"
	"fmt"
	"nourishment_20/internal/AIClient"
	"nourishment_20/internal/api"
	"nourishment_20/internal/auth"
	db "nourishment_20/internal/database"
	log "nourishment_20/internal/logging"
	meal "nourishment_20/internal/mealDomain"
	"nourishment_20/internal/mealOptimizer"
	"nourishment_20/kernel"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// TODO dodać logowanie
type MealKernel struct {
	kernel.Kernel

	dbAccess     *sql.DB
	serverEngine *gin.Engine

	permissionsRepo auth.PermissionsIntf
	jwtGen          *auth.JWTGenerator
	authValidator   *api.AuthMiddleware
	mealsRepo       *meal.FirebirdRepoAccess
}

func NewMealKernel() kernel.KernelIntf {
	k := MealKernel{}
	k.Kernel = kernel.NewKernel()
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

func (k *MealKernel) initServer() {
	k.serverEngine = gin.Default()
	gin.DefaultWriter = log.Global.Writer()
	gin.DefaultErrorWriter = log.Global.Writer()
}

func (k *MealKernel) initPermissionsRepo() {
	k.permissionsRepo = &auth.PermissionsRepo{Db: k.dbAccess}
}

func (k *MealKernel) initJWTGenerator() {
	k.jwtGen = &auth.JWTGenerator{Repo: k.permissionsRepo}
}

func (k *MealKernel) initAuthValidator() {
	k.authValidator = &api.AuthMiddleware{JwtGenerator: *k.jwtGen}
}

func (k *MealKernel) initAuthService() kernel.ModuleIntf {
	authServer := &api.AuthServer{UserRepo: &auth.FirebirdUserRepo{Database: k.dbAccess},
		PermRepo:     k.permissionsRepo,
		JWTGenerator: k.jwtGen}
	return ModuleAuth{Engine: k.serverEngine, AuthServer: authServer, PermRepo: k.permissionsRepo}
}

func (k *MealKernel) initMealRepo() {
	k.mealsRepo = &meal.FirebirdRepoAccess{Database: k.dbAccess}
}

func (k *MealKernel) initMealsModule() kernel.ModuleIntf {
	mealsAPI := &api.MealsAPI{
		Repo: k.mealsRepo,
	}
	return &ModuleMeals{Repo: k.mealsRepo, Engine: k.serverEngine, AuthValidator: k.authValidator, PermRepo: k.permissionsRepo, MethodExposer: mealsAPI}
}

func (k *MealKernel) initCategoriesModule() kernel.ModuleIntf {
	categoriesAPI := &api.CategoriesAPI{
		Repo: k.mealsRepo,
	}
	return &ModuleCategories{Repo: k.mealsRepo, Engine: k.serverEngine, AuthValidator: k.authValidator, PermRepo: k.permissionsRepo, MethodExposer: categoriesAPI}
}

func (k *MealKernel) initLooseProductsInDayModule() kernel.ModuleIntf {
	looseProductsInDayAPI := &api.LooseProductsInDayAPI{
		Repo: k.mealsRepo,
	}
	return &ModuleLooseProductsInDay{Repo: k.mealsRepo, Engine: k.serverEngine, AuthValidator: k.authValidator, PermRepo: k.permissionsRepo, MethodExposer: looseProductsInDayAPI}
}

func (k *MealKernel) initMealsInDayModule() kernel.ModuleIntf {
	mealsInDayAPI := &api.MealsInDayAPI{
		Repo: k.mealsRepo,
	}
	return &ModuleMealsInDay{Repo: k.mealsRepo, Engine: k.serverEngine, AuthValidator: k.authValidator, PermRepo: k.permissionsRepo, MethodExposer: mealsInDayAPI}
}

func (k *MealKernel) initOptimizeMealModule() kernel.ModuleIntf {
	maxTokens, err := strconv.Atoi(os.Getenv("OPENROUTER_MAX_TOKENS"))
	if err != nil {
		log.Global.Panicf("Error converting OPENROUTER_MAX_TOKENS to int: %v", err)
	}
	client := AIClient.OpenRouterClient{
		ApiKey:    os.Getenv("OPENROUTER_API_KEY"),
		Model:     os.Getenv("OPENROUTER_MODEL"),
		MaxTokens: maxTokens,
	}
	aiOptimizer := mealOptimizer.Optimizer{AIClient: &client}

	aiOptimizerAPI := &api.AIOptimizerAPI{
		Repo:     k.mealsRepo,
		AIClient: aiOptimizer,
	}
	return &ModuleOptimizeMeal{Engine: k.serverEngine, AuthValidator: k.authValidator, PermRepo: k.permissionsRepo, MethodExposer: aiOptimizerAPI}
}

func (k *MealKernel) initProductsModule() kernel.ModuleIntf {
	productsAPI := &api.ProductsAPI{
		Repo: k.mealsRepo,
	}
	return &ModuleProducts{Repo: k.mealsRepo, Engine: k.serverEngine, AuthValidator: k.authValidator, PermRepo: k.permissionsRepo, MethodExposer: productsAPI}
}

func (k *MealKernel) Init() {
	err := godotenv.Load()
	if err != nil {
		log.Global.Panicf("Error loading .env file: %v", err)
	}
	k.initLogger()
	k.initServer()
	k.initDB()

	k.initPermissionsRepo()
	k.initJWTGenerator()
	k.initAuthValidator()

	k.initMealRepo()
	k.RegisterModule(k.initAuthService())
	k.RegisterModule(k.initMealsModule())
	k.RegisterModule(k.initProductsModule())
	k.RegisterModule(k.initCategoriesModule())
	k.RegisterModule(k.initLooseProductsInDayModule())
	k.RegisterModule(k.initMealsInDayModule())
	k.RegisterModule(k.initOptimizeMealModule())
}

func (k *MealKernel) Run() {
	k.Kernel.Run()
	k.serverEngine.Run(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}
