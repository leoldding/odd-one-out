package models

type JoinRoomRequest struct {
	Name     string `json:"name"`
	RoomCode string `json:"roomCode"`
}
