package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	constOrder "github.com/asciiu/gomo/common/constants/order"
	"github.com/asciiu/gomo/common/constants/side"
	"github.com/asciiu/gomo/common/constants/status"
	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/mattn/anko/vm"
	micro "github.com/micro/go-micro"
)

// Processor will process orders
type PlanProcessor struct {
	DB        *sql.DB
	Env       *vm.Env
	Receiver  *PlanReceiver
	Completed micro.Publisher
	Triggered micro.Publisher
}

// ProcessEvent will process ExchangeEvents. These events are published from the exchange sockets.
func (processor *PlanProcessor) ProcessTradeEvent(ctx context.Context, payload *evt.TradeEvents) error {
	plans := processor.Receiver.Plans

	// TODO rewrite this - needs to be more efficient!
	for _, tradeEvent := range payload.Events {
		// TODO only look at the plans relevant for this market-exchange
		for p, plan := range plans {
			for _, order := range plan.Orders {
				if tradeEvent.MarketName == order.MarketName && tradeEvent.Exchange == order.Exchange {
					for _, trigger := range order.TriggerExs {
						if isTrue, desc := trigger.Evaluate(tradeEvent.Price); isTrue {
							// remove this order from the processor
							processor.Receiver.Plans = append(plans[:p], plans[p+1:]...)

							if order.OrderType == constOrder.PaperOrder {
								completedEvent := evt.CompletedOrderEvent{
									UserID:     plan.UserID,
									PlanID:     plan.PlanID,
									OrderID:    order.OrderID,
									Exchange:   order.Exchange,
									MarketName: order.MarketName,
									Side:       order.Side,
									InitialCurrencyBalance: plan.ActiveCurrencyBalance,
									InitialCurrencySymbol:  plan.ActiveCurrencySymbol,
									TriggerID:              trigger.TriggerID,
									TriggeredPrice:         tradeEvent.Price,
									TriggeredCondition:     desc,
									ExchangeOrderID:        constOrder.PaperOrder,
									ExchangeMarketName:     constOrder.PaperOrder,
									Status:                 status.Filled,
									CloseOnComplete:        plan.CloseOnComplete,
								}

								symbols := strings.Split(order.MarketName, "-")
								// adjust balances for buy
								if order.Side == side.Buy {
									completedEvent.FinalCurrencySymbol = symbols[0]
									completedEvent.FinalCurrencyBalance = plan.ActiveCurrencyBalance / tradeEvent.Price
									completedEvent.InitialCurrencyRemainder = plan.ActiveCurrencyBalance - (completedEvent.FinalCurrencyBalance * tradeEvent.Price)
									completedEvent.InitialCurrencyTraded = plan.ActiveCurrencyBalance - completedEvent.InitialCurrencyRemainder
									completedEvent.Details = fmt.Sprintf("bought %.8f %s with %.8f %s", completedEvent.FinalCurrencyBalance, symbols[0], completedEvent.InitialCurrencyTraded, symbols[1])
								}

								// adjust balances for sell
								if order.Side == side.Sell {
									completedEvent.FinalCurrencySymbol = symbols[1]
									completedEvent.FinalCurrencyBalance = plan.ActiveCurrencyBalance * tradeEvent.Price
									completedEvent.InitialCurrencyRemainder = 0
									completedEvent.InitialCurrencyTraded = plan.ActiveCurrencyBalance
									completedEvent.Details = fmt.Sprintf("sold %.8f %s for %.8f %s", plan.ActiveCurrencyBalance, symbols[0], completedEvent.FinalCurrencyBalance, symbols[1])
								}

								// Never log the secrets contained in the event
								log.Printf("virtual order processed -- orderID: %s, market: %s\n", order.OrderID, order.MarketName)

								if err := processor.Completed.Publish(ctx, &completedEvent); err != nil {
									log.Println("publish warning: ", err, completedEvent)
								}
							} else {
								quantity := 0.0
								switch {
								case order.Side == side.Buy && order.OrderType == constOrder.LimitOrder:
									quantity = plan.ActiveCurrencyBalance / order.LimitPrice

								case order.Side == side.Buy && order.OrderType == constOrder.MarketOrder:
									quantity = plan.ActiveCurrencyBalance / tradeEvent.Price

								default:
									// sell entire active balance
									quantity = plan.ActiveCurrencyBalance
								}

								// convert this active order event to a triggered order event
								triggeredEvent := evt.TriggeredOrderEvent{
									Exchange:           order.Exchange,
									OrderID:            order.OrderID,
									PlanID:             plan.PlanID,
									UserID:             plan.UserID,
									KeyID:              order.KeyID,
									Key:                order.KeyPublic,
									Secret:             order.KeySecret,
									MarketName:         order.MarketName,
									Side:               order.Side,
									OrderType:          order.OrderType,
									Price:              order.LimitPrice,
									Quantity:           quantity,
									TriggeredPrice:     tradeEvent.Price,
									TriggeredCondition: desc,
								}

								// Never log the secrets contained in the event
								log.Printf("triggered order -- orderID: %s, market: %s\n", order.OrderID, order.MarketName)
								// if non simulated trigger buy event - exchange service subscribes to these events
								if err := processor.Triggered.Publish(ctx, &triggeredEvent); err != nil {
									log.Println("publish warning: ", err)
								}
							}
						}
					}
				}
			}
		}
	}

	//fmt.Println(payload)
	// orders := processor.Receiver.Orders

	// // every order check trade event price with order conditions
	// for i, order := range orders {
	// 	for _, event := range payload.Events {
	// 		marketName := order.EventOrigin.MarketName
	// 		//marketName := strings.Replace(order.EventOrigin.MarketName, "-", "", 1)
	// 		// market name and exchange must match
	// 		if marketName != event.MarketName || order.EventOrigin.Exchange != event.Exchange {
	// 			continue
	// 		}

	// 		executionTriggers := order.ExecutionTriggers
	// 		// eval all conditions for this order
	// 		for _, executionTrig := range executionTriggers {

	// 			// does condition of order eval to true?
	// 			if isValid, desc := executionTrig.TriggerFunction(event.Price); isValid {
	// 				// remove this order from the processor
	// 				processor.Receiver.Orders = append(orders[:i], orders[i+1:]...)

	// 				switch {
	// 				case order.EventOrigin.OrderType == types.PaperOrder:
	// 					details := fmt.Sprintf("%s %s -- %s", order.EventOrigin.Side, order.EventOrigin.MarketName, status.Filled)
	// 					completedEvent := evt.CompletedOrderEvent{
	// 						UserID:             order.EventOrigin.UserID,
	// 						PlanID:             order.EventOrigin.PlanID,
	// 						OrderID:            order.EventOrigin.OrderID,
	// 						Side:               order.EventOrigin.Side,
	// 						TriggeredPrice:     event.Price,
	// 						TriggeredCondition: desc,
	// 						ExchangeOrderID:    types.PaperOrder,
	// 						ExchangeMarketName: types.PaperOrder,
	// 						Status:             status.Filled,
	// 						Details:            details,
	// 					}

	// 					// adjust balances for buy
	// 					if order.EventOrigin.Side == side.Buy {
	// 						currencyQuantity := order.EventOrigin.BaseBalance * order.EventOrigin.BalancePercent / event.Price
	// 						spent := currencyQuantity * event.Price
	// 						completedEvent.BaseBalance = order.EventOrigin.BaseBalance - spent
	// 						completedEvent.CurrencyBalance = order.EventOrigin.CurrencyBalance + currencyQuantity
	// 					}

	// 					// adjust balances for sell
	// 					if order.EventOrigin.Side == side.Sell {
	// 						currencyQuantity := order.EventOrigin.CurrencyBalance * order.EventOrigin.BalancePercent
	// 						completedEvent.BaseBalance = currencyQuantity*event.Price + order.EventOrigin.BaseBalance
	// 						completedEvent.CurrencyBalance = order.EventOrigin.CurrencyBalance - currencyQuantity
	// 					}

	// 					// Never log the secrets contained in the event
	// 					log.Printf("virtual order processed -- orderID: %s, market: %s\n", order.EventOrigin.OrderID, order.EventOrigin.MarketName)

	// 					if err := processor.Completed.Publish(ctx, &completedEvent); err != nil {
	// 						log.Println("publish warning: ", err, completedEvent)
	// 					}
	// 				default:

	// 					quantity := 0.0
	// 					switch {
	// 					case order.EventOrigin.Side == side.Buy && order.EventOrigin.OrderType == types.LimitOrder:
	// 						quantity = order.EventOrigin.BaseBalance * order.EventOrigin.BalancePercent / order.EventOrigin.Price

	// 					case order.EventOrigin.Side == side.Buy && order.EventOrigin.OrderType == types.MarketOrder:
	// 						quantity = order.EventOrigin.BaseBalance * order.EventOrigin.BalancePercent / event.Price

	// 					default:
	// 						quantity = order.EventOrigin.CurrencyBalance * order.EventOrigin.BalancePercent
	// 					}

	// 					// convert this active order event to a triggered order event
	// 					triggeredEvent := evt.TriggeredOrderEvent{
	// 						Exchange:           order.EventOrigin.Exchange,
	// 						OrderID:            order.EventOrigin.OrderID,
	// 						PlanID:             order.EventOrigin.PlanID,
	// 						UserID:             order.EventOrigin.UserID,
	// 						BaseBalance:        order.EventOrigin.BaseBalance,
	// 						CurrencyBalance:    order.EventOrigin.CurrencyBalance,
	// 						KeyID:              order.EventOrigin.KeyID,
	// 						Key:                order.EventOrigin.Key,
	// 						Secret:             order.EventOrigin.Secret,
	// 						MarketName:         order.EventOrigin.MarketName,
	// 						Side:               order.EventOrigin.Side,
	// 						OrderType:          order.EventOrigin.OrderType,
	// 						Price:              order.EventOrigin.Price,
	// 						Quantity:           quantity,
	// 						TriggeredPrice:     event.Price,
	// 						TriggeredCondition: desc,
	// 						NextOrderID:        order.EventOrigin.NextOrderID,
	// 					}

	// 					// Never log the secrets contained in the event
	// 					log.Printf("triggered order -- orderID: %s, market: %s\n", order.EventOrigin.OrderID, order.EventOrigin.MarketName)
	// 					// if non simulated trigger buy event - exchange service subscribes to these events
	// 					if err := processor.Triggered.Publish(ctx, &triggeredEvent); err != nil {
	// 						log.Println("publish warning: ", err)
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	return nil
}
