package sql_test

import (
	"log"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

// func TestInsertDevice(t *testing.T) {
// 	db, _ := db.NewDB("postgres://postgres@localhost/gomo_test?&sslmode=disable")

// 	user := user.NewUser("first", "last", "test@email", "hash")
// 	_, error := userRepo.InsertUser(db, user)
// 	checkErr(error)
// 	defer db.Close()

// 	device := pb.AddDeviceRequest{
// 		UserId:           user.Id,
// 		DeviceType:       "ios",
// 		ExternalDeviceId: "device-1234",
// 		DeviceToken:      "tokie-tokie",
// 	}
// 	deviceAdded, error := sql.InsertDevice(db, &device)
// 	checkErr(error)

// 	error = sql.DeleteDevice(db, deviceAdded.DeviceId)
// 	checkErr(error)

// 	error = userRepo.DeleteUserHard(db, user.Id)
// 	checkErr(error)
// }

// func TestFindDevice(t *testing.T) {
// 	db, _ := db.NewDB("postgres://postgres@localhost/gomo_test?&sslmode=disable")

// 	user := user.NewUser("first", "last", "test@email", "hash")
// 	_, error := userRepo.InsertUser(db, user)
// 	checkErr(error)
// 	defer db.Close()

// 	device := pb.AddDeviceRequest{
// 		UserId:           user.Id,
// 		DeviceType:       "ios",
// 		ExternalDeviceId: "device-1234",
// 		DeviceToken:      "tokie-tokie",
// 	}
// 	deviceAdded, error := sql.InsertDevice(db, &device)
// 	checkErr(error)

// 	request := pb.GetUserDeviceRequest{
// 		DeviceId: deviceAdded.DeviceId,
// 		UserId:   user.Id,
// 	}
// 	device2, err := sql.FindDeviceByDeviceId(db, &request)
// 	checkErr(err)
// 	if deviceAdded.DeviceId != device2.DeviceId {
// 		t.Errorf("no id match!")
// 	}

// 	error = sql.DeleteDevice(db, deviceAdded.DeviceId)
// 	checkErr(error)

// 	error = userRepo.DeleteUserHard(db, user.Id)
// 	checkErr(error)
// }

// func TestUpdateDevice(t *testing.T) {
// 	db, _ := db.NewDB("postgres://postgres@localhost/gomo_test?&sslmode=disable")

// 	user := user.NewUser("first", "last", "test@email", "hash")
// 	_, error := userRepo.InsertUser(db, user)
// 	checkErr(error)
// 	defer db.Close()

// 	device := pb.AddDeviceRequest{
// 		UserId:           user.Id,
// 		DeviceType:       "ios",
// 		ExternalDeviceId: "device-1234",
// 		DeviceToken:      "tokie-tokie",
// 	}
// 	deviceAdded, error := sql.InsertDevice(db, &device)
// 	checkErr(error)

// 	deviceAdded.DeviceType = "android"

// 	request := pb.UpdateDeviceRequest{
// 		DeviceId:         deviceAdded.DeviceId,
// 		UserId:           deviceAdded.UserId,
// 		ExternalDeviceId: "1234",
// 		DeviceType:       "android",
// 		DeviceToken:      "tokie",
// 	}
// 	_, error = sql.UpdateDevice(db, &request)
// 	checkErr(error)

// 	requestFind := pb.GetUserDeviceRequest{
// 		DeviceId: deviceAdded.DeviceId,
// 		UserId:   user.Id,
// 	}
// 	device2, err := sql.FindDeviceByDeviceId(db, &requestFind)
// 	checkErr(err)
// 	if device2.DeviceType != "android" {
// 		t.Errorf("did not update!")
// 	}

// 	error = sql.DeleteDevice(db, device2.DeviceId)
// 	checkErr(error)

// 	error = userRepo.DeleteUserHard(db, user.Id)
// 	checkErr(error)
// }

// func TestFindDevices(t *testing.T) {
// 	db, _ := db.NewDB("postgres://postgres@localhost/gomo_test?&sslmode=disable")

// 	user := user.NewUser("first", "last", "test@email", "hash")
// 	_, error := userRepo.InsertUser(db, user)
// 	checkErr(error)
// 	defer db.Close()

// 	device := pb.AddDeviceRequest{
// 		UserId:           user.Id,
// 		DeviceType:       "ios",
// 		ExternalDeviceId: "device-1234",
// 		DeviceToken:      "tokie-tokie",
// 	}
// 	_, error = sql.InsertDevice(db, &device)
// 	checkErr(error)

// 	device = pb.AddDeviceRequest{
// 		UserId:           user.Id,
// 		DeviceType:       "android",
// 		ExternalDeviceId: "device-5678",
// 		DeviceToken:      "token",
// 	}
// 	_, error = sql.InsertDevice(db, &device)
// 	checkErr(error)

// 	request := pb.GetUserDevicesRequest{
// 		UserId: user.Id,
// 	}
// 	devices, err := sql.FindDevicesByUserId(db, &request)
// 	checkErr(err)

// 	if len(devices) != 2 {
// 		t.Errorf("there should be two devices!")
// 	}
// 	if devices[0].DeviceType != "ios" && devices[1].DeviceType != "android" {
// 		t.Errorf("should have found ios and android")
// 	}

// 	error = userRepo.DeleteUserHard(db, user.Id)
// 	checkErr(error)
// }
