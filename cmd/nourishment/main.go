package main

/*
TODO:DODAĆ TESTY W POSTMANIE:
 - update potrawy - czy updatują sie produkty?
DONE: dodać crud, dto, repo dla kategorii
TODO: dodać api dla optymalizacji potraw - kalorycznosc zmienna
DONE: dodać crud, dto, repo dla produktów wolnych w dniu
DONE: testy dla luźnych produktów w dniu
TODO: zwracac w responsie potraw w dniu całkowite makro
TODO: dodać endpoint do wydruku, niech przesyła pdfa (albo w markdown)
TODO: jwt, autoryzacja uwierzytelnianie
TODO: stworzyc gotowego maina, zeby byl wystawialny w prosty sposób
*/

import (
	"fmt"
	"os"
	"strconv"

	"nourishment_20/internal/AIClient"
	"nourishment_20/internal/api"
	db "nourishment_20/internal/database"
	log "nourishment_20/internal/logging"
	"nourishment_20/internal/mealOptimizer"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	repo := &db.FirebirdRepoAccess{Database: database}

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

	r := gin.Default()

	// Podpięcie zerologa do gin
	gin.DefaultWriter = log.Global.Writer()
	gin.DefaultErrorWriter = log.Global.Writer()

	r.GET("/meals", mealServer.GetMeals)
	r.GET("/meals/:id", mealServer.GetMeal)
	r.POST("/meals", mealServer.CreateMeal)
	r.PUT("/meals", mealServer.UpdateMeal)
	r.DELETE("/meals/:id", mealServer.DeleteMeal)

	r.GET("/mealsinday", mealServer.GetMealsInDay)
	r.GET("/mealsinday/:id", mealServer.GetMealInDay)
	r.POST("/mealsinday", mealServer.CreateMealInDay)
	r.PUT("/mealsinday", mealServer.UpdateMealInDay)
	r.DELETE("/mealsinday/:id", mealServer.DeleteMealInDay)

	r.GET("/products", mealServer.GetProducts)
	r.GET("/products/:id", mealServer.GetProduct)
	r.POST("/products", mealServer.CreateProduct)
	r.PUT("/products", mealServer.UpdateProduct)
	r.DELETE("/products/:id", mealServer.DeleteProduct)

	r.GET("/looseproductsinday", mealServer.GetLooseProductsInDay)
	r.GET("/looseproductsinday/:id", mealServer.GetLooseProductInDay)
	r.POST("/looseproductsinday", mealServer.CreateLooseProductInDay)
	r.PUT("/looseproductsinday", mealServer.UpdateLooseProductInDay)
	r.DELETE("/looseproductsinday/:id", mealServer.DeleteLooseProductInDay)

	// Categories endpoints
	r.GET("/categories", mealServer.GetCategories)
	r.GET("/categories/:id", mealServer.GetCategory)
	r.POST("/categories", mealServer.CreateCategory)
	r.PUT("/categories", mealServer.UpdateCategory)
	r.DELETE("/categories/:id", mealServer.DeleteCategory)

	r.POST("/optimizemeal", mealServer.OptimizeMeal)
	r.POST("/optimizemeal/:id", mealServer.OptimizeMealFromRepo)

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
