package pubsub

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
)

type Publisher struct {
	Games    map[string]map[*Subscriber]struct{}
	Waiting  map[string]map[*Subscriber]struct{}
	GameInfo map[string]*info
	mu       sync.RWMutex
}

type info struct {
	leader    *Subscriber
	oddOne    *Subscriber
	question  string
	gameState string
	confirmed map[string]struct{}
}

const (
	GETQUESTION    string = "Get Question"
	REVEALQUESTION string = "Reveal Question"
	REVEALOOO      string = "Reveal Odd One Out"
)

func NewPublisher() *Publisher {
	log.Println("Creating Publisher")
	return &Publisher{
		Games:    map[string]map[*Subscriber]struct{}{},
		Waiting:  map[string]map[*Subscriber]struct{}{},
		GameInfo: map[string]*info{},
	}
}

func (publisher *Publisher) Subscribe(subscriber *Subscriber, game string) {
	publisher.mu.Lock()

	// create game
	if _, ok := publisher.Games[game]; !ok {
		publisher.Games[game] = make(map[*Subscriber]struct{})
		publisher.GameInfo[game] = &info{leader: subscriber, gameState: GETQUESTION, confirmed: make(map[string]struct{})}
		publisher.Waiting[game] = make(map[*Subscriber]struct{})
		message := Message{GameCode: game, Command: "NEW LEADER", Body: publisher.GameInfo[game].gameState}
		subscriber.MessageChannel <- message
	}

	publisher.mu.Unlock()

	switch publisher.GameInfo[game].gameState {
	case GETQUESTION:
		// add subscriber to game
		publisher.Games[game][subscriber] = struct{}{}
		publisher.Broadcast(game, "PLAYERS", publisher.GetPlayersInGame(game))
		log.Println(subscriber.Name + " subscribed to game " + game)
		break
	case REVEALQUESTION, REVEALOOO:
		// add subscriber to wait queue
		publisher.Waiting[game][subscriber] = struct{}{}
		message := Message{GameCode: game, Command: "WAIT"}
		subscriber.MessageChannel <- message
		log.Println(subscriber.Name + " joined waiting queue for game " + game)
		break
	}

}

func (publisher *Publisher) Unsubscribe(subscriber *Subscriber, game string) {
	publisher.mu.Lock()

	if _, ok := publisher.Waiting[game][subscriber]; ok {
		delete(publisher.Waiting[game], subscriber)
		log.Println(subscriber.Name + " left waiting queue for game " + game)
		goto end
	}

	// remove subscriber from game
	delete(publisher.Games[game], subscriber)
	log.Println(subscriber.Name + " unsubscribed from game " + game)

	// elect new leader if previous leader is unsubscribing
	if subscriber == publisher.GameInfo[game].leader {
		random := rand.Intn(len(publisher.Games[game]))
		subscribers := publisher.Games[game]
		for sub := range subscribers {
			if random == 0 {
				publisher.GameInfo[game].leader = sub
				message := Message{GameCode: game, Command: "NEW LEADER", Body: publisher.GameInfo[game].gameState}
				sub.MessageChannel <- message
				break
			}
			random--
		}
	}

	// check if player had confirmed choice
	if _, ok := publisher.GameInfo[game].confirmed[subscriber.Name]; ok {
		delete(publisher.GameInfo[game].confirmed, subscriber.Name)
		if len(publisher.GameInfo[game].confirmed) == len(publisher.Games[game]) {
			message := Message{GameCode: game, Command: "ALL CONFIRMED", Body: strconv.Itoa(len(publisher.GameInfo[game].confirmed))}
			publisher.GameInfo[game].leader.MessageChannel <- message
		}
	}

	publisher.mu.Unlock()

	// end round if odd one out disconnects
	if subscriber == publisher.GameInfo[game].oddOne {
		message := Message{GameCode: game, Command: "ODD ONE LEFT", Body: "Odd One Out has disconnected. Round has ended."}
		subscribers := publisher.Games[game]
		for sub := range subscribers {
			sub.MessageChannel <- message
		}
		publisher.GameInfo[game].gameState = GETQUESTION
		message = Message{GameCode: game, Command: "NEW ROUND", Body: publisher.GameInfo[game].gameState}
		publisher.GameInfo[game].leader.MessageChannel <- message

		publisher.AddWaitingToGame(game)
	}

	// end round if not enough players left
	if len(publisher.Games[game]) < 3 && publisher.GameInfo[game].gameState != GETQUESTION {
		message := Message{GameCode: game, Command: "NOT ENOUGH PLAYERS", Body: "Not enough players left. Round has ended."}
		subscribers := publisher.Games[game]
		for sub := range subscribers {
			sub.MessageChannel <- message
		}
		publisher.GameInfo[game].gameState = GETQUESTION
		message = Message{GameCode: game, Command: "NEW ROUND", Body: publisher.GameInfo[game].gameState}
		publisher.GameInfo[game].leader.MessageChannel <- message

		publisher.AddWaitingToGame(game)
	}

	publisher.Broadcast(game, "PLAYERS", publisher.GetPlayersInGame(game))

end:
	// delete game if empty
	if len(publisher.Games[game]) == 0 {
		delete(publisher.Games, game)
		delete(publisher.Waiting, game)
		delete(publisher.GameInfo, game)
		return
	}
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

func (publisher *Publisher) CheckIfNameExists(game string, name string) bool {
	publisher.mu.RLock()
	defer publisher.mu.RUnlock()

	subscribers := publisher.Games[game]
	for subscriber := range subscribers {
		if subscriber.Name == name {
			return true
		}
	}

	waitingSubscribers := publisher.Waiting[game]
	for subscriber := range waitingSubscribers {
		if subscriber.Name == name {
			return true
		}
	}

	return false
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

func (publisher *Publisher) AddWaitingToGame(game string) {
	publisher.mu.Lock()
	defer publisher.mu.Unlock()

	message := Message{GameCode: game, Command: "DONE WAITING"}
	for subscriber := range publisher.Waiting[game] {
		publisher.Games[game][subscriber] = struct{}{}
		subscriber.MessageChannel <- message
	}

	publisher.Waiting[game] = map[*Subscriber]struct{}{}
}
