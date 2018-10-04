package database

import (
	"2018_2_codeloft/models"
	"github.com/icrowley/fake"
	"log"
	"sort"
	"strconv"
	"time"
)



type DB struct {
	Users  map[string]*models.User
	UsersSlice []*models.User
	Lastid int
}

func (db * DB) SaveUser(u *models.User){
	db.Users[u.Login] = u
	db.UsersSlice = append(db.UsersSlice, u)
	db.Lastid++
}

func (db * DB) DeleteUser(u models.User){
	user,exist := db.Users[u.Login]
	if !exist {
		return
	}
	user.Login = ""
	user.Score = -1
	delete(db.Users, u.Login)
}

func (db *DB) UpdateUser(u *models.User){
	*db.Users[u.Login] = *u
}

func (db DB) GetUserByLogin(login string) (models.User,bool){
	if user,exist := db.Users[login]; exist{
		return *user, true
	}
	return models.User{}, false
}

func (db DB) GetUserByEmail(email string) (models.User, bool){
	for _,u := range db.Users {
		if u.Email == email {
			return *u, true
		}
	}
	return models.User{}, false
}

func (db DB) GetUserByID(id int) (models.User, bool){
	for _,u := range db.Users {
		if u.Id == id {
			return *u, true
		}
	}
	return models.User{}, false
}

func (db *DB) GenerateUsers(num int) {
	for i := 0; i < num; i++ {
		score, _ := strconv.Atoi(fake.DigitsN(8))

		u := models.User{db.Lastid, fake.FirstName(), fake.SimplePassword(), fake.EmailAddress(), score}
		db.SaveUser(&u)
	}
	//for _,v := range(Users) {
	//	fmt.Println(v)
	//}
}


func (db *DB) SortUsersSlice() {

	UserGreater := func(i,j int) bool {
		return db.UsersSlice[i].Score > db.UsersSlice[j].Score
	}
	sort.Slice(db.UsersSlice, UserGreater)
}

func (db *DB) EndlessSortLeaders(){
	go func() {
		c := time.Tick(time.Second*15)
		for t := range c {
			log.Println("Sort LeaderBoard happend:", t)
			db.SortUsersSlice()
		}
	}()
}

func (db DB) GetLeaders (page, pageSize int) []models.User{
	slice := make([]models.User, 0, pageSize)
	usersLength := len(db.Users)
	begin := (page - 1) * pageSize
	if begin >= usersLength {
		begin = usersLength - pageSize;
	}
	end := begin + pageSize
	if end >= usersLength{
		end = usersLength+1
	}
	for _, val := range db.UsersSlice[begin:end] {
		slice = append(slice, *val)
	}
	return slice
}

func CreateDataBase(size int) DB{
	return DB{make(map[string]*models.User,size),make([]*models.User,0,size),0}
}
