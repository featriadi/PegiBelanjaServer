package db

import (
	"database/sql"
	"pb-dev-be/config"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func Init() {
	conf := config.GetConfig()

	connectionString := conf.DB_Username + ":" + conf.DB_PASSWORD + "@tcp(" + conf.DB_HOST + ":" + conf.DB_PORT + ")/" + conf.DB_NAME

	db, err = sql.Open("mysql", connectionString)

	if err != nil {
		panic("connection error..")
	}

	err = db.Ping()

	if err != nil {
		panic(connectionString)
	}
}

func CreateCon() *sql.DB {
	return db
}
