package main

import (
	"nourishment_20/internal/api"
	log "nourishment_20/internal/logging"

	"github.com/gin-gonic/gin"
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

	r.Run(":8080") // [AI REFACTOR] nasłuch na porcie 8080
}

func main() {
	log.SetGlobalLogger(log.NewZerologLogger())
	StartMealServer() // [AI REFACTOR] uruchom serwer
}