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
	BreakfastId sql.NullInt64
	SecondBreakfastId sql.NullInt64
	LunchId sql.NullInt64
	DinnerId sql.NullInt64
	SupperId sql.NullInt64
	AfternoonSnackId sql.NullInt64
	For5Days sql.NullString // CHAR(1) '1' lub '0'
	FactorBreakfast sql.NullFloat64
	FactorSecondBreakfast sql.NullFloat64
	FactorLunch sql.NullFloat64
	FactorDinner sql.NullFloat64
	FactorSupper sql.NullFloat64
	FactorAfternoonSnack sql.NullFloat64
	Name sql.NullString
}

func (mr *FirebirdRepoAccess) CreateMealInDay(m *MealInDay) int64 {
	mealIds := []int{m.Breakfast.Id, m.SecondBreakfast.Id, m.Lunch.Id, m.Dinner.Id, m.Supper.Id, m.AfternoonSnack.Id}
	for _, mealId := range mealIds {
		if mealId > 0 && mr.GetMeal(mealId).Id == 0 {
			log.Println("CreateMealInDay error: referenced meal does not exist, id:", mealId)
			return -1
		}
	}
	cols := []string{
		MEAL_IN_DAY_BREAKFAST, MEAL_IN_DAY_SECOND_BREAKFAST, MEAL_IN_DAY_LUNCH, MEAL_IN_DAY_DINNER, MEAL_IN_DAY_SUPPER, MEAL_IN_DAY_AFTERNOON_SNACK,
		MEAL_IN_DAY_FOR_5_DAYS, MEAL_IN_DAY_FACTOR_BREAKFAST, MEAL_IN_DAY_FACTOR_SECOND_BREAKFAST, MEAL_IN_DAY_FACTOR_LUNCH, MEAL_IN_DAY_FACTOR_DINNER, MEAL_IN_DAY_FACTOR_SUPPER, MEAL_IN_DAY_FACTOR_AFTERNOON_SNACK, MEAL_IN_DAY_NAME,
	}
	placeholders := QuestionMarks(len(cols))
	sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", MEAL_IN_DAY_TAB, strings.Join(cols, ", "), placeholders)
	for5DaysChar := sql.NullString{String: "0", Valid: true}
	if m.For5Days {
		for5DaysChar.String = "1"
	}
	_, err := mr.DbEngine.Exec(sqlStr,
		m.Breakfast.Id, m.SecondBreakfast.Id, m.Lunch.Id, m.Dinner.Id, m.Supper.Id, m.AfternoonSnack.Id,
		for5DaysChar, m.FactorBreakfast, m.FactorSecondBreakfast, m.FactorLunch, m.FactorDinner, m.FactorSupper, m.FactorAfternoonSnack, m.Name)
	if err != nil {
		log.Println("CreateMealInDay error:", err)
		return -1
	}
	var id int64
	err = mr.DbEngine.QueryRow(fmt.Sprintf("SELECT MAX(%s) FROM %s", MEAL_IN_DAY_MEAL_ID, MEAL_IN_DAY_TAB)).Scan(&id)
	if err != nil {
		log.Println("CreateMealInDay get id error:", err)
		return -1
	}
	return id
}

