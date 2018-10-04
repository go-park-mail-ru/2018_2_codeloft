package models

type User struct {
	Id       int    `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"-"`
	Email    string `json:"email"`
	Score    int    `json:"score"`
}