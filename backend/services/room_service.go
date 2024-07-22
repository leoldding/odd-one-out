package services

import "github.com/leoldding/odd-one-out/models"

func CreateRoom(player *models.Player) error {
	player.ID = "tempID"
	player.RoomID = "tempRoomID"
	player.RoomToken = "tempRoomToken"
	return nil
}
