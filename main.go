package main

import (
	"2018_2_codeloft/database"
	"2018_2_codeloft/handlers"
	"fmt"
	"log"
	"net/http"
	"strings"

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

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("panicMiddleware", r.URL.Path)
		fmt.Printf("URL: %v; Method: %v; Origin: %v\n", r.URL.Path, r.Method, r.Header.Get("Origin"))
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.MainPage)
	mux.HandleFunc("/user", handlers.UserHandler)
	mux.HandleFunc("/session", handlers.SessionHandler)
	mux.HandleFunc("/user/", handlers.UserById)

	fmt.Println("starting server on http://127.0.0.1:8080")
	c := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return strings.Contains(origin, "codeloft") ||
				strings.Contains(origin, "localhost") ||
				strings.Contains(origin, "127.0.0.1")
		},
		//AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Content-Type"},
		//Debug:            true,
	})
	logHandler := logMiddleware(mux)
	corsMW := c.Handler(logHandler)
	panicMW := panicMiddleware(corsMW)
	http.ListenAndServe(":8080", panicMW)
}
