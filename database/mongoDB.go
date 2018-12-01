package database

import (
	"log"

	"gopkg.in/mgo.v2"
)

type MongoDB struct {
	Database    *mgo.Database
	Session     *mgo.Session
	DB_URL      string
	DB_USERNAME string
	DB_PASSWORD string
	DB_NAME     string
}

func (db *MongoDB) Connect() error {
	//url := "mongodb://" + db.DB_USERNAME + ":" + db.DB_PASSWORD + db.DB_URL
	//url := "mongodb://127.0.0.1"
	//log.Printf(url)
	log.Println(db.DB_URL)
	session, err := mgo.Dial(db.DB_URL)
	if err != nil {
		return err
	}

	db.Session = session
	db.Database = session.DB(db.DB_NAME)
	return nil
}
