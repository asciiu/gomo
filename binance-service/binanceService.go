package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	binance "github.com/asciiu/go-binance"
	protoBalance "github.com/asciiu/gomo/binance-service/proto/balance"
	protoBinance "github.com/asciiu/gomo/binance-service/proto/binance"
	constExt "github.com/asciiu/gomo/common/constants/exchange"
	constRes "github.com/asciiu/gomo/common/constants/response"
	protoEvt "github.com/asciiu/gomo/common/proto/events"
	"github.com/asciiu/gomo/common/util"
	constPlan "github.com/asciiu/gomo/plan-service/constants"
	kitLog "github.com/go-kit/kit/log"
	"github.com/lib/pq"
	micro "github.com/micro/go-micro"
	"github.com/pkg/errors"
)

type BinanceService struct {
	CompletedPub micro.Publisher
	Info         *BinanceExchangeInfo
}

func (service *BinanceService) FormatTriggerEvent(triggerEvent *protoEvt.TriggeredOrderEvent) (*protoEvt.CompletedOrderEvent, *binance.NewOrderRequest) {

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
		InitialCurrencyBalance: triggerEvent.CommittedCurrencyAmount,
		InitialCurrencySymbol:  triggerEvent.CommittedCurrencySymbol,
		FinalCurrencySymbol:    finalCurrency,
		TriggerID:              triggerEvent.TriggerID,
		TriggeredPrice:         triggerEvent.TriggeredPrice,
		TriggeredCondition:     triggerEvent.TriggeredCondition,
		TriggeredTime:          triggerEvent.TriggeredTime,
		CloseOnComplete:        triggerEvent.CloseOnComplete,
	}

	// qauntity must be of step size
	lotSize := service.Info.LotSize(triggerEvent.MarketName)
	//minNotional := service.Info.MinNotional(triggerEvent.MarketName)

	var qauntity float64
	switch {
	case triggerEvent.Side == constPlan.Buy && triggerEvent.OrderType == constPlan.LimitOrder:
		// buy limit order should use limit price to compute final currency qty
		qauntity = triggerEvent.CommittedCurrencyAmount / triggerEvent.LimitPrice

		steps := math.Floor(qauntity / lotSize.StepSize)
		qauntity = lotSize.StepSize * steps
		qauntity = util.ToFixedRounded(qauntity, 8)

		completedEvent.InitialCurrencyTraded = util.ToFixedRounded(qauntity*triggerEvent.LimitPrice, 8)
		completedEvent.InitialCurrencyRemainder = util.ToFixedRounded(completedEvent.InitialCurrencyBalance-completedEvent.InitialCurrencyTraded, 8)
		//completedEvent.FinalCurrencyBalance = qauntity

	case triggerEvent.Side == constPlan.Buy && triggerEvent.OrderType == constPlan.MarketOrder:
		// buy market should use the triggered price in the event
		qauntity = triggerEvent.CommittedCurrencyAmount / triggerEvent.TriggeredPrice

		steps := math.Floor(qauntity / lotSize.StepSize)
		qauntity = lotSize.StepSize * steps
		qauntity = util.ToFixedRounded(qauntity, 8)

		completedEvent.InitialCurrencyTraded = util.ToFixedRounded(qauntity*triggerEvent.TriggeredPrice, 8)
		completedEvent.InitialCurrencyRemainder = util.ToFixedRounded(completedEvent.InitialCurrencyBalance-completedEvent.InitialCurrencyTraded, 8)
		//completedEvent.FinalCurrencyBalance = finalCurrencyQty

	default:
		// assume sell entire active balance
		qauntity = triggerEvent.CommittedCurrencyAmount

		steps := math.Floor(qauntity / lotSize.StepSize)
		qauntity = lotSize.StepSize * steps
		qauntity = util.ToFixedRounded(qauntity, 8)

		completedEvent.InitialCurrencyTraded = qauntity
		completedEvent.InitialCurrencyRemainder = util.ToFixedRounded(completedEvent.InitialCurrencyBalance-completedEvent.InitialCurrencyTraded, 8)

		// assume limit order
		//completedEvent.FinalCurrencyBalance = finalCurrencyQty * triggerEvent.LimitPrice

		//if triggerEvent.OrderType == constPlan.MarketOrder {
		//	// when not limit order compute the final balance as triggered price * qty
		//	completedEvent.FinalCurrencyBalance = finalCurrencyQty * triggerEvent.TriggeredPrice
		//}
	}

	binanceRequest := binance.NewOrderRequest{
		Symbol:           marketName,
		Quantity:         qauntity,
		Side:             ellado,
		Price:            triggerEvent.LimitPrice,
		NewClientOrderID: triggerEvent.OrderID,
		TimeInForce:      binance.GTC,
		Type:             orderType,
		Timestamp:        time.Now(),
	}

	return &completedEvent, &binanceRequest
}

