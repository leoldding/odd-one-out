package models

import "log"

type Publisher struct {
	Broadcast  chan Message
	Register   chan Subscriber
	Deregister chan Subscriber
	Games      map[string]map[Subscriber]bool
}

func (publisher Publisher) Publish() {
	log.Println("Starting Publisher")
	for {
		select {
		case subscriber := <-publisher.Register:
			log.Println("Registering player " + subscriber.Name + " to game " + subscriber.RoomCode)
			subscribers := publisher.Games[subscriber.RoomCode]
			if subscribers == nil {
				subscribers = make(map[Subscriber]bool)
				publisher.Games[subscriber.RoomCode] = subscribers
			}
			publisher.Games[subscriber.RoomCode][subscriber] = true
		case subscriber := <-publisher.Deregister:
			log.Println(subscriber)
			log.Println("Deregistering player " + subscriber.Name + " from game " + subscriber.RoomCode)
			subscribers := publisher.Games[subscriber.RoomCode]
			if subscribers != nil {
				if _, ok := subscribers[subscriber]; ok {
					delete(subscribers, subscriber)
					close(subscriber.MessageChannel)
					if len(subscribers) == 0 {
						delete(publisher.Games, subscriber.RoomCode)
					}
				}
			}
		case message := <-publisher.Broadcast:
			subscribers := publisher.Games[message.RoomCode]
			for subscriber := range subscribers {
				select {
				case subscriber.MessageChannel <- message.Text:
				default:
					close(subscriber.MessageChannel)
					delete(subscribers, subscriber)
					if len(subscribers) == 0 {
						delete(publisher.Games, message.RoomCode)
					}
				}
			}
		}
	}
}
