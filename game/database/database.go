package database

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
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
	// postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]
	host := "127.0.0.1"
	if os.Getenv("ENV") == "production" {
		host = "db"
	}
	db.DB_URL = fmt.Sprintf("postgresql://%s:%s@%s:5432/%s?sslmode=disable", db.DB_USERNAME, db.DB_PASSWORD, host, db.DB_NAME)
	//db.DB_URL = fmt.Sprintf("postgresql://%s:%s@db/%s?sslmode=disable", db.DB_USERNAME, db.DB_PASSWORD,db.DB_NAME)
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
	fmt.Println(db.DB_USERNAME, db.DB_PASSWORD, db.DB_NAME)
	db.DataBase = database
}
