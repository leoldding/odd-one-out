package models

type JoinRoomRequest struct {
	RoomCode string `json:"roomCode"`
	Player   Player `json:"player"`
}
