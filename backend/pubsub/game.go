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
	publisher.GameInfo[game].confirmed = 0

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

func (publisher *Publisher) ConfirmChoices(game string) {
	publisher.mu.Lock()
	defer publisher.mu.Unlock()

	publisher.GameInfo[game].confirmed++
	message := Message{GameCode: game, Command: "CONFIRMED CHOICES", Body: strconv.Itoa(publisher.GameInfo[game].confirmed)}
	publisher.GameInfo[game].leader.MessageChannel <- message
}
