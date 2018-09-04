package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

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

// rot32768 rotates utf8 string
// rot5map rotates digits
func rot32768(input string) string {
	var result []string
	rot5map := map[rune]rune{'0': '5', '1': '6', '2': '7', '3': '8', '4': '9', '5': '0', '6': '1', '7': '2', '8': '3', '9': '4'}

	for _, i := range input {
		switch {
		case unicode.IsSpace(i):
			result = append(result, " ")
		case i >= '0' && i <= '9':
			result = append(result, string(rot5map[i]))
		case utf8.ValidRune(i):
			//result = append(result, string(rune(i) ^ 0x80))
			result = append(result, string(rune(i)^utf8.RuneSelf))
		}

	}

	return strings.Join(result, "")
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

	// user specified balances will be ignored if a
	// valid public/secret is send in with request
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
		KeySecret:   rot32768(req.KeySecret),
		Description: req.Description,
		Status:      constAccount.AccountValid,
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
			KeySecret: rot32768(account.KeySecret),
		}
		resBal, _ := service.BinanceClient.GetBalances(ctx, &reqBal)

		// reponse to client on invalid key
		if resBal.Status != constRes.Success {
			res.Status = resBal.Status
			res.Message = resBal.Message
			return nil
		}

		exBalances := make([]*protoBalance.Balance, 0)
		for _, b := range resBal.Data.Balances {
			total := b.Free + b.Locked

			// only add non-zero balances
			if total > 0 {
				balance := protoBalance.Balance{
					UserID:            account.UserID,
					AccountID:         account.AccountID,
					CurrencySymbol:    b.CurrencySymbol,
					Available:         b.Free,
					Locked:            0.0,
					ExchangeTotal:     total,
					ExchangeAvailable: b.Free,
					ExchangeLocked:    b.Locked,
					CreatedOn:         now,
					UpdatedOn:         now,
				}

				exBalances = append(exBalances, &balance)
			}
		}
		account.Balances = exBalances
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

func (service *AccountService) AddBalance(ctx context.Context, req *protoBalance.NewBalanceRequest, res *protoBalance.BalanceResponse) error {
	balance, err := repoAccount.FindAccountBalance(service.DB, req.UserID, req.AccountID, req.CurrencySymbol)
	switch {
	case err == sql.ErrNoRows:
		now := string(pq.FormatTimestamp(time.Now().UTC()))
		balance := protoBalance.Balance{
			UserID:         req.UserID,
			AccountID:      req.AccountID,
			CurrencySymbol: req.CurrencySymbol,
			Available:      req.Available,
			Locked:         0,
			CreatedOn:      now,
			UpdatedOn:      now,
		}
		if err := repoAccount.InsertBalance(service.DB, &balance); err != nil {
			log.Println("AddBalance on InsertBalance error: ", err.Error())
			res.Status = constRes.Error
			res.Message = err.Error()
		}
		res.Status = constRes.Success
		res.Data = &protoBalance.BalanceData{
			Balance: &balance,
		}

	case err != nil:
		// this should not happen log it
		log.Println("AddBalance on FindAccountBalance error: ", err.Error())
		res.Status = constRes.Error
		res.Message = err.Error()
	default:
		available := balance.Available + req.Available

		balance, err = repoAccount.UpdateAvailableBalance(service.DB, req.UserID, req.AccountID, req.CurrencySymbol, available)
		if err != nil {
			log.Println("AddBalance on UpdateAvailableBalance error: ", err.Error())
			res.Status = constRes.Error
			res.Message = err.Error()
		}

		res.Status = constRes.Success
		res.Data = &protoBalance.BalanceData{
			Balance: balance,
		}
	}
	return nil
}

