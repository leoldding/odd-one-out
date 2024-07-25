package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leoldding/odd-one-out/models"
	"github.com/leoldding/odd-one-out/services"
	"github.com/leoldding/odd-one-out/utils"
)

func RegisterRoomHandlers(router *mux.Router) {
	log.Println("Registering Room Handlers")
	router.HandleFunc("/room/create", CreateRoom).Methods("POST")
	router.HandleFunc("/room/join", JoinRoom).Methods("POST")
}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	var createRoomRequest models.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&createRoomRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := utils.IsStructFull(createRoomRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var createRoomResponse models.CreateRoomResponse
	if err := services.CreateRoom(createRoomRequest, &createRoomResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(createRoomResponse)
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	var joinRoomRequest models.JoinRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&joinRoomRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := utils.IsStructFull(joinRoomRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var joinRoomResponse models.JoinRoomResponse
	if err := services.JoinRoom(joinRoomRequest, &joinRoomResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(joinRoomResponse)
}
