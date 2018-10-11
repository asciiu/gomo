package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	constAccount "github.com/asciiu/gomo/account-service/constants"
	constExt "github.com/asciiu/gomo/common/constants/exchange"
	constRes "github.com/asciiu/gomo/common/constants/response"
	protoEvt "github.com/asciiu/gomo/common/proto/events"
	commonUtil "github.com/asciiu/gomo/common/util"
	protoEngine "github.com/asciiu/gomo/execution-engine/proto/engine"
	constPlan "github.com/asciiu/gomo/plan-service/constants"
	"github.com/lib/pq"
	"github.com/mattn/anko/vm"
	micro "github.com/micro/go-micro"
)

const precision = 8

// Order has conditions
type Plan struct {
	PlanID                  string
	UserID                  string
	CommittedCurrencySymbol string
	committedCurrencyAmount float64
	CloseOnComplete         bool
	Orders                  []*Order
}

type Order struct {
	OrderID     string
	Exchange    string
	MarketName  string
	Side        string
	LimitPrice  float64
	OrderType   string
	OrderStatus string
	AccountID   string
	AccountType string
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

	Completed   micro.Publisher
	Triggered   micro.Publisher
	FillBinance micro.Publisher

	PriceLine map[string]float64
	Plans     []*Plan
}

func (engine *Engine) HandleAccountDeleted(ctx context.Context, evt *protoEvt.DeletedAccountEvent) error {
	plans := engine.Plans

	for i, plan := range plans {
		orders := plan.Orders
		for j, order := range orders {
			if order.AccountID == evt.AccountID {
				// remove this order
				orders = append(orders[:j], orders[j+1:]...)
			}
		}
		if len(orders) == 0 {
			engine.Plans = append(plans[:i], plans[i+1:]...)
		}
	}

	return nil
}

