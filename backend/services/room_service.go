package services

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/leoldding/odd-one-out/models"
)

func CreateRoom(createRoomRequest models.CreateRoomRequest, createRoomResponse *models.CreateRoomResponse) error {
	createRoomResponse.Player.Name = createRoomRequest.Name
	createRoomResponse.Player.Leader = true
	b := make([]byte, 3)
	rand.Read(b)
	roomCode := hex.EncodeToString(b)
	createRoomResponse.Player.RoomCode = roomCode
	return nil
}

func JoinRoom(joinRoomRequest models.JoinRoomRequest, joinRoomResponse *models.JoinRoomResponse) error {
	joinRoomResponse.Player.Name = joinRoomRequest.Name
	joinRoomResponse.Player.RoomCode = joinRoomRequest.RoomCode
	joinRoomResponse.Player.Leader = false
	return nil
}
