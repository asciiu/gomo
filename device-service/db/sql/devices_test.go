package sql_test

import (
	"log"
	"testing"

	"github.com/asciiu/gomo/common/db"
	"github.com/asciiu/gomo/device-service/db/sql"
	"github.com/asciiu/gomo/device-service/models"
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

	device := models.NewDevice(user.Id, "device-1234", "ios", "tokie-tokie")
	_, error = sql.InsertDevice(db, device)
	checkErr(error)

	error = sql.DeleteDevice(db, device.Id)
	checkErr(error)

	error = userRepo.DeleteUserHard(db, user.Id)
	checkErr(error)
}

func TestFindDevice(t *testing.T) {
	db, _ := db.NewDB("postgres://postgres@localhost/gomo_test?&sslmode=disable")

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := userRepo.InsertUser(db, user)
	checkErr(error)
	defer db.Close()

	device := models.NewDevice(user.Id, "device-1234", "ios", "tokie-tokie")
	_, error = sql.InsertDevice(db, device)
	checkErr(error)

	device2, err := sql.FindDevice(db, device.Id)
	checkErr(err)
	if device.Id != device2.Id {
		t.Errorf("no id match!")
	}

	error = sql.DeleteDevice(db, device.Id)
	checkErr(error)

	error = userRepo.DeleteUserHard(db, user.Id)
	checkErr(error)
}

func TestUpdateDevice(t *testing.T) {
	db, _ := db.NewDB("postgres://postgres@localhost/gomo_test?&sslmode=disable")

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := userRepo.InsertUser(db, user)
	checkErr(error)
	defer db.Close()

	device := models.NewDevice(user.Id, "device-1234", "ios", "tokie-tokie")
	_, error = sql.InsertDevice(db, device)
	checkErr(error)

	device.DeviceType = "android"

	device, error = sql.UpdateDevice(db, device)
	checkErr(error)

	device2, err := sql.FindDevice(db, device.Id)
	checkErr(err)
	if device2.DeviceType != "android" {
		t.Errorf("did not update!")
	}

	error = sql.DeleteDevice(db, device.Id)
	checkErr(error)

	error = userRepo.DeleteUserHard(db, user.Id)
	checkErr(error)
}
