package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	binance "github.com/asciiu/go-binance"
	kp "github.com/asciiu/gomo/apikey-service/proto/apikey"
	bp "github.com/asciiu/gomo/balance-service/proto/balance"
	msg "github.com/asciiu/gomo/common/messages"
	"github.com/go-kit/kit/log"
	micro "github.com/micro/go-micro"
)

type KeyValidator struct {
	DB    *sql.DB
	Micro micro.Service
}

func (service *KeyValidator) Process(ctx context.Context, key *kp.ApiKey) error {
	if key.Exchange != "Binance" {
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
		key.Status = "verified"

		// publish verify key event
		// publish a new key event
		keyPublisher := micro.NewPublisher(msg.TopicKeyVerified, service.Micro.Client())
		if err := keyPublisher.Publish(context.Background(), key); err != nil {
			logger.Log("could not publish verified key event: ", err)
		}

		balPublisher := micro.NewPublisher(msg.TopicBalanceUpdate, service.Micro.Client())
		balances := make([]*bp.Balance, 0)
		for _, balance := range account.Balances {
			bal := bp.Balance{
				Currency: balance.Asset,
				Free:     balance.Free,
				Locked:   balance.Locked,
			}
			balances = append(balances, &bal)
		}

		accountBalances := bp.AccountBalances{
			ApiKeyId: key.ApiKeyId,
			UserId:   key.UserId,
			Exchange: key.Exchange,
			Balances: balances,
		}
		if err := balPublisher.Publish(context.Background(), &accountBalances); err != nil {
			logger.Log("could not publish account balances event: ", err)
		}

		// publish balances here an an event
		for i, balance := range account.Balances {
			fmt.Printf("%d %s :: FREE: %f LOCKED %f\n", i, balance.Asset, balance.Free, balance.Locked)
		}

		logger.Log("verified keyId: ", key.ApiKeyId)
	}()
	return nil
}
