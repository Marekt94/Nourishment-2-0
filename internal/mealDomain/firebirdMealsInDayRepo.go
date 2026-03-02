package meal

import (
	"database/sql"
	"fmt"
	"nourishment_20/internal/database"
	"nourishment_20/internal/logging"
	"strings"
)

// [API GEN] DTO do tabeli POTRAWY_W_DNIU
// Odpowiada kolumnom zdefiniowanym w firebirdDatabase.go
// Każde pole odpowiada jednej kolumnie w bazie

type MealInDayDb struct { // [API GEN]
	Id                    sql.NullInt64
	BreakfastId           sql.NullInt64
	SecondBreakfastId     sql.NullInt64
	LunchId               sql.NullInt64
	DinnerId              sql.NullInt64
	SupperId              sql.NullInt64
	AfternoonSnackId      sql.NullInt64
	For5Days              sql.NullString // CHAR(1) '1' lub '0'
	FactorBreakfast       sql.NullFloat64
	FactorSecondBreakfast sql.NullFloat64
	FactorLunch           sql.NullFloat64
	FactorDinner          sql.NullFloat64
	FactorSupper          sql.NullFloat64
	FactorAfternoonSnack  sql.NullFloat64
	Name                  sql.NullString
}

func (mr *FirebirdRepoAccess) CreateMealsInDay(m *MealInDay) int64 {
	// Najpierw sprawdź czy wszystkie niepuste posiłki istnieją w bazie
	mealIds := []int{}
	if m.Breakfast.Id > EMPTY_ID {
		mealIds = append(mealIds, m.Breakfast.Id)
	}
	if m.SecondBreakfast.Id > EMPTY_ID {
		mealIds = append(mealIds, m.SecondBreakfast.Id)
	}
	if m.Lunch.Id > EMPTY_ID {
		mealIds = append(mealIds, m.Lunch.Id)
	}
	if m.Dinner.Id > EMPTY_ID {
		mealIds = append(mealIds, m.Dinner.Id)
	}
	if m.Supper.Id > EMPTY_ID {
		mealIds = append(mealIds, m.Supper.Id)
	}
	if m.AfternoonSnack.Id > EMPTY_ID {
		mealIds = append(mealIds, m.AfternoonSnack.Id)
	}

	for _, mealId := range mealIds {
		if mealId > 0 && mr.GetMeal(mealId).Id == 0 {
			logging.Global.Debugf("CreateMealInDay error: referenced meal does not exist, id: %v", mealId)
			return -1
		}
	}

	// Zawsze wstawiaj wszystkie kolumny - ustaw NULL dla EMPTY_ID, wartość dla pozostałych
	cols := []string{
		MEAL_IN_DAY_BREAKFAST, MEAL_IN_DAY_SECOND_BREAKFAST, MEAL_IN_DAY_LUNCH,
		MEAL_IN_DAY_DINNER, MEAL_IN_DAY_SUPPER, MEAL_IN_DAY_AFTERNOON_SNACK,
		MEAL_IN_DAY_FOR_5_DAYS, MEAL_IN_DAY_FACTOR_BREAKFAST, MEAL_IN_DAY_FACTOR_SECOND_BREAKFAST,
		MEAL_IN_DAY_FACTOR_LUNCH, MEAL_IN_DAY_FACTOR_DINNER, MEAL_IN_DAY_FACTOR_SUPPER,
		MEAL_IN_DAY_FACTOR_AFTERNOON_SNACK, MEAL_IN_DAY_NAME,
	}

	for5DaysChar := sql.NullString{String: "0", Valid: true}
	if m.For5Days {
		for5DaysChar.String = "1"
	}

	args := []interface{}{}

	// Dodaj ID posiłków lub NULL
	if m.Breakfast.Id > EMPTY_ID {
		args = append(args, m.Breakfast.Id)
	} else {
		args = append(args, nil)
	}

	if m.SecondBreakfast.Id > EMPTY_ID {
		args = append(args, m.SecondBreakfast.Id)
	} else {
		args = append(args, nil)
	}

	if m.Lunch.Id > EMPTY_ID {
		args = append(args, m.Lunch.Id)
	} else {
		args = append(args, nil)
	}

	if m.Dinner.Id > EMPTY_ID {
		args = append(args, m.Dinner.Id)
	} else {
		args = append(args, nil)
	}

	if m.Supper.Id > EMPTY_ID {
		args = append(args, m.Supper.Id)
	} else {
		args = append(args, nil)
	}

	if m.AfternoonSnack.Id > EMPTY_ID {
		args = append(args, m.AfternoonSnack.Id)
	} else {
		args = append(args, nil)
	}

	args = append(args, for5DaysChar, m.FactorBreakfast, m.FactorSecondBreakfast, m.FactorLunch, m.FactorDinner, m.FactorSupper,
		m.FactorAfternoonSnack, m.Name)

	sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", MEAL_IN_DAY_TAB, strings.Join(cols, ", "), database.QuestionMarks(len(cols)))
	_, err := mr.Database.Exec(sqlStr, args...)
	if err != nil {
		logging.Global.Debugf("CreateMealInDay error: %v", err)
		return -1
	}
	var id int64
	err = mr.Database.QueryRow(fmt.Sprintf("SELECT MAX(%s) FROM %s", MEAL_IN_DAY_ID, MEAL_IN_DAY_TAB)).Scan(&id)
	if err != nil {
		logging.Global.Debugf("CreateMealInDay get id error: %v", err)
		return -1
	}
	m.Id = int(id)
	return id
}

