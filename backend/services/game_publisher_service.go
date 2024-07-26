package services

import (
	"github.com/leoldding/odd-one-out/models"
)

var Publisher = models.Publisher{
	Broadcast:  make(chan models.Message),
	Register:   make(chan models.Subscriber),
	Deregister: make(chan models.Subscriber),
	Games:      make(map[string]map[models.Subscriber]bool),
}
