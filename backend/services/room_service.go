package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/leoldding/odd-one-out/models"
)

func CreateRoom(createRoomRequest models.CreateRoomRequest, createRoomResponse *models.CreateRoomResponse) error {
	createRoomResponse.Player.Name = createRoomRequest.Name
	b := make([]byte, 3)
	rand.Read(b)
	gameCode := hex.EncodeToString(b)
	createRoomResponse.Player.GameCode = gameCode
	return nil
}

func JoinRoom(joinRoomRequest models.JoinRoomRequest, joinRoomResponse *models.JoinRoomResponse) error {
	if Publisher.CheckIfNameExists(joinRoomRequest.GameCode, joinRoomRequest.Name) {
		return errors.New("Name exists in game already.")
	}
	joinRoomResponse.Player.Name = joinRoomRequest.Name
	joinRoomResponse.Player.GameCode = joinRoomRequest.GameCode
	return nil
}
