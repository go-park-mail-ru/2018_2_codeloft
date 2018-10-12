package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"2018_2_codeloft/validator"
)

func checkAuth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if exist := dataBase.CheckCookie(cookie.Value); !exist {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func signIn(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error while reading body in /session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var u struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	err = json.Unmarshal(body, &u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(MyError{"wrong requst format"}))
		return
	}
	err = validator.ValidateLogin(u.Login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(MyError{"bad login"}))
		return
	}
	dbUser, exist := dataBase.GetUserByLogin(u.Login)
	if !exist {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(MyError{"User does not exist"}))
		return
	}
	if dbUser.Password != u.Password {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(MyError{"wrong password"}))
		return
	}
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    "testCookie",
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: false,
	}
	dataBase.AddCookie(cookie.Value)
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)
	dataBase.DelCookie(cookie.Value)
	w.WriteHeader(http.StatusOK)
}

var SessionHandler = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch r.Method {

	case http.MethodGet:
		checkAuth(w, r)
	case http.MethodPost:
		signIn(w, r)
	case http.MethodDelete:
		logout(w, r)

	}
}
