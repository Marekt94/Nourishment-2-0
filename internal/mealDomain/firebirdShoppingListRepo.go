package meal

import (
	"database/sql"
	"fmt"
	"nourishment_20/internal/database"
	"github.com/Marekt94/go-kernel-mt/logging"
	"strings"
)

const ShoppingListPrefix = "sl"
const ProductInShoppingListPrefix = "psl"

type shoppingListDb struct {
	Id        sql.NullInt64
	Name      sql.NullString
	CreatedAt sql.NullTime
	EditDate  sql.NullTime
}

type productInShoppingListDb struct {
	Id           sql.NullInt64
	ListId       sql.NullInt64
	IdProduct    sql.NullInt64
	Weight       sql.NullFloat64 // ILOSC
	EditDate     sql.NullTime
	Bought       sql.NullInt64 // KUPIONE (0 or 1)
	ProductName  sql.NullString
	CategoryName sql.NullString
	ProductUnit  sql.NullString
	UnitWeight   sql.NullFloat64
}

func (s *shoppingListDb) ConvertToShoppingList(out *ShoppingList) {
	out.Id = database.NullInt64ToInt(&s.Id)
	out.Name = database.NullStringToString(&s.Name)
	if s.CreatedAt.Valid {
		out.CreatedAt = s.CreatedAt.Time
	}
	if s.EditDate.Valid {
		out.EditDate = s.EditDate.Time
	}
}

func (p *productInShoppingListDb) ConvertToProductInShoppingList(out *ProductInShoppingList) {
	out.Id = database.NullInt64ToInt(&p.Id)
	out.ListId = database.NullInt64ToInt(&p.ListId)
	out.ProductId = database.NullInt64ToInt(&p.IdProduct)
	out.ProductName = database.NullStringToString(&p.ProductName)
	out.CategoryName = database.NullStringToString(&p.CategoryName)
	out.ProductUnit = database.NullStringToString(&p.ProductUnit)
	out.Weight = database.NullFloat64ToFloat(&p.Weight)
	if p.EditDate.Valid {
		out.EditDate = p.EditDate.Time
	}
	out.Bought = database.NullInt64ToInt(&p.Bought) == 1

	// Calculate Quantity: Weight / UnitWeight
	unitWeight := database.NullFloat64ToFloat(&p.UnitWeight)
	if unitWeight > 0 {
		out.Quantity = out.Weight / unitWeight
	} else {
		out.Quantity = 0
	}
}

