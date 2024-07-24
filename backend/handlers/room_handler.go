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
	router.HandleFunc("/room/player", JoinPlayer).Methods("POST")
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

func JoinPlayer(w http.ResponseWriter, r *http.Request) {
	var joinPlayerRequest models.JoinPlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&joinPlayerRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := utils.IsStructFull(joinPlayerRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var joinPlayerResponse models.JoinPlayerResponse
	if err := services.JoinPlayer(joinPlayerRequest, &joinPlayerResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(joinPlayerResponse)
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
	if err := services.JoinRoom(joinRoomRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
