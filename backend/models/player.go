package models

type Player struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	RoomID    string `json:"roomId"`
	RoomToken string `json:"roomToken"`
}
