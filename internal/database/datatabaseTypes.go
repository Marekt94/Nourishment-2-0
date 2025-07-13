package database

import (
	"database/sql"
	"nourishment_20/internal/logging"
	"os"

	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

type DBEngine interface {
	Connect(c *DBConf) *sql.DB
}

type BaseEngineIntf interface {
	Connect(dbType string, connstr string, dbName string) (*sql.DB)
}

type DBConf struct {
	User       string
	Password   string
	Address    string
	PathOrName string
}

type BaseEngine struct {
}

func (e *BaseEngine) Connect(dbType string, connstr string, dbName string) (*sql.DB) {
	logging.Global.Debugf("connection string: %s", connstr)
	db, err := sql.Open(dbType, connstr)
	if err != nil {
		logging.Global.Panicf("connection error: %v", err)
	}
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: `2006/01/02 15:04:05`,
	}
	zlogger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	db = sqldblogger.OpenDriver(connstr, db.Driver(), zerologadapter.New(zlogger))

	logging.Global.Debugf("connected to %s", dbName)
	return db
}

