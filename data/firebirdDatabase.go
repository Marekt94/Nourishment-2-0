package database

import (
	"database/sql"
	"fmt"

	_ "github.com/nakagami/firebirdsql"
)

const MEAL_TAB = "POTRAWY"
const MEAL_ID  = "ID"
const MEAL_NAME = "NAZWA"
const MEAL_RECIPE = "PRZEPIS"

const PRODUCT_TAB = "PRODUKTY"
const PRODUCT_ID = "ID"
const PRODUCT_NAME = "NAZWA"
const PRODUCT_KCAL_PER_100 = "KCAL_NA_100G"
const PRODUCT_UNIT_WEIGHT = "WAGA_JEDNOSTKI"
const PRODUCT_PROTEINS = "BIALKO"
const PRODUCT_FATS = "TLUSZCZ"
const PRODUCT_SUGAR = "CUKRY_PROSTE"
const PRODUCT_CARBOHYDRATES = "CUKRY_ZLOZONE"
const PRODUCT_SUGAR_AND_CARBO = "CUKRY_SUMA"
const PRODUCT_FIBER = "BLONNIK"
const PRODUCT_SALT = "SOL"
const PRODUCT_UNIT = "JEDNOSTKA"
const PRODUCT_CATEGORY = "KATEGORIA"

const CATEGORY_TAB = "KATEGORIA_PRODUKTU"
const CATEGORY_ID = "ID"
const CATEGORY_NAME = "NAZWA_KATEGORII"

const PRODUCTS_IN_MEAL_TAB = `PRODUKTY_W_POTRAWIE`
const PRODUCTS_IN_MEAL_PRODUCT_ID = `ID_PRODUKTU`
const PRODUCTS_IN_MEAL_MEAL_ID = `ID_POTRAWY`
const PRODUCTS_IN_MEAL_WEIGHT = `ILOSC_W_G`
const PRODUCTS_IN_MEAL_ID = `ID`



type FBDBEngine struct {
	BaseEngineIntf
}

func (e *FBDBEngine) Connect(c *DBConf) (*sql.DB){
	connString := fmt.Sprintf(`%s:%s@%s/%s`, c.User, c.Password, c.Address, c.PathOrName)
	return e.BaseEngineIntf.Connect(`firebirdsql`, connString, c.PathOrName)
}

