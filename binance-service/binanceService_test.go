package main

import (
	"context"
	"database/sql"
	"log"
	"testing"

	protoBalance "github.com/asciiu/gomo/binance-service/proto/balance"
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

func setupService() (*BinanceService, *sql.DB, *user.User) {
	dbUrl := "postgres://postgres@localhost:5432/gomo_test?&sslmode=disable"
	db, _ := db.NewDB(dbUrl)
	service := new(BinanceService)

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := repoUser.InsertUser(db, user)
	checkErr(error)

	return service, db, user
}

func TestInvalidKey1(t *testing.T) {
	service, db, user := setupService()

	defer db.Close()

	request := protoBalance.BalanceRequest{
		UserID:    user.ID,
		KeyPublic: "public",
		KeySecret: "secret",
	}

	response := protoBalance.BalancesResponse{}
	service.GetBalances(context.Background(), &request, &response)

	assert.Equal(t, "fail", response.Status, response.Message)

	repoUser.DeleteUserHard(db, user.ID)
}

func TestInvalidKey2(t *testing.T) {
	service, db, user := setupService()

	defer db.Close()

	request := protoBalance.BalanceRequest{
		UserID:    user.ID,
		KeyPublic: "Sn54bfgy5FILCvhtXSAlqqPhCgF74VLDlLYpJFNyVYeDMRiFCAo6g0F96CPb6xml",
		KeySecret: "AWxEQXvLQuyx218tZeeEHEWbfvdVXZ0zKjQgYEM3aDutkVIxQmtUeJWQVfHkPT1I",
	}

	response := protoBalance.BalancesResponse{}
	service.GetBalances(context.Background(), &request, &response)

	assert.Equal(t, "fail", response.Status, response.Message)

	repoUser.DeleteUserHard(db, user.ID)
}

func TestInvalidKey3(t *testing.T) {
	service, db, user := setupService()

	defer db.Close()

	request := protoBalance.BalanceRequest{
		UserID:    user.ID,
		KeyPublic: "O5oYc5b2TFSdcdWFqjQz8DnvVExeJUFeiGshmSVFet8WLHFVk3Iy1sQ5c",
		KeySecret: "cudE4yw1fQxrk5BfPXb4X5lKLZC3ypmIIWNOzUvE8e9p8sX40CnACo24nxID",
	}

	response := protoBalance.BalancesResponse{}
	service.GetBalances(context.Background(), &request, &response)

	assert.Equal(t, "fail", response.Status, response.Message)

	repoUser.DeleteUserHard(db, user.ID)
}
