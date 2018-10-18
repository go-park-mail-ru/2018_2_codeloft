package models

import (
	"database/sql"
	"log"
)

type User struct {
	Id       int64    `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"-"`
	Email    string `json:"email"`
}


func (user *User) GetUserByID(db *sql.DB, id int64) (bool) {
	row := db.QueryRow("select * from users where id = $1", id)
	err := row.Scan(&user.Id,&user.Login,&user.Password,&user.Email)
	if err != nil {
		//log.Printf("can't scan user with ID: %v. Err: %v\n",id, err)
		return false
	}
	return true
}

func (user *User) GetUserByLogin(db *sql.DB, login string) (bool) {
	row := db.QueryRow("select * from users where login = $1", login)
	err := row.Scan(&user.Id,&user.Login,&user.Password,&user.Email)
	if err != nil {
		//log.Printf("can't scan user with Login: %v. %v\n", login,err)
		return false
	}
	return true
}

func (user *User) GetUserByEmail(db *sql.DB, email string) (bool) {
	row := db.QueryRow("select * from users where email = $1", email)
	err := row.Scan(&user.Id,&user.Login,&user.Password,&user.Email)
	if err != nil {
		//log.Printf("can't scan user with Email: %v. Err: %v\n",email, err)
		return false
	}
	return true
}

func (user *User) AddUser(db *sql.DB) (error) {
	var u User
	if u.GetUserByLogin(db, user.Login) {
		return UserAlreadyExist(user.Login)
	}
	_, err := db.Exec("insert into users(login, password,email) values ($1, $2, $3)", user.Login, user.Password, user.Email)
	if err != nil {
		log.Printf("cant AddUser: %v\n", user)
		return err
	}
	user.GetUserByLogin(db, user.Login)
	return nil
}

type leaders struct {
	Id       int64    `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"-"`
	Email    string `json:"email"`
	Score int64 `json:"score"`
}

func GetLeaders (db *sql.DB, page int, pageSize int) ([]leaders) {
	slice := make([]leaders, 0, pageSize)
	rows, _ := db.Query(`select * from users join 
						(select * from game order by -score limit $1 offset $2) as HS on 
						HS.id = users.id order by -HS.score;`,pageSize, (page-1)*pageSize)
	if rows != nil {
		for rows.Next() {
			var id int64
			var login string
			var password string
			var email string
			var score int64
			var game_id int64
			rows.Scan(&id, &login, &password, &email, &score, &game_id)
			user := leaders{id, login, password, email,score}
			slice = append(slice, user)
		}
	}

	return slice
}

func (user *User) DeleteUser(db *sql.DB) (error) {
	var u User
	if !u.GetUserByLogin(db, user.Login) {
		return UserDoesNotExist(user.Login)
	}
	_, err := db.Exec("delete from users where login = $1", user.Login)
	if err != nil {
		log.Printf("cant DeleteUser: %v. Err %v\n", user,err)
		return err
	}
	return nil
}

func (user *User) UpdateUser(db *sql.DB) (error) {
	var u User
	if !u.GetUserByLogin(db, user.Login) {
		return UserDoesNotExist(user.Login)
	}
	_, err := db.Exec("update users set password=$1, email=$2 where login = $3", user.Password,user.Email, user.Login)
	if err != nil {
		log.Printf("cant UpdateUser: %v. Err %v\n", user, err)
		return err
	}
	return nil
}