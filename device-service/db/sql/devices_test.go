package sql_test

import (
	"log"
	"testing"

	"github.com/asciiu/gomo/common/db"
	repoDevice "github.com/asciiu/gomo/device-service/db/sql"
	protoDevice "github.com/asciiu/gomo/device-service/proto/device"
	repoUser "github.com/asciiu/gomo/user-service/db/sql"
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
	_, error := repoUser.InsertUser(db, user)
	checkErr(error)
	defer db.Close()

	device := protoDevice.AddDeviceRequest{
		UserID:           user.ID,
		DeviceType:       "ios",
		ExternalDeviceID: "device-1234",
		DeviceToken:      "tokie-tokie",
	}
	deviceAdded, error := repoDevice.InsertDevice(db, &device)
	checkErr(error)

	error = repoDevice.DeleteDevice(db, deviceAdded.DeviceID)
	checkErr(error)

	error = repoUser.DeleteUserHard(db, user.ID)
	checkErr(error)
}

func TestFindDevice(t *testing.T) {
	db, _ := db.NewDB("postgres://postgres@localhost/gomo_test?&sslmode=disable")

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := repoUser.InsertUser(db, user)
	checkErr(error)
	defer db.Close()

	device := protoDevice.AddDeviceRequest{
		UserID:           user.ID,
		DeviceType:       "ios",
		ExternalDeviceID: "device-1234",
		DeviceToken:      "tokie-tokie",
	}
	deviceAdded, error := repoDevice.InsertDevice(db, &device)
	checkErr(error)

	request := protoDevice.GetUserDeviceRequest{
		DeviceID: deviceAdded.DeviceID,
		UserID:   user.ID,
	}
	device2, err := repoDevice.FindDeviceByDeviceID(db, &request)
	checkErr(err)
	if deviceAdded.DeviceID != device2.DeviceID {
		t.Errorf("no id match!")
	}

	error = repoDevice.DeleteDevice(db, deviceAdded.DeviceID)
	checkErr(error)

	error = repoUser.DeleteUserHard(db, user.ID)
	checkErr(error)
}

func TestUpdateDevice(t *testing.T) {
	db, _ := db.NewDB("postgres://postgres@localhost/gomo_test?&sslmode=disable")

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := repoUser.InsertUser(db, user)
	checkErr(error)
	defer db.Close()

	device := protoDevice.AddDeviceRequest{
		UserID:           user.ID,
		DeviceType:       "ios",
		ExternalDeviceID: "device-1234",
		DeviceToken:      "tokie-tokie",
	}
	deviceAdded, error := repoDevice.InsertDevice(db, &device)
	checkErr(error)

	deviceAdded.DeviceType = "android"

	request := protoDevice.UpdateDeviceRequest{
		DeviceID:         deviceAdded.DeviceID,
		UserID:           deviceAdded.UserID,
		ExternalDeviceID: "1234",
		DeviceType:       "android",
		DeviceToken:      "tokie",
	}
	_, error = repoDevice.UpdateDevice(db, &request)
	checkErr(error)

	requestFind := protoDevice.GetUserDeviceRequest{
		DeviceID: deviceAdded.DeviceID,
		UserID:   user.ID,
	}
	device2, err := repoDevice.FindDeviceByDeviceID(db, &requestFind)
	checkErr(err)
	if device2.DeviceType != "android" {
		t.Errorf("did not update!")
	}

	error = repoDevice.DeleteDevice(db, device2.DeviceID)
	checkErr(error)

	error = repoUser.DeleteUserHard(db, user.ID)
	checkErr(error)
}

func TestFindDevices(t *testing.T) {
	db, _ := db.NewDB("postgres://postgres@localhost/gomo_test?&sslmode=disable")

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := repoUser.InsertUser(db, user)
	checkErr(error)
	defer db.Close()

	device := protoDevice.AddDeviceRequest{
		UserID:           user.ID,
		DeviceType:       "ios",
		ExternalDeviceID: "device-1234",
		DeviceToken:      "tokie-tokie",
	}
	_, error = repoDevice.InsertDevice(db, &device)
	checkErr(error)

	device = protoDevice.AddDeviceRequest{
		UserID:           user.ID,
		DeviceType:       "android",
		ExternalDeviceID: "device-5678",
		DeviceToken:      "token",
	}
	_, error = repoDevice.InsertDevice(db, &device)
	checkErr(error)

	request := protoDevice.GetUserDevicesRequest{
		UserID: user.ID,
	}
	devices, err := repoDevice.FindDevicesByUserID(db, &request)
	checkErr(err)

	if len(devices) != 2 {
		t.Errorf("there should be two devices!")
	}
	if devices[0].DeviceType != "ios" && devices[1].DeviceType != "android" {
		t.Errorf("should have found ios and android")
	}

	error = repoUser.DeleteUserHard(db, user.ID)
	checkErr(error)
}
