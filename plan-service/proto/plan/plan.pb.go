// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/plan/plan.proto

/*
Package fomo_plans is a generated protocol buffer package.

It is generated from these files:
	proto/plan/plan.proto

It has these top-level messages:
	PlanRequest
	OrderRequest
	TriggerRequest
	GetUserPlanRequest
	GetUserPlansRequest
	DeletePlanRequest
	UpdatePlanRequest
	Plan
	PlanWithPagedOrders
	Order
	Trigger
	PlanData
	PlanResponse
	PlanWithPagedOrdersResponse
	PlansPageResponse
	OrdersPage
	PlansPage
*/
package fomo_plans

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "golang.org/x/net/context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Requests
type PlanRequest struct {
	PlanID          string          `protobuf:"bytes,1,opt,name=planID" json:"planID"`
	UserID          string          `protobuf:"bytes,2,opt,name=userID" json:"userID"`
	KeyID           string          `protobuf:"bytes,3,opt,name=keyID" json:"keyID"`
	PlanTemplateID  string          `protobuf:"bytes,4,opt,name=planTemplateID" json:"planTemplateID"`
	Exchange        string          `protobuf:"bytes,5,opt,name=exchange" json:"exchange"`
	MarketName      string          `protobuf:"bytes,6,opt,name=marketName" json:"marketName"`
	BaseBalance     float64         `protobuf:"fixed64,7,opt,name=baseBalance" json:"baseBalance"`
	CurrencyBalance float64         `protobuf:"fixed64,8,opt,name=currencyBalance" json:"currencyBalance"`
	Status          string          `protobuf:"bytes,9,opt,name=status" json:"status"`
	Orders          []*OrderRequest `protobuf:"bytes,10,rep,name=orders" json:"orders"`
}

func (m *PlanRequest) Reset()                    { *m = PlanRequest{} }
func (m *PlanRequest) String() string            { return proto.CompactTextString(m) }
func (*PlanRequest) ProtoMessage()               {}
func (*PlanRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *PlanRequest) GetPlanID() string {
	if m != nil {
		return m.PlanID
	}
	return ""
}

func (m *PlanRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *PlanRequest) GetKeyID() string {
	if m != nil {
		return m.KeyID
	}
	return ""
}

func (m *PlanRequest) GetPlanTemplateID() string {
	if m != nil {
		return m.PlanTemplateID
	}
	return ""
}

func (m *PlanRequest) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *PlanRequest) GetMarketName() string {
	if m != nil {
		return m.MarketName
	}
	return ""
}

func (m *PlanRequest) GetBaseBalance() float64 {
	if m != nil {
		return m.BaseBalance
	}
	return 0
}

func (m *PlanRequest) GetCurrencyBalance() float64 {
	if m != nil {
		return m.CurrencyBalance
	}
	return 0
}

func (m *PlanRequest) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *PlanRequest) GetOrders() []*OrderRequest {
	if m != nil {
		return m.Orders
	}
	return nil
}

type OrderRequest struct {
	Side            string            `protobuf:"bytes,1,opt,name=side" json:"side"`
	OrderType       string            `protobuf:"bytes,2,opt,name=orderType" json:"orderType"`
	OrderTemplateID string            `protobuf:"bytes,3,opt,name=orderTemplateID" json:"orderTemplateID"`
	BasePercent     float64           `protobuf:"fixed64,4,opt,name=basePercent" json:"basePercent"`
	CurrencyPercent float64           `protobuf:"fixed64,5,opt,name=currencyPercent" json:"currencyPercent"`
	LimitPrice      float64           `protobuf:"fixed64,6,opt,name=limitPrice" json:"limitPrice"`
	Triggers        []*TriggerRequest `protobuf:"bytes,7,rep,name=triggers" json:"triggers"`
}

func (m *OrderRequest) Reset()                    { *m = OrderRequest{} }
func (m *OrderRequest) String() string            { return proto.CompactTextString(m) }
func (*OrderRequest) ProtoMessage()               {}
func (*OrderRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *OrderRequest) GetSide() string {
	if m != nil {
		return m.Side
	}
	return ""
}

func (m *OrderRequest) GetOrderType() string {
	if m != nil {
		return m.OrderType
	}
	return ""
}

func (m *OrderRequest) GetOrderTemplateID() string {
	if m != nil {
		return m.OrderTemplateID
	}
	return ""
}

func (m *OrderRequest) GetBasePercent() float64 {
	if m != nil {
		return m.BasePercent
	}
	return 0
}

func (m *OrderRequest) GetCurrencyPercent() float64 {
	if m != nil {
		return m.CurrencyPercent
	}
	return 0
}

func (m *OrderRequest) GetLimitPrice() float64 {
	if m != nil {
		return m.LimitPrice
	}
	return 0
}

func (m *OrderRequest) GetTriggers() []*TriggerRequest {
	if m != nil {
		return m.Triggers
	}
	return nil
}

type TriggerRequest struct {
	Name    string   `protobuf:"bytes,1,opt,name=name" json:"name"`
	Code    string   `protobuf:"bytes,2,opt,name=code" json:"code"`
	Actions []string `protobuf:"bytes,3,rep,name=actions" json:"actions"`
}

func (m *TriggerRequest) Reset()                    { *m = TriggerRequest{} }
func (m *TriggerRequest) String() string            { return proto.CompactTextString(m) }
func (*TriggerRequest) ProtoMessage()               {}
func (*TriggerRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *TriggerRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *TriggerRequest) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *TriggerRequest) GetActions() []string {
	if m != nil {
		return m.Actions
	}
	return nil
}

