package main

import (
	"context"
	"log"
	"testing"

	constAccount "github.com/asciiu/gomo/account-service/constants"
	protoAccount "github.com/asciiu/gomo/account-service/proto/account"
	protoBalance "github.com/asciiu/gomo/account-service/proto/balance"
	testBinance "github.com/asciiu/gomo/binance-service/test"
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
	service := AccountService{
		DB: db,
	}

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
				CurrencySymbol: "BTC",
				Available:      1.0,
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
				CurrencySymbol: "USDT",
				Available:      10.00,
			},
		},
	}

	response1 := protoAccount.AccountResponse{}
	service.AddAccount(context.Background(), &request, &response1)

	assert.Equal(t, "success", response1.Status, response1.Message)

	getRequest := protoAccount.AccountRequest{
		UserID:    user.ID,
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
				CurrencySymbol: "USDT",
				Available:      10.00,
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
				CurrencySymbol: "BTC",
				Available:      5.0,
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
				CurrencySymbol: "BTC",
				Available:      1.0,
			},
		},
	}

	response := protoAccount.AccountResponse{}
	service.AddAccount(context.Background(), &request, &response)

	assert.Equal(t, "success", response.Status, response.Message)

	delRequest := protoAccount.AccountRequest{
		UserID:    user.ID,
		AccountID: response.Data.Account.AccountID,
	}
	response2 := protoAccount.AccountResponse{}
	service.DeleteAccount(context.Background(), &delRequest, &response2)

	assert.Equal(t, "success", response2.Status, response2.Message)
	assert.Equal(t, constAccount.AccountDeleted, response2.Data.Account.Status)

	repoUser.DeleteUserHard(service.DB, user.ID)
}

// The AccountService should response with nonentity status when the accountID does not exist yet.
func TestGetAccountFailOnID(t *testing.T) {
	service, user := setupService()

	defer service.DB.Close()
	getRequest := protoAccount.AccountRequest{
		UserID:    user.ID,
		AccountID: uuid.New().String(),
	}
	response := protoAccount.AccountResponse{}
	service.GetAccount(context.Background(), &getRequest, &response)

	assert.Equal(t, constRes.Nonentity, response.Status, response.Message)

	repoUser.DeleteUserHard(service.DB, user.ID)
}

// The AccountService should response with nonentity status when the accountID does not exist during a delete request.
func TestDeleteAccountFailOnID(t *testing.T) {
	service, user := setupService()

	defer service.DB.Close()
	getRequest := protoAccount.AccountRequest{
		AccountID: uuid.New().String(),
		UserID:    user.ID,
	}
	response := protoAccount.AccountResponse{}
	service.DeleteAccount(context.Background(), &getRequest, &response)

	assert.Equal(t, constRes.Nonentity, response.Status, response.Message)

	repoUser.DeleteUserHard(service.DB, user.ID)
}

// The AccountService should should succeed with a valid update account request.
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
				CurrencySymbol: "USDT",
				Available:      10.00,
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

func TestSyncAccount(t *testing.T) {
	service, user := setupService()
	service.BinanceClient = testBinance.MockBinanceServiceClient()

	defer service.DB.Close()

	request := protoAccount.NewAccountRequest{
		UserID:      user.ID,
		Exchange:    "binance",
		KeyPublic:   "public",
		KeySecret:   "secret",
		Description: "shit test again!",
		Balances: []*protoBalance.NewBalanceRequest{
			&protoBalance.NewBalanceRequest{
				CurrencySymbol: "BTC",
				Available:      10.0,
			},
		},
	}

	response := protoAccount.AccountResponse{}
	service.AddAccount(context.Background(), &request, &response)

	btc1 := response.Data.Account.Balances[0]
	assert.Equal(t, "success", response.Status, response.Message)
	assert.Equal(t, "binance", response.Data.Account.Exchange, "exchange should be binance")

	assert.Equal(t, "BTC", btc1.CurrencySymbol, "currency symbol should be BTC")
	assert.Equal(t, 1.0, btc1.Available, "available should be 1.0")
	assert.Equal(t, 1.0, btc1.ExchangeAvailable, "available should be 1.0")
	assert.Equal(t, 2.0, btc1.ExchangeTotal, "exchange total after add account should be 2.0")

	// sync the account which will result in a binance service call for balances
	requestAccounts := protoAccount.AccountsRequest{UserID: user.ID}
	responseAccounts := protoAccount.AccountsResponse{}
	service.ResyncAccounts(context.Background(), &requestAccounts, &responseAccounts)

	btc := responseAccounts.Data.Accounts[0].Balances[0]
	assert.Equal(t, "BTC", btc.CurrencySymbol, "currency symbol should be BTC")
	assert.Equal(t, 1.0, btc.Available, "available should be 1.0")
	assert.Equal(t, 0.5, btc.ExchangeAvailable, "BTC exchange available should be 0.5")
	assert.Equal(t, 2.0, btc.ExchangeTotal, "BTC exchange total should be 2.0")

	usdt := responseAccounts.Data.Accounts[0].Balances[1]
	assert.Equal(t, "USDT", usdt.CurrencySymbol, "currency symbol should be USDT")
	assert.Equal(t, 100.0, usdt.Available, "usdt available should be 100.0")
	assert.Equal(t, 0.0, usdt.ExchangeAvailable, "usdt exchange available should be 0")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

func TestGetAccountBalance(t *testing.T) {
	service, user := setupService()
	service.BinanceClient = testBinance.MockBinanceServiceClient()

	defer service.DB.Close()

	request := protoAccount.NewAccountRequest{
		UserID:      user.ID,
		Exchange:    "binance",
		KeyPublic:   "public",
		KeySecret:   "secret",
		Description: "shit test again!",
		Balances: []*protoBalance.NewBalanceRequest{
			&protoBalance.NewBalanceRequest{
				CurrencySymbol: "BTC",
				Available:      10.0,
			},
		},
	}

	response := protoAccount.AccountResponse{}
	service.AddAccount(context.Background(), &request, &response)

	assert.Equal(t, "success", response.Status, response.Message)
	assert.Equal(t, "binance", response.Data.Account.Exchange, "exchange should be binance")

	getBalanceReq := protoBalance.BalanceRequest{
		UserID:         user.ID,
		AccountID:      response.Data.Account.AccountID,
		CurrencySymbol: "BTC",
	}
	balResponse := protoBalance.BalanceResponse{}
	service.GetAccountBalance(context.Background(), &getBalanceReq, &balResponse)

	balance := balResponse.Data.Balance
	assert.Equal(t, "BTC", balance.CurrencySymbol, "currency symbol should be BTC")
	assert.Equal(t, 1.0, balance.Available, "available should be 1.0")
	assert.Equal(t, 0.0, balance.Locked, "locked should be 0")
	assert.Equal(t, 1.0, balance.ExchangeAvailable, "available should be 1.0")
	assert.Equal(t, 2.0, balance.ExchangeTotal, "exchange total after add account should be 2.0")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

func TestValidateAccountBalance(t *testing.T) {
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
				CurrencySymbol: "BTC",
				Available:      10.0,
			},
		},
	}

	response := protoAccount.AccountResponse{}
	service.AddAccount(context.Background(), &request, &response)

	assert.Equal(t, "success", response.Status, response.Message)
	assert.Equal(t, "binance paper", response.Data.Account.Exchange, "exchange should be binance")

	vReq := protoBalance.ValidateBalanceRequest{
		UserID:          user.ID,
		AccountID:       response.Data.Account.AccountID,
		CurrencySymbol:  "BTC",
		RequestedAmount: 9,
	}
	res := protoBalance.ValidateBalanceResponse{}
	service.ValidateAccountBalance(context.Background(), &vReq, &res)

	assert.Equal(t, true, res.Data, "should be false")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

