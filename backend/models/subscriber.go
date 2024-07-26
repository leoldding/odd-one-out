package models

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Subscriber struct {
	Websocket      *websocket.Conn
	Name           string `json:"name"`
	RoomCode       string `json:"roomCode"`
	MessageChannel chan []byte
}

func (subscriber Subscriber) ReadMessages(publisher Publisher) {
	defer func() {
		publisher.Deregister <- subscriber
		subscriber.Websocket.Close()
	}()

	for {
		_, message, err := subscriber.Websocket.ReadMessage()
		if err != nil {
			log.Println("Subscriber error reading message.")
			break
		}

		publisher.Broadcast <- Message{message, subscriber.RoomCode}
	}
}

func (subscriber Subscriber) WriteMessages() {
	defer func() {
		subscriber.Websocket.Close()
	}()

	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case message, ok := <-subscriber.MessageChannel:
			if !ok {
				subscriber.Websocket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := subscriber.Websocket.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := subscriber.Websocket.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
