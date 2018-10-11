package main

import (
	"fmt"
	"net/http"

	"2018_2_codeloft/database"
	"2018_2_codeloft/handlers"

	"github.com/rs/cors"
)

var dataBase *database.DB

func init() {
	dataBase = database.CreateDataBase(20)
	dataBase.GenerateUsers(20)
	dataBase.SortUsersSlice()
	dataBase.EndlessSortLeaders()
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.MainPage)
	mux.HandleFunc("/user", handlers.UserHandler)
	mux.HandleFunc("/session", handlers.SessionHandler)
	mux.HandleFunc("/user/", handlers.UserById)

	fmt.Println("starting server on http://127.0.0.1:8080")
	c := cors.New(cors.Options{
		AllowedOrigins:[]string{"*"},
		AllowCredentials: true,
		AllowedMethods:[]string{"GET", "POST", "DELETE", "PUT"},
	})
	corsMW := c.Handler(mux)
	http.ListenAndServe(":8080", corsMW)
}