type GetUserPlanRequest struct {
	PlanID   string `protobuf:"bytes,1,opt,name=planID" json:"planID"`
	UserID   string `protobuf:"bytes,2,opt,name=userID" json:"userID"`
	Page     uint32 `protobuf:"varint,3,opt,name=page" json:"page"`
	PageSize uint32 `protobuf:"varint,4,opt,name=pageSize" json:"pageSize"`
}

func (m *GetUserPlanRequest) Reset()                    { *m = GetUserPlanRequest{} }
func (m *GetUserPlanRequest) String() string            { return proto.CompactTextString(m) }
func (*GetUserPlanRequest) ProtoMessage()               {}
func (*GetUserPlanRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *GetUserPlanRequest) GetPlanID() string {
	if m != nil {
		return m.PlanID
	}
	return ""
}

func (m *GetUserPlanRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *GetUserPlanRequest) GetPage() uint32 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *GetUserPlanRequest) GetPageSize() uint32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

type GetUserPlansRequest struct {
	UserID     string `protobuf:"bytes,1,opt,name=userID" json:"userID"`
	Exchange   string `protobuf:"bytes,2,opt,name=exchange" json:"exchange"`
	MarketName string `protobuf:"bytes,3,opt,name=marketName" json:"marketName"`
	Status     string `protobuf:"bytes,4,opt,name=status" json:"status"`
	Page       uint32 `protobuf:"varint,5,opt,name=page" json:"page"`
	PageSize   uint32 `protobuf:"varint,6,opt,name=pageSize" json:"pageSize"`
}

func (m *GetUserPlansRequest) Reset()                    { *m = GetUserPlansRequest{} }
func (m *GetUserPlansRequest) String() string            { return proto.CompactTextString(m) }
func (*GetUserPlansRequest) ProtoMessage()               {}
func (*GetUserPlansRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *GetUserPlansRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *GetUserPlansRequest) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *GetUserPlansRequest) GetMarketName() string {
	if m != nil {
		return m.MarketName
	}
	return ""
}

func (m *GetUserPlansRequest) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *GetUserPlansRequest) GetPage() uint32 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *GetUserPlansRequest) GetPageSize() uint32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

type DeletePlanRequest struct {
	PlanID string `protobuf:"bytes,1,opt,name=planID" json:"planID"`
	UserID string `protobuf:"bytes,2,opt,name=userID" json:"userID"`
}

func (m *DeletePlanRequest) Reset()                    { *m = DeletePlanRequest{} }
func (m *DeletePlanRequest) String() string            { return proto.CompactTextString(m) }
func (*DeletePlanRequest) ProtoMessage()               {}
func (*DeletePlanRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *DeletePlanRequest) GetPlanID() string {
	if m != nil {
		return m.PlanID
	}
	return ""
}

func (m *DeletePlanRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

type UpdatePlanRequest struct {
	PlanID          string  `protobuf:"bytes,1,opt,name=planID" json:"planID"`
	UserID          string  `protobuf:"bytes,2,opt,name=userID" json:"userID"`
	Status          string  `protobuf:"bytes,3,opt,name=status" json:"status"`
	BaseBalance     float64 `protobuf:"fixed64,4,opt,name=baseBalance" json:"baseBalance"`
	CurrencyBalance float64 `protobuf:"fixed64,5,opt,name=currencyBalance" json:"currencyBalance"`
}

func (m *UpdatePlanRequest) Reset()                    { *m = UpdatePlanRequest{} }
func (m *UpdatePlanRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdatePlanRequest) ProtoMessage()               {}
func (*UpdatePlanRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *UpdatePlanRequest) GetPlanID() string {
	if m != nil {
		return m.PlanID
	}
	return ""
}

func (m *UpdatePlanRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *UpdatePlanRequest) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *UpdatePlanRequest) GetBaseBalance() float64 {
	if m != nil {
		return m.BaseBalance
	}
	return 0
}

func (m *UpdatePlanRequest) GetCurrencyBalance() float64 {
	if m != nil {
		return m.CurrencyBalance
	}
	return 0
}

// Responses
type Plan struct {
	PlanID             string   `protobuf:"bytes,1,opt,name=planID" json:"planID"`
	PlanTemplateID     string   `protobuf:"bytes,2,opt,name=planTemplateID" json:"planTemplateID"`
	UserID             string   `protobuf:"bytes,3,opt,name=userID" json:"userID"`
	KeyID              string   `protobuf:"bytes,4,opt,name=keyID" json:"keyID"`
	Key                string   `protobuf:"bytes,5,opt,name=key" json:"key"`
	KeySecret          string   `protobuf:"bytes,6,opt,name=keySecret" json:"keySecret"`
	KeyDescription     string   `protobuf:"bytes,7,opt,name=keyDescription" json:"keyDescription"`
	ActiveOrderNumber  uint32   `protobuf:"varint,8,opt,name=activeOrderNumber" json:"activeOrderNumber"`
	Exchange           string   `protobuf:"bytes,9,opt,name=exchange" json:"exchange"`
	ExchangeMarketName string   `protobuf:"bytes,10,opt,name=exchangeMarketName" json:"exchangeMarketName"`
	MarketName         string   `protobuf:"bytes,11,opt,name=marketName" json:"marketName"`
	BaseBalance        float64  `protobuf:"fixed64,12,opt,name=baseBalance" json:"baseBalance"`
	CurrencyBalance    float64  `protobuf:"fixed64,13,opt,name=currencyBalance" json:"currencyBalance"`
	Status             string   `protobuf:"bytes,14,opt,name=status" json:"status"`
	Orders             []*Order `protobuf:"bytes,15,rep,name=orders" json:"orders"`
}

