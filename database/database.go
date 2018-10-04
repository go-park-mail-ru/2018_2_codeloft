package database

import (
	"2018_2_codeloft/models"
	"github.com/icrowley/fake"
	"log"
	"sort"
	"strconv"
	"sync"
	"time"
)

//singleton
var instance *DB
var once sync.Once
var mu *sync.Mutex = &sync.Mutex{}
//singleton


type DB struct {
	Users  map[string]*models.User
	UsersSlice []*models.User
	CookiesBase map[string]bool
	Lastid int
}

func (db *DB) CheckCookie(val string) bool {
	mu.Lock()
	defer mu.Unlock()
	return db.CookiesBase[val]

}

func (db *DB) AddCookie(value string) {
	mu.Lock()
	defer mu.Unlock()
	db.CookiesBase[value] = true
}

func (db *DB) DelCookie(value string) {
	mu.Lock()
	defer mu.Unlock()
	delete(db.CookiesBase,value)
}

func (db * DB) SaveUser(u *models.User){
	mu.Lock()
	db.Users[u.Login] = u
	db.UsersSlice = append(db.UsersSlice, u)
	mu.Unlock()
	db.Lastid++
}

func (db * DB) DeleteUser(u models.User){
	mu.Lock()
	user,exist := db.Users[u.Login]
	if !exist {
		return
	}
	user.Login = ""
	user.Score = -1
	delete(db.Users, u.Login)
	mu.Unlock()
}

func (db *DB) UpdateUser(u *models.User){
	mu.Lock()
	*db.Users[u.Login] = *u
	mu.Unlock()
}

func (db DB) GetUserByLogin(login string) (models.User,bool){
	mu.Lock()
	defer mu.Unlock()
	if user,exist := db.Users[login]; exist{
		return *user, true
	}
	return models.User{}, false
}

func (db DB) GetUserByEmail(email string) (models.User, bool){
	mu.Lock()
	defer mu.Unlock()
	for _,u := range db.Users {
		if u.Email == email {
			return *u, true
		}
	}
	return models.User{}, false
}

func (db DB) GetUserByID(id int) (models.User, bool){
	mu.Lock()
	defer mu.Unlock()
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
	mu.Lock()
	sort.Slice(db.UsersSlice, UserGreater)
	mu.Unlock()
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
	mu.Lock()
	for _, val := range db.UsersSlice[begin:end] {
		slice = append(slice, *val)
	}
	mu.Unlock()
	return slice
}

func CreateDataBase(size int) *DB{
	once.Do(func() {
		instance = &DB{
			make(map[string]*models.User,size),
		make([]*models.User,0,size),
		make(map[string]bool),
		0,
		}
	})
	return instance
}




