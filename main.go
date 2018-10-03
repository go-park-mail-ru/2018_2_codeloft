package main

import (
	"2018_2_codeloft/validator"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	PORT    = ":8080"
	APP_ID  = "6707792"
	APP_KEY = "gQuY2Y2aFVdy9tsIwOAL"
	//APP_SECRET = []byte("fdba0e9ffdba0e9ffdba0e9fc8fddc54cfffdbafdba0e9fa60b49899f33652ed2c03c5f")
	APP_DISPLAY  = "page"
	APP_REDIRECT = "http://127.0.0.1" + PORT
)

var (
	APP_SECRET = "fdba0e9ffdba0e9ffdba0e9fc8fddc54cfffdbafdba0e9fa60b49899f33652ed2c03c5f"
	AUTH_URL   = fmt.Sprintf("https://oauth.vk.com/authorize?client_id=%s&display=%s&redirect_uri=%s", APP_ID, APP_DISPLAY, APP_REDIRECT)
	API_URL    = ""
)

type MyError struct {
	//Code int `json:"ErrorCode"`
	What string `json:"What"`
}

func generateError(err MyError) []byte {
	result, e := json.Marshal(&err)
	if e != nil {
		log.Fatal("Error while MarshalJson while generating Error")
	}
	return result
}

type User struct {
	Id       int    `json:"user_id,omitempty"`
	Login    string `json:"login"`
	Password string `json:"-"`
	Email    string `json:"email"`
	Score    int    `json:"score"`
}

func init() {
	dataBase.generateUsers(20)
}