func (m *Plan) Reset()                    { *m = Plan{} }
func (m *Plan) String() string            { return proto.CompactTextString(m) }
func (*Plan) ProtoMessage()               {}
func (*Plan) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Plan) GetPlanID() string {
	if m != nil {
		return m.PlanID
	}
	return ""
}

func (m *Plan) GetPlanTemplateID() string {
	if m != nil {
		return m.PlanTemplateID
	}
	return ""
}

func (m *Plan) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *Plan) GetKeyID() string {
	if m != nil {
		return m.KeyID
	}
	return ""
}

func (m *Plan) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Plan) GetKeySecret() string {
	if m != nil {
		return m.KeySecret
	}
	return ""
}

func (m *Plan) GetKeyDescription() string {
	if m != nil {
		return m.KeyDescription
	}
	return ""
}

func (m *Plan) GetActiveOrderNumber() uint32 {
	if m != nil {
		return m.ActiveOrderNumber
	}
	return 0
}

func (m *Plan) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *Plan) GetExchangeMarketName() string {
	if m != nil {
		return m.ExchangeMarketName
	}
	return ""
}

func (m *Plan) GetMarketName() string {
	if m != nil {
		return m.MarketName
	}
	return ""
}

func (m *Plan) GetBaseBalance() float64 {
	if m != nil {
		return m.BaseBalance
	}
	return 0
}

func (m *Plan) GetCurrencyBalance() float64 {
	if m != nil {
		return m.CurrencyBalance
	}
	return 0
}

func (m *Plan) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Plan) GetOrders() []*Order {
	if m != nil {
		return m.Orders
	}
	return nil
}

type PlanWithPagedOrders struct {
	PlanID             string      `protobuf:"bytes,1,opt,name=planID" json:"planID"`
	PlanTemplateID     string      `protobuf:"bytes,2,opt,name=planTemplateID" json:"planTemplateID"`
	UserID             string      `protobuf:"bytes,3,opt,name=userID" json:"userID"`
	KeyID              string      `protobuf:"bytes,4,opt,name=keyID" json:"keyID"`
	Key                string      `protobuf:"bytes,5,opt,name=key" json:"key"`
	KeySecret          string      `protobuf:"bytes,6,opt,name=keySecret" json:"keySecret"`
	KeyDescription     string      `protobuf:"bytes,7,opt,name=keyDescription" json:"keyDescription"`
	Exchange           string      `protobuf:"bytes,8,opt,name=exchange" json:"exchange"`
	ExchangeMarketName string      `protobuf:"bytes,9,opt,name=exchangeMarketName" json:"exchangeMarketName"`
	MarketName         string      `protobuf:"bytes,10,opt,name=marketName" json:"marketName"`
	BaseBalance        float64     `protobuf:"fixed64,11,opt,name=baseBalance" json:"baseBalance"`
	CurrencyBalance    float64     `protobuf:"fixed64,12,opt,name=currencyBalance" json:"currencyBalance"`
	Status             string      `protobuf:"bytes,13,opt,name=status" json:"status"`
	OrdersPage         *OrdersPage `protobuf:"bytes,14,opt,name=ordersPage" json:"ordersPage"`
}

func (m *PlanWithPagedOrders) Reset()                    { *m = PlanWithPagedOrders{} }
func (m *PlanWithPagedOrders) String() string            { return proto.CompactTextString(m) }
func (*PlanWithPagedOrders) ProtoMessage()               {}
func (*PlanWithPagedOrders) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *PlanWithPagedOrders) GetPlanID() string {
	if m != nil {
		return m.PlanID
	}
	return ""
}

func (m *PlanWithPagedOrders) GetPlanTemplateID() string {
	if m != nil {
		return m.PlanTemplateID
	}
	return ""
}

func (m *PlanWithPagedOrders) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *PlanWithPagedOrders) GetKeyID() string {
	if m != nil {
		return m.KeyID
	}
	return ""
}

func (m *PlanWithPagedOrders) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *PlanWithPagedOrders) GetKeySecret() string {
	if m != nil {
		return m.KeySecret
	}
	return ""
}

func (m *PlanWithPagedOrders) GetKeyDescription() string {
	if m != nil {
		return m.KeyDescription
	}
	return ""
}

func (m *PlanWithPagedOrders) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *PlanWithPagedOrders) GetExchangeMarketName() string {
	if m != nil {
		return m.ExchangeMarketName
	}
	return ""
}

func (m *PlanWithPagedOrders) GetMarketName() string {
	if m != nil {
		return m.MarketName
	}
	return ""
}

func (m *PlanWithPagedOrders) GetBaseBalance() float64 {
	if m != nil {
		return m.BaseBalance
	}
	return 0
}

func (m *PlanWithPagedOrders) GetCurrencyBalance() float64 {
	if m != nil {
		return m.CurrencyBalance
	}
	return 0
}

func (m *PlanWithPagedOrders) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *PlanWithPagedOrders) GetOrdersPage() *OrdersPage {
	if m != nil {
		return m.OrdersPage
	}
	return nil
}

type Order struct {
	OrderID         string     `protobuf:"bytes,1,opt,name=orderID" json:"orderID"`
	Side            string     `protobuf:"bytes,2,opt,name=side" json:"side"`
	OrderTemplateID string     `protobuf:"bytes,3,opt,name=orderTemplateID" json:"orderTemplateID"`
	OrderNumber     uint32     `protobuf:"varint,4,opt,name=orderNumber" json:"orderNumber"`
	OrderType       string     `protobuf:"bytes,5,opt,name=orderType" json:"orderType"`
	LimitPrice      float64    `protobuf:"fixed64,6,opt,name=limitPrice" json:"limitPrice"`
	BasePercent     float64    `protobuf:"fixed64,7,opt,name=basePercent" json:"basePercent"`
	CurrencyPercent float64    `protobuf:"fixed64,8,opt,name=currencyPercent" json:"currencyPercent"`
	Status          string     `protobuf:"bytes,9,opt,name=status" json:"status"`
	NextOrderID     string     `protobuf:"bytes,10,opt,name=nextOrderID" json:"nextOrderID"`
	Triggers        []*Trigger `protobuf:"bytes,11,rep,name=triggers" json:"triggers"`
}

