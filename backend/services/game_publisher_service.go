package services

import (
	"github.com/leoldding/odd-one-out/pubsub"
)

var Publisher *pubsub.Publisher

func CreatePublisher() {
	Publisher = pubsub.NewPublisher()
}
