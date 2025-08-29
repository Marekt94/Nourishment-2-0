package database

import (
	"database/sql"
	"fmt"

	_ "github.com/nakagami/firebirdsql"
)

type FBDBEngine struct {
	BaseEngineIntf
}

func (e *FBDBEngine) Connect(c *DBConf) *sql.DB {
	connString := fmt.Sprintf(`%s:%s@%s/%s`, c.User, c.Password, c.Address, c.PathOrName)
	return e.BaseEngineIntf.Connect(`firebirdsql`, connString, c.PathOrName)
}
