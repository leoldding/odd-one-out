package services

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/leoldding/odd-one-out/models"
)

func CreateRoom(createRoomRequest models.CreateRoomRequest, createRoomResponse *models.CreateRoomResponse) error {
	createRoomResponse.Player.ID = uuid.New().String()
	createRoomResponse.Player.Leader = true
	b := make([]byte, 3)
	rand.Read(b)
	roomCode := hex.EncodeToString(b)
	createRoomResponse.RoomCode = roomCode
	return nil
}
