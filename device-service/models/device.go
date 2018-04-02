package models

import (
	"time"

	"github.com/google/uuid"
)

func NewDevice(userId, deviceId, deviceType, deviceToken string) *UserDevice {
	newId := uuid.New()

	device := UserDevice{
		Id:          newId.String(),
		UserId:      userId,
		DeviceId:    deviceId,
		DeviceType:  deviceType,
		DeviceToken: deviceToken,
	}
	return &device
}

type UserDevice struct {
	Id          string    `json:"id"`
	UserId      string    `json:"userId"`
	DeviceId    string    `json:"deviceId"`
	DeviceType  string    `json:"deviceType"`
	DeviceToken string    `json:"deviceToken"`
	CreatedOn   time.Time `json:"createdOn"`
	UpdatedOn   time.Time `json:"UpdatedOn"`
}
