package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
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
	MealRepo *FirebirdRepoAccess // Injected for meal lookups
}

func (repo *FirebirdMealsInDayRepo) CreateMealInDay(m *MealInDayDb) int64 {
	// Check if all referenced meals exist
	mealIds := []sql.NullInt64{
		m.Breakfast.Id, m.SecondBreakfast.Id, m.Lunch.Id, m.Dinner.Id, m.Supper.Id, m.AfternoonSnack.Id,
	}
	for _, mealId := range mealIds {
		if mealId.Valid && repo.MealRepo.GetMeal(int(mealId.Int64)).Id == 0 {
			log.Println("CreateMealInDay error: referenced meal does not exist, id:", mealId.Int64)
			return -1
		}
	}
	cols := []string{
		MEAL_IN_DAY_BREAKFAST, MEAL_IN_DAY_SECOND_BREAKFAST, MEAL_IN_DAY_LUNCH, MEAL_IN_DAY_DINNER, MEAL_IN_DAY_SUPPER, MEAL_IN_DAY_AFTERNOON_SNACK,
		MEAL_IN_DAY_FOR_5_DAYS, MEAL_IN_DAY_FACTOR_BREAKFAST, MEAL_IN_DAY_FACTOR_SECOND_BREAKFAST, MEAL_IN_DAY_FACTOR_LUNCH, MEAL_IN_DAY_FACTOR_DINNER, MEAL_IN_DAY_FACTOR_SUPPER, MEAL_IN_DAY_FACTOR_AFTERNOON_SNACK, MEAL_IN_DAY_NAME,
	}
	placeholders := QuestionMarks(len(cols))
	sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", MEAL_IN_DAY_TAB, strings.Join(cols, ", "), placeholders)
	_, err := repo.DbEngine.Exec(sqlStr,
		m.Breakfast.Id, m.SecondBreakfast.Id, m.Lunch.Id, m.Dinner.Id, m.Supper.Id, m.AfternoonSnack.Id,
		m.For5Days, m.FactorBreakfast, m.FactorSecondBreakfast, m.FactorLunch, m.FactorDinner, m.FactorSupper, m.FactorAfternoonSnack, m.Name)
	if err != nil {
		log.Println("CreateMealInDay error:", err)
		return -1
	}
	var id int64
	err = repo.DbEngine.QueryRow(fmt.Sprintf("SELECT MAX(%s) FROM %s", MEAL_IN_DAY_MEAL_ID, MEAL_IN_DAY_TAB)).Scan(&id)
	if err != nil {
		log.Println("CreateMealInDay get id error:", err)
		return -1
	}
	return id
}

func (repo *FirebirdMealsInDayRepo) GetMealInDay(id int) MealInDayDb {
	cols := []string{
		MEAL_IN_DAY_MEAL_ID, MEAL_IN_DAY_BREAKFAST, MEAL_IN_DAY_SECOND_BREAKFAST, MEAL_IN_DAY_LUNCH, MEAL_IN_DAY_DINNER, MEAL_IN_DAY_SUPPER, MEAL_IN_DAY_AFTERNOON_SNACK,
		MEAL_IN_DAY_FOR_5_DAYS, MEAL_IN_DAY_FACTOR_BREAKFAST, MEAL_IN_DAY_FACTOR_SECOND_BREAKFAST, MEAL_IN_DAY_FACTOR_LUNCH, MEAL_IN_DAY_FACTOR_DINNER, MEAL_IN_DAY_FACTOR_SUPPER, MEAL_IN_DAY_FACTOR_AFTERNOON_SNACK, MEAL_IN_DAY_NAME,
	}
	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ?", strings.Join(cols, ", "), MEAL_IN_DAY_TAB, MEAL_IN_DAY_MEAL_ID)
	row := repo.DbEngine.QueryRow(sqlStr, id)
	var m MealInDayDb
	var breakfastId, secondBreakfastId, lunchId, dinnerId, supperId, snackId sql.NullInt64
	err := row.Scan(&m.Id, &breakfastId, &secondBreakfastId, &lunchId, &dinnerId, &supperId, &snackId, &m.For5Days, &m.FactorBreakfast, &m.FactorSecondBreakfast, &m.FactorLunch, &m.FactorDinner, &m.FactorSupper, &m.FactorAfternoonSnack, &m.Name)
	if err != nil {
		log.Println("GetMealInDay error:", err)
		return m
	}
	m.Breakfast = repo.MealRepo.GetMealDb(int(breakfastId.Int64))
	m.SecondBreakfast = repo.MealRepo.GetMealDb(int(secondBreakfastId.Int64))
	m.Lunch = repo.MealRepo.GetMealDb(int(lunchId.Int64))
	m.Dinner = repo.MealRepo.GetMealDb(int(dinnerId.Int64))
	m.Supper = repo.MealRepo.GetMealDb(int(supperId.Int64))
	m.AfternoonSnack = repo.MealRepo.GetMealDb(int(snackId.Int64))
	return m
}

