package main

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2018_2_codeloft/models"
	"log"
	_ "github.com/lib/pq"
)


type User struct {
	Id       int64    `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"-"`
	Email    string `json:"email"`
}

func main()  {
	username := "kexibq"
	password := "22121996"
	dbname := "codeloft"
		dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s host=127.0.0.1 port=5432 sslmode=disable", username, password, dbname)

	database, err := sql.Open("postgres", dbInfo)
	if err != nil {
		fmt.Println("Can't connect to database",err)
	}
	err = database.Ping()
	if err != nil {
		log.Println("error in ping", err)
	}
	var i int64 = 1
	row := database.QueryRow("select * from users where id = $1",i)
	var user User
	err = row.Scan(&user.Id, &user.Login, &user.Password,&user.Email)
	fmt.Println(user)
	if err != nil {
		log.Printf("can't scan user with ID: %v\n", err)
	}
	var s models.Session
	row = database.QueryRow("select * from sessions where value = $1", "as")
	err = row.Scan(&s.Value,&s.User_id)
	if err != nil {
		log.Printf("No session %v. Err %v\n", s, err)
	}
}
