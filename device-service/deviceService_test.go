package main

import (
	"context"
	"log"
	"testing"

	"github.com/asciiu/gomo/common/db"
	pb "github.com/asciiu/gomo/device-service/proto/device"
	userRepo "github.com/asciiu/gomo/user-service/db/sql"
	user "github.com/asciiu/gomo/user-service/models"
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
	_, error := userRepo.InsertUser(db, user)
	checkErr(error)

	return &service, user
}

func TestInsertDevice(t *testing.T) {
	service, user := setupService()

	defer service.DB.Close()

	request := pb.AddDeviceRequest{
		UserID:           user.ID,
		DeviceToken:      "tokie tokie tokie",
		DeviceType:       "test",
		ExternalDeviceID: "1234",
	}

	response := pb.DeviceResponse{}

	service.AddDevice(context.Background(), &request, &response)

	if response.Status != "success" {
		t.Errorf(response.Message)
	}

	if response.Data.Device.UserID != request.UserID {
		t.Errorf("user IDs do not match")
	}

	requestRemove := pb.RemoveDeviceRequest{
		UserID:   user.ID,
		DeviceID: response.Data.Device.DeviceID,
	}

	responseDel := pb.DeviceResponse{}
	service.RemoveDevice(context.Background(), &requestRemove, &responseDel)

	if responseDel.Status != "success" {
		t.Errorf(responseDel.Message)
	}
}
