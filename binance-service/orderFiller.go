package main

import (
	"context"
	"os"

	"github.com/asciiu/gomo/common/constants/exchange"
	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/go-kit/kit/log"
	micro "github.com/micro/go-micro"
)

type OrderFiller struct {
	FilledPub micro.Publisher
}

func (filler *OrderFiller) Process(ctx context.Context, orderEvent *evt.OrderEvent) error {
	// ignore events not binance
	if orderEvent.Exchange != exchange.Binance {
		return nil
	}

	go func() {
		var logger log.Logger
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

		// hmacSigner := &binance.HmacSigner{
		// 	Key: []byte(key.Secret),
		// }

		// binanceService := binance.NewAPIService(
		// 	"https://www.binance.com",
		// 	key.Key,
		// 	hmacSigner,
		// 	logger,
		// 	ctx,
		// )
		// b := binance.NewBinance(binanceService)
		// request := binance.AccountRequest{
		// 	RecvWindow: time.Duration(2) * time.Second,
		// 	Timestamp:  time.Now(),
		// }

		// account, error := b.Account(request)
		// if error != nil {
		// 	fmt.Printf("error encountered: %s", error)
		// }

		// // publish verify key event
		// if err := service.KeyVerifiedPub.Publish(context.Background(), key); err != nil {
		// 	logger.Log("could not publish verified key event: ", err)
		// }

		logger.Log("order filled: ", orderEvent.OrderID)
	}()
	return nil
}
