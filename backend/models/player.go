package models

type Player struct {
	Name     string `json:"name"`
	RoomCode string `json:"roomCode"`
	Leader   bool   `json:"leader"`
}
