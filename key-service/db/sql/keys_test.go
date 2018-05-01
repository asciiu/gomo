package sql_test

import (
	"log"
	"testing"

	"github.com/asciiu/gomo/common/db"
	keyRepo "github.com/asciiu/gomo/key-service/db/sql"
	pb "github.com/asciiu/gomo/key-service/proto/key"
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

	key := pb.ApiKeyRequest{
		ApiKeyID:    user.ID,
		UserID:      user.ID,
		Exchange:    "test",
		Key:         "key",
		Secret:      "secret",
		Description: "Hey this worked!",
	}
	_, error = keyRepo.InsertApiKey(db, &key)
	checkErr(error)

	error = keyRepo.DeleteApiKey(db, key.ApiKeyID)
	checkErr(error)

	error = userRepo.DeleteUserHard(db, user.ID)
	checkErr(error)
}

func TestFindApiKey(t *testing.T) {
	db, _ := db.NewDB("postgres://postgres@localhost/gomo_test?&sslmode=disable")

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := userRepo.InsertUser(db, user)
	checkErr(error)
	defer db.Close()

	key := pb.ApiKeyRequest{
		UserID:      user.ID,
		Exchange:    "test",
		Key:         "key",
		Secret:      "secret",
		Description: "Hey this worked!",
	}
	apiKey, error := keyRepo.InsertApiKey(db, &key)
	checkErr(error)

	request := pb.GetUserApiKeyRequest{
		ApiKeyID: apiKey.ApiKeyID,
		UserID:   user.ID,
	}
	apiKey2, err := keyRepo.FindApiKeyByID(db, &request)
	checkErr(err)
	if apiKey.ApiKeyID != apiKey2.ApiKeyID {
		t.Errorf("no id match!")
	}
	if apiKey2.Exchange != "test" {
		t.Errorf("no exchange matchie!")
	}
	if apiKey2.Description != "Hey this worked!" {
		t.Errorf("whaaa?!")
	}

	error = userRepo.DeleteUserHard(db, user.ID)
	checkErr(error)
}
