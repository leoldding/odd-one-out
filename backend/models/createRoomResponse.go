package models

type CreateRoomResponse struct {
	RoomCode string `json:"roomCode"`
	Player   Player `json:"player"`
}