func (repo *FirebirdMealsInDayRepo) GetMealsInDay() []MealInDayDb {
	cols := []string{
		MEAL_IN_DAY_MEAL_ID, MEAL_IN_DAY_BREAKFAST, MEAL_IN_DAY_SECOND_BREAKFAST, MEAL_IN_DAY_LUNCH, MEAL_IN_DAY_DINNER, MEAL_IN_DAY_SUPPER, MEAL_IN_DAY_AFTERNOON_SNACK,
		MEAL_IN_DAY_FOR_5_DAYS, MEAL_IN_DAY_FACTOR_BREAKFAST, MEAL_IN_DAY_FACTOR_SECOND_BREAKFAST, MEAL_IN_DAY_FACTOR_LUNCH, MEAL_IN_DAY_FACTOR_DINNER, MEAL_IN_DAY_FACTOR_SUPPER, MEAL_IN_DAY_FACTOR_AFTERNOON_SNACK, MEAL_IN_DAY_NAME,
	}
	sqlStr := fmt.Sprintf("SELECT %s FROM %s", strings.Join(cols, ", "), MEAL_IN_DAY_TAB)
	rows, err := repo.DbEngine.Query(sqlStr)
	if err != nil {
		log.Println("GetMealsInDay error:", err)
		return nil
	}
	var res []MealInDayDb
	for rows.Next() {
		var m MealInDayDb
		var breakfastId, secondBreakfastId, lunchId, dinnerId, supperId, snackId sql.NullInt64
		err := rows.Scan(&m.Id, &breakfastId, &secondBreakfastId, &lunchId, &dinnerId, &supperId, &snackId, &m.For5Days, &m.FactorBreakfast, &m.FactorSecondBreakfast, &m.FactorLunch, &m.FactorDinner, &m.FactorSupper, &m.FactorAfternoonSnack, &m.Name)
		if err == nil {
			m.Breakfast = repo.MealRepo.GetMealDb(int(breakfastId.Int64))
			m.SecondBreakfast = repo.MealRepo.GetMealDb(int(secondBreakfastId.Int64))
			m.Lunch = repo.MealRepo.GetMealDb(int(lunchId.Int64))
			m.Dinner = repo.MealRepo.GetMealDb(int(dinnerId.Int64))
			m.Supper = repo.MealRepo.GetMealDb(int(supperId.Int64))
			m.AfternoonSnack = repo.MealRepo.GetMealDb(int(snackId.Int64))
			res = append(res, m)
		}
	}
	return res
}

func (repo *FirebirdMealsInDayRepo) DeleteMealInDay(id int) bool {
	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", MEAL_IN_DAY_TAB, MEAL_IN_DAY_MEAL_ID)
	_, err := repo.DbEngine.Exec(sqlStr, id)
	if err != nil {
		log.Println("DeleteMealInDay error:", err)
		return false
	}
	return true
}

func (repo *FirebirdMealsInDayRepo) UpdateMealInDay(m *MealInDayDb) {
	cols := []string{
		MEAL_IN_DAY_BREAKFAST, MEAL_IN_DAY_SECOND_BREAKFAST, MEAL_IN_DAY_LUNCH, MEAL_IN_DAY_DINNER, MEAL_IN_DAY_SUPPER, MEAL_IN_DAY_AFTERNOON_SNACK,
		MEAL_IN_DAY_FOR_5_DAYS, MEAL_IN_DAY_FACTOR_BREAKFAST, MEAL_IN_DAY_FACTOR_SECOND_BREAKFAST, MEAL_IN_DAY_FACTOR_LUNCH, MEAL_IN_DAY_FACTOR_DINNER, MEAL_IN_DAY_FACTOR_SUPPER, MEAL_IN_DAY_FACTOR_AFTERNOON_SNACK, MEAL_IN_DAY_NAME,
	}
	setExprs := make([]string, len(cols))
	for i, col := range cols {
		setExprs[i] = col + "=?"
	}
	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE %s=?", MEAL_IN_DAY_TAB, strings.Join(setExprs, ", "), MEAL_IN_DAY_MEAL_ID)
	_, err := repo.DbEngine.Exec(sqlStr,
		m.Breakfast.Id, m.SecondBreakfast.Id, m.Lunch.Id, m.Dinner.Id, m.Supper.Id, m.AfternoonSnack.Id,
		m.For5Days, m.FactorBreakfast, m.FactorSecondBreakfast, m.FactorLunch, m.FactorDinner, m.FactorSupper, m.FactorAfternoonSnack, m.Name, m.Id)
	if err != nil {
		log.Println("UpdateMealInDay error:", err)
	}
}