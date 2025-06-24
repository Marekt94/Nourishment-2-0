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

var	cols = []string{
		MEAL_IN_DAY_BREAKFAST, MEAL_IN_DAY_SECOND_BREAKFAST, MEAL_IN_DAY_LUNCH, MEAL_IN_DAY_DINNER, MEAL_IN_DAY_SUPPER, MEAL_IN_DAY_AFTERNOON_SNACK,
		MEAL_IN_DAY_FOR_5_DAYS, MEAL_IN_DAY_FACTOR_BREAKFAST, MEAL_IN_DAY_FACTOR_SECOND_BREAKFAST, MEAL_IN_DAY_FACTOR_LUNCH, MEAL_IN_DAY_FACTOR_DINNER, MEAL_IN_DAY_FACTOR_SUPPER, MEAL_IN_DAY_FACTOR_AFTERNOON_SNACK, MEAL_IN_DAY_NAME,
	}

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


func (mr *FirebirdRepoAccess) CreateMealsInDay(m *MealInDay) int64 {
	mealIds := []int{m.Breakfast.Id, m.SecondBreakfast.Id, m.Lunch.Id, m.Dinner.Id, m.Supper.Id, m.AfternoonSnack.Id}
	for _, mealId := range mealIds {
		if mealId > 0 && mr.GetMeal(mealId).Id == 0 {
			log.Println("CreateMealInDay error: referenced meal does not exist, id:", mealId)
			return -1
		}
	}
	sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", MEAL_IN_DAY_TAB, strings.Join(cols, ", "), QuestionMarks(len(cols)))
	for5DaysChar := sql.NullString{String: "0", Valid: true}
	if m.For5Days {
		for5DaysChar.String = "1"
	}
	_, err := mr.DbEngine.Exec(sqlStr,
		m.Breakfast.Id, m.SecondBreakfast.Id, m.Lunch.Id, m.Dinner.Id, m.Supper.Id, m.AfternoonSnack.Id,
		for5DaysChar, m.FactorBreakfast, m.FactorSecondBreakfast, m.FactorLunch, m.FactorDinner, m.FactorSupper,
		m.FactorAfternoonSnack, m.Name)
	if err != nil {
		log.Println("CreateMealInDay error:", err)
		return -1
	}
	var id int64
	err = mr.DbEngine.QueryRow(fmt.Sprintf("SELECT MAX(%s) FROM %s", MEAL_IN_DAY_ID, MEAL_IN_DAY_TAB)).Scan(&id)
	if err != nil {
		log.Println("CreateMealInDay get id error:", err)
		return -1
	}
	m.Id = int(id)
	return id
}

func (mr *FirebirdRepoAccess) ConvertMealsInDayDbToMealsInDay(m *MealInDayDb) MealInDay {	
	mealInDay := MealInDay{
		Id: int(m.Id.Int64),
		For5Days: m.For5Days.String == "1",
		FactorBreakfast: m.FactorBreakfast.Float64,
		FactorSecondBreakfast: m.FactorSecondBreakfast.Float64,
		FactorLunch: m.FactorLunch.Float64,
		FactorDinner: m.FactorDinner.Float64,
		FactorSupper: m.FactorSupper.Float64,
		FactorAfternoonSnack: m.FactorAfternoonSnack.Float64,
		Name: m.Name.String,
	}
	if m.BreakfastId.Valid {
		mealInDay.Breakfast = mr.GetMeal(int(m.BreakfastId.Int64))
	}
	if m.SecondBreakfastId.Valid {
		mealInDay.SecondBreakfast = mr.GetMeal(int(m.SecondBreakfastId.Int64))	
	}
	if m.LunchId.Valid {
		mealInDay.Lunch = mr.GetMeal(int(m.LunchId.Int64))
	}
	if m.DinnerId.Valid {
		mealInDay.Dinner = mr.GetMeal(int(m.DinnerId.Int64))
	}
	if m.SupperId.Valid {
		mealInDay.Supper = mr.GetMeal(int(m.SupperId.Int64))
	}
	if m.AfternoonSnackId.Valid {
		mealInDay.AfternoonSnack = mr.GetMeal(int(m.AfternoonSnackId.Int64))
	}
	return mealInDay
}

func returnMealsInDayDbFieldsToRetriveFromDb(m *MealInDayDb) (*sql.NullInt64, *sql.NullInt64, *sql.NullInt64,
		*sql.NullInt64, *sql.NullInt64, *sql.NullInt64, *sql.NullInt64, *sql.NullString, *sql.NullFloat64,
		*sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullString) {
	return &m.Id, &m.BreakfastId, &m.SecondBreakfastId, &m.LunchId, &m.DinnerId, &m.SupperId, &m.AfternoonSnackId,
		   &m.For5Days, &m.FactorBreakfast, &m.FactorSecondBreakfast, &m.FactorLunch, &m.FactorDinner, &m.FactorSupper,
		   &m.FactorAfternoonSnack, &m.Name
}

