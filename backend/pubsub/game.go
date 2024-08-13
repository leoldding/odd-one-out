package pubsub

import (
	"math/rand"
	"strconv"
)

func (publisher *Publisher) GetQuestions(game string) {
	publisher.mu.Lock()
	random := rand.Intn(len(publisher.Games[game]))
	subscribers := publisher.Games[game]
	defer publisher.mu.Unlock()

	// reset confirmed
	publisher.GameInfo[game].confirmed = make(map[string]struct{})

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
	publisher.GameInfo[game].gameState = REVEALQUESTION

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

	publisher.GameInfo[game].gameState = REVEALOOO
	message := Message{GameCode: game, Command: "REVEAL QUESTION", Body: publisher.GameInfo[game].question}
	publisher.GameInfo[game].oddOne.MessageChannel <- message
}

func (publisher *Publisher) RevealOddOneOut(game string) {
	publisher.mu.Lock()

	publisher.GameInfo[game].gameState = GETQUESTION
	oddOne := publisher.GameInfo[game].oddOne.Name
	subscribers := publisher.Games[game]

	publisher.mu.Unlock()
	publisher.mu.RLock()

	for subscriber := range subscribers {
		message := Message{GameCode: game, Command: "REVEAL ODD ONE OUT", Body: "real"}
		if subscriber.Name == oddOne {
			message.Body = "fake"
		}
		subscriber.MessageChannel <- message
	}

	publisher.mu.RUnlock()

	publisher.AddWaitingToGame(game)
	publisher.Broadcast(game, "PLAYERS", publisher.GetPlayersInGame(game))
}

func (publisher *Publisher) ConfirmChoices(game string, name string) {
	publisher.mu.Lock()
	defer publisher.mu.Unlock()

	publisher.GameInfo[game].confirmed[name] = struct{}{}
	if len(publisher.GameInfo[game].confirmed) == len(publisher.Games[game]) {
		message := Message{GameCode: game, Command: "ALL CONFIRMED", Body: strconv.Itoa(len(publisher.GameInfo[game].confirmed))}
		publisher.GameInfo[game].leader.MessageChannel <- message
	}
}