func (service *AccountService) ChangeAvailableBalance(ctx context.Context, req *protoBalance.ChangeBalanceRequest, res *protoBalance.BalanceResponse) error {
	balance, err := repoAccount.FindAccountBalance(service.DB, req.UserID, req.AccountID, req.CurrencySymbol)
	switch {
	case err == sql.ErrNoRows:
		now := string(pq.FormatTimestamp(time.Now().UTC()))
		balance := protoBalance.Balance{
			UserID:         req.UserID,
			AccountID:      req.AccountID,
			CurrencySymbol: req.CurrencySymbol,
			Available:      req.Amount,
			Locked:         0,
			CreatedOn:      now,
			UpdatedOn:      now,
		}
		if err := repoAccount.InsertBalance(service.DB, &balance); err != nil {
			log.Println("ChangeAvailableBalance on InsertBalance error: ", err.Error())
			res.Status = constRes.Error
			res.Message = err.Error()
		}
		res.Status = constRes.Success
		res.Data = &protoBalance.BalanceData{
			Balance: &balance,
		}

	case err != nil:
		log.Println("ChangeAvailableBalance on FindAccountBalance error: ", err.Error())
		res.Status = constRes.Error
		res.Message = err.Error()
	default:
		available := balance.Available + req.Amount

		balance, err = repoAccount.UpdateAvailableBalance(service.DB, req.UserID, req.AccountID, req.CurrencySymbol, available)
		if err != nil {
			log.Println("ChangeAvailableBalance on UpdateAvailableBalance error: ", err.Error())
			res.Status = constRes.Error
			res.Message = err.Error()
		}

		res.Status = constRes.Success
		res.Data = &protoBalance.BalanceData{
			Balance: balance,
		}
	}
	return nil
}

func (service *AccountService) ChangeLockedBalance(ctx context.Context, req *protoBalance.ChangeBalanceRequest, res *protoBalance.BalanceResponse) error {
	balance, err := repoAccount.FindAccountBalance(service.DB, req.UserID, req.AccountID, req.CurrencySymbol)
	switch {
	case err == sql.ErrNoRows:
		// currency not found yet create a balance for it
		now := string(pq.FormatTimestamp(time.Now().UTC()))
		balance := protoBalance.Balance{
			UserID:         req.UserID,
			AccountID:      req.AccountID,
			CurrencySymbol: req.CurrencySymbol,
			Available:      0,
			Locked:         req.Amount,
			CreatedOn:      now,
			UpdatedOn:      now,
		}
		if err := repoAccount.InsertBalance(service.DB, &balance); err != nil {
			log.Println("ChangeLockedBalance on InsertBalance error: ", err.Error())
			res.Status = constRes.Error
			res.Message = err.Error()
		}
		res.Status = constRes.Success
		res.Data = &protoBalance.BalanceData{
			Balance: &balance,
		}
	case err != nil:
		log.Println("ChangeLockedBalance on FindAccountBalance error: ", err.Error())
		res.Status = constRes.Error
		res.Message = err.Error()
	default:
		locked := balance.Locked + req.Amount

		balance, err = repoAccount.UpdateLockedBalance(service.DB, req.UserID, req.AccountID, req.CurrencySymbol, locked)
		if err != nil {
			log.Println("ChangeLockedBalance on UpdateLockedBalance error: ", err.Error())
			res.Status = constRes.Error
			res.Message = err.Error()
		}

		res.Status = constRes.Success
		res.Data = &protoBalance.BalanceData{
			Balance: balance,
		}
	}
	return nil
}

func (service *AccountService) LockBalance(ctx context.Context, req *protoBalance.ChangeBalanceRequest, res *protoBalance.BalanceResponse) error {
	balance, err := repoAccount.FindAccountBalance(service.DB, req.UserID, req.AccountID, req.CurrencySymbol)

	switch {
	case err == sql.ErrNoRows:
		now := string(pq.FormatTimestamp(time.Now().UTC()))
		balance := protoBalance.Balance{
			UserID:         req.UserID,
			AccountID:      req.AccountID,
			CurrencySymbol: req.CurrencySymbol,
			Available:      0,
			Locked:         req.Amount,
			CreatedOn:      now,
			UpdatedOn:      now,
		}
		if err := repoAccount.InsertBalance(service.DB, &balance); err != nil {
			log.Println("LockBalance on InsertBalance error: ", err.Error())
			res.Status = constRes.Error
			res.Message = err.Error()
		}
		res.Status = constRes.Success
		res.Data = &protoBalance.BalanceData{
			Balance: &balance,
		}
	case err != nil:
		log.Println("LockBalance error on FindAccountBalance: ", err.Error())
		res.Status = constRes.Error
		res.Message = err.Error()
	default:
		available := balance.Available - req.Amount
		locked := balance.Locked + req.Amount

		balance, err = repoAccount.UpdateInternalBalance(service.DB, req.UserID, req.AccountID, req.CurrencySymbol, available, locked)
		if err != nil {
			log.Println("LockBalance error on UpdateInternalBalance: ", err.Error())
			res.Status = constRes.Error
			res.Message = err.Error()
		}

		res.Status = constRes.Success
		res.Data = &protoBalance.BalanceData{
			Balance: balance,
		}
	}
	return nil
}

