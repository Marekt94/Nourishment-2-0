package meal

import (
	"database/sql"
	"fmt"
	"nourishment_20/internal/database"
	"nourishment_20/internal/logging"
	"strings"
)

// Own the categoryDb model here so it can be reused across the package
type categoryDb struct {
	Id   sql.NullInt64
	Name sql.NullString
}

func (c *categoryDb) ConvertToCategory(out *Category) {
	out.Id = database.NullInt64ToInt(&c.Id)
	out.Name = database.NullStringToString(&c.Name)
}

func (mr *FirebirdRepoAccess) GetCategory(i int) Category {
	var c Category
	sqlStr := fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s = ?", CATEGORY_ID, CATEGORY_NAME, CATEGORY_TAB, CATEGORY_ID)
	row := mr.Database.QueryRow(sqlStr, i)
	var dbRow categoryDb
	if err := row.Scan(&dbRow.Id, &dbRow.Name); err == nil {
		dbRow.ConvertToCategory(&c)
	} else if err != sql.ErrNoRows {
		logging.Global.Panicf("%v", err)
	}
	return c
}

func (mr *FirebirdRepoAccess) GetCategories() []Category {
	sqlStr := fmt.Sprintf("SELECT %s, %s FROM %s ORDER BY %s", CATEGORY_ID, CATEGORY_NAME, CATEGORY_TAB, CATEGORY_ID)
	rows, err := mr.Database.Query(sqlStr)
	if err != nil {
		logging.Global.Panicf("%v", err)
	}
	var res []Category
	for rows.Next() {
		var dbRow categoryDb
		if err := rows.Scan(&dbRow.Id, &dbRow.Name); err == nil {
			var c Category
			dbRow.ConvertToCategory(&c)
			res = append(res, c)
		} else {
			logging.Global.Panicf("%v", err)
		}
	}
	return res
}

func (mr *FirebirdRepoAccess) CreateCategory(c *Category) int64 {
	sqlStr := fmt.Sprintf("INSERT INTO %s (%s) VALUES (?)", CATEGORY_TAB, CATEGORY_NAME)
	if _, err := mr.Database.Exec(sqlStr, &c.Name); err != nil {
		logging.Global.Panicf("%v", err)
		return -1
	}
	var id int64
	if err := mr.Database.QueryRow(fmt.Sprintf("SELECT MAX(%s) FROM %s", CATEGORY_ID, CATEGORY_TAB)).Scan(&id); err != nil {
		logging.Global.Panicf("%v", err)
		return -1
	}
	c.Id = int(id)
	return id
}

func (mr *FirebirdRepoAccess) DeleteCategory(i int) bool {
	if _, err := mr.Database.Exec("DELETE FROM "+CATEGORY_TAB+" WHERE "+CATEGORY_ID+" = ?", i); err != nil {
		logging.Global.Panicf("%v", err)
	}
	row := mr.Database.QueryRow("SELECT "+CATEGORY_ID+" FROM "+CATEGORY_TAB+" WHERE "+CATEGORY_ID+" = ?", i)
	var check int
	return row.Scan(&check) == sql.ErrNoRows
}

func (mr *FirebirdRepoAccess) UpdateCategory(c *Category) {
	sqlStr := fmt.Sprintf("UPDATE %s SET %s=? WHERE %s=?", CATEGORY_TAB, strings.Join([]string{CATEGORY_NAME}, ", "), CATEGORY_ID)
	if _, err := mr.Database.Exec(sqlStr, &c.Name, &c.Id); err != nil {
		logging.Global.Panicf("%v", err)
	}
}
