package database

import (
	"database/sql"
	"fmt"
	"nourishment_20/internal/logging"
	"strings"
)

const MealPrefix = `m`
const ProductPrefix = `p`
const ProductInMealPrefix = `pm`
const CategoryPrefix = `c`

type FirebirdRepoAccess struct {
	Database *sql.DB
}

type MealDb struct {
	Id            sql.NullInt64
	Name          sql.NullString
	Recipe        sql.NullString
	ProductInMeal productInMealDb
}

type productInMealDb struct {
	Id      sql.NullInt64
	Product productDb
	Weight  sql.NullFloat64
}

func returnProdInMealFieldsForDbRetriving(pm *productInMealDb) (*sql.NullInt64, *sql.NullFloat64, *sql.NullInt64, *sql.NullString, *sql.NullFloat64, *sql.NullFloat64,
	*sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64,
	*sql.NullString, *sql.NullInt64, *sql.NullString) {
	id, name, kcalPer100, unitWeight, proteins, fat, sugar, carbohydrates, sugarAndCarno, fiber, salt, unit, categoryId,
		categoryName := ReturnProductFieldsForDbRetriving(&pm.Product)
	return &pm.Id, &pm.Weight, id, name, kcalPer100, unitWeight, proteins, fat, sugar, carbohydrates, sugarAndCarno, fiber,
		salt, unit, categoryId, categoryName
}

func returnMealFieldsForDbRetriving(m *MealDb) (*sql.NullInt64, *sql.NullString, *sql.NullString, *sql.NullInt64, *sql.NullFloat64, *sql.NullInt64, *sql.NullString, *sql.NullFloat64, *sql.NullFloat64,
	*sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64,
	*sql.NullString, *sql.NullInt64, *sql.NullString) {
	idProdInMeal, weight, idProd, name, kcalPer100, unitWeight, proteins, fat, sugar, carbohydrates, sugarAndCarno, fiber, salt, unit, categoryId,
		categoryName := returnProdInMealFieldsForDbRetriving(&m.ProductInMeal)
	return &m.Id, &m.Name, &m.Recipe, idProdInMeal, weight, idProd, name, kcalPer100, unitWeight, proteins, fat, sugar, carbohydrates, sugarAndCarno, fiber, salt, unit, categoryId,
		categoryName
}

func generateGetMealsQuery() string {
	colsForMeal := []string{`ID`, `NAZWA`, `PRZEPIS`}
	colsForProduct := ProductTabs
	colsForCategory := CategoryTabs
	colsForProductInMeal := []string{PRODUCTS_IN_MEAL_ID, PRODUCTS_IN_MEAL_WEIGHT}

	colsForMealStr := CreateColsToSelect(MealPrefix, colsForMeal)
	colsForProductStr := CreateColsToSelect(ProductPrefix, colsForProduct)
	colsForProductInMealStr := CreateColsToSelect(ProductInMealPrefix, colsForProductInMeal)
	colsForCategoryStr := CreateColsToSelect(CategoryPrefix, colsForCategory)
	colsToRetrive := strings.Join([]string{colsForMealStr, colsForProductInMealStr, colsForProductStr, colsForCategoryStr}, `, `)
	logging.Global.Debugf(`cols to retive %s`, colsToRetrive)

	sql := "SELECT %s FROM %s LEFT JOIN %s ON %s=%s LEFT JOIN %s ON %s=%s LEFT JOIN %s ON %s=%s"
	return fmt.Sprintf(sql, colsToRetrive, MEAL_TAB+` `+MealPrefix, PRODUCTS_IN_MEAL_TAB+` `+ProductInMealPrefix,
		MealPrefix+`.`+MEAL_ID, ProductInMealPrefix+`.`+PRODUCTS_IN_MEAL_MEAL_ID,
		PRODUCT_TAB+` `+ProductPrefix, ProductPrefix+`.`+PRODUCT_ID,
		ProductInMealPrefix+`.`+PRODUCTS_IN_MEAL_PRODUCT_ID,
		CATEGORY_TAB+` `+CategoryPrefix,
		CategoryPrefix+`.`+CATEGORY_ID,
		ProductPrefix+`.`+PRODUCT_CATEGORY)
}