func (service *AccountService) UnlockBalance(ctx context.Context, req *protoBalance.ChangeBalanceRequest, res *protoBalance.BalanceResponse) error {
	balance, err := repoAccount.FindAccountBalance(service.DB, req.UserID, req.AccountID, req.CurrencySymbol)

	switch {
	case err == sql.ErrNoRows:
		now := string(pq.FormatTimestamp(time.Now().UTC()))
		balance := protoBalance.Balance{
			UserID:         req.UserID,
			AccountID:      req.AccountID,
			CurrencySymbol: req.CurrencySymbol,
			Available:      req.Amount,
			Locked:         0,
			CreatedOn:      now,
			UpdatedOn:      now,
		}
		if err := repoAccount.InsertBalance(service.DB, &balance); err != nil {
			log.Println("UnlockBalance on InsertBalance error: ", err.Error())
			res.Status = constRes.Error
			res.Message = err.Error()
		}
		res.Status = constRes.Success
		res.Data = &protoBalance.BalanceData{
			Balance: &balance,
		}
	case err != nil:
		log.Println("FindAccountBalance error: ", err.Error())
		res.Status = constRes.Error
		res.Message = err.Error()
	default:
		available := balance.Available + req.Amount
		locked := balance.Locked - req.Amount

		balance, err = repoAccount.UpdateInternalBalance(service.DB, req.UserID, req.AccountID, req.CurrencySymbol, available, locked)
		if err != nil {
			log.Println("UpdateInternalBalance error: ", err.Error())
			res.Status = constRes.Error
			res.Message = err.Error()
		}

		res.Status = constRes.Success
		res.Data = &protoBalance.BalanceData{
			Balance: balance,
		}
	}
	return nil
}

