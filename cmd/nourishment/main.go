package main

import (
	"fmt"
	"nourishment_20/internal/api"
	database "nourishment_20/internal/data"

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
	var conf database.DBConf
	conf.User = `sysdba`
	conf.Password = `masterkey`
	conf.Address = `localhost:3050`
	conf.PathOrName = `C:\Users\marek\Documents\nourishment_backup_db\NOURISHMENT.FDB`

	fDbEngine := database.FBDBEngine{BaseEngineIntf: &database.BaseEngine{}}

	engine := fDbEngine.Connect(&conf)

	var mealsRepo database.MealsRepo = &database.FirebirdRepoAccess{DbEngine: engine}
	fmt.Println(mealsRepo.GetMeal(88))

	StartMealServer() // [AI REFACTOR] uruchom serwer
}