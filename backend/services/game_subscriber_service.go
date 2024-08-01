package services

import (
	"encoding/json"
	"log"
	"net/http"

	ws "github.com/gorilla/websocket"
	"github.com/leoldding/odd-one-out/pubsub"
)

type Subscriber pubsub.Subscriber

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func JoinGame(w http.ResponseWriter, r *http.Request) {
	// upgrade to websocket
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection to websocket:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get player name and game code
	_, register, err := websocket.ReadMessage()
	if err != nil {
		log.Println("Error reading initial websocket message:", err)
	}

	// create subscriber
	var subscriber pubsub.Subscriber
	json.Unmarshal(register, &subscriber)
	subscriber.Websocket = websocket
	subscriber.MessageChannel = make(chan pubsub.Message, 100)
	go subscriber.Run(Publisher)

	// subscribe to publisher
	Publisher.Broadcast(subscriber.GameCode, "PLAYER JOINING", subscriber.Name)
	Publisher.Subscribe(&subscriber, subscriber.GameCode)
}
