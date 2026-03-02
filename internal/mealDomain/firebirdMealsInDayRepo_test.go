package meal

import (
	"testing"

	"nourishment_20/internal/database"
	db "nourishment_20/internal/database"
)

func createTestMealInDay(repo MealsInDayRepoIntf) (MealInDay, int) {
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

func TestCreateMealsInDayWithoutLunch(t *testing.T) {
	repo := initMealsInDayRepo()

	// Stwórz testowy posiłek
	mealApi := Meal{Name: "test meal", Recipe: "przepis"}
	mealsRepo, ok := interface{}(repo).(MealsRepoIntf)
	if !ok {
		t.Fatal("MealsInDayRepo does not support MealsRepo interface")
	}
	mealId := mealsRepo.CreateMeal(&mealApi)
	mealApi.Id = int(mealId)

	// Utwórz MealInDay bez Lunchu (użyj pustego Meal z EMPTY_ID)
	emptyMeal := Meal{Id: EMPTY_ID}
	mealInDay := MealInDay{
		Breakfast:             mealApi,
		SecondBreakfast:       mealApi,
		Lunch:                 emptyMeal, // Brak lunchu
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
		Name:                  "test day without lunch",
	}

	// Dodaj do bazy
	id := repo.CreateMealsInDay(&mealInDay)
	if id <= 0 {
		t.Fatal("CreateMealInDay without lunch failed")
	}

	// Pobierz z bazy i sprawdź
	retrieved := repo.GetMealsInDay(int(id))
	if retrieved.Name != "test day without lunch" {
		t.Errorf("Expected name 'test day without lunch', got '%s'", retrieved.Name)
	}

	// Sprawdź czy Lunch jest pusty (ma EMPTY_ID)
	if retrieved.Lunch.Id != 0 {
		t.Errorf("Expected Lunch.Id to be %d, got %d", EMPTY_ID, retrieved.Lunch.Id)
	}

	// Sprawdź czy pozostałe posiłki są ustawione
	if retrieved.Breakfast.Id != int(mealId) {
		t.Error("Breakfast should be set")
	}
	if retrieved.Dinner.Id != int(mealId) {
		t.Error("Dinner should be set")
	}

	// Cleanup
	repo.DeleteMealsInDay(int(id))
}

func TestUpdateMealsInDayRemoveMeal(t *testing.T) {
	repo := initMealsInDayRepo()

	// Stwórz testowy posiłek
	mealApi := Meal{Name: "test meal for removal", Recipe: "przepis"}
	mealsRepo, ok := interface{}(repo).(MealsRepoIntf)
	if !ok {
		t.Fatal("MealsInDayRepo does not support MealsRepo interface")
	}
	mealId := mealsRepo.CreateMeal(&mealApi)
	mealApi.Id = int(mealId)

	// Utwórz MealInDay z WSZYSTKIMI posiłkami
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
		Name:                  "test day with all meals",
	}

	// Dodaj do bazy
	id := repo.CreateMealsInDay(&mealInDay)
	if id <= 0 {
		t.Fatal("CreateMealInDay failed")
	}

	// Pobierz z bazy i sprawdź że wszystkie posiłki są ustawione
	retrieved := repo.GetMealsInDay(int(id))
	if retrieved.Breakfast.Id != int(mealId) {
		t.Error("Breakfast should be set initially")
	}
	if retrieved.Lunch.Id != int(mealId) {
		t.Error("Lunch should be set initially")
	}
	if retrieved.Dinner.Id != int(mealId) {
		t.Error("Dinner should be set initially")
	}

	// Usuń Lunch i Dinner z obiektu (ustaw na EMPTY_ID)
	retrieved.Lunch = Meal{Id: EMPTY_ID}
	retrieved.Dinner = Meal{Id: EMPTY_ID}

	// Zaktualizuj w bazie
	repo.UpdateMealsInDay(&retrieved)

	// Pobierz ponownie z bazy
	updated := repo.GetMealsInDay(int(id))

	// Sprawdź czy Lunch i Dinner są teraz NULL (Id = 0 po konwersji z NULL)
	if updated.Lunch.Id != 0 {
		t.Errorf("Expected Lunch.Id to be 0 (NULL), got %d", updated.Lunch.Id)
	}
	if updated.Dinner.Id != 0 {
		t.Errorf("Expected Dinner.Id to be 0 (NULL), got %d", updated.Dinner.Id)
	}

	// Sprawdź czy pozostałe posiłki nadal są ustawione
	if updated.Breakfast.Id != int(mealId) {
		t.Error("Breakfast should still be set")
	}
	if updated.SecondBreakfast.Id != int(mealId) {
		t.Error("SecondBreakfast should still be set")
	}
	if updated.Supper.Id != int(mealId) {
		t.Error("Supper should still be set")
	}
	if updated.AfternoonSnack.Id != int(mealId) {
		t.Error("AfternoonSnack should still be set")
	}

	// Cleanup
	repo.DeleteMealsInDay(int(id))
	mealsRepo.DeleteMeal(int(mealId))
}

func initMealsInDayRepo() MealsInDayRepoIntf {
	var conf database.DBConf
	conf.User = `sysdba`
	conf.Password = `masterkey`
	conf.Address = `localhost:3050`
	conf.PathOrName = `C:/Users/marek/Documents/nourishment_backup_db/NOURISHMENT.FDB`

	fDbEngine := db.FBDBEngine{BaseEngineIntf: &database.BaseEngine{}}
	engine := fDbEngine.Connect(&conf)
	return &FirebirdRepoAccess{Database: engine}
}
