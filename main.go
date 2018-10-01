package main

import (
	"encoding/json"
	"fmt"
	"github.com/icrowley/fake"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	PORT = ":8080"
	APP_ID     = "6707792"
	APP_KEY    = "gQuY2Y2aFVdy9tsIwOAL"
	//APP_SECRET = []byte("fdba0e9ffdba0e9ffdba0e9fc8fddc54cfffdbafdba0e9fa60b49899f33652ed2c03c5f")
	APP_DISPLAY = "page"
	APP_REDIRECT = "http://127.0.0.1" + PORT
)

var (
	APP_SECRET = "fdba0e9ffdba0e9ffdba0e9fc8fddc54cfffdbafdba0e9fa60b49899f33652ed2c03c5f"
	AUTH_URL    = fmt.Sprintf("https://oauth.vk.com/authorize?client_id=%s&display=%s&redirect_uri=%s",APP_ID, APP_DISPLAY,APP_REDIRECT)
	API_URL = ""
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
	id int `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Age      int    `json:"age"`
	Score    int  `json:"score"`
}




var users []User = make([]User,0,20)

type BD struct {
	users []User
	lastid int
}

var dataBase BD =BD{make([]User,0,20),0}

func (bd *BD) saveUser(u User) {
	bd.users = append(bd.users, u)
	bd.lastid++
}

func (bd *BD) getUserByEmail(email string) (User,bool){
	for _ , u := range bd.users {
		if u.Email == email {
			return u,true
		}
	}
	return User{},false
}

func (bd *BD) getUserByID(id int) (User,bool){
	for _ , u := range bd.users {
		if u.id == id {
			return u,true
		}
	}
	return User{},false
}

func generateUsers(num int){
	for i:=0; i < num; i++ {
		age,_ := strconv.Atoi(fake.DigitsN(2))
		score,_ := strconv.Atoi(fake.DigitsN(8))
		users = append(users, User{i,fake.EmailAddress(),fake.SimplePassword(),age,score})
	}
	//for _,v := range(users) {
	//	fmt.Println(v)
	//}
}

func init(){
	generateUsers(20)
	dataBase.users = users
	dataBase.lastid = 20
}

func main() {
	//generateUsers(20)

	fmt.Println("AUTH_URL:",AUTH_URL)

	http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("this is backend server API\n"))
		//fmt.Fprintf(w,"<a href=%s>click</a>",AUTH_URL)
		//http.Redirect(w,r,AUTH_URL,http.StatusSeeOther)

	})

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("content-type", "application/json")
		switch r.Method {
		case http.MethodGet:

			slice := make([]User, 0, 20)
			for _, val := range users {
				slice = append(slice, val)
			}
			resp, _ := json.Marshal(&slice)

			w.Write(resp)
		case http.MethodPost:
			//email, exist := r.Form["email"]
			email := r.FormValue("email")
			//u,e := dataBase.getUserByEmail(email[0])
			u,e := dataBase.getUserByEmail(email)
			if e{
				w.Write(generateError(MyError{"User Already exist"}))
				return
			}
			//if !exist {
			//	w.Write(generateError(MyError{"NoGetParam Email"}))
			//	return
			//}
			//password, exist := r.Form["password"]
			password := r.FormValue("password")
			//if !exist {
			//	w.Write(generateError(MyError{"NoGetParam password"}))
			//	return
			//}
			//if u.Password != password {
			//	w.Write(generateError(MyError{"WrongPassword"}))
			//	return
			//}
			u = User{dataBase.lastid,email,password,20,0}
			dataBase.saveUser(u)
			res , err := json.Marshal(&u)
			if err != nil{
				log.Println("error while Marshaling in /user")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(res)
		}
	})

	http.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		switch r.Method{
		case http.MethodGet:
			cookie, err := r.Cookie("session_id")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w,"%v", cookie)
		case http.MethodPost:
			err := r.ParseForm()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				//w.Write([]byte("cant parse form"))

			}
			email := r.FormValue("email")
			u,e := dataBase.getUserByEmail(email)
			if !e{
				w.Write(generateError(MyError{"DoNotExist"}))
				return
			}
			password := r.FormValue("password")
			if u.Password != password {
				w.Write(generateError(MyError{"WrongPassword"}))
				return
			}
			cookie := &http.Cookie{
				Name:  "session_id",
				Value: "testCookie",
			}
			http.SetCookie(w,cookie)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w,"okey post cookie:%v",cookie)
		}

	})

	http.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("content-type", "application/json")
		url := r.URL.Path
		//fmt.Fprintf(w,"url: %s\n",url)
		url = strings.Trim(url,"/user/")
		//fmt.Fprintf(w,"url: %s\n",url)
		id,err := strconv.Atoi(url)
		if err != nil {
			w.Write(generateError(MyError{"Bad URL"}))
			w.WriteHeader(http.StatusBadRequest)
		}
		u,exist := dataBase.getUserByID(id)
		if !exist {
			w.Write(generateError(MyError{"user does not exist"}))
			return
		}
		user,err := json.Marshal(&u)
		if err != nil {
			log.Println("error while Marshaling in /user/")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(user)
	})


	http.HandleFunc("/vkapi", func(w http.ResponseWriter, r *http.Request){
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
