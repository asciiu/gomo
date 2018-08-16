package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/asciiu/gomo/common/constants/response"

	constOrder "github.com/asciiu/gomo/common/constants/order"
	"github.com/asciiu/gomo/common/constants/side"
	"github.com/asciiu/gomo/common/constants/status"
	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/asciiu/gomo/common/util"
	protoEngine "github.com/asciiu/gomo/execution-engine/proto/engine"
	"github.com/mattn/anko/vm"
	micro "github.com/micro/go-micro"
)

const precision = 8

// Order has conditions
type Plan struct {
	PlanID                string
	UserID                string
	ActiveCurrencySymbol  string
	ActiveCurrencyBalance float64
	CloseOnComplete       bool
	Orders                []*Order
}

type Order struct {
	OrderID     string
	Exchange    string
	MarketName  string
	Side        string
	LimitPrice  float64
	OrderType   string
	OrderStatus string
	KeyID       string
	KeyPublic   string
	KeySecret   string
	TriggerExs  []*TriggerEx
}

type TriggerEx struct {
	TriggerID string
	OrderID   string
	Name      string
	Triggered bool
	Actions   []string
	Evaluate  Expression
}

// Processor will process orders
type Engine struct {
	DB  *sql.DB
	Env *vm.Env

	Aborted   micro.Publisher
	Completed micro.Publisher
	Triggered micro.Publisher

	PriceLine map[string]float64
	Plans     []*Plan
}

