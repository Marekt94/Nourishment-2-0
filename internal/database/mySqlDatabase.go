package database

import (
	"database/sql"
	"nourishment_20/internal/logging"
)

type MySQLDBEngine struct {
	BaseEngineIntf
}

func (e *MySQLDBEngine) Connect(c *DBConf) *sql.DB {
	logging.Global.Panicf(`not implemented`)
	connString := ""
	return e.BaseEngineIntf.Connect(`mySQL`, connString, c.PathOrName)
}