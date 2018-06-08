package main

import (
	"context"
	"fmt"
	"os"
	"time"

	binance "github.com/asciiu/go-binance"
	bp "github.com/asciiu/gomo/balance-service/proto/balance"
	"github.com/asciiu/gomo/common/constants/exchange"
	keyconstants "github.com/asciiu/gomo/common/constants/key"
	kp "github.com/asciiu/gomo/key-service/proto/key"
	"github.com/go-kit/kit/log"
	micro "github.com/micro/go-micro"
)

type KeyValidator struct {
	KeyVerifiedPub micro.Publisher
	BalancePub     micro.Publisher
}

func (service *KeyValidator) Process(ctx context.Context, key *kp.Key) error {
	if key.Exchange != exchange.Binance {
		return nil
	}

	go func() {
		var logger log.Logger
		logger = log.NewLogfmtLogger(os.Stdout)
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

		// TODO this should be enum
		key.Status = keyconstants.Verified

		// publish verify key event
		if err := service.KeyVerifiedPub.Publish(context.Background(), key); err != nil {
			logger.Log("could not publish verified key event: ", err)
		}

		balances := make([]*bp.Balance, 0)
		for _, balance := range account.Balances {
			total := balance.Free + balance.Locked

			bal := bp.Balance{
				KeyID:             key.KeyID,
				UserID:            key.UserID,
				ExchangeName:      key.Exchange,
				CurrencyName:      balance.Asset,
				Available:         balance.Free,
				Locked:            balance.Locked,
				ExchangeTotal:     total,
				ExchangeAvailable: balance.Free,
				ExchangeLocked:    balance.Locked,
			}
			balances = append(balances, &bal)
		}

		accountBalances := bp.AccountBalances{
			Balances: balances,
		}

		// publish account balances
		if err := service.BalancePub.Publish(context.Background(), &accountBalances); err != nil {
			logger.Log("could not publish account balances event: ", err)
		}

		logger.Log("verified keyID: ", key.KeyID)
	}()
	return nil
}
