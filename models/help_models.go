package models

//easyjson:json
type HelpUser struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password,omitempty"`
	Email       string `json:"email,omitempty"`
	Score       int64  `json:"score,omitempty"`
}
