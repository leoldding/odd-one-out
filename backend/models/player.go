package models

type Player struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Leader bool   `json:"leader"`
}