func (mr *FirebirdRepoAccess) GetMealInDay(id int) MealInDay {
	cols := []string{
		MEAL_IN_DAY_MEAL_ID, MEAL_IN_DAY_BREAKFAST, MEAL_IN_DAY_SECOND_BREAKFAST, MEAL_IN_DAY_LUNCH, MEAL_IN_DAY_DINNER, MEAL_IN_DAY_SUPPER, MEAL_IN_DAY_AFTERNOON_SNACK,
		MEAL_IN_DAY_FOR_5_DAYS, MEAL_IN_DAY_FACTOR_BREAKFAST, MEAL_IN_DAY_FACTOR_SECOND_BREAKFAST, MEAL_IN_DAY_FACTOR_LUNCH, MEAL_IN_DAY_FACTOR_DINNER, MEAL_IN_DAY_FACTOR_SUPPER, MEAL_IN_DAY_FACTOR_AFTERNOON_SNACK, MEAL_IN_DAY_NAME,
	}
	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ?", strings.Join(cols, ", "), MEAL_IN_DAY_TAB, MEAL_IN_DAY_MEAL_ID)
	row := mr.DbEngine.QueryRow(sqlStr, id)
	var m MealInDayDb
	var for5DaysChar sql.NullString
	err := row.Scan(&m.Id, &m.BreakfastId, &m.SecondBreakfastId, &m.LunchId, &m.DinnerId, &m.SupperId, &m.AfternoonSnackId, &for5DaysChar, &m.FactorBreakfast, &m.FactorSecondBreakfast, &m.FactorLunch, &m.FactorDinner, &m.FactorSupper, &m.FactorAfternoonSnack, &m.Name)
	m.For5Days = for5DaysChar
	if err != nil {
		log.Println("GetMealInDay error:", err)
		return MealInDay{}
	}
	return MealInDay{
		Id: int(m.Id.Int64),
		Breakfast: mr.GetMeal(int(m.BreakfastId.Int64)),
		SecondBreakfast: mr.GetMeal(int(m.SecondBreakfastId.Int64)),
		Lunch: mr.GetMeal(int(m.LunchId.Int64)),
		Dinner: mr.GetMeal(int(m.DinnerId.Int64)),
		Supper: mr.GetMeal(int(m.SupperId.Int64)),
		AfternoonSnack: mr.GetMeal(int(m.AfternoonSnackId.Int64)),
		For5Days: m.For5Days.String == "1",
		FactorBreakfast: m.FactorBreakfast.Float64,
		FactorSecondBreakfast: m.FactorSecondBreakfast.Float64,
		FactorLunch: m.FactorLunch.Float64,
		FactorDinner: m.FactorDinner.Float64,
		FactorSupper: m.FactorSupper.Float64,
		FactorAfternoonSnack: m.FactorAfternoonSnack.Float64,
		Name: m.Name.String,
	}
}

func (mr *FirebirdRepoAccess) GetMealsInDay() []MealInDay {
	cols := []string{
		MEAL_IN_DAY_MEAL_ID, MEAL_IN_DAY_BREAKFAST, MEAL_IN_DAY_SECOND_BREAKFAST, MEAL_IN_DAY_LUNCH, MEAL_IN_DAY_DINNER, MEAL_IN_DAY_SUPPER, MEAL_IN_DAY_AFTERNOON_SNACK,
		MEAL_IN_DAY_FOR_5_DAYS, MEAL_IN_DAY_FACTOR_BREAKFAST, MEAL_IN_DAY_FACTOR_SECOND_BREAKFAST, MEAL_IN_DAY_FACTOR_LUNCH, MEAL_IN_DAY_FACTOR_DINNER, MEAL_IN_DAY_FACTOR_SUPPER, MEAL_IN_DAY_FACTOR_AFTERNOON_SNACK, MEAL_IN_DAY_NAME,
	}
	sqlStr := fmt.Sprintf("SELECT %s FROM %s", strings.Join(cols, ", "), MEAL_IN_DAY_TAB)
	rows, err := mr.DbEngine.Query(sqlStr)
	if err != nil {
		log.Println("GetMealsInDay error:", err)
		return nil
	}
	var res []MealInDay
	for rows.Next() {
		var m MealInDayDb
		var for5DaysChar sql.NullString
		err := rows.Scan(&m.Id, &m.BreakfastId, &m.SecondBreakfastId, &m.LunchId, &m.DinnerId, &m.SupperId, &m.AfternoonSnackId, &for5DaysChar, &m.FactorBreakfast, &m.FactorSecondBreakfast, &m.FactorLunch, &m.FactorDinner, &m.FactorSupper, &m.FactorAfternoonSnack, &m.Name)
		m.For5Days = for5DaysChar
		if err == nil {
			res = append(res, MealInDay{
				Id: int(m.Id.Int64),
				Breakfast: mr.GetMeal(int(m.BreakfastId.Int64)),
				SecondBreakfast: mr.GetMeal(int(m.SecondBreakfastId.Int64)),
				Lunch: mr.GetMeal(int(m.LunchId.Int64)),
				Dinner: mr.GetMeal(int(m.DinnerId.Int64)),
				Supper: mr.GetMeal(int(m.SupperId.Int64)),
				AfternoonSnack: mr.GetMeal(int(m.AfternoonSnackId.Int64)),
				For5Days: m.For5Days.String == "1",
				FactorBreakfast: m.FactorBreakfast.Float64,
				FactorSecondBreakfast: m.FactorSecondBreakfast.Float64,
				FactorLunch: m.FactorLunch.Float64,
				FactorDinner: m.FactorDinner.Float64,
				FactorSupper: m.FactorSupper.Float64,
				FactorAfternoonSnack: m.FactorAfternoonSnack.Float64,
				Name: m.Name.String,
			})
		}
	}
	return res
}

