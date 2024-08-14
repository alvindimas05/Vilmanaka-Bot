package models

type Member struct {
	ID       string `json:"id"`
	Whatsapp string `json:"whatsapp"`
	Discord  string `json:"discord"`
}
