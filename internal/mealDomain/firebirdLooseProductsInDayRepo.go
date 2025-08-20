package meal

import (
	"database/sql"
	"fmt"
	"nourishment_20/internal/database"
	"nourishment_20/internal/logging"
	"strings"
)

type LooseProductsInDayDb struct {
	Id      sql.NullInt64
	DayId   sql.NullInt64
	Product productDb
	Weight  sql.NullFloat64
}

func (p *LooseProductsInDayDb) ConvertToLooseProductInDay() LooseProductInDay {
	var product LooseProductInDay
	product.Id = database.NullInt64ToInt(&p.Id)
	product.DayId = database.NullInt64ToInt(&p.DayId)
	p.Product.ConvertToProduct(&product.Product)
	product.Weight = database.NullFloat64ToFloat(&p.Weight)
	return product
}

func returnLooseProductFields(p *LooseProductsInDayDb) (
	*sql.NullInt64, // pld.ID
	*sql.NullInt64, // pld.ID_DNIA
	*sql.NullFloat64, // pld.ILOSC_W_G (weight)
	*sql.NullInt64, // p.ID
	*sql.NullString, // p.NAZWA
	*sql.NullFloat64, // p.KCAL_NA_100G
	*sql.NullFloat64, // p.WAGA_JEDNOSTKI
	*sql.NullFloat64, // p.BIALKO
	*sql.NullFloat64, // p.TLUSZCZ
	*sql.NullFloat64, // p.CUKRY_PROSTE
	*sql.NullFloat64, // p.CUKRY_ZLOZONE
	*sql.NullFloat64, // p.CUKRY_SUMA
	*sql.NullFloat64, // p.BLONNIK
	*sql.NullFloat64, // p.SOL
	*sql.NullString, // p.JEDNOSTKA
	*sql.NullInt64, // c.ID
	*sql.NullString, // c.NAZWA_KATEGORII
) {
	return &p.Id, &p.DayId, &p.Weight,
		&p.Product.Id, &p.Product.Name, &p.Product.KcalPer100, &p.Product.UnitWeight, &p.Product.Proteins,
		&p.Product.Fat, &p.Product.Sugar, &p.Product.Carbohydrates, &p.Product.SugarAndCarbo, &p.Product.Fiber,
		&p.Product.Salt, &p.Product.Unit, &p.Product.Category.Id, &p.Product.Category.Name
}

func (repo *FirebirdRepoAccess) CreateLooseProductInDay(p *LooseProductInDay) int64 {
	// Use RETURNING to reliably capture the newly inserted ID
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s) VALUES (?, ?, ?) RETURNING %s`,
		LOOSE_PRODUCTS_IN_DAY_TAB,
		LOOSE_PRODUCTS_IN_DAY_DAY_ID,
		LOOSE_PRODUCTS_IN_DAY_PRODUCT_ID,
		LOOSE_PRODUCTS_IN_DAY_WEIGHT,
		LOOSE_PRODUCTS_IN_DAY_ID)
	var id int64
	if err := repo.Database.QueryRow(query, p.DayId, p.Product.Id, p.Weight).Scan(&id); err != nil {
		logging.Global.Panicf("%v", err)
		return -1
	}
	p.Id = int(id)
	return id
}

func (repo *FirebirdRepoAccess) GetLooseProductInDay(id int) LooseProductInDay {
	sqlStr := generateGetLooseProductsQuery()
	sqlStr = fmt.Sprintf("%s WHERE pld.ID = ?", sqlStr)
	logging.Global.Debugf("SQL: %s", sqlStr)

	row := repo.Database.QueryRow(sqlStr, id)
	var product LooseProductsInDayDb
	if err := row.Scan(returnLooseProductFields(&product)); err != nil {
		logging.Global.Panicf("%v", err)
	}
	return product.ConvertToLooseProductInDay()
}

func (repo *FirebirdRepoAccess) GetLooseProductsInDay(dayId int) []LooseProductInDay {
	var products []LooseProductInDay
	sqlStr := generateGetLooseProductsQuery()
	sqlStr = fmt.Sprintf("%s WHERE pld.ID_DNIA = ?", sqlStr)
	logging.Global.Debugf("SQL: %s", sqlStr)

	rows, err := repo.Database.Query(sqlStr, dayId)
	if err != nil {
		logging.Global.Panicf("%v", err)
	}
	for rows.Next() {
		var product LooseProductsInDayDb
		if err := rows.Scan(returnLooseProductFields(&product)); err != nil {
			logging.Global.Panicf("%v", err)
		}
		products = append(products, product.ConvertToLooseProductInDay())
	}
	return products
}

func (repo *FirebirdRepoAccess) UpdateLooseProductInDay(p *LooseProductInDay) {
	sql := fmt.Sprintf(`UPDATE %s SET %s=?, %s=?, %s=? WHERE %s=?`,
		LOOSE_PRODUCTS_IN_DAY_TAB,
		LOOSE_PRODUCTS_IN_DAY_DAY_ID,
		LOOSE_PRODUCTS_IN_DAY_PRODUCT_ID,
		LOOSE_PRODUCTS_IN_DAY_WEIGHT,
		LOOSE_PRODUCTS_IN_DAY_ID)
	if _, err := repo.Database.Exec(sql, p.DayId, p.Product.Id, p.Weight, p.Id); err != nil {
		logging.Global.Panicf("%v", err)
	}
}

func (repo *FirebirdRepoAccess) DeleteLooseProductInDay(id int) bool {
	if _, err := repo.Database.Exec(fmt.Sprintf(`DELETE FROM %s WHERE %s = ?`,
		LOOSE_PRODUCTS_IN_DAY_TAB, LOOSE_PRODUCTS_IN_DAY_ID), id); err != nil {
		logging.Global.Panicf("%v", err)
	}
	row := repo.Database.QueryRow(fmt.Sprintf(`SELECT %s FROM %s WHERE %s = ?`,
		LOOSE_PRODUCTS_IN_DAY_ID, LOOSE_PRODUCTS_IN_DAY_TAB, LOOSE_PRODUCTS_IN_DAY_ID), id)
	var checkId int
	return row.Scan(&checkId) == sql.ErrNoRows
}

func generateGetLooseProductsQuery() string {
	colsForLooseProduct := []string{LOOSE_PRODUCTS_IN_DAY_ID, LOOSE_PRODUCTS_IN_DAY_DAY_ID, LOOSE_PRODUCTS_IN_DAY_WEIGHT}
	colsForProduct := ProductTabs
	colsForCategory := CategoryTabs

	colsLooseProductStr := database.CreateColsToSelect("pld", colsForLooseProduct)
	colsProductStr := database.CreateColsToSelect("p", colsForProduct)
	colsCategoryStr := database.CreateColsToSelect("c", colsForCategory)

	cols := strings.Join([]string{colsLooseProductStr, colsProductStr, colsCategoryStr}, ", ")

	return fmt.Sprintf(`SELECT %s FROM %s pld 
		LEFT JOIN %s p ON p.%s = pld.%s 
		LEFT JOIN %s c ON c.%s = p.%s`,
		cols,
		LOOSE_PRODUCTS_IN_DAY_TAB,
		PRODUCT_TAB, PRODUCT_ID, LOOSE_PRODUCTS_IN_DAY_PRODUCT_ID,
		CATEGORY_TAB, CATEGORY_ID, PRODUCT_CATEGORY)
}
