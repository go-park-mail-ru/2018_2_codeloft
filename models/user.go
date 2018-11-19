package models

import (
	"database/sql"
	"log"

	"go.uber.org/zap"
)

type User struct {
	Id       int64  `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"-"`
	Email    string `json:"email"`
	Score    int64  `json:"score"`
	Lang     string `json:"lang"`
}

func (user *User) GetUserByID(db *sql.DB, id int64) bool {
	row := db.QueryRow("select * from users where id = $1", id)

	err := row.Scan(&user.Id, &user.Login, &user.Password, &user.Email, &user.Score, &user.Lang)

	if err != nil {
		//log.Printf("can't scan user with ID: %v. Err: %v\n",id, err)
		return false
	}
	return true
}

func (user *User) GetUserByLogin(db *sql.DB, login string) bool {
	row := db.QueryRow("select * from users where login = $1", login)

	err := row.Scan(&user.Id, &user.Login, &user.Password, &user.Email, &user.Score, &user.Lang)

	if err != nil {
		//log.Printf("can't scan user with Login: %v. %v\n", login,err)
		return false
	}
	return true
}

func (user *User) GetUserByEmail(db *sql.DB, email string) bool {
	row := db.QueryRow("select * from users where email = $1", email)

	err := row.Scan(&user.Id, &user.Login, &user.Password, &user.Email, &user.Score, &user.Lang)

	if err != nil {
		//log.Printf("can't scan user with Email: %v. Err: %v\n",email, err)
		return false
	}
	return true
}

func (user *User) AddUser(db *sql.DB) error {
	var u User
	if u.GetUserByLogin(db, user.Login) {
		return UserAlreadyExist(user.Login)
	}
	_, err := db.Exec("insert into users(login, password,email) values ($1, $2, $3)", user.Login, user.Password, user.Email)
	if err != nil {
		return err
	}
	user.GetUserByLogin(db, user.Login)
	return nil
}

func GetLeaders(db *sql.DB, page int, pageSize int) []User {
	slice := make([]User, 0, pageSize)
	rows, _ := db.Query(`select * from users order by -score limit $1 offset $2`, pageSize, (page-1)*pageSize)

	if rows != nil {
		for rows.Next() {
			var id int64
			var login string
			var password string
			var email string
			var score int64
			var lang string

			rows.Scan(&id, &login, &password, &email, &score, &lang)
			user := User{id, login, password, email, score, lang}

			slice = append(slice, user)
		}
	}

	return slice
}

func (user *User) DeleteUser(db *sql.DB) error {

	u := &User{}

	if !u.GetUserByLogin(db, user.Login) {
		return UserDoesNotExist(user.Login)
	}
	_, err := db.Exec("delete from users where login = $1", user.Login)
	if err != nil {

		log.Printf("cant DeleteUser: %v. Err %v\n", user, err)

		return err
	}
	return nil
}

func (user *User) UpdateUser(db *sql.DB) error {
	var u User
	if !u.GetUserByLogin(db, user.Login) {
		return UserDoesNotExist(user.Login)
	}
	_, err := db.Exec("update users set password=$1, email=$2 where login = $3", user.Password, user.Email, user.Login)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) UpdateScore(db *sql.DB) error {
	_, err := db.Exec("update users set score=$1 where id = $2", user.Score, user.Id)
	if err != nil {
		log.Printf("cant UpdateScore: %v\n", user)
		return err
	}
	return nil
}

func (user *User) UpdateLang(db *sql.DB) error {
	_, err := db.Exec("update users set lang=$1 where id = $2;", user.Lang, user.Id)
	if err != nil {
		zap.S().Infow("Can not update lang", "err", err)
		return err
	}
	return nil
}
