package database

import (
	"database/sql"
	"log"
)

// [API GEN] DTO do tabeli POTRAWY_W_DNIU
// Odpowiada kolumnom zdefiniowanym w firebirdDatabase.go
// Każde pole odpowiada jednej kolumnie w bazie

type MealInDayDb struct { // [API GEN]
	Id sql.NullInt64
	Breakfast MealDb
	SecondBreakfast MealDb
	Lunch MealDb
	Dinner MealDb
	Supper MealDb
	AfternoonSnack MealDb
	For5Days sql.NullBool
	FactorBreakfast sql.NullFloat64
	FactorSecondBreakfast sql.NullFloat64
	FactorLunch sql.NullFloat64
	FactorDinner sql.NullFloat64
	FactorSupper sql.NullFloat64
	FactorAfternoonSnack sql.NullFloat64
	Name sql.NullString
}

// [API GEN] Interfejs MealsInDayRepo

type MealsInDayRepo interface {
	CreateMealInDay(m *MealInDayDb) int64
	GetMealInDay(id int) MealInDayDb
	GetMealsInDay() []MealInDayDb
	DeleteMealInDay(id int) bool
	UpdateMealInDay(m *MealInDayDb)
}

// [API GEN] Implementacja repozytorium MealsInDayRepo dla Firebird

type FirebirdMealsInDayRepo struct {
	DbEngine *sql.DB
}

func (repo *FirebirdMealsInDayRepo) CreateMealInDay(m *MealInDayDb) int64 {
	// TODO: implementacja zapisu do bazy
	log.Println("CreateMealInDay not implemented")
	return -1
}

func (repo *FirebirdMealsInDayRepo) GetMealInDay(id int) MealInDayDb {
	// TODO: implementacja pobierania z bazy
	log.Println("GetMealInDay not implemented")
	return MealInDayDb{}
}

func (repo *FirebirdMealsInDayRepo) GetMealsInDay() []MealInDayDb {
	// TODO: implementacja pobierania wszystkich rekordów
	log.Println("GetMealsInDay not implemented")
	return nil
}

func (repo *FirebirdMealsInDayRepo) DeleteMealInDay(id int) bool {
	// TODO: implementacja usuwania z bazy
	log.Println("DeleteMealInDay not implemented")
	return false
}

func (repo *FirebirdMealsInDayRepo) UpdateMealInDay(m *MealInDayDb) {
	// TODO: implementacja aktualizacji w bazie
	log.Println("UpdateMealInDay not implemented")
}
