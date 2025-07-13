package api

import (
	"net/http"
	"strconv"

	db "nourishment_20/internal/database"

	"github.com/gin-gonic/gin"
)

func getRepo() db.MealsRepo { // [AI REFACTOR]
	conf := db.DBConf{
		User:       "sysdba",
		Password:   "masterkey",
		Address:    "localhost:3050",
		PathOrName: "C:/Users/marek/Documents/nourishment_backup_db/NOURISHMENT.FDB",
	}
	fDbEngine := db.FBDBEngine{BaseEngineIntf: &db.BaseEngine{}}
	engine := fDbEngine.Connect(&conf)
	return &db.FirebirdRepoAccess{DbEngine: engine}
}

func GetMeal(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err == nil {
		repo := getRepo() // [AI REFACTOR]
		c.IndentedJSON(http.StatusOK, repo.GetMeal(id))
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func GetMeals(c *gin.Context) {
	repo := getRepo() // [AI REFACTOR]
	c.IndentedJSON(http.StatusOK, repo.GetMeals())
}

func CreateMeal(c *gin.Context) {
	var m db.Meal // [AI REFACTOR]
	if err := c.ShouldBindJSON(&m); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo := getRepo() // [AI REFACTOR]
	id := repo.CreateMeal(&m)
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func UpdateMeal(c *gin.Context) {
	var m db.Meal // [AI REFACTOR]
	if err := c.ShouldBindJSON(&m); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo := getRepo() // [AI REFACTOR]
	repo.UpdateMeal(&m)
	c.Status(http.StatusOK)
}

func DeleteMeal(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo := getRepo() // [AI REFACTOR]
	ok := repo.DeleteMeal(id)
	if ok {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}

// CRUD dla MealsInDay
func GetMealInDay(c *gin.Context) {
	repo := getRepo().(db.MealsInDayRepo)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	meal := repo.GetMealsInDay(id)
	if meal.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, meal)
}

func GetMealsInDay(c *gin.Context) {
	repo := getRepo().(db.MealsInDayRepo)
	meals := repo.GetMealsInDays()
	c.IndentedJSON(http.StatusOK, meals)
}

func CreateMealInDay(c *gin.Context) {
	repo := getRepo().(db.MealsInDayRepo)
	var m db.MealInDay
	if err := c.ShouldBindJSON(&m); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := repo.CreateMealsInDay(&m)
	if id <= 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "CreateMealInDay failed"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func UpdateMealInDay(c *gin.Context) {
	repo := getRepo().(db.MealsInDayRepo)
	var m db.MealInDay
	if err := c.ShouldBindJSON(&m); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo.UpdateMealsInDay(&m)
	c.Status(http.StatusOK)
}

func DeleteMealInDay(c *gin.Context) {
	repo := getRepo().(db.MealsInDayRepo)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok := repo.DeleteMealsInDay(id)
	if ok {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}
