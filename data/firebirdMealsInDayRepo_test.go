package database

import (
	"database/sql"
	"testing"
)

func TestFirebirdMealsInDayRepo(t *testing.T) {
	repo := initMealsInDayRepo()
	mealApi := Meal{Name: "test meal", Recipe: "przepis"}
	mealId := repo.MealRepo.CreateMeal(&mealApi)
	meal := repo.MealRepo.GetMealDb(int(mealId))
	mealInDay := MealInDayDb{
		Breakfast: meal,
		SecondBreakfast: meal,
		Lunch: meal,
		Dinner: meal,
		Supper: meal,
		AfternoonSnack: meal,
		For5Days: sql.NullString{String: `1`, Valid: true},
		FactorBreakfast: sql.NullFloat64{Float64: 1.0, Valid: true},
		FactorSecondBreakfast: sql.NullFloat64{Float64: 1.0, Valid: true},
		FactorLunch: sql.NullFloat64{Float64: 1.0, Valid: true},
		FactorDinner: sql.NullFloat64{Float64: 1.0, Valid: true},
		FactorSupper: sql.NullFloat64{Float64: 1.0, Valid: true},
		FactorAfternoonSnack: sql.NullFloat64{Float64: 1.0, Valid: true},
		Name: sql.NullString{String: "test day", Valid: true},
	}
	//TODO: wydzielić do osobnego
	id := repo.CreateMealInDay(&mealInDay)
	if id <= 0 {
		t.Error("CreateMealInDay failed")
	}
	got := repo.GetMealInDay(int(id))
	if got.Name.String != "test day" {
		t.Error("GetMealInDay failed")
	}
	all := repo.GetMealsInDay()
	if len(all) != 1 {
		t.Error("GetMealsInDay failed")
	}
	mealInDay.Name = sql.NullString{String: "updated", Valid: true}
	repo.UpdateMealInDay(&mealInDay)
	got2 := repo.GetMealInDay(int(id))
	if got2.Name.String != "updated" {
		t.Error("UpdateMealInDay failed")
	}
	ok := repo.DeleteMealInDay(int(id))
	if !ok || len(repo.GetMealsInDay()) != 0 {
		t.Error("DeleteMealInDay failed")
	}
}

func initMealsInDayRepo() *FirebirdMealsInDayRepo {
	var conf DBConf
	conf.User = `sysdba`
	conf.Password = `masterkey`
	conf.Address = `localhost:3050`
	conf.PathOrName = `C:/Users/marek/Documents/nourishment_backup_db/NOURISHMENT.FDB`

	fDbEngine := FBDBEngine{BaseEngineIntf: &BaseEngine{}}
	engine := fDbEngine.Connect(&conf)
	mealRepo := &FirebirdRepoAccess{DbEngine: engine}
	return &FirebirdMealsInDayRepo{DbEngine: engine, MealRepo: mealRepo}
}
