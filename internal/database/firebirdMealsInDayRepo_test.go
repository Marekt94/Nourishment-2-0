package database

import (
	"testing"
)

func createTestMealInDay(repo MealsInDayRepo) (MealInDay, int) {
	mealApi := Meal{Name: "test meal", Recipe: "przepis"}
	mealsRepo, ok := interface{}(repo).(MealsRepoIntf)
	if !ok {
		panic("MealsInDayRepo does not support MealsRepo interface")
	}
	mealId := mealsRepo.CreateMeal(&mealApi)
	mealApi.Id = int(mealId)
	mealInDay := MealInDay{
		Breakfast:             mealApi,
		SecondBreakfast:       mealApi,
		Lunch:                 mealApi,
		Dinner:                mealApi,
		Supper:                mealApi,
		AfternoonSnack:        mealApi,
		For5Days:              true,
		FactorBreakfast:       1.0,
		FactorSecondBreakfast: 1.0,
		FactorLunch:           1.0,
		FactorDinner:          1.0,
		FactorSupper:          1.0,
		FactorAfternoonSnack:  1.0,
		Name:                  "test day",
	}
	id := repo.CreateMealsInDay(&mealInDay)
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
	got := repo.GetMealsInDay(id)
	if got.Name != mealInDay.Name {
		t.Error("GetMealInDay failed")
	}
}

func TestGetMealsInDays(t *testing.T) {
	repo := initMealsInDayRepo()
	createTestMealInDay(repo)
	all := repo.GetMealsInDays()
	if len(all) == 1 || len(all) == 0 {
		t.Error("GetMealsInDay failed")
	}
}

func TestUpdateMealsInDay(t *testing.T) {
	repo := initMealsInDayRepo()
	mealInDay, id := createTestMealInDay(repo)
	mealInDay.Name = "updated"
	repo.UpdateMealsInDay(&mealInDay)
	got2 := repo.GetMealsInDay(id)
	if got2.Name != "updated" {
		t.Error("UpdateMealInDay failed")
	}
}

func TestDeleteMealsInDay(t *testing.T) {
	repo := initMealsInDayRepo()
	_, id := createTestMealInDay(repo)
	ok2 := repo.DeleteMealsInDay(id)
	if !ok2 || (repo.GetMealsInDay(id).Id != 0) {
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
	return &FirebirdRepoAccess{Database: engine}
}
