package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
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
	if db.DB_URL != "" {
		dbInfo = db.DB_URL
	} else {
		dbInfo = fmt.Sprintf("user=%s password=%s dbname=%s host=127.0.0.1 port=5432 sslmode=disable", db.DB_USERNAME, db.DB_PASSWORD, db.DB_NAME)
	}
	database, err := sql.Open("postgres", dbInfo)
	if err != nil {
		fmt.Println("Can't connect to database",err)
	}
	err = database.Ping()
	if err != nil {
		log.Println("error in ping", err)
	}
	db.DataBase = database
}

func (db *DB) Init(filename string) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Cant read file:", filename, "error", err)
		return
	}
	str := string(bs)
	_, err = db.DataBase.Exec(str)
	if err != nil {
		log.Println("error while db Init Executing script",err)
		return
	}

}
