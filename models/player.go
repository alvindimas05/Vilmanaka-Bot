package models

type Player struct {
	Username  string `json:"username"`
	Realname  string `json:"realname"`
	Lastlogin string `json:"lastlogin"`
}
