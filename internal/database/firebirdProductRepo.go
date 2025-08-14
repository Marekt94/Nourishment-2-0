package database

import (
	"database/sql"
	"fmt"
	"nourishment_20/internal/logging"
	"strings"
)

const productSQLPrefix = `p`
const categorySQLPrefix = `c`

var ProductTabs = []string{PRODUCT_ID, PRODUCT_NAME, PRODUCT_KCAL_PER_100, PRODUCT_UNIT_WEIGHT, PRODUCT_PROTEINS, PRODUCT_FATS,
	PRODUCT_SUGAR, PRODUCT_CARBOHYDRATES, PRODUCT_SUGAR_AND_CARBO, PRODUCT_FIBER, PRODUCT_SALT, PRODUCT_UNIT}
var CategoryTabs = []string{CATEGORY_ID, CATEGORY_NAME}

type productDb struct {
	Id            sql.NullInt64
	Name          sql.NullString
	KcalPer100    sql.NullFloat64
	UnitWeight    sql.NullFloat64
	Proteins      sql.NullFloat64
	Fat           sql.NullFloat64
	Sugar         sql.NullFloat64
	Carbohydrates sql.NullFloat64
	SugarAndCarbo sql.NullFloat64
	Fiber         sql.NullFloat64
	Salt          sql.NullFloat64
	Unit          sql.NullString
	Category      categoryDb
}

// categoryDb type is defined in firebirdCategoryRepo.go and reused here.

func ReturnCategoryFieldsForDbRetriving(c *categoryDb) (*sql.NullInt64, *sql.NullString) {
	return &c.Id, &c.Name
}

func ReturnProductFieldsForDbRetriving(p *productDb) (*sql.NullInt64, *sql.NullString, *sql.NullFloat64, *sql.NullFloat64,
	*sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64, *sql.NullFloat64,
	*sql.NullString, *sql.NullInt64, *sql.NullString) {
	id, name := ReturnCategoryFieldsForDbRetriving(&p.Category)
	return &p.Id, &p.Name, &p.KcalPer100, &p.UnitWeight, &p.Proteins, &p.Fat, &p.Sugar, &p.Carbohydrates,
		&p.SugarAndCarbo, &p.Fiber, &p.Salt, &p.Unit, id, name
}

func serializeDBToProduct(r *sql.Rows) Product { // [AI REFACTOR]
	var p productDb
	r.Scan(ReturnProductFieldsForDbRetriving(&p))
	prod := Product{} // [AI REFACTOR]
	p.ConvertToProduct(&prod)
	return prod
}

func (s *productDb) ConvertToProduct(p *Product) { // [AI REFACTOR]
	p.Id = NullInt64ToInt(&s.Id)
	p.Name = NullStringToString(&s.Name)
	p.KcalPer100 = NullFloat64ToFloat(&s.KcalPer100)
	p.UnitWeight = NullFloat64ToFloat(&s.UnitWeight)
	p.Proteins = NullFloat64ToFloat(&s.Proteins)
	p.Fat = NullFloat64ToFloat(&s.Fat)
	p.Sugar = NullFloat64ToFloat(&s.Sugar)
	p.Carbohydrates = NullFloat64ToFloat(&s.Carbohydrates)
	p.Fiber = NullFloat64ToFloat(&s.Fiber)
	p.Salt = NullFloat64ToFloat(&s.Salt)
	p.Unit = NullStringToString(&s.Unit)
	p.Category.Id = NullInt64ToInt(&s.Category.Id)
	p.Category.Name = NullStringToString(&s.Category.Name)
}

func (mr *FirebirdRepoAccess) createSQLForProducts() string {
	tabs := productSQLPrefix + `.` + strings.Join(ProductTabs[:], ", "+productSQLPrefix+".") + `, ` + categorySQLPrefix + `.` + strings.Join(CategoryTabs[:], `, `+categorySQLPrefix+`.`)
	logging.Global.Debugf("%v", tabs)

	sqltempl := `SELECT %s FROM %s LEFT JOIN %s ON %s=%s`
	return fmt.Sprintf(sqltempl, tabs, PRODUCT_TAB+` `+productSQLPrefix, CATEGORY_TAB+` `+categorySQLPrefix, categorySQLPrefix+`.`+CATEGORY_ID, productSQLPrefix+`.`+PRODUCT_CATEGORY)
}

func (mr *FirebirdRepoAccess) GetProduct(i int) Product { // [AI REFACTOR]
	var prod Product // [AI REFACTOR]

	sql := mr.createSQLForProducts()
	sql = sql + ` WHERE ` + productSQLPrefix + `.ID = ?`

	row, err := mr.Database.Query(sql, i)
	if err != nil {
		logging.Global.Panicf("%v", err)
	}
	if row.Next() {
		prod = serializeDBToProduct(row)
	}
	logging.Global.Debugf("%v", prod)
	return prod
}

func (mr *FirebirdRepoAccess) GetProducts() []Product { // [AI REFACTOR]
	prods := []Product{} // [AI REFACTOR]

	sql := mr.createSQLForProducts()
	rows, err := mr.Database.Query(sql)
	if err != nil {
		logging.Global.Panicf("%v", err)
	}
	for rows.Next() {
		prod := serializeDBToProduct(rows)
		prods = append(prods, prod)
	}
	return prods
}

func (mr *FirebirdRepoAccess) CreateProduct(p *Product) int64 { // [AI REFACTOR]
	insertTabs := append(ProductTabs[1:], PRODUCT_CATEGORY)
	sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, PRODUCT_TAB, strings.Join(insertTabs[:], `, `), QuestionMarks(len(insertTabs)))
	if _, err := mr.Database.Exec(sql, &p.Name, &p.KcalPer100, &p.UnitWeight, &p.Proteins, &p.Fat, &p.Sugar, &p.Carbohydrates,
		&p.SugarAndCarbo, &p.Fiber, &p.Salt, &p.Unit, &p.Category.Id); err != nil {
		logging.Global.Panicf("%v", err)
	} else {
		var id int
		err := mr.Database.QueryRow(`SELECT MAX(` + PRODUCT_ID + `) FROM ` + PRODUCT_TAB).Scan(&id)
		if err != nil {
			logging.Global.Panicf("%v", err)
		}
		return int64(id)
	}
	return -1
}

func (mr *FirebirdRepoAccess) DeleteProduct(i int) bool {
	if _, err := mr.Database.Exec(`DELETE FROM `+PRODUCT_TAB+` WHERE ID = ?`, i); err != nil {
		logging.Global.Panicf("%v", err)
	}
	row, err := mr.Database.Query(`SELECT ID FROM `+PRODUCT_TAB+` WHERE ID = ?`, i)
	if err != nil {
		logging.Global.Panicf("%v", err)
	}
	return !row.Next()
}

func (mr *FirebirdRepoAccess) UpdateProduct(p *Product) { // [AI REFACTOR]
	sql := fmt.Sprintf(`UPDATE %s SET %s WHERE ID=?`, PRODUCT_TAB, UpdateValues(append(ProductTabs[1:], PRODUCT_CATEGORY)))
	_, err := mr.Database.Exec(sql, &p.Name, &p.KcalPer100, &p.UnitWeight, &p.Proteins, &p.Fat, &p.Sugar, &p.Carbohydrates,
		&p.SugarAndCarbo, &p.Fiber, &p.Salt, &p.Unit, &p.Category.Id, &p.Id)
	if err != nil {
		logging.Global.Panicf("%v", err)
	}
}
