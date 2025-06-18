package database

import (
	"database/sql"
	"testing"
)

const PROD_NAME_1 = `test prod 1`
const PROD_NAME_2 = `test prod 2`

func initMealRepo() MealsRepo {
	var conf DBConf
	conf.User = `sysdba`
	conf.Password = `masterkey`
	conf.Address = `localhost:3050`
	conf.PathOrName = `C:\Users\marek\Documents\nourishment_backup_db\NOURISHMENT.FDB`

	fDbEngine := FBDBEngine{BaseEngineIntf: &BaseEngine{}}

	engine := fDbEngine.Connect(&conf)

	return &FirebirdRepoAccess{DbEngine: engine}
}

func createMeal() meal{
	repo := initMealRepo()
	products := []productInMeal{
		productInMeal{
			Product: product{
				Name: PROD_NAME_1,
				Proteins: 12,
			},
			Weight: 12,
		},
		productInMeal{
			Product: product{
				Name: PROD_NAME_2,
				Proteins: 15,
			},
			Weight: 12,
		},
	}
	meal := meal{Name: "test meal", Recipe: "test recipe", ProductsInMeal: products}
	id := repo.CreateMeal(&meal)
	meal.Id = int(id)
	return meal
}

func TestGetMeals(t *testing.T) {
	repo := initMealRepo()
	meals := repo.GetMeals();
	if len(meals) < 1 {
		t.Error(`Meals list is empty`)
	}		
	if len(meals[0].ProductsInMeal) == 0{
		t.Error(`products in meals not retrived`)
	}
}

func TestGetMeal(t *testing.T) {
	repo := initMealRepo()
	meal := repo.GetMeal(1);
	if len(meal.ProductsInMeal) < 1{
		t.Error(`products in meal not retrived`)
	}
}

func TestCreateMeal(t *testing.T){
	var meal meal
	if meal = createMeal(); meal.Id < 0 {
		t.Error(`meal creation error`)
	}

	repo := initMealRepo();
	mealRetrived := repo.GetMeal(meal.Id)
	if len(mealRetrived.ProductsInMeal) != 2 {
		t.Error(`products in meal not saved`)
	}
	if mealRetrived.ProductsInMeal[0].Id <= 0 {
		t.Error(`product in meal not saved (no id)`)
	}
}

func TestDeleteMeal(t *testing.T){
	meal := createMeal()
	repo := initMealRepo()

	if !repo.DeleteMeal(meal.Id){
		t.Errorf("meal with id: %d not deleted\n", meal.Id)
	}
}

func TestUpdateMeal(t *testing.T){
	prodNameAfterUpdate := `product after update name`
	meal := createMeal();
	repo := initMealRepo()

	meal.Name = `test name after update`
	meal.ProductsInMeal[0].Product.Name = prodNameAfterUpdate
	repo.UpdateMeal(&meal);

	res := repo.GetMeal(meal.Id)
	if res.Name != meal.Name {
		t.Errorf(`meal with id %d not updated`, res.Id)
	}
	var res2 bool = false;
	for _, prod := range meal.ProductsInMeal{
		res2 = res2 || (prod.Product.Name == prodNameAfterUpdate) 
	}
	if !res2 {
		t.Errorf(`product in meal not updated`)
	}
}

func TestUpdateMealWhenOneDeletedAndOneAdded(t *testing.T){
	const prodName3 = `test prod 3`
	meal := createMeal()
	repo := initMealRepo()
	newProd := productInMeal{Product: product{Name: prodName3, Proteins: 31}}
	meal.ProductsInMeal = append(meal.ProductsInMeal, newProd)
	idToDel := -1
	for i, prod := range meal.ProductsInMeal{
		if prod.Product.Name == PROD_NAME_2{
			idToDel = i
			break
		}
	}
	if idToDel == -1{
		t.Error(`init data error`)
	}
	meal.ProductsInMeal = append(meal.ProductsInMeal[:idToDel], meal.ProductsInMeal[idToDel+1:]...)
	t.Log(meal.ProductsInMeal)
	
	repo.UpdateMeal(&meal)
	resMeal := repo.GetMeal(meal.Id)

	res := false
	for _, prod := range resMeal.ProductsInMeal{
		res = res || (prod.Product.Name == prodName3)
	}
	if !res {
		t.Error(`product not added`)
	}

	res = true
	for _, prod := range resMeal.ProductsInMeal{
		res = res && (prod.Product.Name != PROD_NAME_2)
	}
	if !res {
		t.Error(`product not deleted`)
	}

	if len(resMeal.ProductsInMeal) != 2{
		t.Error(`product number mismatch`)
	}
}

func TestUpdateProductInMeal(t *testing.T){
	repo := initMealRepo()
	meal := createMeal()

	meal.ProductsInMeal[0].Product.Name = `product name after update`
	repo.UpdateMeal(&meal)

	res := repo.GetMeal(meal.Id)
	resTemp := false
	for _, prodInMeal := range res.ProductsInMeal{
		resTemp = prodInMeal.Product.Name == `product name after update`
		if resTemp {
			break
		}
	}
	if !resTemp{
		t.Error(`product not updated`)
	}
}

func TestConvertToMealWhenMealsDBEmpty(t *testing.T){
	mealsDB := []MealDb{};
	meals := ConvertToMeals(mealsDB)
	if meals != nil {
		t.Error(`meal should be empty`)
	}
}

func TestConvertToMealWhenMealsDBWithOneMeal(t *testing.T){
	mealDB := []MealDb{MealDb{Name: sql.NullString{`meal_1`, true},
							  ProductInMeal: productInMealDb{Id: sql.NullInt64{412, true},
															 Weight: sql.NullFloat64{100, true},
															 Product: productDb{Name: sql.NullString{`test_prod_1`, true}}}}} 
	meals := ConvertToMeals(mealDB)
	if len(meals) != 1 {
		t.Error(`too many meals`)
	}
	if len(meals[0].ProductsInMeal) != 1 {
		t.Error(`too many products in meal`)
	}
}

func TestConvertToMealWhenNoProducts(t *testing.T){
	mealDB := []MealDb{MealDb{Name: sql.NullString{`meal_1`, true}}} 
	meals := ConvertToMeals(mealDB)
	if len(meals) != 1 {
		t.Error(`too many meals`)
	}
	if meals[0].ProductsInMeal != nil {
		t.Error(`too many products in meal`)
	}
}


