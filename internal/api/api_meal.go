package api

import (
	"net/http"
	"strconv"

	meal "nourishment_20/internal/mealDomain"
	"nourishment_20/internal/mealOptimizer"

	"github.com/gin-gonic/gin"
)

type MealServer struct {
	Repo     meal.MealsRepoIntf
	AIClient mealOptimizer.Optimizer
}

// DONE do refactoringu - trzeba wyciagnac pliki konfiguracyjne wyzej, repo musi byc orzekazywane jako parametr, byc moze do obiektu
func (ms *MealServer) OptimizeMeal(c *gin.Context) {
	kcalStr := c.Query("kcal")
	kcal, err := strconv.ParseFloat(kcalStr, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var meal meal.Meal
	if err := c.ShouldBindJSON(&meal); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := ms.AIClient.OptimizeMeal(&meal, kcal)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, res)
}

func (ms *MealServer) OptimizeMealFromRepo(c *gin.Context) {
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
	meal := ms.Repo.GetMeal(id)
	if meal.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	res, err := ms.AIClient.OptimizeMeal(&meal, kcal)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, res)
}

func (ms *MealServer) GetMeal(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err == nil {
		c.IndentedJSON(http.StatusOK, ms.Repo.GetMeal(id))
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (ms *MealServer) GetMeals(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, ms.Repo.GetMeals())
}

func (ms *MealServer) CreateMeal(c *gin.Context) {
	var m meal.Meal // [AI REFACTOR]
	if err := c.ShouldBindJSON(&m); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := ms.Repo.CreateMeal(&m)
	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func (ms *MealServer) UpdateMeal(c *gin.Context) {
	var m meal.Meal // [AI REFACTOR]
	if err := c.ShouldBindJSON(&m); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ms.Repo.UpdateMeal(&m)
	c.Status(http.StatusOK)
}

func (ms *MealServer) DeleteMeal(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok := ms.Repo.DeleteMeal(id)
	if ok {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}

// CRUD dla MealsInDay
func (ms *MealServer) GetMealInDay(c *gin.Context) {
	repo := ms.Repo.(meal.MealsInDayRepo)
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

func (ms *MealServer) GetMealsInDay(c *gin.Context) {
	repo := ms.Repo.(meal.MealsInDayRepo)
	meals := repo.GetMealsInDays()
	c.IndentedJSON(http.StatusOK, meals)
}

func (ms *MealServer) CreateMealInDay(c *gin.Context) {
	repo := ms.Repo.(meal.MealsInDayRepo)
	var m meal.MealInDay
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

func (ms *MealServer) UpdateMealInDay(c *gin.Context) {
	repo := ms.Repo.(meal.MealsInDayRepo)
	var m meal.MealInDay
	if err := c.ShouldBindJSON(&m); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo.UpdateMealsInDay(&m)
	c.Status(http.StatusOK)
}

func (ms *MealServer) DeleteMealInDay(c *gin.Context) {
	repo := ms.Repo.(meal.MealsInDayRepo)
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
func (ms *MealServer) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo := ms.Repo.(meal.ProductsRepo)
	product := repo.GetProduct(id)
	if product.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, product)
}

func (ms *MealServer) GetProducts(c *gin.Context) {
	repo := ms.Repo.(meal.ProductsRepo)
	products := repo.GetProducts()
	c.IndentedJSON(http.StatusOK, products)
}

func (ms *MealServer) CreateProduct(c *gin.Context) {
	repo := ms.Repo.(meal.ProductsRepo)
	var p meal.Product
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

func (ms *MealServer) UpdateProduct(c *gin.Context) {
	repo := ms.Repo.(meal.ProductsRepo)
	var p meal.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo.UpdateProduct(&p)
	c.Status(http.StatusOK)
}

func (ms *MealServer) DeleteProduct(c *gin.Context) {
	repo := ms.Repo.(meal.ProductsRepo)
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
func (ms *MealServer) GetLooseProductInDay(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo := ms.Repo.(meal.LooseProductsInDayRepo)
	product := repo.GetLooseProductInDay(id)
	if product.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, product)
}

func (ms *MealServer) GetLooseProductsInDay(c *gin.Context) {
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
	repo := ms.Repo.(meal.LooseProductsInDayRepo)
	products := repo.GetLooseProductsInDay(dayId)
	c.IndentedJSON(http.StatusOK, products)
}

func (ms *MealServer) CreateLooseProductInDay(c *gin.Context) {
	repo := ms.Repo.(meal.LooseProductsInDayRepo)
	var p meal.LooseProductInDay
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

func (ms *MealServer) UpdateLooseProductInDay(c *gin.Context) {
	repo := ms.Repo.(meal.LooseProductsInDayRepo)
	var p meal.LooseProductInDay
	if err := c.ShouldBindJSON(&p); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo.UpdateLooseProductInDay(&p)
	c.Status(http.StatusOK)
}

func (ms *MealServer) DeleteLooseProductInDay(c *gin.Context) {
	repo := ms.Repo.(meal.LooseProductsInDayRepo)
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
func (ms *MealServer) GetCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo := ms.Repo.(meal.CategoriesRepo)
	cat := repo.GetCategory(id)
	if cat.Id == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	c.IndentedJSON(http.StatusOK, cat)
}

func (ms *MealServer) GetCategories(c *gin.Context) {
	repo := ms.Repo.(meal.CategoriesRepo)
	cats := repo.GetCategories()
	c.IndentedJSON(http.StatusOK, cats)
}

func (ms *MealServer) CreateCategory(c *gin.Context) {
	repo := ms.Repo.(meal.CategoriesRepo)
	var cat meal.Category
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

func (ms *MealServer) UpdateCategory(c *gin.Context) {
	repo := ms.Repo.(meal.CategoriesRepo)
	var cat meal.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repo.UpdateCategory(&cat)
	c.Status(http.StatusOK)
}

func (ms *MealServer) DeleteCategory(c *gin.Context) {
	repo := ms.Repo.(meal.CategoriesRepo)
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
