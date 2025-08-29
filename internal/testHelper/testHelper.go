package testHelper

import (
	db "nourishment_20/internal/database"
	log "nourishment_20/internal/logging"
	"os"

	"database/sql"

	"github.com/joho/godotenv"
)

func InitTestUnit(path string) *sql.DB {
	err := godotenv.Load(path)
	if err != nil {
		log.Global.Panicf("Error loading .env file: %v", err)
	}
	engine := db.FBDBEngine{BaseEngineIntf: &db.BaseEngine{}}

	conf := db.DBConf{
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		Address:    os.Getenv("DB_ADDRESS"),
		PathOrName: os.Getenv("DB_NAME"),
	}

	return engine.Connect(&conf)
}
