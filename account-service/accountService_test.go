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

func TestNewAccount(t *testing.T) {
	service, user := setupService()

	defer service.DB.Close()

	request := protoAccount.NewAccountRequest{
		UserID:      user.ID,
		Exchange:    "binance paper",
		KeyPublic:   "public",
		KeySecret:   "secret",
		Description: "shit test again!",
		Balances: []*protoBalance.NewBalanceRequest{
			&protoBalance.NewBalanceRequest{
				CurrencySymbol:  "BTC",
				CurrencyBalance: 1.0,
			},
		},
	}

	response := protoAccount.AccountResponse{}
	service.AddAccount(context.Background(), &request, &response)

	assert.Equal(t, "success", response.Status, response.Message)
	assert.Equal(t, 1.0, response.Data.Account.Balances[0].Available, "available should be 1.0")
	assert.Equal(t, "BTC", response.Data.Account.Balances[0].CurrencySymbol, "currency symbol should be BTC")
	assert.Equal(t, "binance paper", response.Data.Account.Exchange, "exchange should be binance paper")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

func TestGetAccount(t *testing.T) {
	service, user := setupService()

	defer service.DB.Close()

	request := protoAccount.NewAccountRequest{
		UserID:      user.ID,
		Exchange:    "binance paper",
		KeyPublic:   "public",
		KeySecret:   "secret",
		Description: "shit test again!",
		Balances: []*protoBalance.NewBalanceRequest{
			&protoBalance.NewBalanceRequest{
				CurrencySymbol:  "USDT",
				CurrencyBalance: 10.00,
			},
		},
	}

	response1 := protoAccount.AccountResponse{}
	service.AddAccount(context.Background(), &request, &response1)

	assert.Equal(t, "success", response1.Status, response1.Message)

	getRequest := protoAccount.AccountRequest{
		AccountID: response1.Data.Account.AccountID,
	}
	response2 := protoAccount.AccountResponse{}
	service.GetAccount(context.Background(), &getRequest, &response2)

	assert.Equal(t, "success", response2.Status, response2.Message)
	assert.Equal(t, response1.Data.Account.AccountID, response2.Data.Account.AccountID, "account IDs do not match")
	assert.Equal(t, response1.Data.Account.UserID, response2.Data.Account.UserID, "user IDs do not match")
	assert.Equal(t, response1.Data.Account.Description, response2.Data.Account.Description, "descriptions do not match")
	assert.Equal(t, response1.Data.Account.Exchange, response2.Data.Account.Exchange, "exchanges do not match")
	assert.Equal(t, "USDT", response2.Data.Account.Balances[0].CurrencySymbol, "currency did not match")

	repoUser.DeleteUserHard(service.DB, user.ID)
}
