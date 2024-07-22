package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leoldding/odd-one-out/models"
	"github.com/leoldding/odd-one-out/services"
)

func RegisterRoomHandlers(router *mux.Router) {
	log.Println("Registering Room Handlers")
	router.HandleFunc("/room/create", CreateRoom).Methods("POST")
}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	var player models.Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err := services.CreateRoom(&player); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(player)
}
