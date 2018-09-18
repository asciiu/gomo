package main

import (
	"context"
	"log"
	"math"
	"os"
	"strings"
	"time"

	binance "github.com/asciiu/go-binance"
	protoBalance "github.com/asciiu/gomo/binance-service/proto/balance"
	constExt "github.com/asciiu/gomo/common/constants/exchange"
	constRes "github.com/asciiu/gomo/common/constants/response"
	protoEvt "github.com/asciiu/gomo/common/proto/events"
	"github.com/asciiu/gomo/common/util"
	constPlan "github.com/asciiu/gomo/plan-service/constants"
	kitLog "github.com/go-kit/kit/log"
	micro "github.com/micro/go-micro"
)

type BinanceService struct {
	CompletedPub micro.Publisher
	Info         *BinanceExchangeInfo
}

func (service *BinanceService) HandleFillOrder(ctx context.Context, triggerEvent *protoEvt.TriggeredOrderEvent) error {
	// ignore events not binance
	// perhaps we can have this handler only receive binance triggers but for the sake of
	// simplicity when adding new exchanges let's just have each exchange service do a check
	// on the exchange
	if triggerEvent.Exchange != constExt.Binance {
		return nil
	}

	go func() {
		var logger kitLog.Logger
		logger = kitLog.NewLogfmtLogger(os.Stdout)
		logger = kitLog.With(logger, "time", kitLog.DefaultTimestampUTC, "caller", kitLog.DefaultCaller)

		hmacSigner := &binance.HmacSigner{
			Key: []byte(util.Rot32768(triggerEvent.KeySecret)),
		}

		binanceService := binance.NewAPIService(
			"https://www.binance.com",
			triggerEvent.KeyPublic,
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
		symbols := strings.Split(triggerEvent.MarketName, "-")
		marketName := strings.Replace(triggerEvent.MarketName, "-", "", 1)
		// buy or sell
		ellado := binance.SideBuy
		finalCurrency := symbols[0]
		if triggerEvent.Side == constPlan.Sell {
			finalCurrency = symbols[1]
			ellado = binance.SideSell
		}
		// order type can be market or limit
		orderType := binance.TypeMarket
		if triggerEvent.OrderType == constPlan.LimitOrder {
			orderType = binance.TypeLimit
		}

		completedEvent := protoEvt.CompletedOrderEvent{
			UserID:                 triggerEvent.UserID,
			PlanID:                 triggerEvent.PlanID,
			OrderID:                triggerEvent.OrderID,
			Exchange:               constExt.Binance,
			MarketName:             triggerEvent.MarketName,
			Side:                   triggerEvent.Side,
			AccountID:              triggerEvent.AccountID,
			InitialCurrencyBalance: triggerEvent.ActiveCurrencyBalance,
			InitialCurrencySymbol:  triggerEvent.ActiveCurrencySymbol,
			FinalCurrencySymbol:    finalCurrency,
			TriggerID:              triggerEvent.TriggerID,
			TriggeredPrice:         triggerEvent.TriggeredPrice,
			TriggeredCondition:     triggerEvent.TriggeredCondition,
		}

		// qauntity must be of step size
		lotSize := service.Info.LotSize(triggerEvent.MarketName)
		//minNotional := service.Info.MinNotional(triggerEvent.MarketName)

		var qauntity float64
		switch {
		case triggerEvent.Side == constPlan.Buy && triggerEvent.OrderType == constPlan.LimitOrder:
			// buy limit order should use limit price to compute final currency qty
			qauntity = triggerEvent.ActiveCurrencyBalance / triggerEvent.LimitPrice
			qauntity = math.Floor(qauntity/lotSize.StepSize) * lotSize.StepSize

			completedEvent.InitialCurrencyTraded = qauntity * triggerEvent.LimitPrice
			completedEvent.InitialCurrencyRemainder = completedEvent.InitialCurrencyBalance - completedEvent.InitialCurrencyTraded
			//completedEvent.FinalCurrencyBalance = qauntity

		case triggerEvent.Side == constPlan.Buy && triggerEvent.OrderType == constPlan.MarketOrder:
			// buy market should use the triggered price in the event
			qauntity = triggerEvent.ActiveCurrencyBalance / triggerEvent.TriggeredPrice
			qauntity = math.Floor(qauntity/lotSize.StepSize) * lotSize.StepSize

			completedEvent.InitialCurrencyTraded = qauntity * triggerEvent.TriggeredPrice
			completedEvent.InitialCurrencyRemainder = completedEvent.InitialCurrencyBalance - completedEvent.InitialCurrencyTraded
			//completedEvent.FinalCurrencyBalance = finalCurrencyQty

		default:
			// assume sell entire active balance
			qauntity = triggerEvent.ActiveCurrencyBalance
			qauntity = math.Floor(qauntity/lotSize.StepSize) * lotSize.StepSize

			completedEvent.InitialCurrencyTraded = qauntity
			completedEvent.InitialCurrencyRemainder = completedEvent.InitialCurrencyBalance - completedEvent.InitialCurrencyTraded
			// assume limit order
			//completedEvent.FinalCurrencyBalance = finalCurrencyQty * triggerEvent.LimitPrice

			//if triggerEvent.OrderType == constPlan.MarketOrder {
			//	// when not limit order compute the final balance as triggered price * qty
			//	completedEvent.FinalCurrencyBalance = finalCurrencyQty * triggerEvent.TriggeredPrice
			//}
		}
		//fmt.Println(minNotional)
		//qauntity = util.ToFixed(quantity, minNotional.MinNotional)

		//fmt.Printf("%+v\n", completedEvent)

		// https://github.com/binance-exchange/binance-official-api-docs/blob/master/rest-api.md
		// Limit type orders require a price
		processedOrder, err := b.NewOrder(binance.NewOrderRequest{
			Symbol:           marketName,
			Quantity:         qauntity,
			Side:             ellado,
			Price:            triggerEvent.LimitPrice,
			NewClientOrderID: triggerEvent.OrderID,
			TimeInForce:      binance.GTC,
			Type:             orderType,
			Timestamp:        time.Now(),
		})

		if err != nil {
			log.Printf("failed binance order call -- orderID: %s\n", triggerEvent.OrderID)
			completedEvent.Status = constPlan.Failed
			completedEvent.Details = err.Error()
		} else {

			log.Printf("processed order -- %+v\n", processedOrder)

			executedOrder, err := b.QueryOrder(binance.QueryOrderRequest{
				Symbol:            marketName,
				OrderID:           processedOrder.OrderID,
				OrigClientOrderID: triggerEvent.OrderID,
				RecvWindow:        time.Duration(2) * time.Second,
				Timestamp:         time.Now(),
			})

			log.Printf("order results -- %+v\n", executedOrder)

			if err != nil {
				completedEvent.Status = constPlan.Filled
				completedEvent.ExchangeOrderID = processedOrder.OrderID
				completedEvent.FinalCurrencyBalance = executedOrder.ExecutedQty
			}
		}

		if err := service.CompletedPub.Publish(ctx, &completedEvent); err != nil {
			log.Println("publish err: ", err.Error())
		}
	}()
	return nil
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

	account, err := b.Account(request)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "nvalid"):
			res.Status = constRes.Fail
			res.Message = "invalid keys"
			return nil
		default:
			res.Status = constRes.Error
			res.Message = err.Error()
			return nil
		}
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