func (mr *FirebirdRepoAccess) GetShoppingList(id int) ShoppingList {
	var sl ShoppingList
	slCols := []string{SHOPPING_LIST_ID, SHOPPING_LIST_NAME, SHOPPING_LIST_CREATED_AT, SHOPPING_LIST_EDIT_DATE}
	slColsStr := database.CreateColsToSelect(ShoppingListPrefix, slCols)

	pslCols := []string{PRODUCTS_IN_SHOPPING_LIST_ID, PRODUCTS_IN_SHOPPING_LIST_LIST_ID, PRODUCTS_IN_SHOPPING_LIST_PRODUCT_ID, PRODUCTS_IN_SHOPPING_LIST_WEIGHT, PRODUCTS_IN_SHOPPING_LIST_EDIT_DATE, PRODUCTS_IN_SHOPPING_LIST_BOUGHT}
	pslColsStr := database.CreateColsToSelect(ProductInShoppingListPrefix, pslCols)

	pCols := []string{PRODUCT_NAME, PRODUCT_UNIT, PRODUCT_UNIT_WEIGHT}
	pColsStr := database.CreateColsToSelect(productSQLPrefix, pCols)

	cCols := []string{CATEGORY_NAME}
	cColsStr := database.CreateColsToSelect(categorySQLPrefix, cCols)

	allCols := strings.Join([]string{slColsStr, pslColsStr, pColsStr, cColsStr}, ", ")

	sql := fmt.Sprintf("SELECT %s FROM %s %s "+
		"LEFT JOIN %s %s ON %s.%s = %s.%s "+
		"LEFT JOIN %s %s ON %s.%s = %s.%s "+
		"LEFT JOIN %s %s ON %s.%s = %s.%s "+
		"WHERE %s.%s = ? ORDER BY %s.%s, %s.%s",
		allCols, SHOPPING_LIST_TAB, ShoppingListPrefix,
		PRODUCTS_IN_SHOPPING_LIST_TAB, ProductInShoppingListPrefix, ShoppingListPrefix, SHOPPING_LIST_ID, ProductInShoppingListPrefix, PRODUCTS_IN_SHOPPING_LIST_LIST_ID,
		PRODUCT_TAB, productSQLPrefix, ProductInShoppingListPrefix, PRODUCTS_IN_SHOPPING_LIST_PRODUCT_ID, productSQLPrefix, PRODUCT_ID,
		CATEGORY_TAB, categorySQLPrefix, productSQLPrefix, PRODUCT_CATEGORY, categorySQLPrefix, CATEGORY_ID,
		ShoppingListPrefix, SHOPPING_LIST_ID, ShoppingListPrefix, SHOPPING_LIST_ID, ProductInShoppingListPrefix, PRODUCTS_IN_SHOPPING_LIST_ID)

	rows, err := mr.Database.Query(sql, id)
	if err != nil {
		logging.Global.Panicf("%v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var sDb shoppingListDb
		var pslDb productInShoppingListDb
		
		if err := rows.Scan(&sDb.Id, &sDb.Name, &sDb.CreatedAt, &sDb.EditDate,
			&pslDb.Id, &pslDb.ListId, &pslDb.IdProduct, &pslDb.Weight, &pslDb.EditDate, &pslDb.Bought,
			&pslDb.ProductName, &pslDb.ProductUnit, &pslDb.UnitWeight, &pslDb.CategoryName); err != nil {
			logging.Global.Panicf("%v", err)
		}

		if sl.Id == 0 {
			sDb.ConvertToShoppingList(&sl)
		}

		if pslDb.Id.Valid {
			var psl ProductInShoppingList
			pslDb.ConvertToProductInShoppingList(&psl)
			sl.Products = append(sl.Products, psl)
		}
	}

	return sl
}

func (mr *FirebirdRepoAccess) GetShoppingLists() []ShoppingList {
	slCols := []string{SHOPPING_LIST_ID, SHOPPING_LIST_NAME, SHOPPING_LIST_CREATED_AT, SHOPPING_LIST_EDIT_DATE}
	slColsStr := database.CreateColsToSelect(ShoppingListPrefix, slCols)

	pslCols := []string{PRODUCTS_IN_SHOPPING_LIST_ID, PRODUCTS_IN_SHOPPING_LIST_LIST_ID, PRODUCTS_IN_SHOPPING_LIST_PRODUCT_ID, PRODUCTS_IN_SHOPPING_LIST_WEIGHT, PRODUCTS_IN_SHOPPING_LIST_EDIT_DATE, PRODUCTS_IN_SHOPPING_LIST_BOUGHT}
	pslColsStr := database.CreateColsToSelect(ProductInShoppingListPrefix, pslCols)

	pCols := []string{PRODUCT_NAME, PRODUCT_UNIT, PRODUCT_UNIT_WEIGHT}
	pColsStr := database.CreateColsToSelect(productSQLPrefix, pCols)

	cCols := []string{CATEGORY_NAME}
	cColsStr := database.CreateColsToSelect(categorySQLPrefix, cCols)

	allCols := strings.Join([]string{slColsStr, pslColsStr, pColsStr, cColsStr}, ", ")

	sql := fmt.Sprintf("SELECT %s FROM %s %s "+
		"LEFT JOIN %s %s ON %s.%s = %s.%s "+
		"LEFT JOIN %s %s ON %s.%s = %s.%s "+
		"LEFT JOIN %s %s ON %s.%s = %s.%s "+
		"ORDER BY %s.%s DESC, %s.%s",
		allCols, SHOPPING_LIST_TAB, ShoppingListPrefix,
		PRODUCTS_IN_SHOPPING_LIST_TAB, ProductInShoppingListPrefix, ShoppingListPrefix, SHOPPING_LIST_ID, ProductInShoppingListPrefix, PRODUCTS_IN_SHOPPING_LIST_LIST_ID,
		PRODUCT_TAB, productSQLPrefix, ProductInShoppingListPrefix, PRODUCTS_IN_SHOPPING_LIST_PRODUCT_ID, productSQLPrefix, PRODUCT_ID,
		CATEGORY_TAB, categorySQLPrefix, productSQLPrefix, PRODUCT_CATEGORY, categorySQLPrefix, CATEGORY_ID,
		ShoppingListPrefix, SHOPPING_LIST_ID, ProductInShoppingListPrefix, PRODUCTS_IN_SHOPPING_LIST_ID)

	rows, err := mr.Database.Query(sql)
	if err != nil {
		logging.Global.Panicf("%v", err)
	}
	defer rows.Close()

	listMap := make(map[int]*ShoppingList)
	var listOrder []int

	for rows.Next() {
		var sDb shoppingListDb
		var pslDb productInShoppingListDb
		
		if err := rows.Scan(&sDb.Id, &sDb.Name, &sDb.CreatedAt, &sDb.EditDate,
			&pslDb.Id, &pslDb.ListId, &pslDb.IdProduct, &pslDb.Weight, &pslDb.EditDate, &pslDb.Bought,
			&pslDb.ProductName, &pslDb.ProductUnit, &pslDb.UnitWeight, &pslDb.CategoryName); err != nil {
			logging.Global.Panicf("%v", err)
		}

		id := int(sDb.Id.Int64)
		if _, exists := listMap[id]; !exists {
			var sl ShoppingList
			sDb.ConvertToShoppingList(&sl)
			sl.Products = []ProductInShoppingList{}
			listMap[id] = &sl
			listOrder = append(listOrder, id)
		}

		if pslDb.Id.Valid {
			var psl ProductInShoppingList
			pslDb.ConvertToProductInShoppingList(&psl)
			listMap[id].Products = append(listMap[id].Products, psl)
		}
	}

	var res []ShoppingList
	for _, id := range listOrder {
		res = append(res, *listMap[id])
	}
	return res
}

func (mr *FirebirdRepoAccess) CreateShoppingList(s *ShoppingList) int64 {
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (?)", SHOPPING_LIST_TAB, SHOPPING_LIST_NAME)
	if _, err := mr.Database.Exec(sql, s.Name); err != nil {
		logging.Global.Panicf("%v", err)
		return -1
	}

	var id int64
	if err := mr.Database.QueryRow(fmt.Sprintf("SELECT MAX(%s) FROM %s", SHOPPING_LIST_ID, SHOPPING_LIST_TAB)).Scan(&id); err != nil {
		logging.Global.Panicf("%v", err)
		return -1
	}
	s.Id = int(id)

	for i := range s.Products {
		s.Products[i].ListId = int(id)
		mr.AddProductToShoppingList(&s.Products[i])
	}

	return id
}

func (mr *FirebirdRepoAccess) syncProductsInShoppingList(s *ShoppingList) {
	for i, prod := range s.Products {
		if prod.Id <= 0 {
			prod.ListId = s.Id
			id := mr.AddProductToShoppingList(&prod)
			s.Products[i].Id = int(id)
			s.Products[i].ListId = s.Id
		} else {
			prod.ListId = s.Id
			mr.UpdateProductInShoppingList(&prod)
		}
	}

	res, err := mr.Database.Query(fmt.Sprintf(`SELECT %s FROM %s WHERE %s = ?`, PRODUCTS_IN_SHOPPING_LIST_ID, PRODUCTS_IN_SHOPPING_LIST_TAB, PRODUCTS_IN_SHOPPING_LIST_LIST_ID), s.Id)
	if err != nil {
		logging.Global.Panicf("%v", err)
	}

	prodInListIds := make(map[int]bool)
	var value sql.NullInt64
	for res.Next() {
		err := res.Scan(&value)
		if err != nil {
			logging.Global.Panicf("%v", err)
		}
		prodInListIds[int(value.Int64)] = true
	}
	res.Close()
	
	for _, el := range s.Products {
		if el.Id > 0 {
			delete(prodInListIds, el.Id)
		}
	}

	for k := range prodInListIds {
		mr.DeleteProductFromShoppingList(k)
	}
}

func (mr *FirebirdRepoAccess) UpdateShoppingList(s *ShoppingList) {
	sql := fmt.Sprintf("UPDATE %s SET %s=?, %s=CURRENT_TIMESTAMP WHERE %s=?",
		SHOPPING_LIST_TAB, SHOPPING_LIST_NAME, SHOPPING_LIST_EDIT_DATE, SHOPPING_LIST_ID)
	if _, err := mr.Database.Exec(sql, s.Name, s.Id); err != nil {
		logging.Global.Panicf("%v", err)
	}
	mr.syncProductsInShoppingList(s)
}

func (mr *FirebirdRepoAccess) DeleteShoppingList(id int) bool {
	if _, err := mr.Database.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s = ?", SHOPPING_LIST_TAB, SHOPPING_LIST_ID), id); err != nil {
		logging.Global.Panicf("%v", err)
		return false
	}
	return true
}

