package models

type JoinRoomRequest struct {
	Name     string `json:"name"`
	GameCode string `json:"gameCode"`
}
