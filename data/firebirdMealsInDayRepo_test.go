package database

import (
	"testing"
)

func createTestMealInDay(repo MealsInDayRepo) (MealInDay, int) {
	mealApi := Meal{Name: "test meal", Recipe: "przepis"}
	mealsRepo, ok := interface{}(repo).(MealsRepo)
	if !ok {
		panic("MealsInDayRepo does not support MealsRepo interface")
	}
	mealId := mealsRepo.CreateMeal(&mealApi)
	mealInDay := MealInDay{
		Breakfast: Meal{Id: int(mealId)},
		SecondBreakfast: Meal{Id: int(mealId)},
		Lunch: Meal{Id: int(mealId)},
		Dinner: Meal{Id: int(mealId)},
		Supper: Meal{Id: int(mealId)},
		AfternoonSnack: Meal{Id: int(mealId)},
		For5Days: true,
		FactorBreakfast: 1.0,
		FactorSecondBreakfast: 1.0,
		FactorLunch: 1.0,
		FactorDinner: 1.0,
		FactorSupper: 1.0,
		FactorAfternoonSnack: 1.0,
		Name: "test day",
	}
	id := repo.CreateMealInDay(&mealInDay)
	return mealInDay, int(id)
}

func TestCreateMealsInDay(t *testing.T) {
	repo := initMealsInDayRepo()
	_, id := createTestMealInDay(repo)
	if id <= 0 {
		t.Error("CreateMealInDay failed")
	}
}

func TestGetMealsInDay(t *testing.T) {
	repo := initMealsInDayRepo()
	mealInDay, id := createTestMealInDay(repo)
	got := repo.GetMealInDay(id)
	if got.Name != mealInDay.Name {
		t.Error("GetMealInDay failed")
	}
}

func TestGetMealsInDays(t *testing.T) {
	repo := initMealsInDayRepo()
	createTestMealInDay(repo)
	all := repo.GetMealsInDay()
	if len(all) != 1 {
		t.Error("GetMealsInDay failed")
	}
}

func TestUpdateMealsInDay(t *testing.T) {
	repo := initMealsInDayRepo()
	mealInDay, id := createTestMealInDay(repo)
	mealInDay.Name = "updated"
	repo.UpdateMealInDay(&mealInDay)
	got2 := repo.GetMealInDay(id)
	if got2.Name != "updated" {
		t.Error("UpdateMealInDay failed")
	}
}

func TestDeleteMealsInDay(t *testing.T) {
	repo := initMealsInDayRepo()
	_, id := createTestMealInDay(repo)
	ok2 := repo.DeleteMealInDay(id)
	if !ok2 || len(repo.GetMealsInDay()) != 0 {
		t.Error("DeleteMealInDay failed")
	}
}

func initMealsInDayRepo() MealsInDayRepo {
	var conf DBConf
	conf.User = `sysdba`
	conf.Password = `masterkey`
	conf.Address = `localhost:3050`
	conf.PathOrName = `C:/Users/marek/Documents/nourishment_backup_db/NOURISHMENT.FDB`

	fDbEngine := FBDBEngine{BaseEngineIntf: &BaseEngine{}}
	engine := fDbEngine.Connect(&conf)
	return &FirebirdRepoAccess{DbEngine: engine}
}
