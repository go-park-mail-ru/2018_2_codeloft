package main

import (
	"2018_2_codeloft/database"
	"2018_2_codeloft/handlers"
	"fmt"
	"net/http"
)

var dataBase *database.DB

func init() {
	dataBase = database.CreateDataBase(20)
	dataBase.GenerateUsers(20)
	dataBase.SortUsersSlice()
	dataBase.EndlessSortLeaders()
}

func main() {
	http.HandleFunc("/", handlers.MainPage)
	http.HandleFunc("/user", handlers.UserHandler)
	http.HandleFunc("/session", handlers.SessionHandler)
	http.HandleFunc("/user/", handlers.UserById)

	fmt.Println("starting server on http://127.0.0.1:8080")

	http.ListenAndServe(":8080", nil)
}

