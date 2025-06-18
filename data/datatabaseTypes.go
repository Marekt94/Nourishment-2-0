package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

type DBEngine interface {
	Connect(c *DBConf) *sql.DB
}


type BaseEngineIntf interface {
	Connect(dbType string, connstr string, dbName string)(*sql.DB)
}

type DBConf struct {
	User       string
	Password   string
	Address    string
	PathOrName string
}

type BaseEngine struct {

}

func (e *BaseEngine) Connect(dbType string, connstr string, dbName string)(*sql.DB){
	log.Printf("connection string: %s\n", connstr)
	db, err := sql.Open(dbType, connstr)
	if err != nil {
		log.Fatal(err)
	}
    consoleWriter := zerolog.ConsoleWriter{
        Out:        os.Stdout,
        TimeFormat: `2006/01/02 15:04:05`,
    }
    zlogger := zerolog.New(consoleWriter).With().Timestamp().Logger()
    db = sqldblogger.OpenDriver(connstr, db.Driver(), zerologadapter.New(zlogger))	

	log.Printf("connected to %s\n", dbName)
	return db
}

