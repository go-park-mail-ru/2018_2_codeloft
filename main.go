package main

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2018_2_codeloft/database"
	"github.com/go-park-mail-ru/2018_2_codeloft/handlers"
	"github.com/go-park-mail-ru/2018_2_codeloft/models"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rs/cors"

	_ "github.com/lib/pq"
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



//func GenerateUsers(num int,db *sql.DB) {
//	for i := 0; i < num; i++ {
//		score, _ := strconv.Atoi(fake.DigitsN(8))
//		login := fake.FirstName()
//		for {
//			if _, exist := db.Users[login]; !exist {
//				break
//			}
//			login = fake.FirstName()
//		}
//		u := models.User{0, login, fake.SimplePassword(), fake.EmailAddress(), score}
//		db.Exec("INSERT INTO users()")
//
//	}
//	u := models.User{db.Lastid, "kek", "qwerty12345", "kek@mail.ru", 0}
//	db.SaveUser(&u)
//	//db.ShowUsers()
//}


func main() {
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	DB_NAME := "codeloft"
	dbInfo := os.Getenv("DATABASE_URL")
	if dbInfo == "" {
		dbInfo = fmt.Sprintf("user=%s password=%s dbname=%s host=127.0.0.1 port=5432 sslmode=disable", username, password, DB_NAME)
	}
	fmt.Println(dbInfo)
	db, err := sql.Open("postgres", dbInfo)
	defer db.Close()
	if err != nil {
		fmt.Println("Can't connect to database")
	}
	err = db.Ping()
	if err != nil {
		log.Println("error in ping", err)
	}
	//GenerateUsers(20,db)
	rows, _ := db.Query("select * from users")
	fmt.Println(rows)
	for rows.Next() {
		var id int
		var login string
		var password string
		var email string
		var score int
		rows.Scan(&id,&login,&password,&email,&score)
		user := models.User{id,login,password,email,score}
		fmt.Println(user)
	}
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
	port := os.Getenv("PORT")
	if port != "" {
		log.Println("get port from env: ", port)
	} else {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, panicMW)
}
