package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	constAccount "github.com/asciiu/gomo/account-service/constants"
	repoAccount "github.com/asciiu/gomo/account-service/db/sql"
	protoAccount "github.com/asciiu/gomo/account-service/proto/account"
	protoBalance "github.com/asciiu/gomo/account-service/proto/balance"
	protoBinanceBal "github.com/asciiu/gomo/binance-service/proto/balance"
	protoBinance "github.com/asciiu/gomo/binance-service/proto/binance"
	constExch "github.com/asciiu/gomo/common/constants/exchange"
	constRes "github.com/asciiu/gomo/common/constants/response"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AccountService struct {
	DB            *sql.DB
	BinanceClient protoBinance.BinanceServiceClient
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
		balance := protoBalance.Balance{
			UserID:         req.UserID,
			AccountID:      accountID,
			CurrencySymbol: b.CurrencySymbol,
			Available:      b.Available,
			Locked:         0,
			CreatedOn:      now,
			UpdatedOn:      now,
		}
		balances = append(balances, &balance)
	}

	// assume account verified
	account := protoAccount.Account{
		AccountID:   accountID,
		UserID:      req.UserID,
		Exchange:    req.Exchange,
		KeyPublic:   req.KeyPublic,
		KeySecret:   req.KeySecret,
		Description: req.Description,
		Status:      constAccount.AccountVerified,
		CreatedOn:   now,
		UpdatedOn:   now,
		Balances:    balances,
	}

	// validate account request when keys are present
	switch {
	case account.KeyPublic != "" && account.KeySecret == "":
		res.Status = constRes.Fail
		res.Message = "keySecret required with keyPublic!"
		return nil
	case account.KeyPublic == "" && account.KeySecret != "":
		res.Status = constRes.Fail
		res.Message = "keyPublic required with keySecret!"
		return nil
	}

	switch account.Exchange {
	case constExch.Binance:
		// if api key ask exchange for balances
		if account.KeyPublic == "" || account.KeySecret == "" {
			res.Status = constRes.Fail
			res.Message = "keyPublic and keySecret required!"
			return nil
		}
		reqBal := protoBinanceBal.BalanceRequest{
			UserID:    account.UserID,
			KeyPublic: account.KeyPublic,
			KeySecret: account.KeySecret,
		}
		resBal, _ := service.BinanceClient.GetBalances(ctx, &reqBal)

		// reponse to client on invalid key
		if resBal.Status != constRes.Success {
			res.Status = resBal.Status
			res.Message = resBal.Message
			return nil
		}

		balances := make([]*protoBalance.Balance, 0)
		for _, b := range resBal.Data.Balances {
			total := b.Free + b.Locked
			balance := protoBalance.Balance{
				UserID:            account.UserID,
				AccountID:         account.AccountID,
				CurrencySymbol:    b.CurrencySymbol,
				Available:         b.Free,
				Locked:            0.0,
				ExchangeTotal:     total,
				ExchangeAvailable: b.Free,
				ExchangeLocked:    b.Locked,
			}
			balances = append(balances, &balance)
		}
		account.Balances = balances
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

func (service *AccountService) AddAccountBalance(ctx context.Context, req *protoBalance.NewBalanceRequest, res *protoBalance.BalanceResponse) error {
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

func (service *AccountService) GetAccountBalance(ctx context.Context, req *protoBalance.BalanceRequest, res *protoBalance.BalanceResponse) error {
	return nil
}

func (service *AccountService) UpdateAccount(ctx context.Context, req *protoAccount.UpdateAccountRequest, res *protoAccount.AccountResponse) error {
	account, err := repoAccount.FindAccount(service.DB, req.AccountID)

	if err == sql.ErrNoRows {
		res.Status = constRes.Nonentity
		res.Message = "no account by that ID"
		return nil
	}

	txn, err := service.DB.Begin()
	if err != nil {
		res.Status = constRes.Error
		res.Message = err.Error()
		log.Println("UpdateAccount on begin transaction: ", err.Error())
		return nil
	}

	if err := repoAccount.UpdateAccountTxn(txn, ctx, req.AccountID, req.KeyPublic, req.KeySecret, req.Description); err != nil {
		txn.Rollback()
		res.Status = constRes.Error
		res.Message = "error encountered while updating account: " + err.Error()
		log.Println("UpdateAccountTxn error: ", err.Error())
		return nil
	}

	txn.Commit()
	account.KeyPublic = req.KeyPublic
	account.KeySecret = req.KeySecret
	account.Description = req.Description

	res.Status = constRes.Success
	res.Data = &protoAccount.UserAccount{
		Account: account,
	}

	return nil
}

func (service *AccountService) UpdateAccountBalance(ctx context.Context, req *protoBalance.UpdateBalanceRequest, res *protoBalance.BalanceResponse) error {
	return nil
}

func (service *AccountService) ValidateAccountBalance(ctx context.Context, req *protoBalance.ValidateBalanceRequest, res *protoBalance.ValidateBalanceResponse) error {
	return nil
}
