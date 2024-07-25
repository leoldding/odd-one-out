package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	ws "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func RegisterGameHandlers(router *mux.Router) {
	log.Println("Registering Game Handlers")
	router.HandleFunc("/game", JoinGame)
}

func JoinGame(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("RoomWS: Error upgrading connection to WebSocket.\nERROR:%f", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	websocket.WriteMessage(ws.TextMessage, []byte("Websocket Connected"))
}
