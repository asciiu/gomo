package main

import (
	"context"
	"fmt"
	"os"
	"time"

	binance "github.com/asciiu/go-binance"
	kp "github.com/asciiu/gomo/apikey-service/proto/apikey"
	"github.com/go-kit/kit/log"
)

type KeyValidator struct{}

func (sub *KeyValidator) Process(ctx context.Context, key *kp.ApiKey) error {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	hmacSigner := &binance.HmacSigner{
		Key: []byte(key.Secret),
	}
	//ctx, _ := context.WithCancel(context.Background())
	// use second return value for cancelling request when shutting down the app

	binanceService := binance.NewAPIService(
		"https://www.binance.com",
		key.Key,
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
		fmt.Printf("error encountered: %s", error)
	}

	for i, balance := range account.Balances {
		fmt.Printf("%d %s :: FREE: %f LOCKED %f\n", i, balance.Asset, balance.Free, balance.Locked)
	}
	return nil
}
