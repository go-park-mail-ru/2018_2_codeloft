package database

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var SERV_DB *gorm.DB
var err error

type Post struct {
	gorm.Model
	Author  string
	Message string
}

func addDatabase(dbname string) error {
	// create database with dbname, won't do anything if db already exists
	SERV_DB.Exec("CREATE DATABASE " + dbname)

	// connect to newly created DB (now has dbname param)
	connectionParams := "dbname=" + dbname + " user=docker password=docker sslmode=disable host=db"
	SERV_DB, err = gorm.Open("postgres", connectionParams)
	if err != nil {
		return err
	}

	return nil
}

func Init() (*gorm.DB, error) {
	// set up DB connection and then attempt to connect 5 times over 25 seconds
	connectionParams := "user=docker password=docker sslmode=disable host=db"
	for i := 0; i < 5; i++ {
		SERV_DB, err = gorm.Open("postgres", connectionParams) // gorm checks Ping on Open
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return SERV_DB, err
	}

	// create table if it does not exist
	if !SERV_DB.HasTable(&Post{}) {
		SERV_DB.CreateTable(&Post{})
	}

	testPost := Post{Author: "Dorper", Message: "GoDoRP is Dope"}
	SERV_DB.Create(&testPost)

	return SERV_DB, err
}
