// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/events/events.proto

/*
Package common_events is a generated protocol buffer package.

It is generated from these files:
	proto/events/events.proto

It has these top-level messages:
	TradeEvent
	TradeEvents
	Auth
	CandleDataRequest
	NewPlanEvent
	Order
	Trigger
	TriggeredOrderEvent
	CompletedOrderEvent
	AbortedOrderEvent
	EngineStartEvent
	DeletedAccountEvent
*/
package common_events

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Events
type TradeEvent struct {
	Exchange   string  `protobuf:"bytes,1,opt,name=exchange" json:"exchange,omitempty"`
	Type       string  `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	EventTime  string  `protobuf:"bytes,3,opt,name=eventTime" json:"eventTime,omitempty"`
	MarketName string  `protobuf:"bytes,4,opt,name=marketName" json:"marketName,omitempty"`
	TradeID    string  `protobuf:"bytes,5,opt,name=tradeID" json:"tradeID,omitempty"`
	Price      float64 `protobuf:"fixed64,6,opt,name=price" json:"price,omitempty"`
	Quantity   float64 `protobuf:"fixed64,7,opt,name=quantity" json:"quantity,omitempty"`
	Total      float64 `protobuf:"fixed64,8,opt,name=total" json:"total,omitempty"`
}

func (m *TradeEvent) Reset()                    { *m = TradeEvent{} }
func (m *TradeEvent) String() string            { return proto.CompactTextString(m) }
func (*TradeEvent) ProtoMessage()               {}
func (*TradeEvent) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *TradeEvent) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *TradeEvent) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *TradeEvent) GetEventTime() string {
	if m != nil {
		return m.EventTime
	}
	return ""
}

func (m *TradeEvent) GetMarketName() string {
	if m != nil {
		return m.MarketName
	}
	return ""
}

func (m *TradeEvent) GetTradeID() string {
	if m != nil {
		return m.TradeID
	}
	return ""
}

func (m *TradeEvent) GetPrice() float64 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *TradeEvent) GetQuantity() float64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

func (m *TradeEvent) GetTotal() float64 {
	if m != nil {
		return m.Total
	}
	return 0
}

type TradeEvents struct {
	Events []*TradeEvent `protobuf:"bytes,1,rep,name=events" json:"events,omitempty"`
}

func (m *TradeEvents) Reset()                    { *m = TradeEvents{} }
func (m *TradeEvents) String() string            { return proto.CompactTextString(m) }
func (*TradeEvents) ProtoMessage()               {}
func (*TradeEvents) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *TradeEvents) GetEvents() []*TradeEvent {
	if m != nil {
		return m.Events
	}
	return nil
}

type Auth struct {
	Key    string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Secret string `protobuf:"bytes,2,opt,name=secret" json:"secret,omitempty"`
}

func (m *Auth) Reset()                    { *m = Auth{} }
func (m *Auth) String() string            { return proto.CompactTextString(m) }
func (*Auth) ProtoMessage()               {}
func (*Auth) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Auth) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Auth) GetSecret() string {
	if m != nil {
		return m.Secret
	}
	return ""
}

type CandleDataRequest struct {
	Exchange   string `protobuf:"bytes,1,opt,name=exchange" json:"exchange,omitempty"`
	MarketName string `protobuf:"bytes,2,opt,name=marketName" json:"marketName,omitempty"`
	Interval   string `protobuf:"bytes,3,opt,name=interval" json:"interval,omitempty"`
}

func (m *CandleDataRequest) Reset()                    { *m = CandleDataRequest{} }
func (m *CandleDataRequest) String() string            { return proto.CompactTextString(m) }
func (*CandleDataRequest) ProtoMessage()               {}
func (*CandleDataRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *CandleDataRequest) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *CandleDataRequest) GetMarketName() string {
	if m != nil {
		return m.MarketName
	}
	return ""
}

func (m *CandleDataRequest) GetInterval() string {
	if m != nil {
		return m.Interval
	}
	return ""
}

type NewPlanEvent struct {
	PlanID                string   `protobuf:"bytes,1,opt,name=planID" json:"planID,omitempty"`
	UserID                string   `protobuf:"bytes,2,opt,name=userID" json:"userID,omitempty"`
	ActiveCurrencySymbol  string   `protobuf:"bytes,3,opt,name=activeCurrencySymbol" json:"activeCurrencySymbol,omitempty"`
	ActiveCurrencyBalance float64  `protobuf:"fixed64,4,opt,name=activeCurrencyBalance" json:"activeCurrencyBalance,omitempty"`
	CloseOnComplete       bool     `protobuf:"varint,5,opt,name=closeOnComplete" json:"closeOnComplete,omitempty"`
	Orders                []*Order `protobuf:"bytes,6,rep,name=orders" json:"orders,omitempty"`
}

func (m *NewPlanEvent) Reset()                    { *m = NewPlanEvent{} }
func (m *NewPlanEvent) String() string            { return proto.CompactTextString(m) }
func (*NewPlanEvent) ProtoMessage()               {}
func (*NewPlanEvent) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *NewPlanEvent) GetPlanID() string {
	if m != nil {
		return m.PlanID
	}
	return ""
}

func (m *NewPlanEvent) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *NewPlanEvent) GetActiveCurrencySymbol() string {
	if m != nil {
		return m.ActiveCurrencySymbol
	}
	return ""
}

func (m *NewPlanEvent) GetActiveCurrencyBalance() float64 {
	if m != nil {
		return m.ActiveCurrencyBalance
	}
	return 0
}

func (m *NewPlanEvent) GetCloseOnComplete() bool {
	if m != nil {
		return m.CloseOnComplete
	}
	return false
}

func (m *NewPlanEvent) GetOrders() []*Order {
	if m != nil {
		return m.Orders
	}
	return nil
}

type Order struct {
	OrderID     string     `protobuf:"bytes,1,opt,name=orderID" json:"orderID,omitempty"`
	Exchange    string     `protobuf:"bytes,2,opt,name=exchange" json:"exchange,omitempty"`
	MarketName  string     `protobuf:"bytes,3,opt,name=marketName" json:"marketName,omitempty"`
	Side        string     `protobuf:"bytes,4,opt,name=side" json:"side,omitempty"`
	LimitPrice  float64    `protobuf:"fixed64,5,opt,name=limitPrice" json:"limitPrice,omitempty"`
	OrderType   string     `protobuf:"bytes,6,opt,name=orderType" json:"orderType,omitempty"`
	OrderStatus string     `protobuf:"bytes,7,opt,name=orderStatus" json:"orderStatus,omitempty"`
	AccountID   string     `protobuf:"bytes,8,opt,name=accountID" json:"accountID,omitempty"`
	KeyPublic   string     `protobuf:"bytes,9,opt,name=keyPublic" json:"keyPublic,omitempty"`
	KeySecret   string     `protobuf:"bytes,10,opt,name=keySecret" json:"keySecret,omitempty"`
	Triggers    []*Trigger `protobuf:"bytes,11,rep,name=triggers" json:"triggers,omitempty"`
}

func (m *Order) Reset()                    { *m = Order{} }
func (m *Order) String() string            { return proto.CompactTextString(m) }
func (*Order) ProtoMessage()               {}
func (*Order) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Order) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
}

func (m *Order) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *Order) GetMarketName() string {
	if m != nil {
		return m.MarketName
	}
	return ""
}

func (m *Order) GetSide() string {
	if m != nil {
		return m.Side
	}
	return ""
}

func (m *Order) GetLimitPrice() float64 {
	if m != nil {
		return m.LimitPrice
	}
	return 0
}

func (m *Order) GetOrderType() string {
	if m != nil {
		return m.OrderType
	}
	return ""
}

func (m *Order) GetOrderStatus() string {
	if m != nil {
		return m.OrderStatus
	}
	return ""
}

func (m *Order) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *Order) GetKeyPublic() string {
	if m != nil {
		return m.KeyPublic
	}
	return ""
}

func (m *Order) GetKeySecret() string {
	if m != nil {
		return m.KeySecret
	}
	return ""
}

func (m *Order) GetTriggers() []*Trigger {
	if m != nil {
		return m.Triggers
	}
	return nil
}

type Trigger struct {
	TriggerID string   `protobuf:"bytes,1,opt,name=triggerID" json:"triggerID,omitempty"`
	OrderID   string   `protobuf:"bytes,2,opt,name=orderID" json:"orderID,omitempty"`
	Name      string   `protobuf:"bytes,5,opt,name=name" json:"name,omitempty"`
	Code      string   `protobuf:"bytes,6,opt,name=code" json:"code,omitempty"`
	Triggered bool     `protobuf:"varint,7,opt,name=triggered" json:"triggered,omitempty"`
	Actions   []string `protobuf:"bytes,8,rep,name=actions" json:"actions,omitempty"`
}

func (m *Trigger) Reset()                    { *m = Trigger{} }
func (m *Trigger) String() string            { return proto.CompactTextString(m) }
func (*Trigger) ProtoMessage()               {}
func (*Trigger) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *Trigger) GetTriggerID() string {
	if m != nil {
		return m.TriggerID
	}
	return ""
}

func (m *Trigger) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
}

func (m *Trigger) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Trigger) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Trigger) GetTriggered() bool {
	if m != nil {
		return m.Triggered
	}
	return false
}

func (m *Trigger) GetActions() []string {
	if m != nil {
		return m.Actions
	}
	return nil
}

type TriggeredOrderEvent struct {
	Exchange              string  `protobuf:"bytes,1,opt,name=exchange" json:"exchange,omitempty"`
	OrderID               string  `protobuf:"bytes,2,opt,name=orderID" json:"orderID,omitempty"`
	PlanID                string  `protobuf:"bytes,3,opt,name=planID" json:"planID,omitempty"`
	UserID                string  `protobuf:"bytes,4,opt,name=userID" json:"userID,omitempty"`
	AccountID             string  `protobuf:"bytes,5,opt,name=accountID" json:"accountID,omitempty"`
	ActiveCurrencySymbol  string  `protobuf:"bytes,6,opt,name=activeCurrencySymbol" json:"activeCurrencySymbol,omitempty"`
	ActiveCurrencyBalance float64 `protobuf:"fixed64,7,opt,name=activeCurrencyBalance" json:"activeCurrencyBalance,omitempty"`
	Quantity              float64 `protobuf:"fixed64,8,opt,name=quantity" json:"quantity,omitempty"`
	KeyPublic             string  `protobuf:"bytes,9,opt,name=keyPublic" json:"keyPublic,omitempty"`
	KeySecret             string  `protobuf:"bytes,10,opt,name=keySecret" json:"keySecret,omitempty"`
	MarketName            string  `protobuf:"bytes,11,opt,name=marketName" json:"marketName,omitempty"`
	Side                  string  `protobuf:"bytes,12,opt,name=side" json:"side,omitempty"`
	OrderType             string  `protobuf:"bytes,13,opt,name=orderType" json:"orderType,omitempty"`
	LimitPrice            float64 `protobuf:"fixed64,14,opt,name=limitPrice" json:"limitPrice,omitempty"`
	TriggerID             string  `protobuf:"bytes,15,opt,name=triggerID" json:"triggerID,omitempty"`
	TriggeredPrice        float64 `protobuf:"fixed64,16,opt,name=triggeredPrice" json:"triggeredPrice,omitempty"`
	TriggeredCondition    string  `protobuf:"bytes,17,opt,name=triggeredCondition" json:"triggeredCondition,omitempty"`
	CloseOnComplete       bool    `protobuf:"varint,18,opt,name=closeOnComplete" json:"closeOnComplete,omitempty"`
}

func (m *TriggeredOrderEvent) Reset()                    { *m = TriggeredOrderEvent{} }
func (m *TriggeredOrderEvent) String() string            { return proto.CompactTextString(m) }
func (*TriggeredOrderEvent) ProtoMessage()               {}
func (*TriggeredOrderEvent) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *TriggeredOrderEvent) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *TriggeredOrderEvent) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
}

func (m *TriggeredOrderEvent) GetPlanID() string {
	if m != nil {
		return m.PlanID
	}
	return ""
}

func (m *TriggeredOrderEvent) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *TriggeredOrderEvent) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *TriggeredOrderEvent) GetActiveCurrencySymbol() string {
	if m != nil {
		return m.ActiveCurrencySymbol
	}
	return ""
}

func (m *TriggeredOrderEvent) GetActiveCurrencyBalance() float64 {
	if m != nil {
		return m.ActiveCurrencyBalance
	}
	return 0
}

func (m *TriggeredOrderEvent) GetQuantity() float64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

func (m *TriggeredOrderEvent) GetKeyPublic() string {
	if m != nil {
		return m.KeyPublic
	}
	return ""
}

func (m *TriggeredOrderEvent) GetKeySecret() string {
	if m != nil {
		return m.KeySecret
	}
	return ""
}

func (m *TriggeredOrderEvent) GetMarketName() string {
	if m != nil {
		return m.MarketName
	}
	return ""
}

func (m *TriggeredOrderEvent) GetSide() string {
	if m != nil {
		return m.Side
	}
	return ""
}

func (m *TriggeredOrderEvent) GetOrderType() string {
	if m != nil {
		return m.OrderType
	}
	return ""
}

func (m *TriggeredOrderEvent) GetLimitPrice() float64 {
	if m != nil {
		return m.LimitPrice
	}
	return 0
}

func (m *TriggeredOrderEvent) GetTriggerID() string {
	if m != nil {
		return m.TriggerID
	}
	return ""
}

func (m *TriggeredOrderEvent) GetTriggeredPrice() float64 {
	if m != nil {
		return m.TriggeredPrice
	}
	return 0
}

func (m *TriggeredOrderEvent) GetTriggeredCondition() string {
	if m != nil {
		return m.TriggeredCondition
	}
	return ""
}

func (m *TriggeredOrderEvent) GetCloseOnComplete() bool {
	if m != nil {
		return m.CloseOnComplete
	}
	return false
}

type CompletedOrderEvent struct {
	OrderID                  string  `protobuf:"bytes,1,opt,name=orderID" json:"orderID,omitempty"`
	PlanID                   string  `protobuf:"bytes,2,opt,name=planID" json:"planID,omitempty"`
	UserID                   string  `protobuf:"bytes,3,opt,name=userID" json:"userID,omitempty"`
	Exchange                 string  `protobuf:"bytes,4,opt,name=exchange" json:"exchange,omitempty"`
	MarketName               string  `protobuf:"bytes,5,opt,name=marketName" json:"marketName,omitempty"`
	Side                     string  `protobuf:"bytes,6,opt,name=side" json:"side,omitempty"`
	AccountID                string  `protobuf:"bytes,7,opt,name=accountID" json:"accountID,omitempty"`
	InitialCurrencySymbol    string  `protobuf:"bytes,8,opt,name=initialCurrencySymbol" json:"initialCurrencySymbol,omitempty"`
	InitialCurrencyBalance   float64 `protobuf:"fixed64,9,opt,name=initialCurrencyBalance" json:"initialCurrencyBalance,omitempty"`
	InitialCurrencyTraded    float64 `protobuf:"fixed64,10,opt,name=initialCurrencyTraded" json:"initialCurrencyTraded,omitempty"`
	InitialCurrencyRemainder float64 `protobuf:"fixed64,11,opt,name=initialCurrencyRemainder" json:"initialCurrencyRemainder,omitempty"`
	InitialCurrencyPrice     float64 `protobuf:"fixed64,12,opt,name=initialCurrencyPrice" json:"initialCurrencyPrice,omitempty"`
	FinalCurrencySymbol      string  `protobuf:"bytes,13,opt,name=finalCurrencySymbol" json:"finalCurrencySymbol,omitempty"`
	FinalCurrencyBalance     float64 `protobuf:"fixed64,14,opt,name=finalCurrencyBalance" json:"finalCurrencyBalance,omitempty"`
	FeeCurrencySymbol        string  `protobuf:"bytes,15,opt,name=feeCurrencySymbol" json:"feeCurrencySymbol,omitempty"`
	FeeCurrencyAmount        float64 `protobuf:"fixed64,16,opt,name=feeCurrencyAmount" json:"feeCurrencyAmount,omitempty"`
	TriggerID                string  `protobuf:"bytes,17,opt,name=triggerID" json:"triggerID,omitempty"`
	TriggeredPrice           float64 `protobuf:"fixed64,18,opt,name=triggeredPrice" json:"triggeredPrice,omitempty"`
	TriggeredCondition       string  `protobuf:"bytes,19,opt,name=triggeredCondition" json:"triggeredCondition,omitempty"`
	ExchangeOrderID          string  `protobuf:"bytes,20,opt,name=exchangeOrderID" json:"exchangeOrderID,omitempty"`
	ExchangeMarketName       string  `protobuf:"bytes,21,opt,name=exchangeMarketName" json:"exchangeMarketName,omitempty"`
	ExchangeTime             string  `protobuf:"bytes,22,opt,name=exchangeTime" json:"exchangeTime,omitempty"`
	Status                   string  `protobuf:"bytes,23,opt,name=status" json:"status,omitempty"`
	Details                  string  `protobuf:"bytes,24,opt,name=details" json:"details,omitempty"`
	CloseOnComplete          bool    `protobuf:"varint,25,opt,name=closeOnComplete" json:"closeOnComplete,omitempty"`
}

func (m *CompletedOrderEvent) Reset()                    { *m = CompletedOrderEvent{} }
func (m *CompletedOrderEvent) String() string            { return proto.CompactTextString(m) }
func (*CompletedOrderEvent) ProtoMessage()               {}
func (*CompletedOrderEvent) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *CompletedOrderEvent) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
}

func (m *CompletedOrderEvent) GetPlanID() string {
	if m != nil {
		return m.PlanID
	}
	return ""
}

func (m *CompletedOrderEvent) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *CompletedOrderEvent) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *CompletedOrderEvent) GetMarketName() string {
	if m != nil {
		return m.MarketName
	}
	return ""
}

func (m *CompletedOrderEvent) GetSide() string {
	if m != nil {
		return m.Side
	}
	return ""
}

func (m *CompletedOrderEvent) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *CompletedOrderEvent) GetInitialCurrencySymbol() string {
	if m != nil {
		return m.InitialCurrencySymbol
	}
	return ""
}

func (m *CompletedOrderEvent) GetInitialCurrencyBalance() float64 {
	if m != nil {
		return m.InitialCurrencyBalance
	}
	return 0
}

func (m *CompletedOrderEvent) GetInitialCurrencyTraded() float64 {
	if m != nil {
		return m.InitialCurrencyTraded
	}
	return 0
}

func (m *CompletedOrderEvent) GetInitialCurrencyRemainder() float64 {
	if m != nil {
		return m.InitialCurrencyRemainder
	}
	return 0
}

func (m *CompletedOrderEvent) GetInitialCurrencyPrice() float64 {
	if m != nil {
		return m.InitialCurrencyPrice
	}
	return 0
}

func (m *CompletedOrderEvent) GetFinalCurrencySymbol() string {
	if m != nil {
		return m.FinalCurrencySymbol
	}
	return ""
}

func (m *CompletedOrderEvent) GetFinalCurrencyBalance() float64 {
	if m != nil {
		return m.FinalCurrencyBalance
	}
	return 0
}

func (m *CompletedOrderEvent) GetFeeCurrencySymbol() string {
	if m != nil {
		return m.FeeCurrencySymbol
	}
	return ""
}

func (m *CompletedOrderEvent) GetFeeCurrencyAmount() float64 {
	if m != nil {
		return m.FeeCurrencyAmount
	}
	return 0
}

func (m *CompletedOrderEvent) GetTriggerID() string {
	if m != nil {
		return m.TriggerID
	}
	return ""
}

func (m *CompletedOrderEvent) GetTriggeredPrice() float64 {
	if m != nil {
		return m.TriggeredPrice
	}
	return 0
}

func (m *CompletedOrderEvent) GetTriggeredCondition() string {
	if m != nil {
		return m.TriggeredCondition
	}
	return ""
}

func (m *CompletedOrderEvent) GetExchangeOrderID() string {
	if m != nil {
		return m.ExchangeOrderID
	}
	return ""
}

func (m *CompletedOrderEvent) GetExchangeMarketName() string {
	if m != nil {
		return m.ExchangeMarketName
	}
	return ""
}

func (m *CompletedOrderEvent) GetExchangeTime() string {
	if m != nil {
		return m.ExchangeTime
	}
	return ""
}

func (m *CompletedOrderEvent) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *CompletedOrderEvent) GetDetails() string {
	if m != nil {
		return m.Details
	}
	return ""
}

func (m *CompletedOrderEvent) GetCloseOnComplete() bool {
	if m != nil {
		return m.CloseOnComplete
	}
	return false
}

type AbortedOrderEvent struct {
	OrderID string `protobuf:"bytes,1,opt,name=orderID" json:"orderID,omitempty"`
	PlanID  string `protobuf:"bytes,2,opt,name=planID" json:"planID,omitempty"`
}

func (m *AbortedOrderEvent) Reset()                    { *m = AbortedOrderEvent{} }
func (m *AbortedOrderEvent) String() string            { return proto.CompactTextString(m) }
func (*AbortedOrderEvent) ProtoMessage()               {}
func (*AbortedOrderEvent) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *AbortedOrderEvent) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
}

func (m *AbortedOrderEvent) GetPlanID() string {
	if m != nil {
		return m.PlanID
	}
	return ""
}

type EngineStartEvent struct {
	EngineID string `protobuf:"bytes,1,opt,name=engineID" json:"engineID,omitempty"`
}

func (m *EngineStartEvent) Reset()                    { *m = EngineStartEvent{} }
func (m *EngineStartEvent) String() string            { return proto.CompactTextString(m) }
func (*EngineStartEvent) ProtoMessage()               {}
func (*EngineStartEvent) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *EngineStartEvent) GetEngineID() string {
	if m != nil {
		return m.EngineID
	}
	return ""
}

type DeletedAccountEvent struct {
	AccountID string `protobuf:"bytes,1,opt,name=accountID" json:"accountID,omitempty"`
}

func (m *DeletedAccountEvent) Reset()                    { *m = DeletedAccountEvent{} }
func (m *DeletedAccountEvent) String() string            { return proto.CompactTextString(m) }
func (*DeletedAccountEvent) ProtoMessage()               {}
func (*DeletedAccountEvent) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *DeletedAccountEvent) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func init() {
	proto.RegisterType((*TradeEvent)(nil), "common.events.TradeEvent")
	proto.RegisterType((*TradeEvents)(nil), "common.events.TradeEvents")
	proto.RegisterType((*Auth)(nil), "common.events.Auth")
	proto.RegisterType((*CandleDataRequest)(nil), "common.events.CandleDataRequest")
	proto.RegisterType((*NewPlanEvent)(nil), "common.events.NewPlanEvent")
	proto.RegisterType((*Order)(nil), "common.events.Order")
	proto.RegisterType((*Trigger)(nil), "common.events.Trigger")
	proto.RegisterType((*TriggeredOrderEvent)(nil), "common.events.TriggeredOrderEvent")
	proto.RegisterType((*CompletedOrderEvent)(nil), "common.events.CompletedOrderEvent")
	proto.RegisterType((*AbortedOrderEvent)(nil), "common.events.AbortedOrderEvent")
	proto.RegisterType((*EngineStartEvent)(nil), "common.events.EngineStartEvent")
	proto.RegisterType((*DeletedAccountEvent)(nil), "common.events.DeletedAccountEvent")
}

func init() { proto.RegisterFile("proto/events/events.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 970 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x56, 0xcd, 0x6e, 0xe3, 0x36,
	0x10, 0x86, 0xfc, 0x23, 0xdb, 0xe3, 0xec, 0x6e, 0x4c, 0x27, 0x2e, 0xb3, 0x28, 0x0a, 0x43, 0x87,
	0x22, 0x87, 0x85, 0xbb, 0xcd, 0x2e, 0x7a, 0xe8, 0xa9, 0x69, 0x9c, 0x43, 0x0e, 0xdd, 0x04, 0x4a,
	0x5e, 0x80, 0x91, 0xb8, 0x59, 0x22, 0x12, 0xe5, 0x95, 0xe8, 0xb4, 0xbe, 0xf5, 0x55, 0xfa, 0x5e,
	0x05, 0xfa, 0x04, 0x3d, 0xf5, 0x05, 0x0a, 0x0e, 0x29, 0xc9, 0x92, 0x25, 0x2f, 0xda, 0x9c, 0xcc,
	0xf9, 0x3e, 0x72, 0x68, 0xcd, 0x7c, 0xdf, 0x48, 0x70, 0xb2, 0x4a, 0x13, 0x95, 0x7c, 0xc7, 0x9f,
	0xb8, 0x54, 0x99, 0xfd, 0x59, 0x20, 0x46, 0x5e, 0x04, 0x49, 0x1c, 0x27, 0x72, 0x61, 0x40, 0xef,
	0x4f, 0x07, 0xe0, 0x2e, 0x65, 0x21, 0xbf, 0xd4, 0x31, 0x79, 0x0d, 0x43, 0xfe, 0x5b, 0xf0, 0x89,
	0xc9, 0x07, 0x4e, 0x9d, 0xb9, 0x73, 0x3a, 0xf2, 0x8b, 0x98, 0x10, 0xe8, 0xa9, 0xcd, 0x8a, 0xd3,
	0x0e, 0xe2, 0xb8, 0x26, 0x5f, 0xc3, 0x08, 0x13, 0xdd, 0x89, 0x98, 0xd3, 0x2e, 0x12, 0x25, 0x40,
	0xbe, 0x01, 0x88, 0x59, 0xfa, 0xc8, 0xd5, 0x07, 0x16, 0x73, 0xda, 0x43, 0x7a, 0x0b, 0x21, 0x14,
	0x06, 0x4a, 0xdf, 0x7d, 0xb5, 0xa4, 0x7d, 0x24, 0xf3, 0x90, 0x1c, 0x41, 0x7f, 0x95, 0x8a, 0x80,
	0x53, 0x77, 0xee, 0x9c, 0x3a, 0xbe, 0x09, 0xf4, 0xbf, 0xfb, 0xbc, 0x66, 0x52, 0x09, 0xb5, 0xa1,
	0x03, 0x24, 0x8a, 0x58, 0x9f, 0x50, 0x89, 0x62, 0x11, 0x1d, 0x9a, 0x13, 0x18, 0x78, 0x3f, 0xc1,
	0xb8, 0x7c, 0xba, 0x8c, 0x7c, 0x0f, 0xae, 0x79, 0x6e, 0xea, 0xcc, 0xbb, 0xa7, 0xe3, 0xb3, 0x93,
	0x45, 0xa5, 0x1a, 0x8b, 0x72, 0xaf, 0x6f, 0x37, 0x7a, 0x6f, 0xa1, 0x77, 0xbe, 0x56, 0x9f, 0xc8,
	0x21, 0x74, 0x1f, 0xf9, 0xc6, 0x16, 0x45, 0x2f, 0xc9, 0x0c, 0xdc, 0x8c, 0x07, 0x29, 0x57, 0xb6,
	0x22, 0x36, 0xf2, 0x1e, 0x61, 0x72, 0xc1, 0x64, 0x18, 0xf1, 0x25, 0x53, 0xcc, 0xe7, 0x9f, 0xd7,
	0x3c, 0xdb, 0x5f, 0xd8, 0x6a, 0x99, 0x3a, 0x3b, 0x65, 0x7a, 0x0d, 0x43, 0x21, 0x15, 0x4f, 0x9f,
	0x58, 0x64, 0x6b, 0x5c, 0xc4, 0xde, 0xef, 0x1d, 0x38, 0xf8, 0xc0, 0x7f, 0xbd, 0x89, 0x98, 0x34,
	0x1d, 0x9c, 0x81, 0xbb, 0x8a, 0x98, 0xbc, 0x5a, 0xda, 0x6b, 0x6c, 0xa4, 0xf1, 0x75, 0xc6, 0xd3,
	0xab, 0x65, 0xfe, 0x6f, 0x4d, 0x44, 0xce, 0xe0, 0x88, 0x05, 0x4a, 0x3c, 0xf1, 0x8b, 0x75, 0x9a,
	0x72, 0x19, 0x6c, 0x6e, 0x37, 0xf1, 0x7d, 0x92, 0x5f, 0xd4, 0xc8, 0x91, 0xf7, 0x70, 0x5c, 0xc5,
	0x7f, 0x66, 0x11, 0x93, 0x81, 0x69, 0xb1, 0xe3, 0x37, 0x93, 0xe4, 0x14, 0x5e, 0x05, 0x51, 0x92,
	0xf1, 0x6b, 0x79, 0x91, 0xc4, 0xab, 0x88, 0x2b, 0x8e, 0x5d, 0x1f, 0xfa, 0x75, 0x98, 0xbc, 0x01,
	0x37, 0x49, 0x43, 0x9e, 0x66, 0xd4, 0xc5, 0x36, 0x1d, 0xd5, 0xda, 0x74, 0xad, 0x49, 0xdf, 0xee,
	0xf1, 0xfe, 0xea, 0x40, 0x1f, 0x11, 0xad, 0x27, 0xc4, 0x8a, 0x87, 0xcf, 0xc3, 0x4a, 0xf9, 0x3b,
	0x7b, 0xcb, 0xdf, 0xdd, 0x29, 0x3f, 0x81, 0x5e, 0x26, 0xc2, 0x5c, 0xbf, 0xb8, 0xd6, 0x67, 0x22,
	0x11, 0x0b, 0x75, 0x83, 0x22, 0xed, 0xe3, 0x63, 0x6f, 0x21, 0xda, 0x17, 0x78, 0xf5, 0x9d, 0x36,
	0x8c, 0x6b, 0x7c, 0x51, 0x00, 0x64, 0x0e, 0x63, 0x0c, 0x6e, 0x15, 0x53, 0xeb, 0x0c, 0xa5, 0x3c,
	0xf2, 0xb7, 0x21, 0x7d, 0x9e, 0x05, 0x41, 0xb2, 0x96, 0xea, 0x6a, 0x89, 0x8a, 0x1e, 0xf9, 0x25,
	0xa0, 0xd9, 0x47, 0xbe, 0xb9, 0x59, 0xdf, 0x47, 0x22, 0xa0, 0x23, 0xc3, 0x16, 0x80, 0x65, 0x6f,
	0x8d, 0x34, 0xa1, 0x60, 0x0d, 0x40, 0xce, 0x60, 0xa8, 0x52, 0xf1, 0xf0, 0xa0, 0xab, 0x3b, 0xc6,
	0xea, 0xce, 0x76, 0x4c, 0x80, 0xb4, 0x5f, 0xec, 0xf3, 0xfe, 0x70, 0x60, 0x60, 0x51, 0x9d, 0xdd,
	0xe2, 0x45, 0x95, 0x4b, 0x60, 0xbb, 0x03, 0x9d, 0x6a, 0x07, 0x08, 0xf4, 0xa4, 0xae, 0xaf, 0x31,
	0x3a, 0xae, 0x35, 0x16, 0x24, 0x61, 0x5e, 0x20, 0x5c, 0x6f, 0xe5, 0xe7, 0x21, 0x56, 0x66, 0xe8,
	0x97, 0x80, 0xce, 0xaf, 0xc5, 0x95, 0xc8, 0x8c, 0x0e, 0xe7, 0x5d, 0x9d, 0xdf, 0x86, 0xde, 0xdf,
	0x3d, 0x98, 0xde, 0xe5, 0xfb, 0x50, 0x0e, 0x5f, 0x9e, 0x68, 0xed, 0xff, 0xb6, 0x74, 0x51, 0xb7,
	0xc5, 0x45, 0xbd, 0x8a, 0x8b, 0x2a, 0xfd, 0xea, 0xd7, 0xfb, 0xd5, 0xe6, 0x31, 0xf7, 0xff, 0x78,
	0x6c, 0xb0, 0xcf, 0x63, 0xdb, 0x13, 0x72, 0x58, 0x9b, 0x90, 0xcf, 0x51, 0x4d, 0xd5, 0x23, 0xe3,
	0x56, 0x8f, 0x1c, 0x6c, 0x79, 0xa4, 0xe2, 0x81, 0x17, 0x75, 0x0f, 0x54, 0x1d, 0xf4, 0xb2, 0xc9,
	0x41, 0xa5, 0xce, 0x5e, 0xd5, 0x75, 0xf6, 0x2d, 0xbc, 0x2c, 0x44, 0x61, 0x32, 0x1c, 0x62, 0x86,
	0x1a, 0x4a, 0x16, 0x40, 0x0a, 0xe4, 0x22, 0x91, 0xa1, 0xd0, 0x62, 0xa1, 0x13, 0x4c, 0xd7, 0xc0,
	0x34, 0xcd, 0x28, 0xd2, 0x38, 0xa3, 0xbc, 0x7f, 0x06, 0x30, 0xcd, 0x83, 0x6d, 0xbd, 0xb5, 0xcf,
	0xa0, 0x52, 0x53, 0x9d, 0x16, 0x4d, 0x75, 0x2b, 0x9a, 0xda, 0x56, 0x6e, 0x6f, 0xef, 0xcc, 0xea,
	0xb7, 0xf6, 0xc3, 0xad, 0xf6, 0xa3, 0xd4, 0xe8, 0xa0, 0xae, 0xd1, 0xf7, 0x70, 0x2c, 0xa4, 0x50,
	0x82, 0x45, 0x35, 0x91, 0x9a, 0xe9, 0xd3, 0x4c, 0x92, 0x1f, 0x60, 0x56, 0x23, 0x72, 0x99, 0x8e,
	0xb0, 0x1f, 0x2d, 0x6c, 0xc3, 0x6d, 0xf8, 0xea, 0x0d, 0x51, 0x79, 0x8e, 0xdf, 0x4c, 0x92, 0x1f,
	0x81, 0xd6, 0x08, 0x9f, 0xc7, 0x4c, 0xc8, 0x90, 0xa7, 0xa8, 0x49, 0xc7, 0x6f, 0xe5, 0xb5, 0x07,
	0x6b, 0x9c, 0xd1, 0xcd, 0x01, 0x9e, 0x6b, 0xe4, 0xc8, 0x5b, 0x98, 0x7e, 0x14, 0x72, 0xa7, 0x22,
	0x46, 0xcb, 0x4d, 0x94, 0xbe, 0xa5, 0x02, 0xe7, 0xd5, 0x30, 0xfa, 0x6e, 0xe4, 0xc8, 0x1b, 0x98,
	0x7c, 0xe4, 0xf5, 0xd1, 0x60, 0x14, 0xbf, 0x4b, 0xd4, 0x76, 0x9f, 0xc7, 0xba, 0x7d, 0x56, 0xfc,
	0xbb, 0x44, 0xd5, 0x45, 0x93, 0x2f, 0xbb, 0x88, 0xfc, 0x07, 0x17, 0x4d, 0xf7, 0xb9, 0x28, 0x57,
	0xea, 0xb5, 0xf5, 0xc2, 0x11, 0x6e, 0xae, 0xc3, 0x3a, 0x73, 0x0e, 0xfd, 0x52, 0xea, 0xf9, 0xd8,
	0x64, 0xde, 0x65, 0x88, 0x07, 0x07, 0x39, 0x8a, 0x9f, 0x9c, 0x33, 0xdc, 0x59, 0xc1, 0xf0, 0xbb,
	0xcc, 0xbc, 0x58, 0xbf, 0xb2, 0xdf, 0x65, 0xe6, 0x9d, 0x4a, 0x61, 0x10, 0x72, 0xc5, 0x44, 0x94,
	0x51, 0x6a, 0x9c, 0x69, 0xc3, 0x26, 0xd7, 0x9f, 0x34, 0xbb, 0xfe, 0x12, 0x26, 0xe7, 0xf7, 0x49,
	0xfa, 0x4c, 0xcb, 0x7b, 0x0b, 0x38, 0xbc, 0x94, 0x0f, 0x42, 0xf2, 0x5b, 0xc5, 0x52, 0x55, 0xbe,
	0xa8, 0x10, 0x2b, 0xd2, 0x14, 0xb1, 0xf7, 0x0e, 0xa6, 0x4b, 0x8e, 0x93, 0xe6, 0xdc, 0x18, 0xd6,
	0x1c, 0xa9, 0x38, 0xda, 0xa9, 0x39, 0xfa, 0xde, 0xc5, 0x0f, 0xfe, 0x77, 0xff, 0x06, 0x00, 0x00,
	0xff, 0xff, 0xe3, 0x2d, 0x21, 0x37, 0x0d, 0x0c, 0x00, 0x00,
}
