package main

import (
	"context"
	"log"
	"testing"

	"github.com/asciiu/gomo/common/db"
	keyProto "github.com/asciiu/gomo/key-service/proto/key"
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

func setupService() (*KeyService, *user.User) {
	dbUrl := "postgres://postgres@localhost:5432/gomo_test?&sslmode=disable"
	db, _ := db.NewDB(dbUrl)
	service := KeyService{DB: db}

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := repoUser.InsertUser(db, user)
	checkErr(error)

	return &service, user
}

func TestUnsupportedKey(t *testing.T) {
	service, user := setupService()

	defer service.DB.Close()

	key := keyProto.KeyRequest{
		UserID:      user.ID,
		Exchange:    "monex",
		Key:         "public",
		Secret:      "secret",
		Description: "shit test!",
	}

	response := keyProto.KeyResponse{}
	service.AddKey(context.Background(), &key, &response)

	assert.Equal(t, "fail", response.Status, "return status is expected to be fail due to unsupported exchange")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

func TestSuccessfulKey(t *testing.T) {
	service, user := setupService()

	defer service.DB.Close()

	key := keyProto.KeyRequest{
		UserID:      user.ID,
		Exchange:    "binance",
		Key:         "public",
		Secret:      "secret",
		Description: "shit test 2",
	}

	response := keyProto.KeyResponse{}
	service.AddKey(context.Background(), &key, &response)

	assert.Equal(t, "success", response.Status, response.Message)
	assert.Equal(t, key.UserID, response.Data.Key.UserID, "user ids do not match")

	requestRemove := keyProto.RemoveKeyRequest{
		UserID: user.ID,
		KeyID:  response.Data.Key.KeyID,
	}
	responseDel := keyProto.KeyResponse{}
	service.RemoveKey(context.Background(), &requestRemove, &responseDel)

	assert.Equal(t, "success", responseDel.Status, response.Message)

	repoUser.DeleteUserHard(service.DB, user.ID)
}
