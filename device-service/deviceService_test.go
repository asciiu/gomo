package main

import (
	"context"
	"log"
	"testing"

	"github.com/asciiu/gomo/common/db"
	protoDevice "github.com/asciiu/gomo/device-service/proto/device"
	repoUser "github.com/asciiu/gomo/user-service/db/sql"
	user "github.com/asciiu/gomo/user-service/models"
	"github.com/stretchr/testify/assert"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func setupService() (*DeviceService, *user.User) {
	dbUrl := "postgres://postgres@localhost:5432/gomo_test?&sslmode=disable"
	db, _ := db.NewDB(dbUrl)
	service := DeviceService{db}

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := repoUser.InsertUser(db, user)
	checkErr(error)

	return &service, user
}

func TestInsertDevice(t *testing.T) {
	service, user := setupService()

	defer service.DB.Close()

	request := protoDevice.AddDeviceRequest{
		UserID:           user.ID,
		DeviceToken:      "tokie tokie tokie",
		DeviceType:       "test",
		ExternalDeviceID: "1234",
	}

	response := protoDevice.DeviceResponse{}
	service.AddDevice(context.Background(), &request, &response)

	assert.Equal(t, "success", response.Status, response.Message)
	assert.Equal(t, request.UserID, response.Data.Device.UserID, "user ids do not match")

	requestRemove := protoDevice.RemoveDeviceRequest{
		UserID:   user.ID,
		DeviceID: response.Data.Device.DeviceID,
	}
	responseDel := protoDevice.DeviceResponse{}
	service.RemoveDevice(context.Background(), &requestRemove, &responseDel)
	assert.Equal(t, "success", responseDel.Status, responseDel.Message)

	repoUser.DeleteUserHard(service.DB, user.ID)
}
