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

	repo := &meal.FirebirdRepoAccess{Database: database}

	// Utworzenie instancji AI Client
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

	// Utworzenie instancji MealServer
	mealServer := &api.MealServer{
		Repo:     repo,
		AIClient: aiOptimizer,
	}

	permissionRepo := auth.PermissionsRepo{Db: database}
	jwtGen := auth.JWTGenerator{Repo: &permissionRepo}
	authValidator := api.AuthMiddleware{JwtGenerator: jwtGen}

	authServer := &api.AuthServer{UserRepo: &auth.FirebirdUserRepo{Database: database},
		PermRepo:     &permissionRepo,
		JWTGenerator: &jwtGen}

	r := gin.Default()

	// Podpięcie zerologa do gin
	gin.DefaultWriter = log.Global.Writer()
	gin.DefaultErrorWriter = log.Global.Writer()

	// // Rejestracja uprawnień w systemie
	// permissionRepo.RegisterPermissions(api.RESOURCE_MEALS, []string{"read", "write"})
	// permissionRepo.RegisterPermissions(api.RESOURCE_MEALSINDAY, []string{"read", "write"})
	// permissionRepo.RegisterPermissions(api.RESOURCE_PRODUCTS, []string{"read", "write"})
	// permissionRepo.RegisterPermissions(api.RESOURCE_LOOSEPRODUCTSINDAY, []string{"read", "write"})
	// permissionRepo.RegisterPermissions(api.RESOURCE_CATEGORIES, []string{"read", "write"})
	// permissionRepo.RegisterPermissions(api.RESOURCE_OPTIMIZEMEAL, []string{"write"})

	// // Zarejestruj wszystkie uprawnienia dla użytkownika ADMIN
	// adminUser := "ADMIN"
	// permissionRepo.RegisterUserPermission(adminUser, api.RESOURCE_MEALS, "read")
	// permissionRepo.RegisterUserPermission(adminUser, api.RESOURCE_MEALS, "write")
	// permissionRepo.RegisterUserPermission(adminUser, api.RESOURCE_MEALSINDAY, "read")
	// permissionRepo.RegisterUserPermission(adminUser, api.RESOURCE_MEALSINDAY, "write")
	// permissionRepo.RegisterUserPermission(adminUser, api.RESOURCE_PRODUCTS, "read")
	// permissionRepo.RegisterUserPermission(adminUser, api.RESOURCE_PRODUCTS, "write")
	// permissionRepo.RegisterUserPermission(adminUser, api.RESOURCE_LOOSEPRODUCTSINDAY, "read")
	// permissionRepo.RegisterUserPermission(adminUser, api.RESOURCE_LOOSEPRODUCTSINDAY, "write")
	// permissionRepo.RegisterUserPermission(adminUser, api.RESOURCE_CATEGORIES, "read")
	// permissionRepo.RegisterUserPermission(adminUser, api.RESOURCE_CATEGORIES, "write")
	// permissionRepo.RegisterUserPermission(adminUser, api.RESOURCE_OPTIMIZEMEAL, "write")

	// Zarejestruj uprawnienia do odczytu dla użytkownika READER
	// readerUser := "READER"
	// permissionRepo.RegisterUserPermission(readerUser, api.RESOURCE_MEALS, "read")
	// permissionRepo.RegisterUserPermission(readerUser, api.RESOURCE_MEALSINDAY, "read")
	// permissionRepo.RegisterUserPermission(readerUser, api.RESOURCE_PRODUCTS, "read")
	// permissionRepo.RegisterUserPermission(readerUser, api.RESOURCE_LOOSEPRODUCTSINDAY, "read")
	// permissionRepo.RegisterUserPermission(readerUser, api.RESOURCE_CATEGORIES, "read")

	// TODO: Dodać moduły które jako parametr przekazywałyby gin i w tych podułach byłyby rejestrowane endpointy
	r.POST(api.PATH_LOGIN, authServer.GenerateToken)

	r.GET(api.PATH_MEALS, authValidator.Middleware, mealServer.GetMeals)
	r.GET(api.PATH_MEALS_WITH_ID, authValidator.Middleware, mealServer.GetMeal)
	r.POST(api.PATH_MEALS, authValidator.Middleware, mealServer.CreateMeal)
	r.PUT(api.PATH_MEALS, authValidator.Middleware, mealServer.UpdateMeal)
	r.DELETE(api.PATH_MEALS_WITH_ID, authValidator.Middleware, mealServer.DeleteMeal)

	r.GET(api.PATH_MEALSINDAY, authValidator.Middleware, mealServer.GetMealsInDay)
	r.GET(api.PATH_MEALSINDAY_WITH_ID, authValidator.Middleware, mealServer.GetMealInDay)
	r.POST(api.PATH_MEALSINDAY, authValidator.Middleware, mealServer.CreateMealInDay)
	r.PUT(api.PATH_MEALSINDAY, authValidator.Middleware, mealServer.UpdateMealInDay)
	r.DELETE(api.PATH_MEALSINDAY_WITH_ID, authValidator.Middleware, mealServer.DeleteMealInDay)

	r.GET(api.PATH_PRODUCTS, authValidator.Middleware, mealServer.GetProducts)
	r.GET(api.PATH_PRODUCTS_WITH_ID, authValidator.Middleware, mealServer.GetProduct)
	r.POST(api.PATH_PRODUCTS, authValidator.Middleware, mealServer.CreateProduct)
	r.PUT(api.PATH_PRODUCTS, authValidator.Middleware, mealServer.UpdateProduct)
	r.DELETE(api.PATH_PRODUCTS_WITH_ID, authValidator.Middleware, mealServer.DeleteProduct)

	r.GET(api.PATH_LOOSEPRODUCTSINDAY, authValidator.Middleware, mealServer.GetLooseProductsInDay)
	r.GET(api.PATH_LOOSEPRODUCTSINDAY_WITH_ID, authValidator.Middleware, mealServer.GetLooseProductInDay)
	r.POST(api.PATH_LOOSEPRODUCTSINDAY, authValidator.Middleware, mealServer.CreateLooseProductInDay)
	r.PUT(api.PATH_LOOSEPRODUCTSINDAY, authValidator.Middleware, mealServer.UpdateLooseProductInDay)
	r.DELETE(api.PATH_LOOSEPRODUCTSINDAY_WITH_ID, authValidator.Middleware, mealServer.DeleteLooseProductInDay)

	// Categories endpoints
	r.GET(api.PATH_CATEGORIES, authValidator.Middleware, mealServer.GetCategories)
	r.GET(api.PATH_CATEGORIES_WITH_ID, authValidator.Middleware, mealServer.GetCategory)
	r.POST(api.PATH_CATEGORIES, authValidator.Middleware, mealServer.CreateCategory)
	r.PUT(api.PATH_CATEGORIES, authValidator.Middleware, mealServer.UpdateCategory)
	r.DELETE(api.PATH_CATEGORIES_WITH_ID, authValidator.Middleware, mealServer.DeleteCategory)

	r.POST(api.PATH_OPTIMIZEMEAL, authValidator.Middleware, mealServer.OptimizeMeal)
	r.POST(api.PATH_OPTIMIZEMEAL_WITH_ID, authValidator.Middleware, mealServer.OptimizeMealFromRepo)

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
