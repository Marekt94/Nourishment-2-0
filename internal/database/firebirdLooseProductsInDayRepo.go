package database

import (
	"database/sql"
	"fmt"
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
	product.Id = NullInt64ToInt(&p.Id)
	product.DayId = NullInt64ToInt(&p.DayId)
	p.Product.ConvertToProduct(&product.Product)
	product.Weight = NullFloat64ToFloat(&p.Weight)
	return product
}

func returnLooseProductFields(p *LooseProductsInDayDb) (*sql.NullInt64, *sql.NullInt64, *sql.NullInt64, *sql.NullString,
	*sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64,
	*sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullString,
	*sql.NullInt64, *sql.NullString) {
	return &p.Id, &p.DayId, &p.Product.Id, &p.Product.Name, &p.Product.KcalPer100,
		&p.Product.UnitWeight, &p.Product.Proteins, &p.Product.Fat, &p.Product.Sugar,
		&p.Product.Carbohydrates, &p.Product.SugarAndCarbo, &p.Product.Fiber, &p.Weight,
		&p.Product.Unit, &p.Product.Category.Id, &p.Product.Category.Name
}

func (repo *FirebirdRepoAccess) CreateLooseProductInDay(p *LooseProductInDay) int64 {
	sql := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s) VALUES (?, ?, ?)`,
		LOOSE_PRODUCTS_IN_DAY_TAB,
		LOOSE_PRODUCTS_IN_DAY_DAY_ID,
		LOOSE_PRODUCTS_IN_DAY_PRODUCT_ID,
		LOOSE_PRODUCTS_IN_DAY_WEIGHT)
	if _, err := repo.DbEngine.Exec(sql, p.DayId, p.Product.Id, p.Weight); err != nil {
		logging.Global.Panicf("%v", err)
	} else {
		query := fmt.Sprintf(`SELECT MAX(%s) FROM %s`, LOOSE_PRODUCTS_IN_DAY_ID, LOOSE_PRODUCTS_IN_DAY_TAB)
		repo.DbEngine.QueryRow(query).Scan(&p.Id)
		return int64(p.Id)
	}
	return -1
}

func (repo *FirebirdRepoAccess) GetLooseProductInDay(id int) LooseProductInDay {
	sqlStr := generateGetLooseProductsQuery()
	sqlStr = fmt.Sprintf(sqlStr+` WHERE pld.ID = ?`, id)
	logging.Global.Debugf("SQL: %s", sqlStr)

	row := repo.DbEngine.QueryRow(sqlStr, id)
	var product LooseProductsInDayDb
	if err := row.Scan(returnLooseProductFields(&product)); err != nil {
		logging.Global.Panicf("%v", err)
	}
	return product.ConvertToLooseProductInDay()
}

func (repo *FirebirdRepoAccess) GetLooseProductsInDay(dayId int) []LooseProductInDay {
	var products []LooseProductInDay
	sqlStr := generateGetLooseProductsQuery()
	sqlStr = fmt.Sprintf(sqlStr + ` WHERE pld.ID_DNIA = ?`)
	logging.Global.Debugf("SQL: %s", sqlStr)

	rows, err := repo.DbEngine.Query(sqlStr, dayId)
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
	if _, err := repo.DbEngine.Exec(sql, p.DayId, p.Product.Id, p.Weight, p.Id); err != nil {
		logging.Global.Panicf("%v", err)
	}
}

func (repo *FirebirdRepoAccess) DeleteLooseProductInDay(id int) bool {
	if _, err := repo.DbEngine.Exec(fmt.Sprintf(`DELETE FROM %s WHERE %s = ?`,
		LOOSE_PRODUCTS_IN_DAY_TAB, LOOSE_PRODUCTS_IN_DAY_ID), id); err != nil {
		logging.Global.Panicf("%v", err)
	}
	row := repo.DbEngine.QueryRow(fmt.Sprintf(`SELECT %s FROM %s WHERE %s = ?`,
		LOOSE_PRODUCTS_IN_DAY_ID, LOOSE_PRODUCTS_IN_DAY_TAB, LOOSE_PRODUCTS_IN_DAY_ID), id)
	var checkId int
	return row.Scan(&checkId) == sql.ErrNoRows
}

func generateGetLooseProductsQuery() string {
	colsForLooseProduct := []string{LOOSE_PRODUCTS_IN_DAY_ID, LOOSE_PRODUCTS_IN_DAY_DAY_ID, LOOSE_PRODUCTS_IN_DAY_WEIGHT}
	colsForProduct := ProductTabs
	colsForCategory := CategoryTabs

	colsLooseProductStr := CreateColsToSelect("pld", colsForLooseProduct)
	colsProductStr := CreateColsToSelect("p", colsForProduct)
	colsCategoryStr := CreateColsToSelect("c", colsForCategory)

	cols := strings.Join([]string{colsLooseProductStr, colsProductStr, colsCategoryStr}, ", ")

	return fmt.Sprintf(`SELECT %s FROM %s pld 
		LEFT JOIN %s p ON p.%s = pld.%s 
		LEFT JOIN %s c ON c.%s = p.%s`,
		cols,
		LOOSE_PRODUCTS_IN_DAY_TAB,
		PRODUCT_TAB, PRODUCT_ID, LOOSE_PRODUCTS_IN_DAY_PRODUCT_ID,
		CATEGORY_TAB, CATEGORY_ID, PRODUCT_CATEGORY)
}
