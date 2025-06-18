package database

import (
	"database/sql"
	"log"
)

type MySQLDBEngine struct {
	BaseEngineIntf
}

func (e *MySQLDBEngine) Connect(c *DBConf) *sql.DB {
	log.Fatalln(`not implemented`)
	connString := ""
	return e.BaseEngineIntf.Connect(`mySQL`, connString, c.PathOrName)
}