// This test checks to see if validate returns false when the balance does not exist.
func TestValidateAccountBalance2(t *testing.T) {
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
				CurrencySymbol: "BTC",
				Available:      10.0,
			},
		},
	}

	response := protoAccount.AccountResponse{}
	service.AddAccount(context.Background(), &request, &response)

	assert.Equal(t, "success", response.Status, response.Message)
	assert.Equal(t, "binance paper", response.Data.Account.Exchange, "exchange should be binance")

	vReq := protoBalance.ValidateBalanceRequest{
		UserID:          user.ID,
		AccountID:       response.Data.Account.AccountID,
		CurrencySymbol:  "USDT",
		RequestedAmount: 100,
	}
	res := protoBalance.ValidateBalanceResponse{}
	service.ValidateAccountBalance(context.Background(), &vReq, &res)

	assert.Equal(t, false, res.Data, "should be false")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

func TestAdjustAccountBalance(t *testing.T) {
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
				CurrencySymbol: "BTC",
				Available:      10.0,
			},
		},
	}

	response := protoAccount.AccountResponse{}
	service.AddAccount(context.Background(), &request, &response)

	assert.Equal(t, "success", response.Status, response.Message)
	assert.Equal(t, "binance paper", response.Data.Account.Exchange, "exchange should be binance")

	adjustReq := protoBalance.AdjustBalanceRequest{
		UserID:         user.ID,
		AccountID:      response.Data.Account.AccountID,
		CurrencySymbol: "BTC",
		Delta:          -1.0,
	}
	res := protoBalance.BalanceResponse{}
	service.AdjustBalance(context.Background(), &adjustReq, &res)

	balance := res.Data.Balance
	assert.Equal(t, 9.0, balance.Available, "should be 9 btc left")
	assert.Equal(t, 1.0, balance.Locked, "should be 1 locked")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

func TestAdjustAddBalance(t *testing.T) {
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
				CurrencySymbol: "BTC",
				Available:      10.0,
			},
		},
	}

	response := protoAccount.AccountResponse{}
	service.AddAccount(context.Background(), &request, &response)

	assert.Equal(t, "success", response.Status, response.Message)

	addReq := protoBalance.NewBalanceRequest{
		UserID:         user.ID,
		AccountID:      response.Data.Account.AccountID,
		CurrencySymbol: "USDT",
		Available:      1000.0,
	}
	res := protoBalance.BalanceResponse{}
	service.AddBalance(context.Background(), &addReq, &res)

	balance := res.Data.Balance
	assert.Equal(t, 1000.0, balance.Available, "should be 1000 usdt")
	assert.Equal(t, 0.0, balance.Locked, "should be 0 locked")

	getRequest := protoAccount.AccountRequest{
		UserID:    user.ID,
		AccountID: response.Data.Account.AccountID,
	}
	response2 := protoAccount.AccountResponse{}
	service.GetAccount(context.Background(), &getRequest, &response2)

	account := response2.Data.Account
	balances := account.Balances
	assert.Equal(t, "USDT", balances[1].CurrencySymbol, "currency did not match")
	assert.Equal(t, 1000.0, balances[1].Available, "available did not match")
	assert.Equal(t, 0.0, balances[1].Locked, "locked did not match")

	repoUser.DeleteUserHard(service.DB, user.ID)
}
