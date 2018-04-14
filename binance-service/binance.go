package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	binance "github.com/asciiu/go-binance"
	apikey "github.com/asciiu/gomo/apikey-service/models"
	"github.com/go-kit/kit/log"
)

func VerifyKey(details []byte) {
	var keyDetails apikey.ExchangeKey
	err := json.Unmarshal(details, &keyDetails)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(keyDetails)

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	hmacSigner := &binance.HmacSigner{
		Key: []byte(keyDetails.Secret),
	}
	ctx, _ := context.WithCancel(context.Background())
	// use second return value for cancelling request when shutting down the app

	binanceService := binance.NewAPIService(
		"https://www.binance.com",
		keyDetails.ApiKey,
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
}
