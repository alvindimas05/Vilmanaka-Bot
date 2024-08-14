package handler

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

type Handler struct {
	Client         *whatsmeow.Client
	eventHandlerID uint32
	baseUrl        string
}

func (handler *Handler) register() {
	handler.eventHandlerID = handler.Client.AddEventHandler(handler.EventHandler)
}

func (handler *Handler) Initialize() {
	handler.baseUrl = os.Getenv("BASE_URL")
}

func (handler *Handler) EventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		handler.HandleMessage(v)
		break
	}
}

func (handler *Handler) HandleMessage(e *events.Message) {
	message := e.Message.GetConversation()

	if e.Info.IsFromMe || len(message) == 0 || string(message[0]) != os.Getenv("PREFIX") {
		return
	}
	cmd := strings.Replace(strings.Split(message, " ")[0], os.Getenv("PREFIX"), "", 1)

	switch cmd {
	case "online-players":
		handler.CommandOnlinePlayers(e)
		return
	}

	filename := cmd
	text, err := handler.GetTextMessage(filename)
	if err != nil || filename == "help" || filename == "" {
		handler.SendMessage(e, handler.GetHelpMessage())
		return
	}

	handler.SendMessage(e, text)
}

func (handler *Handler) GetHelpMessage() string {
	help, err := handler.GetTextMessage("help")
	if err != nil {
		panic(err)
	}
	return help + handler.ListMessages()
}

func (handler *Handler) ListMessages() string {
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	files, err := os.ReadDir(path.Join(workdir, "messages"))
	if err != nil {
		panic(err)
	}

	text := ""
	for _, file := range files {
		text += "\n" + os.Getenv("PREFIX") + file.Name()
	}

	return text
}

func (handler *Handler) SendMessage(e *events.Message, message string) {
	fmt.Println(
		"Sending message with reply:",
		&waE2E.ContextInfo{
			StanzaID:    proto.String(e.Info.ID),
			Participant: proto.String(e.Info.Sender.ToNonAD().String()),
		},
	)
	_, err := handler.Client.SendMessage(context.Background(), e.Info.Chat,
		&waE2E.Message{
			Conversation: proto.String(message),
			//ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			//	Text: proto.String(message),
			//	ContextInfo: &waE2E.ContextInfo{
			//		StanzaID:    proto.String(e.Info.ID),
			//		Participant: proto.String(e.Info.Sender.ToNonAD().String()),
			//	},
			//},
		},
	)
	if err != nil {
		panic(err)
	}
}

func (handler *Handler) GetTextMessage(filename string) (string, error) {
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	file, err := os.Open(path.Join(workdir, "messages", filename))
	if err != nil {
		return "", err
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	b, err := io.ReadAll(file)
	return string(b), nil
}
