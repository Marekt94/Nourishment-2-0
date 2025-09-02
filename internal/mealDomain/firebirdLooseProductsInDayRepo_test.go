package meal

import (
	db "nourishment_20/internal/database"
	"testing"
)

func initLooseProductsInDayRepo() LooseProductsInDayRepoIntf {
	var conf db.DBConf
	conf.User = `sysdba`
	conf.Password = `masterkey`
	conf.Address = `localhost:3050`
	conf.PathOrName = `C:/Users/marek/Documents/nourishment_backup_db/NOURISHMENT.FDB`

	fDbEngine := db.FBDBEngine{BaseEngineIntf: &db.BaseEngine{}}
	engine := fDbEngine.Connect(&conf)
	return &FirebirdRepoAccess{Database: engine}
}

func createTestProductForLoose() Product {
	repo := initProductRepo()
	prod := Product{Name: "loose_prod_test", KcalPer100: 100, UnitWeight: 100, Proteins: 10, Fat: 5, Sugar: 5, Carbohydrates: 20, Fiber: 1, Salt: 1, Unit: `g`, Category: Category{Id: 1, Name: `no category`}}
	i := int(repo.CreateProduct(&prod))
	prod.Id = i
	return prod
}

func createTestMealInDayForLoose() (MealInDay, int) {
	repo := initMealsInDayRepo()
	mealsRepo, ok := interface{}(repo).(MealsRepoIntf)
	if !ok {
		panic("MealsInDayRepo does not support MealsRepo interface")
	}
	meal := Meal{Name: "tmp meal for day", Recipe: "tmp"}
	id := mealsRepo.CreateMeal(&meal)
	meal.Id = int(id)
	mid := MealInDay{
		Breakfast:             meal,
		SecondBreakfast:       meal,
		Lunch:                 meal,
		Dinner:                meal,
		Supper:                meal,
		AfternoonSnack:        meal,
		For5Days:              true,
		FactorBreakfast:       1.0,
		FactorSecondBreakfast: 1.0,
		FactorLunch:           1.0,
		FactorDinner:          1.0,
		FactorSupper:          1.0,
		FactorAfternoonSnack:  1.0,
		Name:                  "day for loose",
	}
	dayId := repo.CreateMealsInDay(&mid)
	return mid, int(dayId)
}

func TestCreateLooseProductInDay(t *testing.T) {
	repo := initLooseProductsInDayRepo()
	prod := createTestProductForLoose()
	_, dayId := createTestMealInDayForLoose()

	lp := LooseProductInDay{DayId: dayId, Product: prod, Weight: 150}
	id := repo.CreateLooseProductInDay(&lp)
	if id <= 0 || lp.Id <= 0 {
		t.Error("CreateLooseProductInDay failed")
	}
}

func TestGetLooseProductInDay(t *testing.T) {
	repo := initLooseProductsInDayRepo()
	prod := createTestProductForLoose()
	_, dayId := createTestMealInDayForLoose()

	lp := LooseProductInDay{DayId: dayId, Product: prod, Weight: 200}
	id := repo.CreateLooseProductInDay(&lp)
	got := repo.GetLooseProductInDay(int(id))
	if got.Id != int(id) || got.DayId != dayId || got.Product.Id != prod.Id {
		t.Error("GetLooseProductInDay failed")
	}
}

func TestGetLooseProductsInDay(t *testing.T) {
	repo := initLooseProductsInDayRepo()
	prod := createTestProductForLoose()
	_, dayId := createTestMealInDayForLoose()

	lp1 := LooseProductInDay{DayId: dayId, Product: prod, Weight: 50}
	lp2 := LooseProductInDay{DayId: dayId, Product: prod, Weight: 75}
	repo.CreateLooseProductInDay(&lp1)
	repo.CreateLooseProductInDay(&lp2)

	list := repo.GetLooseProductsInDay(dayId)
	if len(list) == 0 {
		t.Error("GetLooseProductsInDay returned empty list")
	}
	// ensure at least one of the created IDs exists in the list
	found := false
	for _, it := range list {
		if it.Id == lp1.Id || it.Id == lp2.Id {
			found = true
			break
		}
	}
	if !found {
		t.Error("Created loose products not found in day list")
	}
}

func TestUpdateLooseProductInDay(t *testing.T) {
	repo := initLooseProductsInDayRepo()
	prod := createTestProductForLoose()
	_, dayId := createTestMealInDayForLoose()

	lp := LooseProductInDay{DayId: dayId, Product: prod, Weight: 10}
	repo.CreateLooseProductInDay(&lp)
	lp.Weight = 999
	repo.UpdateLooseProductInDay(&lp)
	got := repo.GetLooseProductInDay(lp.Id)
	if got.Weight != 999 {
		t.Error("UpdateLooseProductInDay failed")
	}
}

func TestDeleteLooseProductInDay(t *testing.T) {
	repo := initLooseProductsInDayRepo()
	prod := createTestProductForLoose()
	_, dayId := createTestMealInDayForLoose()

	lp := LooseProductInDay{DayId: dayId, Product: prod, Weight: 10}
	repo.CreateLooseProductInDay(&lp)
	ok := repo.DeleteLooseProductInDay(lp.Id)
	if !ok {
		t.Error("DeleteLooseProductInDay failed")
	}
}
