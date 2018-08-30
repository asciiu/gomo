package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	binance "github.com/asciiu/go-binance"
	constExt "github.com/asciiu/gomo/common/constants/exchange"
	protoEvt "github.com/asciiu/gomo/common/proto/events"
	constPlan "github.com/asciiu/gomo/plan-service/constants"
	gokitlog "github.com/go-kit/kit/log"
	micro "github.com/micro/go-micro"
)

type OrderFulfiller struct {
	CompletedPub micro.Publisher
}

func (filler *OrderFulfiller) FillOrder(ctx context.Context, triggerEvent *protoEvt.TriggeredOrderEvent) error {
	// ignore events not binance
	// perhaps we can have this handler only receive binance triggers but for the sake of
	// simplicity when adding new exchanges let's just have each exchange service do a check
	// on the exchange
	if triggerEvent.Exchange != constExt.Binance {
		return nil
	}

	go func() {
		var logger gokitlog.Logger
		logger = gokitlog.NewLogfmtLogger(os.Stdout)
		logger = gokitlog.With(logger, "time", gokitlog.DefaultTimestampUTC, "caller", gokitlog.DefaultCaller)

		hmacSigner := &binance.HmacSigner{
			Key: []byte(triggerEvent.Secret),
		}

		binanceService := binance.NewAPIService(
			"https://www.binance.com",
			triggerEvent.Key,
			hmacSigner,
			logger,
			ctx,
		)
		b := binance.NewBinance(binanceService)

		// Buy-limit: plan.baseBalance / planOrder.Price
		// Buy-market: plan.baseBalance / trigger.Price (can only determine this at trigger time)
		// Sell-limit: currencyBalance
		// Sell-market: currencyBalance

		// binance expects the symbol to be formatted as a single word: e.g. BNBBTC
		symbol := strings.Replace(triggerEvent.MarketName, "-", "", 1)
		// buy or sell
		ellado := binance.SideBuy
		if triggerEvent.Side == constPlan.Sell {
			ellado = binance.SideSell
		}
		// order type can be market or limit
		orderType := binance.TypeMarket
		if triggerEvent.OrderType == constPlan.LimitOrder {
			orderType = binance.TypeLimit
		}
		// https://github.com/binance-exchange/binance-official-api-docs/blob/master/rest-api.md
		// Limit type orders require a price
		newOrder, err := b.NewOrder(binance.NewOrderRequest{
			Symbol:      symbol,
			Quantity:    triggerEvent.Quantity,
			Side:        ellado,
			Price:       triggerEvent.Price,
			TimeInForce: binance.GTC,
			Type:        orderType,
			Timestamp:   time.Now(),
		})
		if err != nil {
			log.Printf("failed new order binance call -- orderID: %s, market: %s\n", triggerEvent.OrderID, triggerEvent.MarketName)
			//title := fmt.Sprintf("%s %s ordered failed", triggerEvent.MarketName, triggerEvent.Side)
			// notification := notifications.Notification{
			// 	UserID:           triggerEvent.UserID,
			// 	NotificationType: "order",
			// 	ObjectID:         triggerEvent.OrderID,
			// 	Title:            title,
			// 	Description:      err.Error(),
			// 	Timestamp:        time.Now().UTC().Format(time.RFC3339),
			// }

			// // publish verify key event
			// if err := filler.FailedPub.Publish(context.Background(), &notification); err != nil {
			// 	log.Println("could not publish failed order: ", err)
			// }

			completedEvent := protoEvt.CompletedOrderEvent{
				UserID:             triggerEvent.UserID,
				PlanID:             triggerEvent.PlanID,
				OrderID:            triggerEvent.OrderID,
				Side:               triggerEvent.Side,
				TriggeredPrice:     triggerEvent.TriggeredPrice,
				TriggeredCondition: triggerEvent.TriggeredCondition,
				Status:             constPlan.Failed,
				Details:            err.Error(),
			}

			if err := filler.CompletedPub.Publish(ctx, &completedEvent); err != nil {
				log.Println("publish warning: ", err, completedEvent)
			}

		} else {
			//if err := filler.FilledPub.Publish(ctx, orderEvent); err != nil {
			//	log.Println("publish warning: ", err, orderEvent)
			//}

			log.Printf("processed order -- orderID: %s, details: %+v\n", triggerEvent.OrderID, newOrder)
		}
	}()
	return nil
}
