package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leoldding/odd-one-out/services"
)

func RegisterGameHandlers(router *mux.Router) {
	log.Println("Registering Game Handlers")
	router.HandleFunc("/game", JoinGame)
}

func JoinGame(w http.ResponseWriter, r *http.Request) {
	services.UpgradeConnection(w, r)
}
