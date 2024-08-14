package handler

import (
	"VilmanakaBot/models"
	"context"
	"github.com/carlmjohnson/requests"
	"go.mau.fi/whatsmeow/types/events"
	"strconv"
	"time"
)

func (handler *Handler) CommandOnlinePlayers(e *events.Message) {
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

//func (handler *Handler) CommandLastOnline(e *events.Message, args []string) {
//	days, err := strconv.Atoi(args[0])
//
//	if err != nil {
//        handler.CommandLastOnlinePlayer(e, args[0])
//        return
//	}
//}

func (handler *Handler) CommandLastOnlinePlayer(e *events.Message, playername string) {
	var player models.Player
	err := requests.
		URL(handler.baseUrl + "authme-info/" + playername).
		ToJSON(&player).
		Fetch(context.Background())

	if err != nil {
		handler.SendMessage(e, "Can't find player with username : "+playername)
		return
	}

	i, err := strconv.ParseInt(player.Lastlogin, 10, 64)
	if err != nil {
		panic(err)
	}

	tm := time.Unix(i, 0)

	handler.SendMessage(e, "Player Name : "+playername+
		"\nLast Online : "+tm.In(time.Local).Format(" 15:04:05 Monday, 02 January 2006"))
}
