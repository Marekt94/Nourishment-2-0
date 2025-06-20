package database

import (
	"database/sql"
	"testing"
)

// [API GEN] Test inicjalizacji struktury MealInDayDb
func TestCreateMealInDayDb(t *testing.T) {
	meal := MealDb{Id: sql.NullInt64{Int64: 1, Valid: true}, Name: sql.NullString{String: "test meal", Valid: true}}
	mealInDay := MealInDayDb{
		Id: sql.NullInt64{Int64: 10, Valid: true},
		Breakfast: meal,
		SecondBreakfast: meal,
		Lunch: meal,
		Dinner: meal,
		Supper: meal,
		AfternoonSnack: meal,
		For5Days: sql.NullBool{Bool: true, Valid: true},
		FactorBreakfast: sql.NullFloat64{Float64: 1.0, Valid: true},
		FactorSecondBreakfast: sql.NullFloat64{Float64: 1.0, Valid: true},
		FactorLunch: sql.NullFloat64{Float64: 1.0, Valid: true},
		FactorDinner: sql.NullFloat64{Float64: 1.0, Valid: true},
		FactorSupper: sql.NullFloat64{Float64: 1.0, Valid: true},
		FactorAfternoonSnack: sql.NullFloat64{Float64: 1.0, Valid: true},
		Name: sql.NullString{String: "test day", Valid: true},
	}

	if !mealInDay.Id.Valid || mealInDay.Id.Int64 != 10 {
		t.Error("MealInDayDb.Id not set correctly")
	}
	if mealInDay.Breakfast.Name.String != "test meal" {
		t.Error("MealInDayDb.Breakfast not set correctly")
	}
	if !mealInDay.For5Days.Bool {
		t.Error("MealInDayDb.For5Days not set correctly")
	}
	if mealInDay.Name.String != "test day" {
		t.Error("MealInDayDb.Name not set correctly")
	}
}

// [API GEN] Zaslepka repozytorium MealsInDayRepo
// Usunięto interfejs MealsInDayRepo z testu, zostaje tylko FakeMealsInDayRepo

type FakeMealsInDayRepo struct { // [API GEN]
	store map[int]MealInDayDb
	nextId int
}

func NewFakeMealsInDayRepo() *FakeMealsInDayRepo {
	return &FakeMealsInDayRepo{store: make(map[int]MealInDayDb), nextId: 1}
}

func (r *FakeMealsInDayRepo) CreateMealInDay(m *MealInDayDb) int64 {
	m.Id = sql.NullInt64{Int64: int64(r.nextId), Valid: true}
	r.store[r.nextId] = *m
	r.nextId++
	return m.Id.Int64
}

func (r *FakeMealsInDayRepo) GetMealInDay(id int) MealInDayDb {
	return r.store[id]
}

func (r *FakeMealsInDayRepo) GetMealsInDay() []MealInDayDb {
	res := []MealInDayDb{}
	for _, v := range r.store {
		res = append(res, v)
	}
	return res
}

func (r *FakeMealsInDayRepo) DeleteMealInDay(id int) bool {
	if _, ok := r.store[id]; ok {
		delete(r.store, id)
		return true
	}
	return false
}

func (r *FakeMealsInDayRepo) UpdateMealInDay(m *MealInDayDb) {
	if m.Id.Valid {
		r.store[int(m.Id.Int64)] = *m
	}
}

// [API GEN] Testy CRUD dla MealsInDayRepo
func TestFakeMealsInDayRepo_CRUD(t *testing.T) {
	repo := NewFakeMealsInDayRepo()
	meal := MealDb{Id: sql.NullInt64{Int64: 1, Valid: true}, Name: sql.NullString{String: "test meal", Valid: true}}
	mealInDay := MealInDayDb{
		Breakfast: meal,
		SecondBreakfast: meal,
		Lunch: meal,
		Dinner: meal,
		Supper: meal,
		AfternoonSnack: meal,
		For5Days: sql.NullBool{Bool: true, Valid: true},
		FactorBreakfast: sql.NullFloat64{Float64: 1.0, Valid: true},
		FactorSecondBreakfast: sql.NullFloat64{Float64: 1.0, Valid: true},
		FactorLunch: sql.NullFloat64{Float64: 1.0, Valid: true},
		FactorDinner: sql.NullFloat64{Float64: 1.0, Valid: true},
		FactorSupper: sql.NullFloat64{Float64: 1.0, Valid: true},
		FactorAfternoonSnack: sql.NullFloat64{Float64: 1.0, Valid: true},
		Name: sql.NullString{String: "test day", Valid: true},
	}
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
