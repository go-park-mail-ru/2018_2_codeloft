package handlers

import (
	"2018_2_codeloft/database"
	"2018_2_codeloft/models"
	"2018_2_codeloft/validator"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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

var dataBase *database.DB = database.CreateDataBase(0)



func leaders(w http.ResponseWriter,r * http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(MyError{"error while parsing form"}))
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
	slice := dataBase.GetLeaders(page, pageSize)
	resp, _ := json.Marshal(&slice)
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return
}



func signUp(w http.ResponseWriter,r * http.Request) {
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
		w.Write(generateError(MyError{"wrong request format"}))
		return
	}
	if _, exist := dataBase.GetUserByLogin(u.Login); exist {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(MyError{"User already exist"}))
		return
	}
	err = validator.ValidateEmail(u.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(MyError{"bad email"}))
		return
	}
	err = validator.ValidateLogin(u.Login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(MyError{"bad login"}))
		return
	}
	err = validator.ValidatePassword(u.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(MyError{"bad password"}))
		return
	}
	var user models.User = models.User{Id: dataBase.Lastid, Login: u.Login, Email: u.Email, Password: u.Password, Score: 0}
	dataBase.SaveUser(&user)
	res, err := json.Marshal(&user)
	if err != nil {
		log.Println("error while Marshaling in /user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(res)
}



func deleteUser(w http.ResponseWriter, r* http.Request){
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
		w.Write(generateError(MyError{"wrong request format"}))
		return
	}
	err = validator.ValidateLogin(u.Login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(MyError{"bad login"}))
		return
	}
	user, exist := dataBase.GetUserByLogin(u.Login)
	if !exist {
		w.Write(generateError(MyError{"User does not exist"}))
		return
	}
	dataBase.DeleteUser(user)
	w.WriteHeader(http.StatusOK)
}



func updateUser(w http.ResponseWriter, r *http.Request){
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
		NewPassword string `json:"new_password,omitempty"`
		Password    string `json:"password"`
		Email       string `json:"email,omitempty"`
		Score       int    `json:"score,omitempty"`
	}
	err = json.Unmarshal(body, &u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(MyError{"wrong request format"}))
		return
	}
	err = validator.ValidateLogin(u.Login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(MyError{"bad login"}))
		return
	}
	user, exist := dataBase.GetUserByLogin(u.Login)
	if !exist {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(MyError{"User does not exist"}))
		return
	}

	if user.Password != u.Password {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(MyError{"wrong password"}))
		return
	}
	var newPassword string = user.Password
	var newEmail string = user.Email
	var newScore int = user.Score
	if u.NewPassword != "" {
		newPassword = u.NewPassword
	}
	err = validator.ValidatePassword(newPassword)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(MyError{"bad New password"}))
		return
	}
	if u.Email != "" {
		newEmail = u.Email
	}
	err = validator.ValidateEmail(newEmail)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(MyError{"bad New email"}))
		return
	}
	if u.Score != 0 {
		newScore = u.Score
	}

	newUser := models.User{user.Id, u.Login, newPassword, newEmail, newScore}
	dataBase.UpdateUser(&newUser)
	w.WriteHeader(http.StatusOK)
	res, err := json.Marshal(&newUser)
	if err != nil {
		log.Println("error while Marshaling in /user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(res)
}


var UserHandler = func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("content-type", "application/json")
		switch r.Method {

		case http.MethodGet:
			leaders(w, r)
		case http.MethodPost:
			signUp(w, r)
		case http.MethodDelete:
			deleteUser(w, r)
		case http.MethodPut:
			updateUser(w, r)
		}

}

var UserById = func(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type", "application/json")
	url := r.URL.Path
	url = strings.Trim(url, "/user/")
	id, err := strconv.Atoi(url)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(generateError(MyError{"Bad URL"}))
		return
	}
	u, exist := dataBase.GetUserByID(id)
	if !exist {
		w.WriteHeader(http.StatusBadRequest)
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
}
