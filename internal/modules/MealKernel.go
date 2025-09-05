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
	log.Global.Infof("Creating new MealKernel instance...")
	k := MealKernel{}
	k.Kernel = kernel.NewKernel()
	log.Global.Infof("MealKernel instance created successfully")
	return &k
}

func (k *MealKernel) initDB() {
	log.Global.Infof("Initializing database connection...")
	DbEngine := db.FBDBEngine{BaseEngineIntf: &db.BaseEngine{}}

	conf := db.DBConf{
		User:       os.Getenv("DB_USER"),
		Password:   "***", // Don't log password
		Address:    os.Getenv("DB_ADDRESS"),
		PathOrName: os.Getenv("DB_NAME"),
	}
	log.Global.Infof("Connecting to database at %s with user %s", conf.Address, conf.User)

	conf.Password = os.Getenv("DB_PASSWORD") // Set real password
	k.dbAccess = DbEngine.Connect(&conf)
	log.Global.Infof("Database connection established successfully")
}

func (k *MealKernel) initLogger() {
	log.Global.Infof("Initializing logger...")
	log.SetGlobalLogger(log.NewZerologLogger())
	log.Global.Infof("Logger initialized successfully")
}

func (k *MealKernel) initServer() {
	log.Global.Infof("Initializing server engine...")
	k.serverEngine = gin.Default()
	gin.DefaultWriter = log.Global.Writer()
	gin.DefaultErrorWriter = log.Global.Writer()
	log.Global.Infof("Server engine initialized successfully")
}

func (k *MealKernel) initPermissionsRepo() {
	log.Global.Infof("Initializing permissions repository...")
	k.permissionsRepo = &auth.PermissionsRepo{Db: k.dbAccess}
	log.Global.Infof("Permissions repository initialized successfully")
}

func (k *MealKernel) initJWTGenerator() {
	log.Global.Infof("Initializing JWT generator...")
	k.jwtGen = &auth.JWTGenerator{Repo: k.permissionsRepo}
	log.Global.Infof("JWT generator initialized successfully")
}

func (k *MealKernel) initAuthValidator() {
	log.Global.Infof("Initializing auth validator...")
	k.authValidator = &api.AuthMiddleware{JwtGenerator: *k.jwtGen}
	log.Global.Infof("Auth validator initialized successfully")
}

func (k *MealKernel) initAuthService() kernel.ModuleIntf {
	authServer := &api.AuthServer{UserRepo: &auth.FirebirdUserRepo{Database: k.dbAccess},
		PermRepo:     k.permissionsRepo,
		JWTGenerator: k.jwtGen}
	return ModuleAuth{Engine: k.serverEngine, AuthServer: authServer, PermRepo: k.permissionsRepo}
}

func (k *MealKernel) initMealRepo() {
	log.Global.Infof("Initializing meal repository...")
	k.mealsRepo = &meal.FirebirdRepoAccess{Database: k.dbAccess}
	log.Global.Infof("Meal repository initialized successfully")
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
	log.Global.Infof("Initializing AI optimization module...")

	maxTokens, err := strconv.Atoi(os.Getenv("OPENROUTER_MAX_TOKENS"))
	if err != nil {
		log.Global.Panicf("Error converting OPENROUTER_MAX_TOKENS to int: %v", err)
	}
	log.Global.Infof("AI Client configuration - Model: %s, MaxTokens: %d",
		os.Getenv("OPENROUTER_MODEL"), maxTokens)

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
	log.Global.Infof("AI optimization module initialized successfully")
	return &ModuleOptimizeMeal{Engine: k.serverEngine, AuthValidator: k.authValidator, PermRepo: k.permissionsRepo, MethodExposer: aiOptimizerAPI}
}

func (k *MealKernel) initProductsModule() kernel.ModuleIntf {
	productsAPI := &api.ProductsAPI{
		Repo: k.mealsRepo,
	}
	return &ModuleProducts{Repo: k.mealsRepo, Engine: k.serverEngine, AuthValidator: k.authValidator, PermRepo: k.permissionsRepo, MethodExposer: productsAPI}
}

func (k *MealKernel) Init() {
	log.Global.Infof("Starting MealKernel initialization...")

	err := godotenv.Load()
	if err != nil {
		log.Global.Panicf("Error loading .env file: %v", err)
	}
	log.Global.Infof("Environment variables loaded successfully")

	k.initLogger()
	k.initServer()
	k.initDB()
	k.initPermissionsRepo()
	k.initJWTGenerator()
	k.initAuthValidator()
	k.initMealRepo()

	// Register all modules (each will log its own registration)
	k.RegisterModule(k.initAuthService())
	k.RegisterModule(k.initMealsModule())
	k.RegisterModule(k.initProductsModule())
	k.RegisterModule(k.initCategoriesModule())
	k.RegisterModule(k.initLooseProductsInDayModule())
	k.RegisterModule(k.initMealsInDayModule())
	k.RegisterModule(k.initOptimizeMealModule())

	log.Global.Infof("MealKernel initialization completed successfully")
}

func (k *MealKernel) Run() {
	log.Global.Infof("Starting MealKernel run sequence...")

	log.Global.Infof("Running kernel modules...")
	k.Kernel.Run()
	log.Global.Infof("All modules initialized successfully")

	port := os.Getenv("SERVER_PORT")
	log.Global.Infof("Starting HTTP server on port %s", port)
	k.serverEngine.Run(fmt.Sprintf(":%s", port))
}
