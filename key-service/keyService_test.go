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
	service := KeyService{db}

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := userRepo.InsertUser(db, user)
	checkErr(error)

	return &service, user
}

func TestInsertApiKey(t *testing.T) {
	service, user := setupService()

	defer service.DB.Close()

	key := keyProto.ApiKeyRequest{
		UserId:      user.Id,
		Exchange:    "monex",
		Key:         "key",
		Secret:      "sssh",
		Description: "shit dwog!",
	}

	response := keyProto.ApiKeyResponse{}

	service.AddApiKey(context.Background(), &key, &response)

	if response.Status != "success" {
		t.Errorf(response.Message)
	}

	if response.Data.ApiKey.UserId != key.UserId {
		t.Errorf("user IDs do not match")
	}

	requestRemove := keyProto.RemoveApiKeyRequest{
		UserId:   user.Id,
		ApiKeyId: response.Data.ApiKey.ApiKeyId,
	}

	responseDel := keyProto.ApiKeyResponse{}
	service.RemoveApiKey(context.Background(), &requestRemove, &responseDel)

	if responseDel.Status != "success" {
		t.Errorf(responseDel.Message)
	}

	userRepo.DeleteUserHard(service.DB, user.Id)
}
