package services

import (
	"encoding/json"
	"log"
	"net/http"

	ws "github.com/gorilla/websocket"
	"github.com/leoldding/odd-one-out/models"
)

type Subscriber models.Subscriber

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func UpgradeConnection(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection to WebSocket.\nERROR:%f", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, register, err := websocket.ReadMessage()
	if err != nil {
		log.Println("Error reading register message")
	}
	var subscriber models.Subscriber
	json.Unmarshal(register, &subscriber)
	subscriber.Websocket = websocket
	subscriber.MessageChannel = make(chan []byte, 256)

	Publisher.Register <- subscriber

	go subscriber.ReadMessages(Publisher)
	go subscriber.WriteMessages()
}