// ProcessEvent will process TradeEvents. These events are published from the exchange sockets.
// Whether or not a trigger for an order will execute will be deteremined here.
func (engine *Engine) HandleTradeEvents(ctx context.Context, payload *protoEvt.TradeEvents) error {
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
						// evalute the trigger - is it true?
						if isTrue, desc := trigger.Evaluate(tradeEvent.Price); isTrue {
							// remove this order from the processor
							engine.Plans = append(plans[:p], plans[p+1:]...)
							now := string(pq.FormatTimestamp(time.Now().UTC()))

							// assume if there is no key that it is paper
							// paper accounts should NOT have a public key.
							if order.AccountType == constAccount.AccountPaper {
								completedEvent := protoEvt.CompletedOrderEvent{
									UserID:                 plan.UserID,
									PlanID:                 plan.PlanID,
									OrderID:                order.OrderID,
									Exchange:               order.Exchange,
									MarketName:             order.MarketName,
									Side:                   order.Side,
									AccountID:              order.AccountID,
									InitialCurrencyBalance: plan.committedCurrencyAmount,
									InitialCurrencySymbol:  plan.CommittedCurrencySymbol,
									TriggerID:              trigger.TriggerID,
									TriggeredPrice:         tradeEvent.Price,
									TriggeredCondition:     desc,
									TriggeredTime:          now,
									ExchangeMarketName:     constPlan.PaperOrder,
									ExchangePrice:          tradeEvent.Price,
									ExchangeTime:           now,
									Status:                 constPlan.Filled,
									CloseOnComplete:        plan.CloseOnComplete,
								}

								symbols := strings.Split(order.MarketName, "-")
								// adjust balances for buy
								if order.Side == constPlan.Buy {
									qty := commonUtil.ToFixedFloor(plan.committedCurrencyAmount/tradeEvent.Price, precision)

									completedEvent.ExchangePrice = tradeEvent.Price
									completedEvent.FinalCurrencySymbol = symbols[0]
									completedEvent.FinalCurrencyBalance = qty
									completedEvent.InitialCurrencyTraded = commonUtil.ToFixedFloor(completedEvent.FinalCurrencyBalance*tradeEvent.Price, precision)
									completedEvent.InitialCurrencyRemainder = commonUtil.ToFixedFloor(plan.committedCurrencyAmount-completedEvent.InitialCurrencyTraded, precision)
									completedEvent.Details = fmt.Sprintf("orderID: %s, bought %.8f %s with %.8f %s", completedEvent.OrderID, completedEvent.FinalCurrencyBalance, symbols[0], completedEvent.InitialCurrencyTraded, symbols[1])
								}

								// adjust balances for sell
								if order.Side == constPlan.Sell {
									completedEvent.FinalCurrencySymbol = symbols[1]
									completedEvent.FinalCurrencyBalance = commonUtil.ToFixedFloor(plan.committedCurrencyAmount*tradeEvent.Price, precision)
									completedEvent.InitialCurrencyRemainder = 0
									completedEvent.InitialCurrencyTraded = commonUtil.ToFixedFloor(plan.committedCurrencyAmount, precision)
									completedEvent.Details = fmt.Sprintf("orderID: %s, sold %.8f %s for %.8f %s", completedEvent.OrderID, plan.committedCurrencyAmount, symbols[0], completedEvent.FinalCurrencyBalance, symbols[1])
								}

								// Never log the secrets contained in the event
								log.Printf("virtual order processed -- orderID: %s, market: %s\n", order.OrderID, order.MarketName)

								if err := engine.Completed.Publish(ctx, &completedEvent); err != nil {
									log.Println("publish warning: ", err, completedEvent)
								}
							} else {
								quantity := 0.0
								switch {
								case order.Side == constPlan.Buy && order.OrderType == constPlan.LimitOrder:
									quantity = plan.committedCurrencyAmount / order.LimitPrice

								case order.Side == constPlan.Buy && order.OrderType == constPlan.MarketOrder:
									quantity = plan.committedCurrencyAmount / tradeEvent.Price

								default:
									// sell entire active balance
									quantity = plan.committedCurrencyAmount
								}

								// convert this active order event to a triggered order event
								triggeredEvent := protoEvt.TriggeredOrderEvent{
									Exchange:                order.Exchange,
									OrderID:                 order.OrderID,
									PlanID:                  plan.PlanID,
									UserID:                  plan.UserID,
									AccountID:               order.AccountID,
									CommittedCurrencySymbol: plan.CommittedCurrencySymbol,
									CommittedCurrencyAmount: plan.committedCurrencyAmount,
									KeyPublic:               order.KeyPublic,
									KeySecret:               order.KeySecret,
									MarketName:              order.MarketName,
									Side:                    order.Side,
									OrderType:               order.OrderType,
									LimitPrice:              order.LimitPrice,
									Quantity:                quantity,
									TriggerID:               trigger.TriggerID,
									TriggeredPrice:          tradeEvent.Price,
									TriggeredCondition:      desc,
									TriggeredTime:           now,
								}

								// Never log the secrets contained in the event
								log.Printf("triggered order -- %+v\n", triggeredEvent)

								switch triggeredEvent.Exchange {
								case constExt.Binance:
									// send out fill request to binance subscribers
									if err := engine.FillBinance.Publish(ctx, &triggeredEvent); err != nil {
										log.Println("publish warning: ", err)
									}
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
				now := string(pq.FormatTimestamp(time.Now().UTC()))

				// execute this order right NOW! It's FOMO TIME! but only when we have the last price
				// for the market in question
				if order.AccountType == constAccount.AccountPaper && lastPrice != 0 {
					completedEvent := protoEvt.CompletedOrderEvent{
						UserID:                 req.UserID,
						PlanID:                 req.PlanID,
						OrderID:                order.OrderID,
						Exchange:               order.Exchange,
						MarketName:             order.MarketName,
						Side:                   order.Side,
						AccountID:              order.AccountID,
						InitialCurrencyBalance: req.CommittedCurrencyAmount,
						InitialCurrencySymbol:  req.CommittedCurrencySymbol,
						TriggerID:              trigger.TriggerID,
						TriggeredPrice:         lastPrice,
						TriggeredCondition:     "Immeadiate!",
						TriggeredTime:          now,
						ExchangeMarketName:     constPlan.PaperOrder,
						ExchangeTime:           now,
						Status:                 constPlan.Filled,
						CloseOnComplete:        req.CloseOnComplete,
					}

					symbols := strings.Split(order.MarketName, "-")
					// adjust balances for buy
					if order.Side == constPlan.Buy {
						qty := commonUtil.ToFixedFloor(req.CommittedCurrencyAmount/lastPrice, precision)

						completedEvent.FinalCurrencySymbol = symbols[0]
						completedEvent.FinalCurrencyBalance = qty
						completedEvent.InitialCurrencyTraded = commonUtil.ToFixedFloor(completedEvent.FinalCurrencyBalance*lastPrice, precision)
						completedEvent.InitialCurrencyRemainder = commonUtil.ToFixedFloor(req.CommittedCurrencyAmount-completedEvent.InitialCurrencyTraded, precision)
						completedEvent.Details = fmt.Sprintf("bought %.8f %s with %.8f %s", completedEvent.FinalCurrencyBalance, symbols[0], completedEvent.InitialCurrencyTraded, symbols[1])
					}

					// adjust balances for sell
					if order.Side == constPlan.Sell {
						completedEvent.FinalCurrencySymbol = symbols[1]
						completedEvent.FinalCurrencyBalance = commonUtil.ToFixedFloor(req.CommittedCurrencyAmount*lastPrice, precision)
						completedEvent.InitialCurrencyRemainder = 0
						completedEvent.InitialCurrencyTraded = commonUtil.ToFixedFloor(req.CommittedCurrencyAmount, precision)
						completedEvent.Details = fmt.Sprintf("sold %.8f %s for %.8f %s", req.CommittedCurrencyAmount, symbols[0], completedEvent.FinalCurrencyBalance, symbols[1])
					}

					// Never log the secrets contained in the event
					log.Printf("virtual order processed -- orderID: %s, market: %s\n", order.OrderID, order.MarketName)

					if err := engine.Completed.Publish(ctx, &completedEvent); err != nil {
						log.Println("publish warning: ", err, completedEvent)
					}
					// no need to add this plan just return
					res.Status = constRes.Success
					res.Data = &protoEngine.PlanList{
						Plans: []*protoEngine.Plan{
							&protoEngine.Plan{
								PlanID: req.PlanID,
							},
						},
					}
					return nil
				}

				immediate := new(Immediate)
				expression = immediate.evaluate

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
				AccountID:   order.AccountID,
				AccountType: order.AccountType,
				KeyPublic:   order.KeyPublic,
				KeySecret:   order.KeySecret,
				TriggerExs:  expressions,
			})
		}
	}

	engine.Plans = append(engine.Plans, &Plan{
		PlanID:                  req.PlanID,
		UserID:                  req.UserID,
		CommittedCurrencySymbol: req.CommittedCurrencySymbol,
		committedCurrencyAmount: req.CommittedCurrencyAmount,
		CloseOnComplete:         req.CloseOnComplete,
		Orders:                  orders,
	})

	res.Status = constRes.Success
	res.Data = &protoEngine.PlanList{
		Plans: []*protoEngine.Plan{
			&protoEngine.Plan{
				PlanID: req.PlanID,
			},
		},
	}

	return nil
}

