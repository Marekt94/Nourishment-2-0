package database

import (
	"database/sql"
	"github.com/Marekt94/go-kernel-mt/logging"
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
	
	// Instruct the adapter to be controlled by the `LOG_LEVEL` configuration
	logAdapter := zerologadapter.New(zlogger)
	logOptions := []sqldblogger.Option{
		sqldblogger.WithMinimumLevel(sqldblogger.LevelDebug), // Route whatever comes below the adapter
	}
	db = sqldblogger.OpenDriver(connstr, db.Driver(), logAdapter, logOptions...)

	logging.Global.Debugf("connected to %s", dbName)
	return db
}