func (m *Order) Reset()                    { *m = Order{} }
func (m *Order) String() string            { return proto.CompactTextString(m) }
func (*Order) ProtoMessage()               {}
func (*Order) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *Order) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
}

func (m *Order) GetSide() string {
	if m != nil {
		return m.Side
	}
	return ""
}

func (m *Order) GetOrderTemplateID() string {
	if m != nil {
		return m.OrderTemplateID
	}
	return ""
}

func (m *Order) GetOrderNumber() uint32 {
	if m != nil {
		return m.OrderNumber
	}
	return 0
}

func (m *Order) GetOrderType() string {
	if m != nil {
		return m.OrderType
	}
	return ""
}

func (m *Order) GetLimitPrice() float64 {
	if m != nil {
		return m.LimitPrice
	}
	return 0
}

func (m *Order) GetBasePercent() float64 {
	if m != nil {
		return m.BasePercent
	}
	return 0
}

func (m *Order) GetCurrencyPercent() float64 {
	if m != nil {
		return m.CurrencyPercent
	}
	return 0
}

func (m *Order) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Order) GetNextOrderID() string {
	if m != nil {
		return m.NextOrderID
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
	TriggerID     string   `protobuf:"bytes,1,opt,name=triggerID" json:"triggerID"`
	OrderID       string   `protobuf:"bytes,2,opt,name=orderID" json:"orderID"`
	TriggerNumber uint32   `protobuf:"varint,3,opt,name=triggerNumber" json:"triggerNumber"`
	Name          string   `protobuf:"bytes,4,opt,name=name" json:"name"`
	Code          string   `protobuf:"bytes,5,opt,name=code" json:"code"`
	Triggered     bool     `protobuf:"varint,6,opt,name=triggered" json:"triggered"`
	Actions       []string `protobuf:"bytes,7,rep,name=actions" json:"actions"`
}

func (m *Trigger) Reset()                    { *m = Trigger{} }
func (m *Trigger) String() string            { return proto.CompactTextString(m) }
func (*Trigger) ProtoMessage()               {}
func (*Trigger) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

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

func (m *Trigger) GetTriggerNumber() uint32 {
	if m != nil {
		return m.TriggerNumber
	}
	return 0
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

type PlanData struct {
	Plan *Plan `protobuf:"bytes,1,opt,name=plan" json:"plan"`
}

func (m *PlanData) Reset()                    { *m = PlanData{} }
func (m *PlanData) String() string            { return proto.CompactTextString(m) }
func (*PlanData) ProtoMessage()               {}
func (*PlanData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *PlanData) GetPlan() *Plan {
	if m != nil {
		return m.Plan
	}
	return nil
}

type PlanResponse struct {
	Status  string    `protobuf:"bytes,1,opt,name=status" json:"status"`
	Message string    `protobuf:"bytes,2,opt,name=message" json:"message"`
	Data    *PlanData `protobuf:"bytes,3,opt,name=data" json:"data"`
}

func (m *PlanResponse) Reset()                    { *m = PlanResponse{} }
func (m *PlanResponse) String() string            { return proto.CompactTextString(m) }
func (*PlanResponse) ProtoMessage()               {}
func (*PlanResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *PlanResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *PlanResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *PlanResponse) GetData() *PlanData {
	if m != nil {
		return m.Data
	}
	return nil
}

type PlanWithPagedOrdersResponse struct {
	Status  string               `protobuf:"bytes,1,opt,name=status" json:"status"`
	Message string               `protobuf:"bytes,2,opt,name=message" json:"message"`
	Data    *PlanWithPagedOrders `protobuf:"bytes,3,opt,name=data" json:"data"`
}

func (m *PlanWithPagedOrdersResponse) Reset()                    { *m = PlanWithPagedOrdersResponse{} }
func (m *PlanWithPagedOrdersResponse) String() string            { return proto.CompactTextString(m) }
func (*PlanWithPagedOrdersResponse) ProtoMessage()               {}
func (*PlanWithPagedOrdersResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *PlanWithPagedOrdersResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *PlanWithPagedOrdersResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *PlanWithPagedOrdersResponse) GetData() *PlanWithPagedOrders {
	if m != nil {
		return m.Data
	}
	return nil
}

type PlansPageResponse struct {
	Status  string     `protobuf:"bytes,1,opt,name=status" json:"status"`
	Message string     `protobuf:"bytes,2,opt,name=message" json:"message"`
	Data    *PlansPage `protobuf:"bytes,3,opt,name=data" json:"data"`
}

func (m *PlansPageResponse) Reset()                    { *m = PlansPageResponse{} }
func (m *PlansPageResponse) String() string            { return proto.CompactTextString(m) }
func (*PlansPageResponse) ProtoMessage()               {}
func (*PlansPageResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *PlansPageResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *PlansPageResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *PlansPageResponse) GetData() *PlansPage {
	if m != nil {
		return m.Data
	}
	return nil
}

type OrdersPage struct {
	Page     uint32   `protobuf:"varint,1,opt,name=page" json:"page"`
	PageSize uint32   `protobuf:"varint,2,opt,name=pageSize" json:"pageSize"`
	Total    uint32   `protobuf:"varint,3,opt,name=total" json:"total"`
	Orders   []*Order `protobuf:"bytes,4,rep,name=orders" json:"orders"`
}

func (m *OrdersPage) Reset()                    { *m = OrdersPage{} }
func (m *OrdersPage) String() string            { return proto.CompactTextString(m) }
func (*OrdersPage) ProtoMessage()               {}
func (*OrdersPage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *OrdersPage) GetPage() uint32 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *OrdersPage) GetPageSize() uint32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *OrdersPage) GetTotal() uint32 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *OrdersPage) GetOrders() []*Order {
	if m != nil {
		return m.Orders
	}
	return nil
}