func ConvertToMeals(m []MealDb) []Meal { // [AI] poprawka: []Meal zamiast []meal
	var res []Meal
	if len(m) < 1 {
		return nil
	}
	initMeal := m[0]
	initMealsDb := m[1:]
	mealsDb := m[:1]

	for _, lMealDb := range initMealsDb {
		if lMealDb.Id == initMeal.Id {
			mealsDb = append(mealsDb, lMealDb)
		} else {
			lMeal := ConvertToMeal(mealsDb)
			res = append(res, lMeal)
			initMeal = lMealDb
			mealsDb = []MealDb{lMealDb}
		}
	}
	lMeal := ConvertToMeal(mealsDb)
	res = append(res, lMeal)
	return res
}

func ConvertToMeal(m []MealDb) Meal { // [AI] poprawka: Meal zamiast meal
	var meal Meal
	logging.Global.Debugf(`start converting db to meal`)
	meal.Id = NullInt64ToInt(&m[0].Id)
	meal.Name = NullStringToString(&m[0].Name)
	meal.Recipe = NullStringToString(&m[0].Recipe)
	for _, pml := range m {
		if pml.ProductInMeal.Id.Valid {
			var pm ProductInMeal
			pm.Id = NullInt64ToInt(&pml.ProductInMeal.Id)
			pm.Weight = NullFloat64ToFloat(&pml.ProductInMeal.Weight)
			pml.ProductInMeal.Product.ConvertToProduct(&pm.Product)
			meal.ProductsInMeal = append(meal.ProductsInMeal, pm)
		}
	}
	logging.Global.Debugf(`finish converting db to meal`)
	return meal
}

func (mr *FirebirdRepoAccess) GetMeal(i int) Meal { // [AI] poprawka: Meal zamiast meal
	sqlStr := generateGetMealsQuery()
	sqlStr = fmt.Sprintf(sqlStr+` WHERE %s = ?`, MealPrefix+`.`+MEAL_ID)
	logging.Global.Debugf(`SQL: %s`, sqlStr)

	var meals []MealDb
	rows, err := mr.Database.Query(sqlStr, i)
	if err != nil {
		logging.Global.Panicf("%v", err)
	}
	for rows.Next() {
		var meal MealDb
		err := rows.Scan(returnMealFieldsForDbRetriving(&meal))
		logging.Global.Debugf("%v", meal)
		if err == nil {
			meals = append(meals, meal)
		} else {
			logging.Global.Panicf("%v", err)
		}
	}
	res := ConvertToMeal(meals)
	logging.Global.Debugf("%v", res)
	return res
}

func (mr *FirebirdRepoAccess) GetMeals() []Meal { // [AI] poprawka: []Meal zamiast []meal
	var meals []Meal
	sqlStr := generateGetMealsQuery()
	sqlStr = fmt.Sprintf(sqlStr+` ORDER BY %s`, MealPrefix+`.`+MEAL_ID)
	rows, err := mr.Database.Query(sqlStr)
	if err != nil {
		logging.Global.Panicf("%v", err)
	} else {
		var mealsDb []MealDb
		for rows.Next() {
			var mealDb MealDb
			if err := rows.Scan(returnMealFieldsForDbRetriving(&mealDb)); err == nil {
				mealsDb = append(mealsDb, mealDb)
			} else {
				logging.Global.Panicf("%v", err)
			}
		}
		meals = ConvertToMeals(mealsDb)
	}
	for i, meal := range meals {
		logging.Global.Debugf(`%d: %v`, i, meal)
	}
	return meals
}

func (mr *FirebirdRepoAccess) DeleteMeal(i int) bool {
	if _, err := mr.Database.Exec(`DELETE FROM `+MEAL_TAB+` WHERE ID = ?`, i); err != nil {
		logging.Global.Panicf("%v", err)
	}
	row, err := mr.Database.Query(`SELECT ID FROM `+MEAL_TAB+` WHERE ID = ?`, i)
	if err != nil {
		logging.Global.Panicf("%v", err)
	}
	return !row.Next()
}