func (service *BinanceService) HandleFillOrder(ctx context.Context, triggerEvent *protoEvt.TriggeredOrderEvent) error {
	log.Printf("trigger event -- %+v\n", triggerEvent)

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
		completedEvent, binanceRequest := service.FormatTriggerEvent(triggerEvent)

		// https://github.com/binance-exchange/binance-official-api-docs/blob/master/rest-api.md
		// Limit type orders require a price
		processedOrder, err := b.NewOrder(*binanceRequest)

		if err != nil {
			log.Printf("failed binance order call -- orderID: %s\n", triggerEvent.OrderID)
			completedEvent.Status = constPlan.Failed
			completedEvent.Details = err.Error()
		} else {
			log.Printf("processed order -- %+v\n", processedOrder)

			// need to wait for binance to fill the order on their end to avoid race condition
			time.Sleep(1 * time.Second)

			var quantity float64
			var commission float64
			var feesym string
			var latest time.Time
			var totalTrades float64
			var sumPrice float64

			// try pulling the results 3 times
			for i := 0; i < 3; i++ {

				// ask for most recent 200
				trades, err := b.MyTrades(binance.MyTradesRequest{
					Symbol:     binanceRequest.Symbol,
					Limit:      200,
					RecvWindow: time.Duration(2) * time.Second,
					Timestamp:  time.Now(),
				})

				if err != nil {
					log.Printf("could not retrieve recent trades after processed order %+v\n", processedOrder)
					completedEvent.Status = constPlan.Failed
					completedEvent.Details = err.Error()
					break
				}

				// loop trades and find all that match order ID
				for _, trade := range trades {
					if trade.OrderID == processedOrder.OrderID {
						totalTrades++

						if triggerEvent.Side == constPlan.Sell {
							quantity += trade.Price * trade.Qty
						} else {
							quantity += trade.Qty
						}

						// assuming same fee currency
						feesym = trade.CommissionAsset

						if trade.Time.After(latest) {
							latest = trade.Time
						}
						sumPrice += trade.Price
						commission += trade.Commission
					}
				}

				// assume the results were pulled completely if totalTrades > 0
				if totalTrades > 0 && commission > 0 && sumPrice > 0 {
					log.Printf("trade results -- %+v\n", trades)
					break
				}

				// wait for 10 seconds to pull the trade results
				time.Sleep(10 * time.Second)
			}

			switch {
			case totalTrades > 0 && commission > 0 && sumPrice > 0:
				exchangeRate := 0.0
				if totalTrades > 0 {
					exchangeRate = util.ToFixedRounded(sumPrice/totalTrades, 8)
				}
				completedEvent.ExchangeTime = string(pq.FormatTimestamp(latest))
				completedEvent.FeeCurrencySymbol = feesym
				completedEvent.FeeCurrencyAmount = util.ToFixedRounded(commission, 8)
				// compute the exchange rate as an average
				completedEvent.ExchangePrice = exchangeRate

				completedEvent.Status = constPlan.Filled
				completedEvent.ExchangeOrderID = strconv.FormatInt(processedOrder.OrderID, 10)
				completedEvent.FinalCurrencyBalance = util.ToFixedRounded(quantity, 8)

				if triggerEvent.Side == constPlan.Sell {
					completedEvent.Details = fmt.Sprintf("sold %.8f %s in exchange for %.8f %s",
						completedEvent.InitialCurrencyBalance,
						completedEvent.InitialCurrencySymbol,
						completedEvent.FinalCurrencyBalance,
						completedEvent.FinalCurrencySymbol)
				} else {
					completedEvent.Details = fmt.Sprintf("bought %.8f %s and traded %.8f %s",
						completedEvent.FinalCurrencyBalance,
						completedEvent.FinalCurrencySymbol,
						completedEvent.InitialCurrencyTraded,
						completedEvent.InitialCurrencySymbol)
				}
			default:
				log.Printf("3rd attempt to retrieve binance trade results failed for %+v\n", processedOrder)
				completedEvent.Status = constPlan.Failed
				completedEvent.Details = err.Error()
				break
			}
			// I'd prefer to just call this but binance doesn't return the price in
			// the response, so I had to resort to pulling the most recent trades
			//executedOrder, err := b.QueryOrder(binance.QueryOrderRequest{
			//	Symbol:            marketName,
			//	OrderID:           processedOrder.OrderID,
			//	OrigClientOrderID: triggerEvent.OrderID,
			//	RecvWindow:        time.Duration(2) * time.Second,
			//	Timestamp:         time.Now(),
			//})
		}

		if err := service.CompletedPub.Publish(ctx, completedEvent); err != nil {
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

func (service *BinanceService) GetMarketRestrictions(ctx context.Context, req *protoBinance.MarketRestrictionRequest, res *protoBinance.MarketRestrictionResponse) error {
	priceFilter := service.Info.PriceFilter(req.MarketName)
	lotFilter := service.Info.LotSize(req.MarketName)

	if priceFilter == nil || lotFilter == nil {
		res.Status = constRes.Nonentity
		return nil
	}

	res.Status = constRes.Success
	res.Data = &protoBinance.RestrictionData{
		Restrictions: &protoBinance.MarketRestriction{
			MinTradeSize:    lotFilter.MinQty,
			MaxTradeSize:    lotFilter.MaxQty,
			TradeSizeStep:   lotFilter.StepSize,
			MinMarketPrice:  priceFilter.Min,
			MaxMarketPrice:  priceFilter.Max,
			MarketPriceStep: priceFilter.TickSize,
			BasePrecision:   8,
			MarketPrecision: 8,
		},
	}

	return nil
}

func (service *BinanceService) GetCandles(ctx context.Context, req *protoBinance.CandleRequest, res *protoBinance.CandlesResponse) error {

	symbol := strings.Replace(req.MarketName, "-", "", 1)

	url := fmt.Sprintf("https://api.binance.com/api/v1/klines?symbol=%s&interval=%s&limit=1000", symbol, "1m")

	request, _ := http.NewRequest("GET", url, nil)

	response, _ := http.DefaultClient.Do(request)

	defer response.Body.Close()

	textRes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read response from Klines")
	}

	//if res.StatusCode != 200 {
	//	as.handleError(textRes)
	//}

	rawKlines := [][]interface{}{}
	if err := json.Unmarshal(textRes, &rawKlines); err != nil {
		return errors.Wrap(err, "rawKlines unmarshal failed")
	}
	klines := []*protoBinance.Candle{}
	for _, k := range rawKlines {
		ot, err := timeFromUnixTimestampFloat(k[0])
		if err != nil {
			return errors.Wrap(err, "cannot parse Kline.OpenTime")
		}
		open, err := floatFromString(k[1])
		if err != nil {
			return errors.Wrap(err, "cannot parse Kline.Open")
		}
		high, err := floatFromString(k[2])
		if err != nil {
			return errors.Wrap(err, "cannot parse Kline.High")
		}
		low, err := floatFromString(k[3])
		if err != nil {
			return errors.Wrap(err, "cannot parse Kline.Low")
		}
		cls, err := floatFromString(k[4])
		if err != nil {
			return errors.Wrap(err, "cannot parse Kline.Close")
		}
		volume, err := floatFromString(k[5])
		if err != nil {
			return errors.Wrap(err, "cannot parse Kline.Volume")
		}
		ct, err := timeFromUnixTimestampFloat(k[6])
		if err != nil {
			return errors.Wrap(err, "cannot parse Kline.CloseTime")
		}
		qav, err := floatFromString(k[7])
		if err != nil {
			return errors.Wrap(err, "cannot parse Kline.QuoteAssetVolume")
		}
		not, ok := k[8].(float64)
		if !ok {
			return errors.Wrap(err, "cannot parse Kline.NumberOfTrades")
		}
		tbbav, err := floatFromString(k[9])
		if err != nil {
			return errors.Wrap(err, "cannot parse Kline.TakerBuyBaseAssetVolume")
		}
		tbqav, err := floatFromString(k[10])
		if err != nil {
			return errors.Wrap(err, "cannot parse Kline.TakerBuyQuoteAssetVolume")
		}
		klines = append(klines, &protoBinance.Candle{
			OpenTime:                 ot.String(),
			Open:                     open,
			High:                     high,
			Low:                      low,
			Close:                    cls,
			Volume:                   volume,
			CloseTime:                ct.String(),
			QuoteAssetVolume:         qav,
			NumberOfTrades:           int32(not),
			TakerBuyBaseAssetVolume:  tbbav,
			TakerBuyQuoteAssetVolume: tbqav,
		})
	}

	res.Status = constRes.Success
	res.Data = &protoBinance.Candles{
		Candles: klines,
	}
	return nil
}
