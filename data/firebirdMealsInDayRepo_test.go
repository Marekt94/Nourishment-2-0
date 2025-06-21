package database

import (
	"testing"
)

func TestFirebirdMealsInDayRepo(t *testing.T) {
	repo := initMealsInDayRepo()
	mealApi := Meal{Name: "test meal", Recipe: "przepis"}
	res, supp := interface{}(repo).(MealsRepo)
	if !supp {
		t.Fatal("MealsInDayRepo does not support MealsRepo interface")
	}
	mealId := res.CreateMeal(&mealApi)
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
	if id <= 0 {
		t.Error("CreateMealInDay failed")
	}
	got := repo.GetMealInDay(int(id))
	if got.Name != "test day" {
		t.Error("GetMealInDay failed")
	}
	all := repo.GetMealsInDay()
	if len(all) != 1 {
		t.Error("GetMealsInDay failed")
	}
	mealInDay.Name = "updated"
	repo.UpdateMealInDay(&mealInDay)
	got2 := repo.GetMealInDay(int(id))
	if got2.Name != "updated" {
		t.Error("UpdateMealInDay failed")
	}
	ok := repo.DeleteMealInDay(int(id))
	if !ok || len(repo.GetMealsInDay()) != 0 {
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
