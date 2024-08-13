package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leoldding/odd-one-out/pubsub"
	"github.com/leoldding/odd-one-out/services"
)

func RegisterGameHandlers(router *mux.Router, publisher *pubsub.Publisher) {
	log.Println("Registering Game Handlers")
	router.HandleFunc("/game", JoinGame(publisher))
}

func JoinGame(publisher *pubsub.Publisher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services.JoinGame(w, r, publisher)
	}
}
