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

// CRUD dla Products
func GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo := getRepo().(db.ProductsRepo)
	product := repo.GetProduct(id)
	if product.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, product)
}

func GetProducts(c *gin.Context) {
	repo := getRepo().(db.ProductsRepo)
	products := repo.GetProducts()
	c.IndentedJSON(http.StatusOK, products)
}

func CreateProduct(c *gin.Context) {
	repo := getRepo().(db.ProductsRepo)
	var p db.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := repo.CreateProduct(&p)
	if id <= 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "CreateProduct failed"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func UpdateProduct(c *gin.Context) {
	repo := getRepo().(db.ProductsRepo)
	var p db.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo.UpdateProduct(&p)
	c.Status(http.StatusOK)
}

func DeleteProduct(c *gin.Context) {
	repo := getRepo().(db.ProductsRepo)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok := repo.DeleteProduct(id)
	if ok {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}

// CRUD dla LooseProductInDay
func GetLooseProductInDay(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo := getRepo().(db.LooseProductsInDayRepo)
	product := repo.GetLooseProductInDay(id)
	if product.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, product)
}

func GetLooseProductsInDay(c *gin.Context) {
	dayId, err := strconv.Atoi(c.Param("dayId"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo := getRepo().(db.LooseProductsInDayRepo)
	products := repo.GetLooseProductsInDay(dayId)
	c.IndentedJSON(http.StatusOK, products)
}

func CreateLooseProductInDay(c *gin.Context) {
	repo := getRepo().(db.LooseProductsInDayRepo)
	var p db.LooseProductInDay
	if err := c.ShouldBindJSON(&p); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := repo.CreateLooseProductInDay(&p)
	if id <= 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "CreateLooseProductInDay failed"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func UpdateLooseProductInDay(c *gin.Context) {
	repo := getRepo().(db.LooseProductsInDayRepo)
	var p db.LooseProductInDay
	if err := c.ShouldBindJSON(&p); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo.UpdateLooseProductInDay(&p)
	c.Status(http.StatusOK)
}

func DeleteLooseProductInDay(c *gin.Context) {
	repo := getRepo().(db.LooseProductsInDayRepo)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok := repo.DeleteLooseProductInDay(id)
	if ok {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}