// ProcessEvent will process TradeEvents. These events are published from the exchange sockets.
func (engine *Engine) ProcessTradeEvents(ctx context.Context, payload *evt.TradeEvents) error {
	plans := engine.Plans

	// TODO this implementation is fine for prototype but needs to be more efficient for production!
	for _, tradeEvent := range payload.Events {
		// update the last price for this market
		engine.PriceLine[tradeEvent.MarketName] = tradeEvent.Price

		// TODO only look at the plans relevant for this market-exchange
		for p, plan := range plans {
			for _, order := range plan.Orders {
				if tradeEvent.MarketName == order.MarketName && tradeEvent.Exchange == order.Exchange {
					for _, trigger := range order.TriggerExs {
						if isTrue, desc := trigger.Evaluate(tradeEvent.Price); isTrue {
							// remove this order from the processor
							engine.Plans = append(plans[:p], plans[p+1:]...)

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
									qty := util.ToFixed(plan.ActiveCurrencyBalance/tradeEvent.Price, precision)

									completedEvent.FinalCurrencySymbol = symbols[0]
									completedEvent.FinalCurrencyBalance = qty
									completedEvent.InitialCurrencyTraded = util.ToFixed(completedEvent.FinalCurrencyBalance*tradeEvent.Price, precision)
									completedEvent.InitialCurrencyRemainder = util.ToFixed(plan.ActiveCurrencyBalance-completedEvent.InitialCurrencyTraded, precision)
									completedEvent.Details = fmt.Sprintf("bought %.8f %s with %.8f %s", completedEvent.FinalCurrencyBalance, symbols[0], completedEvent.InitialCurrencyTraded, symbols[1])
								}

								// adjust balances for sell
								if order.Side == side.Sell {
									completedEvent.FinalCurrencySymbol = symbols[1]
									completedEvent.FinalCurrencyBalance = util.ToFixed(plan.ActiveCurrencyBalance*tradeEvent.Price, precision)
									completedEvent.InitialCurrencyRemainder = 0
									completedEvent.InitialCurrencyTraded = util.ToFixed(plan.ActiveCurrencyBalance, precision)
									completedEvent.Details = fmt.Sprintf("sold %.8f %s for %.8f %s", plan.ActiveCurrencyBalance, symbols[0], completedEvent.FinalCurrencyBalance, symbols[1])
								}

								// Never log the secrets contained in the event
								log.Printf("virtual order processed -- orderID: %s, market: %s\n", order.OrderID, order.MarketName)

								if err := engine.Completed.Publish(ctx, &completedEvent); err != nil {
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
								if err := engine.Triggered.Publish(ctx, &triggeredEvent); err != nil {
									log.Println("publish warning: ", err)
								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}

// Clients of the execution engine should use this to add plans to the engine. The processing of plans
// will be handled in the ProcessTradeEvents function.
func (engine *Engine) AddPlan(ctx context.Context, req *protoEngine.NewPlanRequest, res *protoEngine.PlanResponse) error {
	log.Printf("received new plan event for plan ID: %s", req.PlanID)

	// convert plan orders to have trigger expressions
	orders := make([]*Order, 0)

	trailingPoint := regexp.MustCompile(`^.*?TrailingStopPoint\((0\.\d{2,}),\s(\d+\.\d+).*?`)
	trailingPercent := regexp.MustCompile(`^.*?TrailingStopPercent\((0\.\d{2,}),\s(\d+\.\d+).*?`)
	for _, order := range req.Orders {

		// convert trigger to trigger expression
		expressions := make([]*TriggerEx, 0)

		for _, trigger := range order.Triggers {
			str := trigger.Code
			var expression Expression
			switch {
			case trailingPoint.MatchString(str):
				rs := trailingPoint.FindStringSubmatch(str)
				top, _ := strconv.ParseFloat(rs[1], 64)
				points, _ := strconv.ParseFloat(rs[2], 64)

				ts := TrailingStopPoint{
					Top:    top,
					Points: points,
				}
				expression = (&ts).evaluate
			case trailingPercent.MatchString(str):
				rs := trailingPercent.FindStringSubmatch(str)
				top, _ := strconv.ParseFloat(rs[1], 64)
				percent, _ := strconv.ParseFloat(rs[2], 64)

				ts := TrailingStopPercent{
					Top:     top,
					Percent: percent,
				}
				expression = (&ts).evaluate

			case str == "immediateMarketPrice":
				lastPrice := engine.PriceLine[order.MarketName]
				// execute this order right NOW! It's FOMO TIME!
				if order.OrderType == constOrder.PaperOrder {
					completedEvent := evt.CompletedOrderEvent{
						UserID:     req.UserID,
						PlanID:     req.PlanID,
						OrderID:    order.OrderID,
						Exchange:   order.Exchange,
						MarketName: order.MarketName,
						Side:       order.Side,
						InitialCurrencyBalance: req.ActiveCurrencyBalance,
						InitialCurrencySymbol:  req.ActiveCurrencySymbol,
						TriggerID:              trigger.TriggerID,
						TriggeredPrice:         lastPrice,
						TriggeredCondition:     "Jordan!",
						ExchangeOrderID:        constOrder.PaperOrder,
						ExchangeMarketName:     constOrder.PaperOrder,
						Status:                 status.Filled,
						CloseOnComplete:        req.CloseOnComplete,
					}

					symbols := strings.Split(order.MarketName, "-")
					// adjust balances for buy
					if order.Side == side.Buy {
						qty := util.ToFixed(req.ActiveCurrencyBalance/lastPrice, precision)

						completedEvent.FinalCurrencySymbol = symbols[0]
						completedEvent.FinalCurrencyBalance = qty
						completedEvent.InitialCurrencyTraded = util.ToFixed(completedEvent.FinalCurrencyBalance*lastPrice, precision)
						completedEvent.InitialCurrencyRemainder = util.ToFixed(req.ActiveCurrencyBalance-completedEvent.InitialCurrencyTraded, precision)
						completedEvent.Details = fmt.Sprintf("bought %.8f %s with %.8f %s", completedEvent.FinalCurrencyBalance, symbols[0], completedEvent.InitialCurrencyTraded, symbols[1])
					}

					// adjust balances for sell
					if order.Side == side.Sell {
						completedEvent.FinalCurrencySymbol = symbols[1]
						completedEvent.FinalCurrencyBalance = util.ToFixed(req.ActiveCurrencyBalance*lastPrice, precision)
						completedEvent.InitialCurrencyRemainder = 0
						completedEvent.InitialCurrencyTraded = util.ToFixed(req.ActiveCurrencyBalance, precision)
						completedEvent.Details = fmt.Sprintf("sold %.8f %s for %.8f %s", req.ActiveCurrencyBalance, symbols[0], completedEvent.FinalCurrencyBalance, symbols[1])
					}

					// Never log the secrets contained in the event
					log.Printf("virtual order processed -- orderID: %s, market: %s\n", order.OrderID, order.MarketName)

					if err := engine.Completed.Publish(ctx, &completedEvent); err != nil {
						log.Println("publish warning: ", err, completedEvent)
					}
				}
				// no need to add this plan just return
				res.Status = response.Success
				res.Data = &protoEngine.PlanList{
					Plans: []*protoEngine.Plan{
						&protoEngine.Plan{
							PlanID: req.PlanID,
						},
					},
				}
				return nil

			case strings.Contains(str, "price"):
				priceCond := PriceCondition{
					Env:       engine.Env,
					Statement: str,
				}
				expression = (&priceCond).evaluate
			default:
				// skip this trigger because it is invalid
				continue
			}

			expressions = append(expressions, &TriggerEx{
				TriggerID: trigger.TriggerID,
				OrderID:   trigger.OrderID,
				Name:      trigger.Name,
				Triggered: trigger.Triggered,
				Actions:   trigger.Actions,
				Evaluate:  expression,
			})
		}

		// only add an order if there are valid expressions
		if len(expressions) > 0 {
			orders = append(orders, &Order{
				OrderID:     order.OrderID,
				Exchange:    order.Exchange,
				MarketName:  order.MarketName,
				Side:        order.Side,
				LimitPrice:  order.LimitPrice,
				OrderType:   order.OrderType,
				OrderStatus: order.OrderStatus,
				KeyID:       order.KeyID,
				KeyPublic:   order.KeyPublic,
				KeySecret:   order.KeySecret,
				TriggerExs:  expressions,
			})
		}
	}

	engine.Plans = append(engine.Plans, &Plan{
		PlanID:                req.PlanID,
		UserID:                req.UserID,
		ActiveCurrencySymbol:  req.ActiveCurrencySymbol,
		ActiveCurrencyBalance: req.ActiveCurrencyBalance,
		CloseOnComplete:       req.CloseOnComplete,
		Orders:                orders,
	})

	res.Status = response.Success
	res.Data = &protoEngine.PlanList{
		Plans: []*protoEngine.Plan{
			&protoEngine.Plan{
				PlanID: req.PlanID,
			},
		},
	}

	return nil
}

func (engine *Engine) GetActivePlans(ctx context.Context, req *protoEngine.ActiveRequest, res *protoEngine.PlanResponse) error {
	return nil
}
func (engine *Engine) KillPlan(ctx context.Context, req *protoEngine.KillRequest, res *protoEngine.PlanResponse) error {
	return nil
}
func (engine *Engine) KillUserPlans(ctx context.Context, req *protoEngine.KillUserRequest, res *protoEngine.PlanResponse) error {
	return nil
}
