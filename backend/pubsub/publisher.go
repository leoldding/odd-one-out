package pubsub

import (
	"log"
	"math/rand"
	"sync"
)

type Publisher struct {
	Games    map[string]map[*Subscriber]struct{}
	GameInfo map[string]*info
	mu       sync.RWMutex
}

type info struct {
	leader   *Subscriber
	oddOne   *Subscriber
	question string
	state    string
}

func NewPublisher() *Publisher {
	log.Println("Creating Publisher")
	return &Publisher{
		Games:    map[string]map[*Subscriber]struct{}{},
		GameInfo: map[string]*info{},
	}
}

func (publisher *Publisher) Subscribe(subscriber *Subscriber, game string) {
	publisher.mu.Lock()
	defer publisher.mu.Unlock()

	// create game
	if _, ok := publisher.Games[game]; !ok {
		publisher.Games[game] = make(map[*Subscriber]struct{})
		publisher.GameInfo[game] = &info{leader: subscriber, state: "Get Question"}
		message := Message{GameCode: game, Command: "NEW LEADER", Body: publisher.GameInfo[game].state}
		subscriber.MessageChannel <- message
	}
	// add subscriber to game
	publisher.Games[game][subscriber] = struct{}{}

	// have subscriber wait for next round
	switch publisher.GameInfo[game].state {
	case "Reveal Question":
		message := Message{GameCode: game, Command: "WAIT", Body: "2"}
		subscriber.MessageChannel <- message
		break
	case "Reveal Odd One Out":
		message := Message{GameCode: game, Command: "WAIT", Body: "1"}
		subscriber.MessageChannel <- message
		break
	}

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
		delete(publisher.GameInfo, game)
	}

	// elect new leader if previous leader is unsubscribing
	if subscriber == publisher.GameInfo[game].leader {
		random := rand.Intn(len(publisher.Games[game]))
		subscribers := publisher.Games[game]
		for sub := range subscribers {
			if random == 0 {
				publisher.GameInfo[game].leader = sub
				message := Message{GameCode: game, Command: "NEW LEADER", Body: publisher.GameInfo[game].state}
				sub.MessageChannel <- message
				break
			}
			random--
		}
	}

	// end round if odd one out disconnects
	if subscriber == publisher.GameInfo[game].oddOne {
		message := Message{GameCode: game, Command: "ODD ONE LEFT", Body: "Odd One Out has disconnected. Round has ended."}
		subscribers := publisher.Games[game]
		for sub := range subscribers {
			sub.MessageChannel <- message
		}
		publisher.GameInfo[game].state = "Get Question"

		message = Message{GameCode: game, Command: "NEW ROUND", Body: publisher.GameInfo[game].state}
		publisher.GameInfo[game].leader.MessageChannel <- message
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

func (publisher *Publisher) GetQuestions(game string) {
	publisher.mu.Lock()
	random := rand.Intn(len(publisher.Games[game]))
	subscribers := publisher.Games[game]
	defer publisher.mu.Unlock()

	// select which subscriber to be the odd one out
	var oddOne string
	for subscriber := range subscribers {
		if random == 0 {
			publisher.GameInfo[game].oddOne = subscriber
			oddOne = subscriber.Name
			break
		}
		random--
	}

	// retrieve to questions from database
	realQuestion := "real question"
	fakeQuestion := "fake question"
	publisher.GameInfo[game].question = realQuestion
	publisher.GameInfo[game].state = "Reveal Question"

	// send questions to subscribers
	for subscriber := range subscribers {
		message := Message{GameCode: game, Command: "GET QUESTION", Body: realQuestion}
		if subscriber.Name == oddOne {
			message.Body = fakeQuestion
		}
		subscriber.MessageChannel <- message
	}
}

func (publisher *Publisher) RevealQuestion(game string) {
	publisher.mu.Lock()
	defer publisher.mu.Unlock()

	publisher.GameInfo[game].state = "Reveal Odd One Out"
	message := Message{GameCode: game, Command: "REVEAL QUESTION", Body: publisher.GameInfo[game].question}
	publisher.GameInfo[game].oddOne.MessageChannel <- message
}

func (publisher *Publisher) RevealOddOneOut(game string) {
	publisher.mu.Lock()
	defer publisher.mu.Unlock()

	publisher.GameInfo[game].state = "Get Question"
	oddOne := publisher.GameInfo[game].oddOne.Name
	subscribers := publisher.Games[game]
	for subscriber := range subscribers {
		message := Message{GameCode: game, Command: "REVEAL ODD ONE OUT", Body: "real"}
		if subscriber.Name == oddOne {
			message.Body = "fake"
		}
		subscriber.MessageChannel <- message
	}
}
