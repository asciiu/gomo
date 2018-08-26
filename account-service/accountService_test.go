package main

import (
	"context"
	"log"
	"testing"

	constAccount "github.com/asciiu/gomo/account-service/constants"
	protoAccount "github.com/asciiu/gomo/account-service/proto/account"
	protoBalance "github.com/asciiu/gomo/account-service/proto/balance"
	constRes "github.com/asciiu/gomo/common/constants/response"
	"github.com/asciiu/gomo/common/db"
	repoUser "github.com/asciiu/gomo/user-service/db/sql"
	user "github.com/asciiu/gomo/user-service/models"
	"github.com/google/uuid"
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

func TestGetAccounts(t *testing.T) {
	service, user := setupService()

	defer service.DB.Close()

	newRequest1 := protoAccount.NewAccountRequest{
		UserID:      user.ID,
		Exchange:    "binance paper",
		KeyPublic:   "public1",
		KeySecret:   "secret1",
		Description: "Test Account 1",
		Balances: []*protoBalance.NewBalanceRequest{
			&protoBalance.NewBalanceRequest{
				CurrencySymbol:  "USDT",
				CurrencyBalance: 10.00,
			},
		},
	}
	response1 := protoAccount.AccountResponse{}
	service.AddAccount(context.Background(), &newRequest1, &response1)
	assert.Equal(t, "success", response1.Status, response1.Message)

	newRequest2 := protoAccount.NewAccountRequest{
		UserID:      user.ID,
		Exchange:    "binance paper",
		KeyPublic:   "public2",
		KeySecret:   "secret2",
		Description: "Test Account 2",
		Balances: []*protoBalance.NewBalanceRequest{
			&protoBalance.NewBalanceRequest{
				CurrencySymbol:  "BTC",
				CurrencyBalance: 5.0,
			},
		},
	}
	response2 := protoAccount.AccountResponse{}
	service.AddAccount(context.Background(), &newRequest2, &response2)
	assert.Equal(t, "success", response2.Status, response2.Message)

	requestAccounts := protoAccount.AccountsRequest{
		UserID: user.ID,
	}
	response3 := protoAccount.AccountsResponse{}
	service.GetAccounts(context.Background(), &requestAccounts, &response3)

	assert.Equal(t, "success", response2.Status, response2.Message)
	assert.Equal(t, 2, len(response3.Data.Accounts), "should have two accounts")
	assert.Equal(t, "BTC", response3.Data.Accounts[1].Balances[0].CurrencySymbol, "currency for second account should be BTC")
	assert.Equal(t, "USDT", response3.Data.Accounts[0].Balances[0].CurrencySymbol, "currency for first account should be USDT")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

func TestDeleteAccount(t *testing.T) {
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

	delRequest := protoAccount.AccountRequest{
		AccountID: response.Data.Account.AccountID,
	}
	response2 := protoAccount.AccountResponse{}
	service.DeleteAccount(context.Background(), &delRequest, &response2)

	assert.Equal(t, "success", response2.Status, response2.Message)
	assert.Equal(t, constAccount.AccountDeleted, response2.Data.Account.Status)

	repoUser.DeleteUserHard(service.DB, user.ID)
}

func TestGetAccountFailOnID(t *testing.T) {
	service, user := setupService()

	defer service.DB.Close()
	getRequest := protoAccount.AccountRequest{
		AccountID: uuid.New().String(),
	}
	response := protoAccount.AccountResponse{}
	service.GetAccount(context.Background(), &getRequest, &response)

	assert.Equal(t, constRes.Nonentity, response.Status, response.Message)

	repoUser.DeleteUserHard(service.DB, user.ID)
}

func TestDeleteAccountFailOnID(t *testing.T) {
	service, user := setupService()

	defer service.DB.Close()
	getRequest := protoAccount.AccountRequest{
		AccountID: uuid.New().String(),
	}
	response := protoAccount.AccountResponse{}
	service.DeleteAccount(context.Background(), &getRequest, &response)

	assert.Equal(t, constRes.Nonentity, response.Status, response.Message)

	repoUser.DeleteUserHard(service.DB, user.ID)
}

func TestAccountUpdate(t *testing.T) {
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

	updateRequest := protoAccount.UpdateAccountRequest{
		AccountID:   response1.Data.Account.AccountID,
		UserID:      user.ID,
		KeyPublic:   "public2",
		KeySecret:   "secret2",
		Description: "description2",
	}

	response2 := protoAccount.AccountResponse{}
	service.UpdateAccount(context.Background(), &updateRequest, &response2)

	assert.Equal(t, "success", response2.Status, response2.Message)
	assert.Equal(t, response1.Data.Account.Exchange, response2.Data.Account.Exchange, "exchanges do not match")
	assert.Equal(t, "public2", response2.Data.Account.KeyPublic, "public keys don't match")
	assert.Equal(t, "description2", response2.Data.Account.Description, "descriptions don't match")
	repoUser.DeleteUserHard(service.DB, user.ID)
}
