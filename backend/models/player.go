package models

type Player struct {
	Name     string `json:"name"`
	GameCode string `json:"gameCode"`
	Leader   bool   `json:"leader"`
}
