package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	"go.uber.org/zap"
)

type DB struct {
	DataBase    *sql.DB
	DB_NAME     string
	DB_URL      string
	DB_USERNAME string
	DB_PASSWORD string
}

func (db *DB) ConnectDataBase() {
	var dbInfo string
	// Если есть DB_URL то мы используем его(для хероку)
	if db.DB_URL != "" {
		dbInfo = db.DB_URL
	} else {
		dbInfo = fmt.Sprintf("user=%s password=%s dbname=%s host=127.0.0.1 port=5432 sslmode=disable", db.DB_USERNAME, db.DB_PASSWORD, db.DB_NAME)
	}
	database, err := sql.Open("postgres", dbInfo)
	if err != nil {
		fmt.Println("Can't connect to database", err)
	}
	err = database.Ping()
	if err != nil {
		zap.L().Error("Error in field",
			zap.Error(err),
		)
	}
	db.DataBase = database
}

func (db *DB) Init(filename string) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		zap.L().Error("Cant read file",
			zap.String("filename", filename),
			zap.Error(err),
		)
		return
	}
	str := string(bs)
	_, err = db.DataBase.Exec(str)
	if err != nil {
		zap.L().Error("Error while db Init Executing script",
			zap.Error(err),
		)
		return
	}

}
