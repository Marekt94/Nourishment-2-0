package api

import (
	"net/http"
	"os"
	"strconv"

	"nourishment_20/internal/AIClient"
	db "nourishment_20/internal/database"
	"nourishment_20/internal/logging"
	"nourishment_20/internal/mealOptimizer"

	"github.com/gin-gonic/gin"
)

// TODO do refactoringu - trzeba wyciagnac pliki konfiguracyjne wyzej, repo musi byc orzekazywane jako parametr, byc moze do obiektu
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

// TODO (ai opmitimizer) do refactoringu - trzeba wyciagnac pliki konfiguracyjne wyzej, repo musi byc orzekazywane jako parametr, byc moze do obiektu
func getAIClient() *mealOptimizer.Optimizer { // [AI REFACTOR]
	maxTokens, err := strconv.Atoi(os.Getenv("OPENROUTER_MAX_TOKENS"))
	if err != nil {
		logging.Global.Panicf("Error converting OPENROUTER_MAX_TOKENS to int: %v", err)
	}
	client := AIClient.OpenRouterClient{
		ApiKey:    os.Getenv("OPENROUTER_API_KEY"),
		Model:     os.Getenv("OPENROUTER_MODEL"),
		MaxTokens: maxTokens,
	}
	res := mealOptimizer.Optimizer{AIClient: &client}
	return &res
}

func OptimizeMeal(c *gin.Context) {
	kcalStr := c.Query("kcal")
	kcal, err := strconv.ParseFloat(kcalStr, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	optimizer := getAIClient() // [AI REFACTOR]
	var meal db.Meal
	if err := c.ShouldBindJSON(&meal); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := optimizer.OptimizeMeal(&meal, kcal)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, res)
}

func OptimizeMealFromRepo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	kcal, err := strconv.ParseFloat(c.Query("kcal"), 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo := getRepo() // [AI REFACTOR]
	meal := repo.GetMeal(id)
	if meal.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	res, err := getAIClient().OptimizeMeal(&meal, kcal)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, res)
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
	dayIdStr := c.Query("dayId")
	if dayIdStr == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "missing query parameter: dayId"})
		return
	}
	dayId, err := strconv.Atoi(dayIdStr)
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

// CRUD dla Categories
func GetCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo := getRepo().(db.CategoriesRepo)
	cat := repo.GetCategory(id)
	if cat.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, cat)
}

func GetCategories(c *gin.Context) {
	repo := getRepo().(db.CategoriesRepo)
	cats := repo.GetCategories()
	c.IndentedJSON(http.StatusOK, cats)
}

func CreateCategory(c *gin.Context) {
	repo := getRepo().(db.CategoriesRepo)
	var cat db.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := repo.CreateCategory(&cat)
	if id <= 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "CreateCategory failed"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func UpdateCategory(c *gin.Context) {
	repo := getRepo().(db.CategoriesRepo)
	var cat db.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo.UpdateCategory(&cat)
	c.Status(http.StatusOK)
}

func DeleteCategory(c *gin.Context) {
	repo := getRepo().(db.CategoriesRepo)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok := repo.DeleteCategory(id)
	if ok {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}
