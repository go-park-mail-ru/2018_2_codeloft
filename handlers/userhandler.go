package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2018_2_codeloft/models"
	"github.com/go-park-mail-ru/2018_2_codeloft/validator"
)

func generateError(err models.MyError) []byte {
	result, e := json.Marshal(&err)
	if e != nil {
		log.Fatal("Error while MarshalJson while generating Error")
	}
	return result
}

func leaders(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path,"error while parsing form",err}))
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
	slice := models.GetLeaders(db,page,pageSize)
	resp, _ := json.Marshal(&slice)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return
}

func signUp(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path,"wrong request format",err}))
		return
	}
	var user models.User
	if exist := user.GetUserByLogin(db, u.Login); exist {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path,"User already exist",fmt.Errorf("")}))
		return
	}
	err = validator.ValidateEmail(u.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path,"bad email",err}))
		return
	}
	err = validator.ValidateLogin(u.Login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path,"bad login",err}))
		return
	}
	err = validator.ValidatePassword(u.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path,"bad password",err}))
		return
	}
	user = models.User{Login: u.Login, Email: u.Email, Password: u.Password}
	user.AddUser(db)
	res, err := json.Marshal(&user)
	if err != nil {
		log.Println("error while Marshaling in /user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// SET COOKIE
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    "testCookie"+user.Login,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
	}
	if os.Getenv("ENV") == "production" {
		cookie.Secure = true
	}
	session := models.Session{cookie.Value, user.Id}
	session.AddCookie(db)
	http.SetCookie(w, &cookie)

	w.Write(res)
}

func deleteUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path,"wrong request format",err}))
		return
	}
	err = validator.ValidateLogin(u.Login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path,"bad login",err}))
		return
	}
	var user models.User
	if !user.GetUserByLogin(db, u.Login) {
		w.Write(generateError(models.MyError{r.URL.Path,"User does not exist",models.UserDoesNotExist(u.Login)}))
		return
	}
	user.DeleteUser(db)
	w.WriteHeader(http.StatusOK)
}

func updateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	_, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(generateError(models.MyError{r.URL.Path,"authorization required",err}))
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("error while reading body in /user",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var u struct {
		Login       string `json:"login"`
		NewPassword string `json:"new_password,omitempty"`
		Password    string `json:"password"`
		Email       string `json:"email,omitempty"`
		Score       int64    `json:"score,omitempty"`
	}
	err = json.Unmarshal(body, &u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path,"wrong request format",err}))
		return
	}
	err = validator.ValidateLogin(u.Login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path,"bad login",err}))
		return
	}
	var user models.User
	if !user.GetUserByLogin(db,u.Login) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path,"User does not exist",err}))
		return
	}

	if user.Password != u.Password {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path,"wrong password",fmt.Errorf("")}))
		return
	}
	var newPassword string = user.Password
	var newEmail string = user.Email
	var newScore int64 = 0
	if u.NewPassword != "" {
		newPassword = u.NewPassword
	}
	err = validator.ValidatePassword(newPassword)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path,"bad New password",err}))
		return
	}
	if u.Email != "" {
		newEmail = u.Email
	}
	err = validator.ValidateEmail(newEmail)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path,"bad New email",err}))
		return
	}
	if u.Score != 0 {
		newScore = u.Score
	}

	newUser := models.User{user.Id, u.Login, newPassword, newEmail}
	newUser.UpdateUser(db)
	game := models.Game{newScore,user.Id}
	game.UpdateScore(db)
	var result struct {
		Id       int64    `json:"user_id"`
		Login    string `json:"login"`
		Email    string `json:"email"`
		Score int64 `json:"score"`
	}
	result.Id = newUser.Id
	result.Login = newUser.Login
	result.Email = newUser.Email
	result.Score = game.Score
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&result)
	if err != nil {
		log.Println("error while Marshaling in /user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(res)
}

//var UserHandler = func(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("content-type", "application/json")
//
//	switch r.Method {
//
//	case http.MethodGet:
//		leaders(w, r)
//	case http.MethodPost:
//		signUp(w, r)
//	case http.MethodDelete:
//		deleteUser(w, r)
//	case http.MethodPut:
//		updateUser(w, r)
//	}
//
//}

type UserHandler struct {
	Db *sql.DB
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	switch r.Method {

	case http.MethodGet:
		leaders(w, r, h.Db)
	case http.MethodPost:
		signUp(w, r, h.Db)
	case http.MethodDelete:
		deleteUser(w, r, h.Db)
	case http.MethodPut:
		updateUser(w, r, h.Db)
	}
}

//var UserById = func(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("content-type", "application/json")
//	url := r.URL.Path
//	url = strings.Trim(url, "/user/")
//	id, err := strconv.Atoi(url)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		w.Write(generateError(MyError{"Bad URL"}))
//		return
//	}
//	u, exist := dataBase.GetUserByID(id)
//	if !exist {
//		w.WriteHeader(http.StatusBadRequest)
//		w.Write(generateError(MyError{"user does not exist"}))
//		return
//	}
//	user, err := json.Marshal(&u)
//	if err != nil {
//		log.Println("error while Marshaling in /user/")
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	w.Write(user)
//}

type UserById struct {
	Db *sql.DB
}

func (h *UserById) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	url := r.URL.Path
	url = strings.Trim(url, "/user/")
	id, err := strconv.ParseInt(url,10,64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path,"Bad URL",err}))
		return
	}
	var u models.User
	if !u.GetUserByID(h.Db,id) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(models.MyError{r.URL.Path,"user does not exist",models.UserDoesNotExist(u.Login)}))
		return
	}
	user, err := json.Marshal(&u)
	if err != nil {
		log.Println("error while Marshaling in /user/")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(user)
}
