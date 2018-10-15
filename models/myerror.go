package models

import "fmt"

type MyError struct {
	URL string `json:"URL"`
	What string `json:"What"`
	Err error `json:"error"`
}

func (e *MyError) Error() string {
	return fmt.Sprintf("MyError: URL: %s, What: %s, Error: %v",e.URL,e.What,e.Err)
}


type UserError struct {
	What string
	Login string
}

func (e * UserError) Error() string {
	return fmt.Sprintf("%s. (Login: %s)",e.What,e.Login)
}

func UserAlreadyExist(login string) (*UserError) {
	return &UserError{fmt.Sprintf("User %s already exist",login),login}
}

func UserDoesNotExist(login string) (*UserError) {
	return &UserError{fmt.Sprintf("User %s does not exist",login),login}
}