func (mr *FirebirdRepoAccess) updateProductsInMeal(m *Meal, r ProductsRepo) { // [AI] poprawka: *Meal zamiast *meal
	for i, prod := range m.ProductsInMeal {
		if prod.Product.Id <= 0 {
			id := r.CreateProduct(&prod.Product)
			prod.Product.Id = int(id)
			m.ProductsInMeal[i].Product.Id = int(id)
		} else {
			r.UpdateProduct(&prod.Product)
		}

		if prod.Id <= 0 {
			prodInMealDb := productInMealDb{Product: productDb{Id: sql.NullInt64{Valid: true, Int64: int64(prod.Product.Id)}},
				Weight: sql.NullFloat64{Float64: prod.Weight, Valid: true}}
			prodInDbTabs := []string{PRODUCTS_IN_MEAL_PRODUCT_ID, PRODUCTS_IN_MEAL_MEAL_ID, PRODUCTS_IN_MEAL_WEIGHT}
			sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, PRODUCTS_IN_MEAL_TAB,
				strings.Join(prodInDbTabs[:], `, `),
				QuestionMarks(len(prodInDbTabs)))
			_, err := mr.Database.Exec(sql, &prodInMealDb.Product.Id, m.Id, prodInMealDb.Weight)
			if err != nil {
				logging.Global.Panicf("%v", err)
			}

			sql = fmt.Sprintf(`SELECT MAX(%s) FROM `+PRODUCTS_IN_MEAL_TAB, PRODUCTS_IN_MEAL_ID)
			row := mr.Database.QueryRow(sql)
			var id int
			row.Scan(&id)
			prod.Id = id
			m.ProductsInMeal[i].Id = id
		}
	}

	//delete producte from db
	tabs := []string{PRODUCTS_IN_MEAL_ID, PRODUCTS_IN_MEAL_PRODUCT_ID}
	res, err := mr.Database.Query(`SELECT `+strings.Join(tabs, ", ")+` FROM `+PRODUCTS_IN_MEAL_TAB+` WHERE `+PRODUCTS_IN_MEAL_MEAL_ID+` = ?`, m.Id)
	if err != nil {
		logging.Global.Panicf("%v", err)
	}

	prodInMealIds := make(map[int]int)
	var key, value sql.NullInt64
	for res.Next() {
		err := res.Scan(&value, &key)
		if err != nil {
			logging.Global.Panicf("%v", err)
		}
		prodInMealIds[int(key.Int64)] = int(value.Int64)
	}
	for _, el := range m.ProductsInMeal {
		delete(prodInMealIds, el.Product.Id)
	}

	for _, v := range prodInMealIds {
		mr.Database.Exec(`DELETE FROM `+PRODUCTS_IN_MEAL_TAB+` WHERE `+PRODUCTS_IN_MEAL_ID+` = ?`, v)
	}
}

func (mr *FirebirdRepoAccess) CreateMeal(m *Meal) int64 { // [AI] poprawka: *Meal zamiast *meal
	if _, err := mr.Database.Exec(`INSERT INTO `+MEAL_TAB+` (`+MEAL_NAME+`, `+MEAL_RECIPE+`) `+`VALUES (?, ?)`, m.Name, m.Recipe); err != nil {
		logging.Global.Panicf("%v", err)
	} else {
		query := fmt.Sprintf(`SELECT MAX(%s) FROM %s`, MEAL_ID, MEAL_TAB)
		mr.Database.QueryRow(query).Scan(&m.Id)

		r, supp := interface{}(mr).(ProductsRepo)
		if !supp {
			logging.Global.Panicf(`object does not support ProductRepo interface`)
		}
		mr.updateProductsInMeal(m, r)
		return int64(m.Id)
	}
	return -1
}

func (mr *FirebirdRepoAccess) UpdateMeal(m *Meal) { // [AI] poprawka: *Meal zamiast *meal
	sql := fmt.Sprintf(`UPDATE %s SET %s=?, %s=? WHERE ID=?`, MEAL_TAB, MEAL_NAME, MEAL_RECIPE)
	logging.Global.Debugf(sql)

	_, err := mr.Database.Exec(sql, m.Name, m.Recipe, m.Id)
	if err == nil {
		r, supp := interface{}(mr).(ProductsRepo)
		if !supp {
			logging.Global.Panicf(`object does not support ProductRepo interface`)
		}
		mr.updateProductsInMeal(m, r)
	} else {
		logging.Global.Panicf("%v", err)
	}
}