func (mr *FirebirdRepoAccess) AddProductToShoppingList(p *ProductInShoppingList) int64 {
	if p.ProductId <= 0 {
		return -1 // Invalid product ID, prevent FK constraint violation
	}

	bought := 0
	if p.Bought {
		bought = 1
	}
	sql := fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s) VALUES (?, ?, ?, ?)",
		PRODUCTS_IN_SHOPPING_LIST_TAB,
		PRODUCTS_IN_SHOPPING_LIST_LIST_ID, PRODUCTS_IN_SHOPPING_LIST_PRODUCT_ID,
		PRODUCTS_IN_SHOPPING_LIST_WEIGHT, PRODUCTS_IN_SHOPPING_LIST_BOUGHT)
	
	if _, err := mr.Database.Exec(sql, p.ListId, p.ProductId, p.Weight, bought); err != nil {
		logging.Global.Panicf("%v", err)
		return -1
	}

	var id int64
	if err := mr.Database.QueryRow(fmt.Sprintf("SELECT MAX(%s) FROM %s", PRODUCTS_IN_SHOPPING_LIST_ID, PRODUCTS_IN_SHOPPING_LIST_TAB)).Scan(&id); err != nil {
		logging.Global.Panicf("%v", err)
		return -1
	}
	p.Id = int(id)
	return id
}

