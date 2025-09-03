package main

/*
TODO:DODAĆ TESTY W POSTMANIE:
 - update potrawy - czy updatują sie produkty?
DONE: dodać crud, dto, repo dla kategorii
DONE: dodać api dla optymalizacji potraw - kalorycznosc zmienna
DONE: dodać crud, dto, repo dla produktów wolnych w dniu
DONE: testy dla luźnych produktów w dniu
TODO: zwracac w responsie potraw w dniu całkowite makro
TODO: dodać endpoint do wydruku, niech przesyła pdfa (albo w markdown)
DONE: jwt, autoryzacja uwierzytelnianie
TODO: dodać weryfikację uprawnień po danych w bazie danych
TODO: stworzyc gotowego maina, zeby byl wystawialny w prosty sposób
*/

import (
	"fmt"
	"os"
	"strconv"

	"nourishment_20/internal/AIClient"
	"nourishment_20/internal/api"
	"nourishment_20/internal/auth"
	log "nourishment_20/internal/logging"
	meal "nourishment_20/internal/mealDomain"
	"nourishment_20/internal/mealOptimizer"
	"nourishment_20/internal/modules"
	"nourishment_20/kernel"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	db "nourishment_20/internal/database"
)

// [AI REFACTOR] Tworzenie i uruchamianie serwera HTTP na porcie 8080
func StartMealServer() {
	// Utworzenie instancji repo
	DbEngine := db.FBDBEngine{BaseEngineIntf: &db.BaseEngine{}}

	conf := db.DBConf{
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		Address:    os.Getenv("DB_ADDRESS"),
		PathOrName: os.Getenv("DB_NAME"),
	}
	database := DbEngine.Connect(&conf)

	permissionRepo := auth.PermissionsRepo{Db: database}
	jwtGen := auth.JWTGenerator{Repo: &permissionRepo}
	authValidator := api.AuthMiddleware{JwtGenerator: jwtGen}

	//TODO: poprawić - powinno byc dla firebirdrepoacces a nie osobny struct
	authServer := &api.AuthServer{UserRepo: &auth.FirebirdUserRepo{Database: database},
		PermRepo:     &permissionRepo,
		JWTGenerator: &jwtGen}

	r := gin.Default()

	// Podpięcie zerologa do gin
	gin.DefaultWriter = log.Global.Writer()
	gin.DefaultErrorWriter = log.Global.Writer()

	kernel := kernel.NewKernel()
	kernel.RegisterModule(modules.ModuleAuth{Engine: r, AuthServer: authServer, PermRepo: &permissionRepo})

	// Wspólne repozytorium dla wszystkich modułów
	repo := meal.FirebirdRepoAccess{Database: database}

	// AI Client dla optymalizacji
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

	// API dla każdego modułu
	mealsAPI := &api.MealsAPI{
		Repo: &repo,
	}
	mealsInDayAPI := &api.MealsInDayAPI{
		Repo: &repo,
	}
	productsAPI := &api.ProductsAPI{
		Repo: &repo,
	}
	looseProductsInDayAPI := &api.LooseProductsInDayAPI{
		Repo: &repo,
	}
	categoriesAPI := &api.CategoriesAPI{
		Repo: &repo,
	}
	aiOptimizerAPI := &api.AIOptimizerAPI{
		Repo:     &repo,
		AIClient: aiOptimizer,
	}

	kernel.RegisterModule(&modules.ModuleMeals{Repo: &repo, Engine: r, AuthValidator: &authValidator, PermRepo: &permissionRepo, MethodExposer: mealsAPI})
	kernel.RegisterModule(&modules.ModuleMealsInDay{Repo: &repo, Engine: r, AuthValidator: &authValidator, PermRepo: &permissionRepo, MethodExposer: mealsInDayAPI})
	kernel.RegisterModule(&modules.ModuleProducts{Repo: &repo, Engine: r, AuthValidator: &authValidator, PermRepo: &permissionRepo, MethodExposer: productsAPI})
	kernel.RegisterModule(&modules.ModuleLooseProductsInDay{Repo: &repo, Engine: r, AuthValidator: &authValidator, PermRepo: &permissionRepo, MethodExposer: looseProductsInDayAPI})
	kernel.RegisterModule(&modules.ModuleCategories{Repo: &repo, Engine: r, AuthValidator: &authValidator, PermRepo: &permissionRepo, MethodExposer: categoriesAPI})
	kernel.RegisterModule(&modules.ModuleOptimizeMeal{Engine: r, AuthValidator: &authValidator, PermRepo: &permissionRepo, MethodExposer: aiOptimizerAPI})

	kernel.Run()

	r.Run(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}

func main() {
	log.SetGlobalLogger(log.NewZerologLogger())
	err := godotenv.Load()
	if err != nil {
		log.Global.Panicf("Error loading .env file: %v", err)
	}
	StartMealServer() // [AI REFACTOR] uruchom serwer
}