// Get active plans
func (engine *Engine) GetActivePlans(ctx context.Context, req *protoEngine.ActiveRequest, res *protoEngine.PlanResponse) error {
	plans := engine.Plans
	livePlans := make([]*protoEngine.Plan, 0)

	for _, plan := range plans {
		livePlans = append(livePlans, &protoEngine.Plan{
			PlanID: plan.PlanID,
		})
	}

	res.Status = constRes.Success
	res.Data = &protoEngine.PlanList{
		Plans: livePlans,
	}
	return nil
}

// kills a single plan. If the plan was not found returns error.
func (engine *Engine) KillPlan(ctx context.Context, req *protoEngine.KillRequest, res *protoEngine.PlanResponse) error {

	plans := engine.Plans

	for i, plan := range plans {
		if plan.PlanID == req.PlanID {
			// remove this plan by appending the slice up until i (:i excludes i)
			// with all elements after i
			engine.Plans = append(plans[:i], plans[i+1:]...)
			res.Status = constRes.Success
			res.Data = &protoEngine.PlanList{
				Plans: []*protoEngine.Plan{
					&protoEngine.Plan{
						PlanID: req.PlanID,
					},
				},
			}
			return nil
		}
	}

	return errors.New(fmt.Sprintf("plan %s not found", req.PlanID))
}

// Kill all plans belonging to a user.
func (engine *Engine) KillUserPlans(ctx context.Context, req *protoEngine.KillUserRequest, res *protoEngine.PlanResponse) error {
	plans := engine.Plans
	killedPlans := make([]*protoEngine.Plan, 0)

	for i, plan := range plans {
		if plan.UserID == req.UserID {
			// remove this plan by appending the slice up until i (:i excludes i)
			// with all elements after i
			engine.Plans = append(plans[:i], plans[i+1:]...)

			killedPlans = append(killedPlans, &protoEngine.Plan{
				PlanID: plan.PlanID,
			})
		}
	}

	res.Status = constRes.Success
	res.Data = &protoEngine.PlanList{
		Plans: killedPlans,
	}
	return nil
}
