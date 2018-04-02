package models

import (
	"time"

	"github.com/google/uuid"
)

func NewDevice(userId, externalDeviceId, deviceType, deviceToken string) *UserDevice {
	newId := uuid.New()

	device := UserDevice{
		Id:               newId.String(),
		UserId:           userId,
		ExternalDeviceId: externalDeviceId,
		DeviceType:       deviceType,
		DeviceToken:      deviceToken,
	}
	return &device
}

type UserDevice struct {
	Id               string    `json:"id"`
	UserId           string    `json:"userId"`
	ExternalDeviceId string    `json:"externalDeviceId"`
	DeviceType       string    `json:"deviceType"`
	DeviceToken      string    `json:"deviceToken"`
	CreatedOn        time.Time `json:"createdOn"`
	UpdatedOn        time.Time `json:"UpdatedOn"`
}
