package main

import (
	"context"
	"log"
	"testing"

	protoAccount "github.com/asciiu/gomo/account-service/proto/account"
	protoBalance "github.com/asciiu/gomo/account-service/proto/balance"
	"github.com/asciiu/gomo/common/db"
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

func setupService() (*AccountService, *user.User) {
	dbUrl := "postgres://postgres@localhost:5432/gomo_test?&sslmode=disable"
	db, _ := db.NewDB(dbUrl)
	service := AccountService{DB: db}

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := repoUser.InsertUser(db, user)
	checkErr(error)

	return &service, user
}

func TestUnsupportedExchange(t *testing.T) {
	service, user := setupService()

	defer service.DB.Close()

	request := protoAccount.NewAccountRequest{
		UserID:      user.ID,
		Exchange:    "monex",
		KeyPublic:   "public",
		KeySecret:   "secret",
		Description: "shit test again!",
		Balances:    make([]*protoBalance.NewBalanceRequest, 0),
	}

	response := protoAccount.AccountResponse{}
	service.AddAccount(context.Background(), &request, &response)

	assert.Equal(t, "fail", response.Status, "return status is expected to be fail due to unsupported exchange")

	repoUser.DeleteUserHard(service.DB, user.ID)
}
