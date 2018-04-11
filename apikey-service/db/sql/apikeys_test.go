package sql_test

import (
	"log"
	"testing"

	keyRepo "github.com/asciiu/gomo/apikey-service/db/sql"
	pb "github.com/asciiu/gomo/apikey-service/proto/apikey"
	"github.com/asciiu/gomo/common/db"
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
		ApiKeyId:    user.Id,
		UserId:      user.Id,
		Exchange:    "test",
		Key:         "key",
		Secret:      "secret",
		Description: "Hey this worked!",
	}
	_, error = keyRepo.InsertApiKey(db, &key)
	checkErr(error)

	error = keyRepo.DeleteApiKey(db, key.ApiKeyId)
	checkErr(error)

	error = userRepo.DeleteUserHard(db, user.Id)
	checkErr(error)
}

func TestFindApiKey(t *testing.T) {
	db, _ := db.NewDB("postgres://postgres@localhost/gomo_test?&sslmode=disable")

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := userRepo.InsertUser(db, user)
	checkErr(error)
	defer db.Close()

	key := pb.ApiKeyRequest{
		UserId:      user.Id,
		Exchange:    "test",
		Key:         "key",
		Secret:      "secret",
		Description: "Hey this worked!",
	}
	apiKey, error := keyRepo.InsertApiKey(db, &key)
	checkErr(error)

	request := pb.GetUserApiKeyRequest{
		ApiKeyId: apiKey.ApiKeyId,
		UserId:   user.Id,
	}
	apiKey2, err := keyRepo.FindApiKeyById(db, &request)
	checkErr(err)
	if apiKey.ApiKeyId != apiKey2.ApiKeyId {
		t.Errorf("no id match!")
	}
	if apiKey2.Exchange != "test" {
		t.Errorf("no exchange matchie!")
	}
	if apiKey2.Description != "Hey this worked!" {
		t.Errorf("whaaa?!")
	}

	error = userRepo.DeleteUserHard(db, user.Id)
	checkErr(error)
}
