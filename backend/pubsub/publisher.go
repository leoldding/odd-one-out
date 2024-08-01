package pubsub

import (
	"log"
	"sync"
)

type Publisher struct {
	Games map[string]map[*Subscriber]struct{}
	mu    sync.RWMutex
}

func NewPublisher() *Publisher {
	log.Println("Creating Publisher")
	return &Publisher{
		Games: map[string]map[*Subscriber]struct{}{},
	}
}

func (publisher *Publisher) Subscribe(subscriber *Subscriber, game string) {
	publisher.mu.Lock()
	defer publisher.mu.Unlock()

	// create game
	if _, ok := publisher.Games[game]; !ok {
		publisher.Games[game] = make(map[*Subscriber]struct{})
	}
	// add subscriber to game
	publisher.Games[game][subscriber] = struct{}{}
	log.Println(subscriber.Name + " subscribed to game " + game)
}

func (publisher *Publisher) Unsubscribe(subscriber *Subscriber, game string) {
	publisher.mu.Lock()
	defer publisher.mu.Unlock()

	// remove subscriber from game
	delete(publisher.Games[game], subscriber)
	// delete game if empty
	if len(publisher.Games[game]) == 0 {
		delete(publisher.Games, game)
	}
	log.Println(subscriber.Name + " unsubscribed from game " + game)
}

func (publisher *Publisher) Broadcast(game string, command string, body string) {
	publisher.mu.RLock()
	subscribers := publisher.Games[game]
	defer publisher.mu.RUnlock()

	// send message to each subscriber in game
	for subscriber := range subscribers {
		message := Message{GameCode: game, Command: command, Body: body}
		subscriber.MessageChannel <- message
	}
}

func (publisher *Publisher) GetPlayersInGame(game string) string {
	publisher.mu.RLock()
	subscribers := publisher.Games[game]
	defer publisher.mu.RUnlock()

	var players string
	for subscriber := range subscribers {
		players += subscriber.Name + ","
	}

	return players[:len(players)-1]
}
