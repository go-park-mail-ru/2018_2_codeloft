package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
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
	host := "127.0.0.1"
	if os.Getenv("ENV") == "production"{
		host = "db"
	}
	db.DB_URL = fmt.Sprintf("postgresql://%s:%s@%s:5432/%s?sslmode=disable", db.DB_USERNAME, db.DB_PASSWORD, host,db.DB_NAME)
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
		log.Println("error, cant ping:", err)
	}
	db.DataBase = database
}
