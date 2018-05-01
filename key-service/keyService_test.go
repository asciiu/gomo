package main

import (
	"context"
	"log"
	"testing"

	"github.com/asciiu/gomo/common/db"
	keyProto "github.com/asciiu/gomo/key-service/proto/key"
	userRepo "github.com/asciiu/gomo/user-service/db/sql"
	user "github.com/asciiu/gomo/user-service/models"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func setupService() (*KeyService, *user.User) {
	dbUrl := "postgres://postgres@localhost:5432/gomo_test?&sslmode=disable"
	db, _ := db.NewDB(dbUrl)
	service := KeyService{DB: db}

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := userRepo.InsertUser(db, user)
	checkErr(error)

	return &service, user
}

func TestInsertKey(t *testing.T) {
	service, user := setupService()

	defer service.DB.Close()

	key := keyProto.KeyRequest{
		UserID:      user.ID,
		Exchange:    "monex",
		Key:         "key",
		Secret:      "sssh",
		Description: "shit dwog!",
	}

	response := keyProto.KeyResponse{}

	service.AddKey(context.Background(), &key, &response)

	if response.Status != "success" {
		t.Errorf(response.Message)
	}

	if response.Data.Key.UserID != key.UserID {
		t.Errorf("user IDs do not match")
	}

	requestRemove := keyProto.RemoveKeyRequest{
		UserID: user.ID,
		KeyID:  response.Data.Key.KeyID,
	}

	responseDel := keyProto.KeyResponse{}
	service.RemoveKey(context.Background(), &requestRemove, &responseDel)

	if responseDel.Status != "success" {
		t.Errorf(responseDel.Message)
	}

	userRepo.DeleteUserHard(service.DB, user.ID)
}
