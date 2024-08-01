package pubsub

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Subscriber struct {
	Websocket      *websocket.Conn
	Name           string `json:"name"`
	GameCode       string `json:"gameCode"`
	MessageChannel chan Message
	mu             sync.RWMutex
}

func (subscriber *Subscriber) Run(publisher *Publisher) {
	stop := make(chan bool, 2)

	go subscriber.ReadFromWebsocket(publisher, stop)
	go subscriber.SendMessagesToWebsocket(publisher, stop)

	<-stop
	<-stop

	close(subscriber.MessageChannel)
	publisher.Unsubscribe(subscriber, subscriber.GameCode)
	subscriber.Websocket.Close()
}

func (subscriber *Subscriber) ReadFromWebsocket(publisher *Publisher, stop chan bool) {
	defer func() {
		stop <- true
	}()

	for {
		_, message, err := subscriber.Websocket.ReadMessage()
		if err != nil {
			log.Println("Error reading messages from websocket:", err)
			break
		}

		publisher.Broadcast(subscriber.GameCode, string(message))
	}
}

func (subscriber *Subscriber) SendMessagesToWebsocket(publisher *Publisher, stop chan bool) {
	defer func() {
		stop <- true
	}()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-subscriber.MessageChannel:
			if !ok {
				subscriber.mu.Lock()
				subscriber.Websocket.WriteMessage(websocket.CloseMessage, []byte{})
				subscriber.mu.Unlock()
				return
			}
			subscriber.mu.Lock()
			if err := subscriber.Websocket.WriteMessage(websocket.TextMessage, []byte(message.Body)); err != nil {
				subscriber.mu.Unlock()
				log.Println("Error sending message to websocket:", err)
				return
			}
			subscriber.mu.Unlock()
		case <-ticker.C:
			subscriber.mu.Lock()
			if err := subscriber.Websocket.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				subscriber.mu.Unlock()
				log.Println("Error sending ping message:", err)
				return
			}
			subscriber.mu.Unlock()
		}
	}
}
