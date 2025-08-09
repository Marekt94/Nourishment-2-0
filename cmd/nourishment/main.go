package main

/*
TODO:DODAĆ TESTY W POSTMANIE:
 - update potrawy - czy updatują sie produkty?
TODO: dodać crud, dto, repo dla kategorii
TODO: dodać api dla optymalizacji potraw - kalorycznosc zmienna
TODO: dodać crud, dto, repo dla produktów wolnych w dniu
TODO: zwracac w responsie potraw w dniu całkowite makro
TODO: dodać endpoint do wydruku, niech przesyła pdfa
TODO: jwt, autoryzacja uwierzytelnianie
TODO: stworzyc gotowego maina, zeby byl wystawialny w prosty sposób
*/

import (
	"nourishment_20/internal/AIClient"
	"nourishment_20/internal/api"
	"nourishment_20/internal/database"
	log "nourishment_20/internal/logging"
	"nourishment_20/internal/mealOptimizer"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// [AI REFACTOR] Tworzenie i uruchamianie serwera HTTP na porcie 8080
func StartMealServer() {
	r := gin.Default()

	r.GET("/meals", api.GetMeals)
	r.GET("/meals/:id", api.GetMeal)
	r.POST("/meals", api.CreateMeal)
	r.PUT("/meals", api.UpdateMeal)
	r.DELETE("/meals/:id", api.DeleteMeal)

	r.GET("/mealsinday", api.GetMealsInDay)
	r.GET("/mealsinday/:id", api.GetMealInDay)
	r.POST("/mealsinday", api.CreateMealInDay)
	r.PUT("/mealsinday", api.UpdateMealInDay)
	r.DELETE("/mealsinday/:id", api.DeleteMealInDay)

	r.GET("/products", api.GetProducts)
	r.GET("/products/:id", api.GetProduct)
	r.POST("/products", api.CreateProduct)
	r.PUT("/products", api.UpdateProduct)
	r.DELETE("/products/:id", api.DeleteProduct)

	r.GET("/looseproductsinday", api.GetLooseProductsInDay)
	r.GET("/looseproductsinday/:id", api.GetLooseProductInDay)
	r.POST("/looseproductsinday", api.CreateLooseProductInDay)
	r.PUT("/looseproductsinday", api.UpdateLooseProductInDay)
	r.DELETE("/looseproductsinday/:id", api.DeleteLooseProductInDay)

	r.Run(":8080") // [AI REFACTOR] nasłuch na porcie 8080
}

func main() {
	log.SetGlobalLogger(log.NewZerologLogger())
	// StartMealServer() // [AI REFACTOR] uruchom serwer
	err := godotenv.Load()
	if err != nil {
		log.Global.Panicf("Error loading .env file: %v", err)
	}
	maxTokensStr := os.Getenv("OPENROUTER_MAX_TOKENS")
	maxTokens, err := strconv.Atoi(maxTokensStr)
	if err != nil {
		log.Global.Panicf("Error converting OPEROUTER_MAX_TOKENS to int: %v", err)
	}
	client := AIClient.OpenRouterClient{ApiKey: os.Getenv("OPENROUTER_API_KEY"), Model: os.Getenv("OPENROUTER_MODEL"), MaxTokens: maxTokens}
	mealOptimizer := mealOptimizer.Optimizer{AIClient: &client}

	var conf database.DBConf
	conf.User = `sysdba`
	conf.Password = `masterkey`
	conf.Address = `localhost:3050`
	conf.PathOrName = `C:\Users\marek\Documents\nourishment_backup_db\NOURISHMENT.FDB`
	fDbEngine := database.FBDBEngine{BaseEngineIntf: &database.BaseEngine{}}
	engine := fDbEngine.Connect(&conf)
	var mealsRepo database.MealsRepo

	mealsRepo = &database.FirebirdRepoAccess{DbEngine: engine}

	meal := mealsRepo.GetMeal(15)
	_, err = mealOptimizer.OptimizeMeal(&meal)
	if err != nil {
		log.Global.Panicf("Error optimizing meal: %v", err)
	}
}