func (mr *FirebirdRepoAccess) ConvertMealsInDayDbToMealsInDay(m *MealInDayDb) MealInDay {
	mealInDay := MealInDay{
		Id:                    int(m.Id.Int64),
		For5Days:              m.For5Days.String == "1",
		FactorBreakfast:       m.FactorBreakfast.Float64,
		FactorSecondBreakfast: m.FactorSecondBreakfast.Float64,
		FactorLunch:           m.FactorLunch.Float64,
		FactorDinner:          m.FactorDinner.Float64,
		FactorSupper:          m.FactorSupper.Float64,
		FactorAfternoonSnack:  m.FactorAfternoonSnack.Float64,
		Name:                  m.Name.String,
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
	row := mr.Database.QueryRow(sqlStr, id)
	var m MealInDayDb
	err := row.Scan(returnMealsInDayDbFieldsToRetriveFromDb(&m))
	if err != nil {
		logging.Global.Debugf("GetMealInDay error: %v", err)
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
	rows, err := mr.Database.Query(sqlStr)
	if err != nil {
		logging.Global.Debugf("GetMealsInDay error: %v", err)
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
			logging.Global.Panicf("%v", err)
		}
	}
	return res
}

func (mr *FirebirdRepoAccess) DeleteMealsInDay(id int) bool {
	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", MEAL_IN_DAY_TAB, MEAL_IN_DAY_ID)
	_, err := mr.Database.Exec(sqlStr, id)
	if err != nil {
		logging.Global.Debugf("DeleteMealInDay error: %v", err)
		return false
	}
	return true
}

func (mr *FirebirdRepoAccess) UpdateMealsInDay(m *MealInDay) {
	cols := []string{}
	args := []interface{}{}

	// Zawsze aktualizuj wszystkie kolumny posiłków - ustaw NULL dla EMPTY_ID, wartość dla pozostałych
	cols = append(cols, MEAL_IN_DAY_BREAKFAST)
	if m.Breakfast.Id > EMPTY_ID {
		args = append(args, m.Breakfast.Id)
	} else {
		args = append(args, nil) // NULL w bazie danych
	}

	cols = append(cols, MEAL_IN_DAY_SECOND_BREAKFAST)
	if m.SecondBreakfast.Id > EMPTY_ID {
		args = append(args, m.SecondBreakfast.Id)
	} else {
		args = append(args, nil) // NULL w bazie danych
	}

	cols = append(cols, MEAL_IN_DAY_LUNCH)
	if m.Lunch.Id > EMPTY_ID {
		args = append(args, m.Lunch.Id)
	} else {
		args = append(args, nil) // NULL w bazie danych
	}

	cols = append(cols, MEAL_IN_DAY_DINNER)
	if m.Dinner.Id > EMPTY_ID {
		args = append(args, m.Dinner.Id)
	} else {
		args = append(args, nil) // NULL w bazie danych
	}

	cols = append(cols, MEAL_IN_DAY_SUPPER)
	if m.Supper.Id > EMPTY_ID {
		args = append(args, m.Supper.Id)
	} else {
		args = append(args, nil) // NULL w bazie danych
	}

	cols = append(cols, MEAL_IN_DAY_AFTERNOON_SNACK)
	if m.AfternoonSnack.Id > EMPTY_ID {
		args = append(args, m.AfternoonSnack.Id)
	} else {
		args = append(args, nil) // NULL w bazie danych
	}

	for5DaysChar := sql.NullString{String: "0", Valid: true}
	if m.For5Days {
		for5DaysChar.String = "1"
	}

	cols = append(cols, MEAL_IN_DAY_FOR_5_DAYS, MEAL_IN_DAY_FACTOR_BREAKFAST, MEAL_IN_DAY_FACTOR_SECOND_BREAKFAST,
		MEAL_IN_DAY_FACTOR_LUNCH, MEAL_IN_DAY_FACTOR_DINNER, MEAL_IN_DAY_FACTOR_SUPPER,
		MEAL_IN_DAY_FACTOR_AFTERNOON_SNACK, MEAL_IN_DAY_NAME)
	args = append(args, for5DaysChar, m.FactorBreakfast, m.FactorSecondBreakfast, m.FactorLunch, m.FactorDinner,
		m.FactorSupper, m.FactorAfternoonSnack, m.Name)

	args = append(args, m.Id)

	sqlStr := fmt.Sprintf("UPDATE %s SET %s WHERE %s=?", MEAL_IN_DAY_TAB, database.UpdateValues(cols), MEAL_IN_DAY_ID)
	_, err := mr.Database.Exec(sqlStr, args...)
	if err != nil {
		logging.Global.Debugf("UpdateMealInDay error: %v", err)
	}
}
