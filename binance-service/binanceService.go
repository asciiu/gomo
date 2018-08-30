package main

import (
	"context"
	"os"
	"time"

	binance "github.com/asciiu/go-binance"
	protoBalance "github.com/asciiu/gomo/binance-service/proto/balance"
	constRes "github.com/asciiu/gomo/common/constants/response"
	kitLog "github.com/go-kit/kit/log"
)

type BinanceService struct {
}

// Retrieve the exchange balances
func (service *BinanceService) GetBalances(ctx context.Context, req *protoBalance.BalanceRequest, res *protoBalance.BalancesResponse) error {

	var logger kitLog.Logger
	logger = kitLog.NewLogfmtLogger(os.Stdout)
	logger = kitLog.With(logger, "time", kitLog.DefaultTimestampUTC, "caller", kitLog.DefaultCaller)

	hmacSigner := &binance.HmacSigner{
		Key: []byte(req.KeySecret),
	}

	binanceService := binance.NewAPIService(
		"https://www.binance.com",
		req.KeyPublic,
		hmacSigner,
		logger,
		ctx,
	)
	b := binance.NewBinance(binanceService)
	request := binance.AccountRequest{
		RecvWindow: time.Duration(2) * time.Second,
		Timestamp:  time.Now(),
	}

	account, error := b.Account(request)
	if error != nil {
		res.Status = constRes.Error
		res.Message = error.Error()
		return nil
	}

	balances := make([]*protoBalance.Balance, 0)
	for _, balance := range account.Balances {
		bal := protoBalance.Balance{
			CurrencySymbol: balance.Asset,
			Free:           balance.Free,
			Locked:         balance.Locked,
		}
		balances = append(balances, &bal)
	}

	res.Status = constRes.Success
	res.Data = &protoBalance.BalanceList{
		Balances: balances,
	}

	return nil

}
