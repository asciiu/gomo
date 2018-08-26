package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	constAccount "github.com/asciiu/gomo/account-service/constants"
	repoAccount "github.com/asciiu/gomo/account-service/db/sql"
	protoAccount "github.com/asciiu/gomo/account-service/proto/account"
	protoBalance "github.com/asciiu/gomo/account-service/proto/balance"
	constExch "github.com/asciiu/gomo/common/constants/exchange"
	constRes "github.com/asciiu/gomo/common/constants/response"
	"github.com/google/uuid"
	"github.com/lib/pq"
	micro "github.com/micro/go-micro"
)

type AccountService struct {
	DB        *sql.DB
	KeyPub    micro.Publisher
	NotifyPub micro.Publisher
}

// IMPORTANT! When adding support for new exchanges we must add there names here!
// TODO: read these from a config or from the DB
var supported = [...]string{
	constExch.Binance,
	constExch.BinancePaper,
}

func supportedExchange(a string) bool {
	for _, b := range supported {
		if b == a {
			return true
		}
	}
	return false
}

// func (service *AccountService) HandleVerifiedKey(ctx context.Context, key *protoKey.Key) error {
// 	log.Println("received verified key ", key.KeyID)

// 	_, error := repoKey.UpdateKeyStatus(service.DB, key)

// 	description := fmt.Sprintf("%s key verified", key.Exchange)
// 	notification := protoActivity.Activity{
// 		UserID:      key.UserID,
// 		Type:        "key",
// 		ObjectID:    key.KeyID,
// 		Title:       "Exchange Setup",
// 		Description: description,
// 		Timestamp:   time.Now().UTC().Format(time.RFC3339),
// 	}

// 	// publish verify key event
// 	if err := service.NotifyPub.Publish(context.Background(), &notification); err != nil {
// 		log.Println("could not publish verified key event: ", err)
// 	}

// 	return error
// }

// Add a new account. An account may be real or paper. Paper accounts do not need to be verfied.
func (service *AccountService) AddAccount(ctx context.Context, req *protoAccount.NewAccountRequest, res *protoAccount.AccountResponse) error {
	// supported exchange keys check
	if !supportedExchange(req.Exchange) {
		res.Status = constRes.Fail
		res.Message = fmt.Sprintf("%s is not supported", req.Exchange)
		return nil
	}

	accountID := uuid.New().String()
	now := string(pq.FormatTimestamp(time.Now().UTC()))
	balances := make([]*protoBalance.Balance, 0, len(req.Balances))

	for _, b := range req.Balances {
		balanceID := uuid.New()
		balance := protoBalance.Balance{
			BalanceID:      balanceID.String(),
			UserID:         req.UserID,
			AccountID:      accountID,
			CurrencySymbol: b.CurrencySymbol,
			Available:      b.CurrencyBalance,
			Locked:         0,
			CreatedOn:      now,
			UpdatedOn:      now,
		}
		balances = append(balances, &balance)
	}

	status := constAccount.AccountUnverified
	if strings.Contains(req.Exchange, "paper") {
		status = constAccount.AccountVerified
	}

	account := protoAccount.Account{
		AccountID:   accountID,
		UserID:      req.UserID,
		Exchange:    req.Exchange,
		KeyPublic:   req.KeyPublic,
		KeySecret:   req.KeySecret,
		Description: req.Description,
		Status:      status,
		CreatedOn:   now,
		UpdatedOn:   now,
		Balances:    balances,
	}

	if err := repoAccount.InsertAccount(service.DB, &account); err != nil {
		msg := fmt.Sprintf("insert account failed %s", err.Error())
		log.Println(msg)

		res.Status = constRes.Error
		res.Message = msg
	}

	res.Status = constRes.Success
	res.Data = &protoAccount.UserAccount{Account: &account}

	return nil
}

func (service *AccountService) DeleteAccount(ctx context.Context, req *protoAccount.AccountRequest, res *protoAccount.AccountResponse) error {
	account, err := repoAccount.UpdateAccountStatus(service.DB, req.AccountID, constAccount.AccountDeleted)

	switch {
	case err == sql.ErrNoRows:
		res.Status = constRes.Nonentity
		res.Message = "account not found"
	case err != nil:
		log.Println("GetAccount error: ", err.Error())
		res.Status = constRes.Error
		res.Message = err.Error()
	default:
		res.Status = constRes.Success
		res.Data = &protoAccount.UserAccount{
			Account: account,
		}
	}

	return nil
}

func (service *AccountService) GetAccounts(ctx context.Context, req *protoAccount.AccountsRequest, res *protoAccount.AccountsResponse) error {
	accounts, err := repoAccount.FindAccounts(service.DB, req.UserID)

	if err != nil {
		log.Println("GetAccount error: ", err.Error())
		res.Status = constRes.Error
		res.Message = err.Error()
	} else {
		res.Status = constRes.Success
		res.Data = &protoAccount.UserAccounts{
			Accounts: accounts,
		}
	}

	return nil
}

// Get a specific account by ID.
func (service *AccountService) GetAccount(ctx context.Context, req *protoAccount.AccountRequest, res *protoAccount.AccountResponse) error {
	account, err := repoAccount.FindAccount(service.DB, req.AccountID)

	switch {
	case err == sql.ErrNoRows:
		res.Status = constRes.Nonentity
		res.Message = "no account by that ID"
	case err != nil:
		log.Println("GetAccount error: ", err.Error())
		res.Status = constRes.Error
		res.Message = err.Error()
	default:
		res.Status = constRes.Success
		res.Data = &protoAccount.UserAccount{
			Account: account,
		}
	}

	return nil
}

func (service *AccountService) UpdateAccount(ctx context.Context, req *protoAccount.UpdateAccountRequest, res *protoAccount.AccountResponse) error {
	return nil
}
