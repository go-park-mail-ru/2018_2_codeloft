package main

import (
	"github.com/icrowley/fake"
	"strconv"
)

var users []User = make([]User, 0, 20)

type BD struct {
	users  []User
	lastid int
}

var dataBase BD = BD{make([]User, 0, 20), 0}

func (bd *BD) saveUser(u User) {
	bd.users = append(bd.users, u)
	bd.lastid++
}

func (db *BD) deleteUser(u User) {
	db.users = append(db.users[:u.Id], db.users[u.Id+1:]...)
}

func (db *BD) updateUser(id int, newUser User) {
	db.users[id] = newUser
}

func (bd *BD) getUserByEmail(email string) (User, bool) {
	for _, u := range bd.users {
		if u.Email == email {
			return u, true
		}
	}
	return User{}, false
}

func (bd *BD) getUserByLogin(login string) (User, bool) {
	for _, u := range bd.users {
		if u.Login == login {
			return u, true
		}
	}
	return User{}, false
}

func (bd *BD) getUserByID(id int) (User, bool) {
	for _, u := range bd.users {
		if u.Id == id {
			return u, true
		}
	}
	return User{}, false
}

func (db *BD) generateUsers(num int) {
	for i := 0; i < num; i++ {
		score, _ := strconv.Atoi(fake.DigitsN(8))

		u := User{db.lastid, fake.FirstName(), fake.SimplePassword(), fake.EmailAddress(), score}
		db.saveUser(u)
	}
	//for _,v := range(users) {
	//	fmt.Println(v)
	//}
}
