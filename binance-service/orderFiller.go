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
	"github.com/asciiu/gomo/common/constants/order"
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

		// binance expects the symbol to be formatted as a single word: e.g. BNBBTC
		symbol := strings.Replace(orderEvent.MarketName, "-", "", 1)
		// buy or sell
		ellado := binance.SideBuy
		if orderEvent.Side == side.Sell {
			ellado = binance.SideSell
		}
		// order type can be market or limit
		orderType := binance.TypeMarket
		if orderEvent.OrderType == order.LimitOrder {
			orderType = binance.TypeLimit
		}
		//orderEvent.BaseQuantity
		// https://github.com/binance-exchange/binance-official-api-docs/blob/master/rest-api.md
		// Limit type orders require a price
		newOrder, err := b.NewOrder(binance.NewOrderRequest{
			Symbol:      symbol,
			Quantity:    orderEvent.Quantity,
			Side:        ellado,
			Price:       orderEvent.Price,
			TimeInForce: binance.GTC,
			Type:        orderType,
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
