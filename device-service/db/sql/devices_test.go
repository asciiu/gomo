package sql_test

import (
	"log"
	"testing"

	"github.com/asciiu/gomo/common/db"
	"github.com/asciiu/gomo/device-service/db/sql"
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

func TestInsertDevice(t *testing.T) {
	db, _ := db.NewDB("postgres://postgres@localhost/gomo_test?&sslmode=disable")

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := userRepo.InsertUser(db, user)
	checkErr(error)
	defer db.Close()

	device := pb.AddDeviceRequest{
		UserID:           user.ID,
		DeviceType:       "ios",
		ExternalDeviceID: "device-1234",
		DeviceToken:      "tokie-tokie",
	}
	deviceAdded, error := sql.InsertDevice(db, &device)
	checkErr(error)

	error = sql.DeleteDevice(db, deviceAdded.DeviceID)
	checkErr(error)

	error = userRepo.DeleteUserHard(db, user.ID)
	checkErr(error)
}

func TestFindDevice(t *testing.T) {
	db, _ := db.NewDB("postgres://postgres@localhost/gomo_test?&sslmode=disable")

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := userRepo.InsertUser(db, user)
	checkErr(error)
	defer db.Close()

	device := pb.AddDeviceRequest{
		UserID:           user.ID,
		DeviceType:       "ios",
		ExternalDeviceID: "device-1234",
		DeviceToken:      "tokie-tokie",
	}
	deviceAdded, error := sql.InsertDevice(db, &device)
	checkErr(error)

	request := pb.GetUserDeviceRequest{
		DeviceID: deviceAdded.DeviceID,
		UserID:   user.ID,
	}
	device2, err := sql.FindDeviceByDeviceID(db, &request)
	checkErr(err)
	if deviceAdded.DeviceID != device2.DeviceID {
		t.Errorf("no id match!")
	}

	error = sql.DeleteDevice(db, deviceAdded.DeviceID)
	checkErr(error)

	error = userRepo.DeleteUserHard(db, user.ID)
	checkErr(error)
}

func TestUpdateDevice(t *testing.T) {
	db, _ := db.NewDB("postgres://postgres@localhost/gomo_test?&sslmode=disable")

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := userRepo.InsertUser(db, user)
	checkErr(error)
	defer db.Close()

	device := pb.AddDeviceRequest{
		UserID:           user.ID,
		DeviceType:       "ios",
		ExternalDeviceID: "device-1234",
		DeviceToken:      "tokie-tokie",
	}
	deviceAdded, error := sql.InsertDevice(db, &device)
	checkErr(error)

	deviceAdded.DeviceType = "android"

	request := pb.UpdateDeviceRequest{
		DeviceID:         deviceAdded.DeviceID,
		UserID:           deviceAdded.UserID,
		ExternalDeviceID: "1234",
		DeviceType:       "android",
		DeviceToken:      "tokie",
	}
	_, error = sql.UpdateDevice(db, &request)
	checkErr(error)

	requestFind := pb.GetUserDeviceRequest{
		DeviceID: deviceAdded.DeviceID,
		UserID:   user.ID,
	}
	device2, err := sql.FindDeviceByDeviceID(db, &requestFind)
	checkErr(err)
	if device2.DeviceType != "android" {
		t.Errorf("did not update!")
	}

	error = sql.DeleteDevice(db, device2.DeviceID)
	checkErr(error)

	error = userRepo.DeleteUserHard(db, user.ID)
	checkErr(error)
}

func TestFindDevices(t *testing.T) {
	db, _ := db.NewDB("postgres://postgres@localhost/gomo_test?&sslmode=disable")

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := userRepo.InsertUser(db, user)
	checkErr(error)
	defer db.Close()

	device := pb.AddDeviceRequest{
		UserID:           user.ID,
		DeviceType:       "ios",
		ExternalDeviceID: "device-1234",
		DeviceToken:      "tokie-tokie",
	}
	_, error = sql.InsertDevice(db, &device)
	checkErr(error)

	device = pb.AddDeviceRequest{
		UserID:           user.ID,
		DeviceType:       "android",
		ExternalDeviceID: "device-5678",
		DeviceToken:      "token",
	}
	_, error = sql.InsertDevice(db, &device)
	checkErr(error)

	request := pb.GetUserDevicesRequest{
		UserID: user.ID,
	}
	devices, err := sql.FindDevicesByUserID(db, &request)
	checkErr(err)

	if len(devices) != 2 {
		t.Errorf("there should be two devices!")
	}
	if devices[0].DeviceType != "ios" && devices[1].DeviceType != "android" {
		t.Errorf("should have found ios and android")
	}

	error = userRepo.DeleteUserHard(db, user.ID)
	checkErr(error)
}