type PlansPage struct {
	Page     uint32  `protobuf:"varint,1,opt,name=page" json:"page"`
	PageSize uint32  `protobuf:"varint,2,opt,name=pageSize" json:"pageSize"`
	Total    uint32  `protobuf:"varint,3,opt,name=total" json:"total"`
	Plans    []*Plan `protobuf:"bytes,4,rep,name=plans" json:"plans"`
}

func (m *PlansPage) Reset()                    { *m = PlansPage{} }
func (m *PlansPage) String() string            { return proto.CompactTextString(m) }
func (*PlansPage) ProtoMessage()               {}
func (*PlansPage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

func (m *PlansPage) GetPage() uint32 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *PlansPage) GetPageSize() uint32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *PlansPage) GetTotal() uint32 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *PlansPage) GetPlans() []*Plan {
	if m != nil {
		return m.Plans
	}
	return nil
}

func init() {
	proto.RegisterType((*PlanRequest)(nil), "fomo.plans.PlanRequest")
	proto.RegisterType((*OrderRequest)(nil), "fomo.plans.OrderRequest")
	proto.RegisterType((*TriggerRequest)(nil), "fomo.plans.TriggerRequest")
	proto.RegisterType((*GetUserPlanRequest)(nil), "fomo.plans.GetUserPlanRequest")
	proto.RegisterType((*GetUserPlansRequest)(nil), "fomo.plans.GetUserPlansRequest")
	proto.RegisterType((*DeletePlanRequest)(nil), "fomo.plans.DeletePlanRequest")
	proto.RegisterType((*UpdatePlanRequest)(nil), "fomo.plans.UpdatePlanRequest")
	proto.RegisterType((*Plan)(nil), "fomo.plans.Plan")
	proto.RegisterType((*PlanWithPagedOrders)(nil), "fomo.plans.PlanWithPagedOrders")
	proto.RegisterType((*Order)(nil), "fomo.plans.Order")
	proto.RegisterType((*Trigger)(nil), "fomo.plans.Trigger")
	proto.RegisterType((*PlanData)(nil), "fomo.plans.PlanData")
	proto.RegisterType((*PlanResponse)(nil), "fomo.plans.PlanResponse")
	proto.RegisterType((*PlanWithPagedOrdersResponse)(nil), "fomo.plans.PlanWithPagedOrdersResponse")
	proto.RegisterType((*PlansPageResponse)(nil), "fomo.plans.PlansPageResponse")
	proto.RegisterType((*OrdersPage)(nil), "fomo.plans.OrdersPage")
	proto.RegisterType((*PlansPage)(nil), "fomo.plans.PlansPage")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for PlanService service

type PlanServiceClient interface {
	AddPlan(ctx context.Context, in *PlanRequest, opts ...client.CallOption) (*PlanResponse, error)
	GetUserPlan(ctx context.Context, in *GetUserPlanRequest, opts ...client.CallOption) (*PlanWithPagedOrdersResponse, error)
	GetUserPlans(ctx context.Context, in *GetUserPlansRequest, opts ...client.CallOption) (*PlansPageResponse, error)
	DeletePlan(ctx context.Context, in *DeletePlanRequest, opts ...client.CallOption) (*PlanResponse, error)
	UpdatePlan(ctx context.Context, in *UpdatePlanRequest, opts ...client.CallOption) (*PlanResponse, error)
}

type planServiceClient struct {
	c           client.Client
	serviceName string
}

func NewPlanServiceClient(serviceName string, c client.Client) PlanServiceClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "fomo.plans"
	}
	return &planServiceClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *planServiceClient) AddPlan(ctx context.Context, in *PlanRequest, opts ...client.CallOption) (*PlanResponse, error) {
	req := c.c.NewRequest(c.serviceName, "PlanService.AddPlan", in)
	out := new(PlanResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *planServiceClient) GetUserPlan(ctx context.Context, in *GetUserPlanRequest, opts ...client.CallOption) (*PlanWithPagedOrdersResponse, error) {
	req := c.c.NewRequest(c.serviceName, "PlanService.GetUserPlan", in)
	out := new(PlanWithPagedOrdersResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *planServiceClient) GetUserPlans(ctx context.Context, in *GetUserPlansRequest, opts ...client.CallOption) (*PlansPageResponse, error) {
	req := c.c.NewRequest(c.serviceName, "PlanService.GetUserPlans", in)
	out := new(PlansPageResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *planServiceClient) DeletePlan(ctx context.Context, in *DeletePlanRequest, opts ...client.CallOption) (*PlanResponse, error) {
	req := c.c.NewRequest(c.serviceName, "PlanService.DeletePlan", in)
	out := new(PlanResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *planServiceClient) UpdatePlan(ctx context.Context, in *UpdatePlanRequest, opts ...client.CallOption) (*PlanResponse, error) {
	req := c.c.NewRequest(c.serviceName, "PlanService.UpdatePlan", in)
	out := new(PlanResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for PlanService service

type PlanServiceHandler interface {
	AddPlan(context.Context, *PlanRequest, *PlanResponse) error
	GetUserPlan(context.Context, *GetUserPlanRequest, *PlanWithPagedOrdersResponse) error
	GetUserPlans(context.Context, *GetUserPlansRequest, *PlansPageResponse) error
	DeletePlan(context.Context, *DeletePlanRequest, *PlanResponse) error
	UpdatePlan(context.Context, *UpdatePlanRequest, *PlanResponse) error
}

func RegisterPlanServiceHandler(s server.Server, hdlr PlanServiceHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&PlanService{hdlr}, opts...))
}

type PlanService struct {
	PlanServiceHandler
}

func (h *PlanService) AddPlan(ctx context.Context, in *PlanRequest, out *PlanResponse) error {
	return h.PlanServiceHandler.AddPlan(ctx, in, out)
}

func (h *PlanService) GetUserPlan(ctx context.Context, in *GetUserPlanRequest, out *PlanWithPagedOrdersResponse) error {
	return h.PlanServiceHandler.GetUserPlan(ctx, in, out)
}

func (h *PlanService) GetUserPlans(ctx context.Context, in *GetUserPlansRequest, out *PlansPageResponse) error {
	return h.PlanServiceHandler.GetUserPlans(ctx, in, out)
}

func (h *PlanService) DeletePlan(ctx context.Context, in *DeletePlanRequest, out *PlanResponse) error {
	return h.PlanServiceHandler.DeletePlan(ctx, in, out)
}

func (h *PlanService) UpdatePlan(ctx context.Context, in *UpdatePlanRequest, out *PlanResponse) error {
	return h.PlanServiceHandler.UpdatePlan(ctx, in, out)
}

func init() { proto.RegisterFile("proto/plan/plan.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 1082 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xdc, 0x57, 0x5f, 0x8f, 0xdb, 0x44,
	0x10, 0xc7, 0xb1, 0x73, 0x49, 0xc6, 0x97, 0x6b, 0xb3, 0xd7, 0x16, 0x2b, 0xc0, 0x35, 0xb2, 0xaa,
	0x92, 0x4a, 0x28, 0xad, 0x52, 0xa9, 0x8f, 0x48, 0x40, 0xa4, 0xd3, 0x3d, 0xd0, 0x46, 0xbe, 0xab,
	0xe0, 0x75, 0xcf, 0x19, 0xd2, 0x90, 0xc4, 0x0e, 0xf6, 0xa6, 0x6a, 0x90, 0x78, 0xe3, 0x0d, 0x89,
	0xaf, 0xc1, 0x33, 0x9f, 0x82, 0x27, 0x3e, 0x05, 0x9f, 0x01, 0x9e, 0xd1, 0xce, 0xfa, 0xcf, 0xda,
	0x4e, 0x2e, 0x81, 0xe3, 0x89, 0x97, 0xbb, 0x9d, 0xdf, 0x4c, 0x76, 0x76, 0x7f, 0x33, 0x3b, 0x33,
	0x86, 0xfb, 0xab, 0x28, 0x14, 0xe1, 0xd3, 0xd5, 0x82, 0x07, 0xf4, 0x67, 0x40, 0x32, 0x83, 0x6f,
	0xc2, 0x65, 0x38, 0x90, 0x40, 0xec, 0xfe, 0x5e, 0x03, 0x7b, 0xbc, 0xe0, 0x81, 0x87, 0xdf, 0xad,
	0x31, 0x16, 0xec, 0x01, 0x1c, 0x49, 0xc5, 0xc5, 0xc8, 0x31, 0x7a, 0x46, 0xbf, 0xe5, 0x25, 0x92,
	0xc4, 0xd7, 0x31, 0x46, 0x17, 0x23, 0xa7, 0xa6, 0x70, 0x25, 0xb1, 0x7b, 0x50, 0x9f, 0xe3, 0xe6,
	0x62, 0xe4, 0x98, 0x04, 0x2b, 0x81, 0x3d, 0x86, 0x13, 0xf9, 0xbb, 0x2b, 0x5c, 0xae, 0x16, 0x5c,
	0xe0, 0xc5, 0xc8, 0xb1, 0x48, 0x5d, 0x42, 0x59, 0x17, 0x9a, 0xf8, 0xce, 0x7f, 0xc3, 0x83, 0x29,
	0x3a, 0x75, 0xb2, 0xc8, 0x64, 0x76, 0x06, 0xb0, 0xe4, 0xd1, 0x1c, 0xc5, 0x4b, 0xbe, 0x44, 0xe7,
	0x88, 0xb4, 0x1a, 0xc2, 0x7a, 0x60, 0x5f, 0xf3, 0x18, 0x3f, 0xe7, 0x0b, 0x1e, 0xf8, 0xe8, 0x34,
	0x7a, 0x46, 0xdf, 0xf0, 0x74, 0x88, 0xf5, 0xe1, 0x8e, 0xbf, 0x8e, 0x22, 0x0c, 0xfc, 0x4d, 0x6a,
	0xd5, 0x24, 0xab, 0x32, 0x2c, 0x6f, 0x17, 0x0b, 0x2e, 0xd6, 0xb1, 0xd3, 0x52, 0xb7, 0x53, 0x12,
	0x7b, 0x06, 0x47, 0x61, 0x34, 0xc1, 0x28, 0x76, 0xa0, 0x67, 0xf6, 0xed, 0xa1, 0x33, 0xc8, 0xa9,
	0x1b, 0xbc, 0x92, 0x9a, 0x84, 0x37, 0x2f, 0xb1, 0x73, 0x7f, 0xae, 0xc1, 0xb1, 0xae, 0x60, 0x0c,
	0xac, 0x78, 0x36, 0xc1, 0x84, 0x4e, 0x5a, 0xb3, 0x0f, 0xa1, 0x45, 0xe6, 0x57, 0x9b, 0x15, 0x26,
	0x7c, 0xe6, 0x80, 0x3c, 0xb6, 0x12, 0x72, 0xf6, 0x14, 0xb9, 0x65, 0x38, 0xa5, 0x60, 0x8c, 0x91,
	0x8f, 0x81, 0x20, 0x8e, 0x13, 0x0a, 0x12, 0x48, 0xa7, 0x20, 0xb5, 0xaa, 0x17, 0x29, 0x48, 0x2d,
	0xcf, 0x00, 0x16, 0xb3, 0xe5, 0x4c, 0x8c, 0xa3, 0x99, 0xaf, 0xe8, 0x36, 0x3c, 0x0d, 0x61, 0x2f,
	0xa0, 0x29, 0xa2, 0xd9, 0x74, 0x2a, 0xc9, 0x68, 0x10, 0x19, 0x5d, 0x9d, 0x8c, 0x2b, 0xa5, 0x4b,
	0xe9, 0xc8, 0x6c, 0x5d, 0x0f, 0x4e, 0x8a, 0x3a, 0xc9, 0x48, 0x20, 0x43, 0x9a, 0x30, 0x22, 0xd7,
	0x12, 0xf3, 0xc3, 0x49, 0x4a, 0x06, 0xad, 0x99, 0x03, 0x0d, 0xee, 0x8b, 0x59, 0x18, 0xc4, 0x8e,
	0xd9, 0x33, 0xfb, 0x2d, 0x2f, 0x15, 0x5d, 0x01, 0xec, 0x1c, 0xc5, 0xeb, 0x18, 0xa3, 0xdb, 0xa4,
	0x2e, 0x03, 0x6b, 0xc5, 0xa7, 0x48, 0xe4, 0xb6, 0x3d, 0x5a, 0xcb, 0x84, 0x94, 0xff, 0x2f, 0x67,
	0xdf, 0x23, 0xd1, 0xd9, 0xf6, 0x32, 0xd9, 0xfd, 0xd5, 0x80, 0x53, 0xcd, 0x6d, 0xac, 0xf9, 0x4d,
	0xf6, 0x37, 0x0a, 0xfb, 0xeb, 0xc9, 0x5d, 0xbb, 0x31, 0xb9, 0xcd, 0x4a, 0x72, 0xe7, 0x09, 0x69,
	0x15, 0x12, 0x32, 0x3d, 0x73, 0x7d, 0xc7, 0x99, 0x8f, 0x4a, 0x67, 0xfe, 0x02, 0x3a, 0x23, 0x5c,
	0xa0, 0xc0, 0x5b, 0x10, 0xe5, 0xfe, 0x62, 0x40, 0xe7, 0xf5, 0x6a, 0xc2, 0x6f, 0xb5, 0x8b, 0x76,
	0x25, 0xb3, 0x70, 0xa5, 0xd2, 0x3b, 0xb6, 0x0e, 0x7a, 0xc7, 0xf5, 0xad, 0xef, 0xd8, 0xfd, 0xd3,
	0x04, 0x4b, 0x9e, 0x71, 0xe7, 0xe1, 0xaa, 0x85, 0xa9, 0xb6, 0xb5, 0x30, 0xe5, 0x97, 0x30, 0xb7,
	0x97, 0x3b, 0x4b, 0x2f, 0x77, 0x77, 0xc1, 0x9c, 0xe3, 0x26, 0xa9, 0x60, 0x72, 0x29, 0x5f, 0xf8,
	0x1c, 0x37, 0x97, 0xe8, 0x47, 0x28, 0x92, 0xda, 0x95, 0x03, 0xf2, 0x14, 0x73, 0xdc, 0x8c, 0x30,
	0xf6, 0xa3, 0xd9, 0x4a, 0xa6, 0x34, 0x55, 0xaf, 0x96, 0x57, 0x42, 0xd9, 0x27, 0xd0, 0x91, 0x29,
	0xff, 0x16, 0xa9, 0xa2, 0xbc, 0x5c, 0x2f, 0xaf, 0x31, 0xa2, 0x12, 0xd6, 0xf6, 0xaa, 0x8a, 0x42,
	0xbe, 0xb5, 0x4a, 0xf9, 0x36, 0x00, 0x96, 0xae, 0xbf, 0xcc, 0xf3, 0x0e, 0xc8, 0x6a, 0x8b, 0xa6,
	0x94, 0x9f, 0xf6, 0xbe, 0xe2, 0x7b, 0x7c, 0x50, 0xd0, 0xda, 0xfb, 0x8a, 0xef, 0x49, 0x21, 0x31,
	0x9e, 0x64, 0xc5, 0xf7, 0x0e, 0xd5, 0x9b, 0x4e, 0xb5, 0xf8, 0xa6, 0x55, 0xf7, 0x0f, 0x13, 0x4e,
	0x65, 0xdc, 0xbf, 0x9a, 0x89, 0x37, 0x63, 0x3e, 0xc5, 0x09, 0xa9, 0xe3, 0xff, 0x49, 0x1a, 0xe8,
	0x81, 0x6d, 0x1e, 0x14, 0xd8, 0xd6, 0x81, 0x81, 0x85, 0x7d, 0x81, 0xb5, 0x0f, 0x0a, 0xec, 0xf1,
	0xbe, 0xc0, 0xb6, 0x0b, 0x81, 0x7d, 0x01, 0xa0, 0xe2, 0x26, 0x43, 0x45, 0x41, 0xb7, 0x87, 0x0f,
	0x2a, 0xc1, 0x25, 0xad, 0xa7, 0x59, 0xba, 0x7f, 0xd5, 0xa0, 0x4e, 0x2a, 0xd9, 0x1a, 0x08, 0xcf,
	0x02, 0x9b, 0x8a, 0x59, 0xbb, 0xad, 0x69, 0xed, 0xf6, 0x1f, 0x35, 0xd4, 0x50, 0x7b, 0x6a, 0xaa,
	0x03, 0xe8, 0x50, 0xb1, 0x75, 0xd7, 0xcb, 0xad, 0x7b, 0x5f, 0x13, 0x2d, 0x35, 0xec, 0xc6, 0x41,
	0x0d, 0xbb, 0xb9, 0xbd, 0x61, 0xef, 0x9a, 0x59, 0x7a, 0x60, 0x07, 0xf8, 0x4e, 0xbc, 0x4a, 0xf8,
	0x51, 0x21, 0xd6, 0x21, 0xf6, 0x54, 0x6b, 0xe5, 0x36, 0x3d, 0xad, 0xd3, 0x6d, 0xad, 0x3c, 0xef,
	0xe1, 0xbf, 0x19, 0xd0, 0x48, 0x50, 0x49, 0x40, 0x82, 0x67, 0xe4, 0xe7, 0x80, 0x1e, 0x98, 0x5a,
	0x31, 0x30, 0x8f, 0xa0, 0x9d, 0x98, 0x25, 0xe4, 0xaa, 0xb6, 0x5b, 0x04, 0xb3, 0xd9, 0xc0, 0xda,
	0x32, 0x1b, 0xd4, 0xb5, 0xd9, 0x20, 0x3f, 0x05, 0x4e, 0x88, 0xe7, 0xa6, 0x97, 0x03, 0xfa, 0xe4,
	0xd0, 0x28, 0x4e, 0x0e, 0xcf, 0xa0, 0x29, 0xeb, 0xc4, 0x88, 0x0b, 0xce, 0x1e, 0x81, 0x25, 0x2f,
	0x4c, 0x97, 0xb0, 0x87, 0x77, 0x75, 0x0a, 0xa8, 0xcf, 0x91, 0xd6, 0xfd, 0x16, 0x8e, 0x55, 0xd7,
	0x8b, 0x57, 0x61, 0x10, 0xeb, 0x49, 0x6d, 0x14, 0x68, 0x77, 0xa0, 0xb1, 0xc4, 0x38, 0xe6, 0x59,
	0xb3, 0x4f, 0x45, 0xd6, 0x07, 0x6b, 0xc2, 0x05, 0xa7, 0x0b, 0xdb, 0xc3, 0x7b, 0x65, 0x3f, 0xf2,
	0x2c, 0x1e, 0x59, 0xb8, 0x3f, 0x1a, 0xf0, 0xc1, 0x96, 0x32, 0x76, 0x0b, 0xdf, 0xcf, 0x0b, 0xbe,
	0x1f, 0x96, 0x7d, 0x97, 0x1d, 0xa9, 0x63, 0xac, 0xa0, 0x43, 0x03, 0x0e, 0x3d, 0xc0, 0x7f, 0xef,
	0xfb, 0x49, 0xc1, 0xf7, 0xfd, 0xb2, 0x6f, 0xb5, 0xbd, 0xf2, 0xf8, 0x03, 0x40, 0xfe, 0xe6, 0xb3,
	0x21, 0xc7, 0xd8, 0x31, 0xe4, 0xd4, 0x8a, 0x43, 0x8e, 0xac, 0xc6, 0x22, 0x14, 0x7c, 0x91, 0xa4,
	0x94, 0x12, 0xb4, 0xf6, 0x61, 0xed, 0x6b, 0x1f, 0x1b, 0x68, 0x65, 0x27, 0xfa, 0x8f, 0xbc, 0x3f,
	0x86, 0x3a, 0x79, 0x4a, 0x9c, 0x57, 0xb3, 0x4b, 0xa9, 0x87, 0x3f, 0x99, 0xea, 0xfb, 0xeb, 0x12,
	0xa3, 0xb7, 0xb2, 0x42, 0x7c, 0x0a, 0x8d, 0xcf, 0x26, 0x13, 0x9a, 0x61, 0xde, 0xaf, 0xfc, 0x46,
	0x4d, 0x5e, 0x5d, 0xa7, 0xaa, 0x50, 0x41, 0x72, 0xdf, 0x63, 0x5f, 0x83, 0xad, 0xcd, 0xa8, 0xec,
	0x4c, 0x37, 0xad, 0xce, 0xcc, 0xdd, 0x8f, 0xf7, 0x65, 0x44, 0xbe, 0xf3, 0x18, 0x8e, 0xf5, 0xe9,
	0x97, 0x3d, 0xdc, 0xb1, 0x75, 0x3a, 0x17, 0x77, 0x3f, 0xda, 0x1e, 0xf1, 0x7c, 0xc7, 0x73, 0x80,
	0x7c, 0x38, 0x65, 0x05, 0xf3, 0xca, 0xd0, 0x7a, 0xe3, 0xa5, 0xcf, 0x01, 0xf2, 0xf9, 0xb4, 0xb8,
	0x51, 0x65, 0x6e, 0xbd, 0x69, 0xa3, 0xeb, 0x23, 0xfa, 0x40, 0x7e, 0xfe, 0x77, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x4b, 0x1e, 0x62, 0xf4, 0x39, 0x0f, 0x00, 0x00,
}
