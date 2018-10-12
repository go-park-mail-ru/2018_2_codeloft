package main

import (
	"2018_2_codeloft/database"
	"2018_2_codeloft/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

var dataBase *database.DB

func init() {
	dataBase = database.CreateDataBase(20)
	dataBase.GenerateUsers(20)
	dataBase.SortUsersSlice()
	dataBase.EndlessSortLeaders()
}

func panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("panicMiddleware", r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				log.Printf("in URL: %v\n\tWith method %v", r.URL.Path, r.Method)
				log.Println("recovered", err)

			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.MainPage)
	mux.HandleFunc("/user", handlers.UserHandler)
	mux.HandleFunc("/session", handlers.SessionHandler)
	mux.HandleFunc("/user/", handlers.UserById)
	panicMW := panicMiddleware(mux)

	fmt.Println("starting server on http://127.0.0.1:8080")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Content-Type"},
	})
	corsMW := c.Handler(panicMW)
	http.ListenAndServe(":8080", corsMW)
}