func (service *AccountService) DeleteAccount(ctx context.Context, req *protoAccount.AccountRequest, res *protoAccount.AccountResponse) error {
	account, err := repoAccount.UpdateAccountStatus(service.DB, req.AccountID, req.UserID, constAccount.AccountDeleted)

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
	account, err := repoAccount.FindAccount(service.DB, req.AccountID, req.UserID)

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

// Get a specific account by ID.
func (service *AccountService) GetAccountKeys(ctx context.Context, req *protoAccount.GetAccountKeysRequest, res *protoAccount.AccountKeysResponse) error {
	keys, err := repoAccount.FindAccountKeys(service.DB, req.AccountIDs)

	switch {
	case err == sql.ErrNoRows:
		res.Status = constRes.Nonentity
		res.Message = "no accounts by those IDs"
	case err != nil:
		log.Println("FindAccountKeys error: ", err.Error())
		res.Status = constRes.Error
		res.Message = err.Error()
	default:
		res.Status = constRes.Success
		res.Data = &protoAccount.KeysList{
			Keys: keys,
		}
	}

	return nil
}

func (service *AccountService) GetAccountBalance(ctx context.Context, req *protoBalance.BalanceRequest, res *protoBalance.BalanceResponse) error {
	balance, err := repoAccount.FindAccountBalance(service.DB, req.UserID, req.AccountID, req.CurrencySymbol)

	switch {
	case err == sql.ErrNoRows:
		res.Status = constRes.Nonentity
		res.Message = "balance not found"
	case err != nil:
		log.Println("FindAccountBalance error: ", err.Error())
		res.Status = constRes.Error
		res.Message = err.Error()
	default:
		res.Status = constRes.Success
		res.Data = &protoBalance.BalanceData{
			Balance: balance,
		}
	}
	return nil
}

func (service *AccountService) resyncBinanceBalances(ctx context.Context, account *protoAccount.Account) *protoAccount.Account {
	reqBal := protoBinanceBal.BalanceRequest{
		UserID:    account.UserID,
		KeyPublic: account.KeyPublic,
		KeySecret: account.KeySecret,
	}
	exBals, _ := service.BinanceClient.GetBalances(ctx, &reqBal)

	// invalid exchange keys?
	if exBals.Status == constRes.Fail {
		_, err := repoAccount.UpdateAccountStatus(service.DB, account.AccountID, account.UserID, constAccount.AccountInvalid)
		if err != nil {
			log.Println("ResyncBinanceBalances on UpdateAccountStatus to invalid: ", err.Error())
		}
		// update status of account to invalid
		account.Status = constAccount.AccountInvalid
		return account
	}
	// not supposed to happen, but if it does log error
	if exBals.Status == constRes.Error {
		log.Println("ResyncBinanceBalances on GetBalances: ", exBals.Message)
		// log error
		return account
	}

	txn, err := service.DB.Begin()
	if err != nil {
		// not supposed to happen, but if it does log error
		log.Println("ResyncBinanceBalances on begin transaction: ", err.Error())
		// log error
		return account
	}

	now := string(pq.FormatTimestamp(time.Now().UTC()))
	newBalances := make([]*protoBalance.Balance, 0)
	for _, exBal := range exBals.Data.Balances {
		total := exBal.Free + exBal.Locked

		found := false
		for _, accBal := range account.Balances {
			// if the balance already exists
			if accBal.CurrencySymbol == exBal.CurrencySymbol {
				// if total, locked, or available has changed on the exchange we need to update our balance for this currency
				if accBal.ExchangeTotal != total || accBal.ExchangeLocked != exBal.Locked || accBal.ExchangeAvailable != exBal.Free {
					// update the user's balance for this account's currency
					if err := repoAccount.UpdateExchangeBalanceTxn(txn, ctx, account.AccountID, account.UserID, accBal.CurrencySymbol, total, exBal.Free, exBal.Locked); err != nil {
						txn.Rollback()
						log.Println("ResyncBinanceBalances on UpdateExchangeBalanceTxn: ", err.Error())
						return account
					}
					accBal.ExchangeTotal = total
					accBal.ExchangeLocked = exBal.Locked
					accBal.ExchangeAvailable = exBal.Free
				}
				found = true
				break
			}
		}
		if !found {
			balance := protoBalance.Balance{
				UserID:            account.UserID,
				AccountID:         account.AccountID,
				CurrencySymbol:    exBal.CurrencySymbol,
				Available:         exBal.Free,
				Locked:            0.0,
				ExchangeTotal:     total,
				ExchangeAvailable: exBal.Free,
				ExchangeLocked:    exBal.Locked,
				CreatedOn:         now,
				UpdatedOn:         now,
			}
			newBalances = append(newBalances, &balance)
		}
	}
	if len(newBalances) > 0 {
		if err := repoAccount.InsertBalances(txn, newBalances); err != nil {
			txn.Rollback()
			log.Println("ResyncBinanceBalances on InsertBalances: ", err.Error())
			return account
		}
		// append new balances to account balances
		account.Balances = append(account.Balances, newBalances...)
	}
	txn.Commit()

	return account
}

func (service *AccountService) ResyncAccounts(ctx context.Context, req *protoAccount.AccountsRequest, res *protoAccount.AccountsResponse) error {
	accounts, err := repoAccount.FindAccounts(service.DB, req.UserID)

	if err != nil {
		log.Println("ResyncAccounts error: ", err.Error())
		res.Status = constRes.Error
		res.Message = err.Error()
		return nil
	}

	// loop through each account
	for i, acc := range accounts {
		if strings.Contains(acc.Exchange, "paper") {
			continue
		}

		switch acc.Exchange {
		case constExch.Binance:
			accounts[i] = service.resyncBinanceBalances(ctx, acc)
		}
	}
	res.Status = constRes.Success
	res.Data = &protoAccount.UserAccounts{
		Accounts: accounts,
	}

	return nil
}

func (service *AccountService) UpdateAccount(ctx context.Context, req *protoAccount.UpdateAccountRequest, res *protoAccount.AccountResponse) error {
	// validate request when keys are present
	switch {
	case req.KeyPublic != "" && req.KeySecret == "":
		res.Status = constRes.Fail
		res.Message = "keySecret required with keyPublic!"
		return nil
	case req.KeyPublic == "" && req.KeySecret != "":
		res.Status = constRes.Fail
		res.Message = "keyPublic required with keySecret!"
		return nil
	}

	account, err := repoAccount.FindAccount(service.DB, req.AccountID, req.UserID)

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

	// if new key
	if req.KeyPublic != "" {
		now := string(pq.FormatTimestamp(time.Now().UTC()))

		switch account.Exchange {
		case constExch.Binance:
			reqBal := protoBinanceBal.BalanceRequest{
				UserID:    account.UserID,
				KeyPublic: req.KeyPublic,
				KeySecret: req.KeySecret,
			}
			resBal, _ := service.BinanceClient.GetBalances(ctx, &reqBal)

			// response to client on invalid key
			if resBal.Status != constRes.Success {
				res.Status = resBal.Status
				res.Message = resBal.Message
				return nil
			}

			newBalances := make([]*protoBalance.Balance, 0)
			for _, b := range resBal.Data.Balances {
				total := b.Free + b.Locked

				// only add non-zero balances
				if total <= 0 {
					continue
				}

				found := false
				for _, a := range account.Balances {
					// if the balance already exists
					if a.CurrencySymbol == b.CurrencySymbol {
						// if total, locked, or available has changed on the exchange we need to update our balance for this currency
						if a.ExchangeTotal != total || a.ExchangeLocked != b.Locked || a.ExchangeAvailable != b.Free {
							if err := repoAccount.UpdateExchangeBalanceTxn(txn, ctx, account.AccountID, account.UserID, a.CurrencySymbol, total, b.Free, b.Locked); err != nil {
								txn.Rollback()
								res.Status = constRes.Error
								res.Message = "error encountered while updating balance: " + err.Error()
								log.Println("UpdateExchangeBalanceTxn error: ", err.Error())
								return nil
							}
						}
						found = true
						break
					}
				}
				if !found {
					balance := protoBalance.Balance{
						UserID:            account.UserID,
						AccountID:         account.AccountID,
						CurrencySymbol:    b.CurrencySymbol,
						Available:         b.Free,
						Locked:            0.0,
						ExchangeTotal:     total,
						ExchangeAvailable: b.Free,
						ExchangeLocked:    b.Locked,
						CreatedOn:         now,
						UpdatedOn:         now,
					}
					newBalances = append(newBalances, &balance)
				}
			}
			if len(newBalances) > 0 {
				if err := repoAccount.InsertBalances(txn, newBalances); err != nil {
					txn.Rollback()
					return errors.New("insert new balances failed within an account update: " + err.Error())
				}
			}
		}
	}

	if err := repoAccount.UpdateAccountTxn(txn, ctx, req.AccountID, req.UserID, req.KeyPublic, req.KeySecret, req.Description); err != nil {
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

func (service *AccountService) ValidateAccountBalance(ctx context.Context, req *protoBalance.ValidateBalanceRequest, res *protoBalance.ValidateBalanceResponse) error {
	balance, err := repoAccount.FindAccountBalance(service.DB, req.UserID, req.AccountID, req.CurrencySymbol)

	switch {
	case err == sql.ErrNoRows:
		// user does not have balance for said currency therefore not valid
		res.Status = constRes.Success
		res.Data = false
	case err != nil:
		log.Println("ValidateAccountBalance on FindAccountBalance error: ", err.Error())
		res.Status = constRes.Error
		res.Message = err.Error()
	default:
		valid := balance.Available > req.RequestedAmount
		res.Status = constRes.Success
		res.Data = valid
	}
	return nil
}
