package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	binance "github.com/asciiu/go-binance"
	kp "github.com/asciiu/gomo/apikey-service/proto/apikey"
	"github.com/go-kit/kit/log"
)

type KeyValidator struct {
	DB *sql.DB
}

func (sub *KeyValidator) Process(ctx context.Context, key *kp.ApiKey) error {
	if key.Exchange != "Binance" {
		return nil
	}

	go func() {
		var logger log.Logger
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		logger = log.With(logger, "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

		hmacSigner := &binance.HmacSigner{
			Key: []byte(key.Secret),
		}

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

		// publish verify key event

		// publish balances here an an event
		for i, balance := range account.Balances {
			fmt.Printf("%d %s :: FREE: %f LOCKED %f\n", i, balance.Asset, balance.Free, balance.Locked)
		}
	}()
	return nil
}