func (mr *FirebirdRepoAccess) GetMealsInDay(id int) MealInDay {
	cols := []string{
		MEAL_IN_DAY_ID, MEAL_IN_DAY_BREAKFAST, MEAL_IN_DAY_SECOND_BREAKFAST, MEAL_IN_DAY_LUNCH, MEAL_IN_DAY_DINNER, MEAL_IN_DAY_SUPPER, MEAL_IN_DAY_AFTERNOON_SNACK,
		MEAL_IN_DAY_FOR_5_DAYS, MEAL_IN_DAY_FACTOR_BREAKFAST, MEAL_IN_DAY_FACTOR_SECOND_BREAKFAST, MEAL_IN_DAY_FACTOR_LUNCH, MEAL_IN_DAY_FACTOR_DINNER, MEAL_IN_DAY_FACTOR_SUPPER, MEAL_IN_DAY_FACTOR_AFTERNOON_SNACK, MEAL_IN_DAY_NAME,
	}
	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ?", strings.Join(cols, ", "), MEAL_IN_DAY_TAB, MEAL_IN_DAY_ID)
	row := mr.DbEngine.QueryRow(sqlStr, id)
	var m MealInDayDb
	err := row.Scan(returnMealsInDayDbFieldsToRetriveFromDb(&m))
	if err != nil {
		log.Println("GetMealInDay error:", err)
		return MealInDay{}
	}
	return mr.ConvertMealsInDayDbToMealsInDay(&m)
}

func (mr *FirebirdRepoAccess) GetMealsInDays() []MealInDay {
	cols := []string{
		MEAL_IN_DAY_ID, MEAL_IN_DAY_BREAKFAST, MEAL_IN_DAY_SECOND_BREAKFAST, MEAL_IN_DAY_LUNCH, MEAL_IN_DAY_DINNER, MEAL_IN_DAY_SUPPER, MEAL_IN_DAY_AFTERNOON_SNACK,
		MEAL_IN_DAY_FOR_5_DAYS, MEAL_IN_DAY_FACTOR_BREAKFAST, MEAL_IN_DAY_FACTOR_SECOND_BREAKFAST, MEAL_IN_DAY_FACTOR_LUNCH, MEAL_IN_DAY_FACTOR_DINNER, MEAL_IN_DAY_FACTOR_SUPPER, MEAL_IN_DAY_FACTOR_AFTERNOON_SNACK, MEAL_IN_DAY_NAME,
	}
	sqlStr := fmt.Sprintf("SELECT %s FROM %s ORDER BY %s ASC", strings.Join(cols, ", "), MEAL_IN_DAY_TAB, MEAL_IN_DAY_ID)
	rows, err := mr.DbEngine.Query(sqlStr)
	if err != nil {
		log.Println("GetMealsInDay error:", err)
		return nil
	}
	var res []MealInDay
	for rows.Next() {
		var m MealInDayDb
		err := rows.Scan(returnMealsInDayDbFieldsToRetriveFromDb(&m))
		if err == nil {
			mealInDay := mr.ConvertMealsInDayDbToMealsInDay(&m)
			res = append(res, mealInDay)
		} else {
			log.Fatalln(err)
		}
	}
	return res
}

func (mr *FirebirdRepoAccess) DeleteMealsInDay(id int) bool {
	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", MEAL_IN_DAY_TAB, MEAL_IN_DAY_ID)
	_, err := mr.DbEngine.Exec(sqlStr, id)
	if err != nil {
		log.Println("DeleteMealInDay error:", err)
		return false
	}
	return true
}

func (mr *FirebirdRepoAccess) UpdateMealsInDay(m *MealInDay) {
	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE %s=?", MEAL_IN_DAY_TAB, QuestionMarks(len(cols)), MEAL_IN_DAY_ID)
	for5DaysChar := sql.NullString{String: "0", Valid: true}
	if m.For5Days {
		for5DaysChar.String = "1"
	}
	_, err := mr.DbEngine.Exec(sqlStr,
		m.Breakfast.Id, m.SecondBreakfast.Id, m.Lunch.Id, m.Dinner.Id, m.Supper.Id, m.AfternoonSnack.Id,
		for5DaysChar, m.FactorBreakfast, m.FactorSecondBreakfast, m.FactorLunch, m.FactorDinner, m.FactorSupper,
		m.FactorAfternoonSnack, m.Name, m.Id)
	if err != nil {
		log.Println("UpdateMealInDay error:", err)
	}
}