func main() {
	//generateUsers(20)

	//fmt.Println("AUTH_URL:",AUTH_URL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("this is backend server API\n"))
		//fmt.Fprintf(w,"<a href=%s>click</a>",AUTH_URL)
		//http.Redirect(w,r,AUTH_URL,http.StatusSeeOther)

	})

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("content-type", "application/json")

		switch r.Method {

		case http.MethodGet:
			err := r.ParseForm()
			if err != nil {
				w.Write(generateError(MyError{"error while parsing form"}))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			page, err := strconv.Atoi(r.FormValue("page"))
			if err != nil {
				page = 1
			}
			pageSize, err := strconv.Atoi(r.FormValue("page_size"))
			if err != nil {
				pageSize = 0
			}
			if page <= 0 {
				page = 1
			}
			if pageSize <= 0 {
				pageSize = 5
			}
			slice := make([]User, 0, pageSize)
			begin := (page - 1) * pageSize
			end := begin + pageSize
			for _, val := range dataBase.users[begin:end] {
				slice = append(slice, val)
			}
			resp, _ := json.Marshal(&slice)
			w.WriteHeader(http.StatusOK)
			w.Write(resp)

		case http.MethodPost:
			body, err := ioutil.ReadAll(r.Body)

			if err != nil {
				log.Println("error while reading body in /user")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var u struct {
				Login    string `json:"login"`
				Password string `json:"password"`
				Email    string `json:"email"`
			}
			err = json.Unmarshal(body, &u)
			if err != nil {
				w.Write(generateError(MyError{"wrong requst format"}))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if _, exist := dataBase.getUserByLogin(u.Login); exist {
				w.Write(generateError(MyError{"User already exist"}))
				return
			}
			err = validator.ValidateEmail(u.Email)
			if err != nil {
				w.Write(generateError(MyError{"bad email"}))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = validator.ValidateLogin(u.Login)
			if err != nil {
				w.Write(generateError(MyError{"bad login"}))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			err = validator.ValidatePassword(u.Password)
			if err != nil {
				w.Write(generateError(MyError{"bad password"}))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			var user User = User{Id: dataBase.lastid, Login: u.Login, Email: u.Email, Password: u.Password, Score: 0}
			dataBase.saveUser(user)

			res, err := json.Marshal(&user)
			if err != nil {
				log.Println("error while Marshaling in /user")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(res)

		case http.MethodDelete:
			body, err := ioutil.ReadAll(r.Body)

			if err != nil {
				log.Println("error while reading body in /user")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var u struct {
				Login    string `json:"login"`
				Password string `json:"password"`
			}
			err = json.Unmarshal(body, &u)
			if err != nil {
				w.Write(generateError(MyError{"wrong requst format"}))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			user, exist := dataBase.getUserByLogin(u.Login)
			if !exist {
				w.Write(generateError(MyError{"User does not exist"}))
				return
			}
			dataBase.deleteUser(user)
			w.WriteHeader(http.StatusOK)

		case http.MethodPut:
			_, err := r.Cookie("session_id")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(generateError(MyError{"authorization required"}))
				return
			}

			body, err := ioutil.ReadAll(r.Body)

			if err != nil {
				log.Println("error while reading body in /user")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var u struct {
				Login       string `json:"login"`
				NewPassword string `json:"password,omitempty"`
				Password    string `json:"old_password"`
				Email       string `json:"email,omitempty"`
				Score       int    `json:"score,omitempty"`
			}
			err = json.Unmarshal(body, &u)
			if err != nil {
				w.Write(generateError(MyError{"wrong requst format"}))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			user, exist := dataBase.getUserByLogin(u.Login)
			if !exist {
				w.Write(generateError(MyError{"User does not exist"}))
				return
			}
			if user.Password != u.Password {
				w.Write(generateError(MyError{"wrong password"}))
				return
			}
			var newPassword string = user.Password
			var newEmail string = user.Email
			var newScore int = user.Score
			if u.Password != "" {
				newPassword = u.NewPassword
			}
			if u.Email != "" {
				newEmail = u.Email
			}
			if u.Score != 0 {
				newScore = u.Score
			}

			newUser := User{user.Id, u.Login, newPassword, newEmail, newScore}
			dataBase.updateUser(user.Id, newUser)
			w.WriteHeader(http.StatusOK)
		}
	})

	http.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		switch r.Method {

		case http.MethodGet:
			_, err := r.Cookie("session_id")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusOK)

		case http.MethodPost:
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println("error while reading body in /session")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var u User
			err = json.Unmarshal(body, &u)
			if err != nil {
				w.Write(generateError(MyError{"wrong requst format"}))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			dbUser, exist := dataBase.getUserByLogin(u.Login)
			if !exist {
				w.Write(generateError(MyError{"User does not exist"}))
				return
			}
			if dbUser.Password != u.Password {
				w.Write(generateError(MyError{"wrong password"}))
				return
			}
			cookie := http.Cookie{
				Name:     "session_id",
				Value:    u.Login + "testCookie" + u.Password,
				Expires:  time.Now().Add(30 * 24 * time.Hour),
				HttpOnly: false,
			}
			http.SetCookie(w, &cookie)
			w.WriteHeader(http.StatusOK)

		case http.MethodDelete:
			cookie, err := r.Cookie("session_id")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			cookie.Expires = time.Now()
			http.SetCookie(w, cookie)
			w.WriteHeader(http.StatusOK)
		}

	})

	http.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("content-type", "application/json")
		url := r.URL.Path
		url = strings.Trim(url, "/user/")
		id, err := strconv.Atoi(url)
		if err != nil {
			w.Write(generateError(MyError{"Bad URL"}))
			w.WriteHeader(http.StatusBadRequest)
		}
		u, exist := dataBase.getUserByID(id)
		if !exist {
			w.Write(generateError(MyError{"user does not exist"}))
			return
		}
		user, err := json.Marshal(&u)
		if err != nil {
			log.Println("error while Marshaling in /user/")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(user)
	})

	http.HandleFunc("/vkapi", func(w http.ResponseWriter, r *http.Request) {
		//http.Redirect(w,r,AUTH_URL,http.StatusSeeOther)

		ctx := r.Context()
		code := r.FormValue("code")
		conf := oauth2.Config{
			ClientID:     APP_ID,
			ClientSecret: APP_KEY,
			RedirectURL:  APP_REDIRECT,
			Endpoint:     vk.Endpoint,
		}

		token, err := conf.Exchange(ctx, code)
		if err != nil {
			log.Println("cannot exchange")
			log.Println(err)
			return
		}

		client := conf.Client(ctx, token)
		resp, err := client.Get(fmt.Sprintf(API_URL, token.AccessToken))
		if err != nil {
			log.Println("cannot request data")
			log.Println(err)
			return
		}
		defer resp.Body.Close()
	})
	fmt.Println("starting server on http://127.0.0.1:8080")

	http.ListenAndServe(":8080", nil)
}

// curl -X POST -d "email=123@mail.ru&passowrd=123" http://127.0.0.1:8080/user
// curl -X POST -d "email=123@mail.ru&passowrd=123" http://127.0.0.1:8080/session
