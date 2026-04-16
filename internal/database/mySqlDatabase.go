package database

import (
	"database/sql"
	"github.com/Marekt94/go-kernel-mt/logging"
)

type MySQLDBEngine struct {
	BaseEngineIntf
}

func (e *MySQLDBEngine) Connect(c *DBConf) *sql.DB {
	logging.Global.Panicf(`not implemented`)
	connString := ""
	return e.BaseEngineIntf.Connect(`mySQL`, connString, c.PathOrName)
}