func (mr *FirebirdRepoAccess) DeleteMealInDay(id int) bool {
	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", MEAL_IN_DAY_TAB, MEAL_IN_DAY_MEAL_ID)
	_, err := mr.DbEngine.Exec(sqlStr, id)
	if err != nil {
		log.Println("DeleteMealInDay error:", err)
		return false
	}
	return true
}

func (mr *FirebirdRepoAccess) UpdateMealInDay(m *MealInDay) {
	cols := []string{
		MEAL_IN_DAY_BREAKFAST, MEAL_IN_DAY_SECOND_BREAKFAST, MEAL_IN_DAY_LUNCH, MEAL_IN_DAY_DINNER, MEAL_IN_DAY_SUPPER, MEAL_IN_DAY_AFTERNOON_SNACK,
		MEAL_IN_DAY_FOR_5_DAYS, MEAL_IN_DAY_FACTOR_BREAKFAST, MEAL_IN_DAY_FACTOR_SECOND_BREAKFAST, MEAL_IN_DAY_FACTOR_LUNCH, MEAL_IN_DAY_FACTOR_DINNER, MEAL_IN_DAY_FACTOR_SUPPER, MEAL_IN_DAY_FACTOR_AFTERNOON_SNACK, MEAL_IN_DAY_NAME,
	}
	setExprs := make([]string, len(cols))
	for i, col := range cols {
		setExprs[i] = col + "=?"
	}
	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE %s=?", MEAL_IN_DAY_TAB, strings.Join(setExprs, ", "), MEAL_IN_DAY_MEAL_ID)
	for5DaysChar := sql.NullString{String: "0", Valid: true}
	if m.For5Days {
		for5DaysChar.String = "1"
	}
	_, err := mr.DbEngine.Exec(sqlStr,
		m.Breakfast.Id, m.SecondBreakfast.Id, m.Lunch.Id, m.Dinner.Id, m.Supper.Id, m.AfternoonSnack.Id,
		for5DaysChar, m.FactorBreakfast, m.FactorSecondBreakfast, m.FactorLunch, m.FactorDinner, m.FactorSupper, m.FactorAfternoonSnack, m.Name, m.Id)
	if err != nil {
		log.Println("UpdateMealInDay error:", err)
	}
}

// Mapowanie DB -> API
func (mr *FirebirdRepoAccess) ConvertToMealInDay(m MealInDayDb) MealInDay {
	return MealInDay{
		Id: int(m.Id.Int64),
		Breakfast: mr.GetMeal(int(m.BreakfastId.Int64)),
		SecondBreakfast: mr.GetMeal(int(m.SecondBreakfastId.Int64)),
		Lunch: mr.GetMeal(int(m.LunchId.Int64)),
		Dinner: mr.GetMeal(int(m.DinnerId.Int64)),
		Supper: mr.GetMeal(int(m.SupperId.Int64)),
		AfternoonSnack: mr.GetMeal(int(m.AfternoonSnackId.Int64)),
		For5Days: m.For5Days.String == "1",
		FactorBreakfast: m.FactorBreakfast.Float64,
		FactorSecondBreakfast: m.FactorSecondBreakfast.Float64,
		FactorLunch: m.FactorLunch.Float64,
		FactorDinner: m.FactorDinner.Float64,
		FactorSupper: m.FactorSupper.Float64,
		FactorAfternoonSnack: m.FactorAfternoonSnack.Float64,
		Name: m.Name.String,
	}
}