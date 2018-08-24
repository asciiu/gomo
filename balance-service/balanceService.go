package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	repoBalance "github.com/asciiu/gomo/balance-service/db/sql"
	protoBalance "github.com/asciiu/gomo/balance-service/proto/balance"
)

type BalanceService struct {
	DB *sql.DB
}

func (service *BalanceService) HandleBalances(ctx context.Context, balances *protoBalance.AccountBalances) error {

	count, err := repoBalance.UpsertBalances(service.DB, balances)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("inserted ", count)
	return nil
}

// GetUserBalance returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *BalanceService) GetUserBalance(ctx context.Context, req *protoBalance.GetUserBalanceRequest, res *protoBalance.BalanceResponse) error {
	balance, error := repoBalance.FindBalance(service.DB, req)

	if error == nil {
		res.Status = "success"
		res.Data = &protoBalance.UserBalanceData{
			Balance: balance,
		}
	} else {
		res.Status = "error"
		res.Message = error.Error()
	}

	return nil
}

// GetUserBalances returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *BalanceService) GetUserBalances(ctx context.Context, req *protoBalance.GetUserBalancesRequest, res *protoBalance.BalancesResponse) error {

	switch {
	case req.Symbol == "":
		balances, error := repoBalance.FindAllBalancesByUserID(service.DB, req)
		if error == nil {
			res.Status = "success"
			res.Data = balances
		} else {
			res.Status = "error"
			res.Message = error.Error()
		}
	default:
		balances, error := repoBalance.FindSymbolBalancesByUserID(service.DB, req)
		if error == nil {
			res.Status = "success"
			res.Data = balances
		} else {
			res.Status = "error"
			res.Message = error.Error()
		}
	}

	return nil
}