func (mr *FirebirdRepoAccess) UpdateProductInShoppingList(p *ProductInShoppingList) {
	bought := 0
	if p.Bought {
		bought = 1
	}
	sql := fmt.Sprintf("UPDATE %s SET %s=?, %s=?, %s=?, %s=CURRENT_TIMESTAMP WHERE %s=?",
		PRODUCTS_IN_SHOPPING_LIST_TAB,
		PRODUCTS_IN_SHOPPING_LIST_PRODUCT_ID, PRODUCTS_IN_SHOPPING_LIST_WEIGHT,
		PRODUCTS_IN_SHOPPING_LIST_BOUGHT, PRODUCTS_IN_SHOPPING_LIST_EDIT_DATE,
		PRODUCTS_IN_SHOPPING_LIST_ID)

	if _, err := mr.Database.Exec(sql, p.ProductId, p.Weight, bought, p.Id); err != nil {
		logging.Global.Panicf("%v", err)
	}
}

func (mr *FirebirdRepoAccess) DeleteProductFromShoppingList(id int) bool {
	if _, err := mr.Database.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s = ?", PRODUCTS_IN_SHOPPING_LIST_TAB, PRODUCTS_IN_SHOPPING_LIST_ID), id); err != nil {
		logging.Global.Panicf("%v", err)
		return false
	}
	return true
}

func (mr *FirebirdRepoAccess) BulkUpdateProductsInShoppingList(products []ProductInShoppingList) {
	if len(products) == 0 {
		return
	}

	tx, err := mr.Database.Begin()
	if err != nil {
		logging.Global.Panicf("Failed to begin transaction: %v", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				logging.Global.Panicf("Failed to commit bulk update: %v", err)
			}
		}
	}()

	sql := fmt.Sprintf("UPDATE %s SET %s=?, %s=?, %s=?, %s=CURRENT_TIMESTAMP WHERE %s=?",
		PRODUCTS_IN_SHOPPING_LIST_TAB,
		PRODUCTS_IN_SHOPPING_LIST_PRODUCT_ID, PRODUCTS_IN_SHOPPING_LIST_WEIGHT,
		PRODUCTS_IN_SHOPPING_LIST_BOUGHT, PRODUCTS_IN_SHOPPING_LIST_EDIT_DATE,
		PRODUCTS_IN_SHOPPING_LIST_ID)

	stmt, err := tx.Prepare(sql)
	if err != nil {
		logging.Global.Panicf("Failed to prepare bulk update statement: %v", err)
	}
	defer stmt.Close()

	for _, p := range products {
		bought := 0
		if p.Bought {
			bought = 1
		}
		if _, err := stmt.Exec(p.ProductId, p.Weight, bought, p.Id); err != nil {
			logging.Global.Panicf("Error during bulk exec: %v", err)
		}
	}
}
