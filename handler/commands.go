package handler

import (
	"context"
	"github.com/carlmjohnson/requests"
	"go.mau.fi/whatsmeow/types/events"
	"log"
)

func (handler *Handler) CommandOnlinePlayers(e *events.Message) {
	log.Printf("CommandOnlinePlayers executed")

	var players []string
	err := requests.
		URL(handler.baseUrl + "online-players").
		ToJSON(&players).
		Fetch(context.Background())

	if err != nil {
		panic(err)
	}

	msg := "List online players :"
	for _, player := range players {
		msg += "\n- " + player
	}
	handler.SendMessage(e, msg)
}
