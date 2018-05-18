package main

import (
	"context"
	"database/sql"

	repo "github.com/asciiu/gomo/balance-service/db/sql"
	bp "github.com/asciiu/gomo/balance-service/proto/balance"
)

type BalanceService struct {
	DB *sql.DB
}

// GetUserBalance returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *BalanceService) GetUserBalance(ctx context.Context, req *bp.GetUserBalanceRequest, res *bp.BalanceResponse) error {
	balance, error := repo.FindBalance(service.DB, req)

	if error == nil {
		res.Status = "success"
		res.Data = &bp.UserBalanceData{
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
func (service *BalanceService) GetUserBalances(ctx context.Context, req *bp.GetUserBalancesRequest, res *bp.BalancesResponse) error {

	switch {
	case req.Symbol == "":
		balances, error := repo.FindAllBalancesByUserID(service.DB, req)
		if error == nil {
			res.Status = "success"
			res.Data = balances
		} else {
			res.Status = "error"
			res.Message = error.Error()
		}
	default:
		balances, error := repo.FindSymbolBalancesByUserID(service.DB, req)
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
