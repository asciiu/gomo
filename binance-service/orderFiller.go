package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	binance "github.com/asciiu/go-binance"
	"github.com/asciiu/gomo/common/constants/exchange"
	"github.com/asciiu/gomo/common/constants/side"
	evt "github.com/asciiu/gomo/common/proto/events"
	notifications "github.com/asciiu/gomo/notification-service/proto"
	gokitlog "github.com/go-kit/kit/log"
	micro "github.com/micro/go-micro"
)

type OrderFiller struct {
	FilledPub micro.Publisher
	FailedPub micro.Publisher
}

func (filler *OrderFiller) FillOrder(ctx context.Context, orderEvent *evt.OrderEvent) error {
	// ignore events not binance
	if orderEvent.Exchange != exchange.Binance {
		return nil
	}

	go func() {
		var logger gokitlog.Logger
		logger = gokitlog.NewLogfmtLogger(os.Stdout)
		logger = gokitlog.With(logger, "time", gokitlog.DefaultTimestampUTC, "caller", gokitlog.DefaultCaller)

		hmacSigner := &binance.HmacSigner{
			Key: []byte(orderEvent.Secret),
		}

		binanceService := binance.NewAPIService(
			"https://www.binance.com",
			orderEvent.Key,
			hmacSigner,
			logger,
			ctx,
		)
		b := binance.NewBinance(binanceService)

		symbol := strings.Replace(orderEvent.MarketName, "-", "", 1)
		sidee := binance.SideBuy
		if orderEvent.Side == side.Sell {
			sidee = binance.SideSell
		}
		//orderEvent.BaseQuantity
		newOrder, err := b.NewOrder(binance.NewOrderRequest{
			Symbol:      symbol,
			Quantity:    100000,
			Price:       0.000010,
			Side:        sidee,
			TimeInForce: binance.GTC,
			Type:        binance.TypeLimit,
			Timestamp:   time.Now(),
		})
		if err != nil {
			log.Printf("failed new order binance call -- orderID: %s, market: %s\n", orderEvent.OrderID, orderEvent.MarketName)
			title := fmt.Sprintf("%s %s ordered failed", orderEvent.MarketName, orderEvent.Side)
			notification := notifications.Notification{
				UserID:           orderEvent.UserID,
				NotificationType: "order",
				ObjectID:         orderEvent.OrderID,
				Title:            title,
				Description:      err.Error(),
				Timestamp:        time.Now().UTC().Format(time.RFC3339),
			}

			// publish verify key event
			if err := filler.FailedPub.Publish(context.Background(), &notification); err != nil {
				log.Println("could not publish failed order: ", err)
			}
		} else {

			fmt.Println("Success: ", newOrder)

			orderEvent.ExchangeMarketName = "binance"
			orderEvent.ExchangeOrderID = "binance"

			if err := filler.FilledPub.Publish(ctx, orderEvent); err != nil {
				log.Println("publish warning: ", err, orderEvent)
			}

			log.Printf("order filled -- orderID: %s, market: %s\n", orderEvent.OrderID, orderEvent.MarketName)
		}
	}()
	return nil
}
