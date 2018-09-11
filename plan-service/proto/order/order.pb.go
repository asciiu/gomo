// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/asciiu/gomo/plan-service/proto/order/order.proto

/*
Package order is a generated protocol buffer package.

It is generated from these files:
	github.com/asciiu/gomo/plan-service/proto/order/order.proto

It has these top-level messages:
	TriggerRequest
	NewOrderRequest
	GetUserOrderRequest
	GetUserOrdersRequest
	RemoveOrderRequest
	Order
	Trigger
	UserOrderData
	UserOrdersData
	OrderResponse
	OrderListResponse
	OrdersPage
*/
package order

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

// Requests
type TriggerRequest struct {
	TriggerID         string   `protobuf:"bytes,1,opt,name=triggerID" json:"triggerID"`
	Code              string   `protobuf:"bytes,2,opt,name=code" json:"code"`
	Index             uint32   `protobuf:"varint,3,opt,name=index" json:"index"`
	Name              string   `protobuf:"bytes,4,opt,name=name" json:"name"`
	Title             string   `protobuf:"bytes,5,opt,name=title" json:"title"`
	TriggerTemplateID string   `protobuf:"bytes,6,opt,name=triggerTemplateID" json:"triggerTemplateID"`
	Actions           []string `protobuf:"bytes,7,rep,name=actions" json:"actions"`
}

func (m *TriggerRequest) Reset()                    { *m = TriggerRequest{} }
func (m *TriggerRequest) String() string            { return proto.CompactTextString(m) }
func (*TriggerRequest) ProtoMessage()               {}
func (*TriggerRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *TriggerRequest) GetTriggerID() string {
	if m != nil {
		return m.TriggerID
	}
	return ""
}

func (m *TriggerRequest) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *TriggerRequest) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *TriggerRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *TriggerRequest) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *TriggerRequest) GetTriggerTemplateID() string {
	if m != nil {
		return m.TriggerTemplateID
	}
	return ""
}

func (m *TriggerRequest) GetActions() []string {
	if m != nil {
		return m.Actions
	}
	return nil
}

type NewOrderRequest struct {
	OrderID                string            `protobuf:"bytes,1,opt,name=orderID" json:"orderID"`
	OrderPriority          uint32            `protobuf:"varint,2,opt,name=orderPriority" json:"orderPriority"`
	OrderType              string            `protobuf:"bytes,3,opt,name=orderType" json:"orderType"`
	OrderTemplateID        string            `protobuf:"bytes,4,opt,name=orderTemplateID" json:"orderTemplateID"`
	AccountID              string            `protobuf:"bytes,5,opt,name=accountID" json:"accountID"`
	ParentOrderID          string            `protobuf:"bytes,6,opt,name=parentOrderID" json:"parentOrderID"`
	Grupo                  string            `protobuf:"bytes,7,opt,name=grupo" json:"grupo"`
	MarketName             string            `protobuf:"bytes,8,opt,name=marketName" json:"marketName"`
	Side                   string            `protobuf:"bytes,9,opt,name=side" json:"side"`
	LimitPrice             float64           `protobuf:"fixed64,10,opt,name=limitPrice" json:"limitPrice"`
	InitialCurrencyBalance float64           `protobuf:"fixed64,11,opt,name=initialCurrencyBalance" json:"initialCurrencyBalance"`
	Triggers               []*TriggerRequest `protobuf:"bytes,12,rep,name=triggers" json:"triggers"`
}

func (m *NewOrderRequest) Reset()                    { *m = NewOrderRequest{} }
func (m *NewOrderRequest) String() string            { return proto.CompactTextString(m) }
func (*NewOrderRequest) ProtoMessage()               {}
func (*NewOrderRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *NewOrderRequest) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
}

func (m *NewOrderRequest) GetOrderPriority() uint32 {
	if m != nil {
		return m.OrderPriority
	}
	return 0
}

func (m *NewOrderRequest) GetOrderType() string {
	if m != nil {
		return m.OrderType
	}
	return ""
}

func (m *NewOrderRequest) GetOrderTemplateID() string {
	if m != nil {
		return m.OrderTemplateID
	}
	return ""
}

func (m *NewOrderRequest) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *NewOrderRequest) GetParentOrderID() string {
	if m != nil {
		return m.ParentOrderID
	}
	return ""
}

func (m *NewOrderRequest) GetGrupo() string {
	if m != nil {
		return m.Grupo
	}
	return ""
}

func (m *NewOrderRequest) GetMarketName() string {
	if m != nil {
		return m.MarketName
	}
	return ""
}

func (m *NewOrderRequest) GetSide() string {
	if m != nil {
		return m.Side
	}
	return ""
}

func (m *NewOrderRequest) GetLimitPrice() float64 {
	if m != nil {
		return m.LimitPrice
	}
	return 0
}

func (m *NewOrderRequest) GetInitialCurrencyBalance() float64 {
	if m != nil {
		return m.InitialCurrencyBalance
	}
	return 0
}

func (m *NewOrderRequest) GetTriggers() []*TriggerRequest {
	if m != nil {
		return m.Triggers
	}
	return nil
}

type GetUserOrderRequest struct {
	OrderID string `protobuf:"bytes,1,opt,name=orderID" json:"orderID"`
	UserID  string `protobuf:"bytes,2,opt,name=userID" json:"userID"`
}

func (m *GetUserOrderRequest) Reset()                    { *m = GetUserOrderRequest{} }
func (m *GetUserOrderRequest) String() string            { return proto.CompactTextString(m) }
func (*GetUserOrderRequest) ProtoMessage()               {}
func (*GetUserOrderRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *GetUserOrderRequest) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
}

func (m *GetUserOrderRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

type GetUserOrdersRequest struct {
	UserID string `protobuf:"bytes,1,opt,name=userID" json:"userID"`
}

func (m *GetUserOrdersRequest) Reset()                    { *m = GetUserOrdersRequest{} }
func (m *GetUserOrdersRequest) String() string            { return proto.CompactTextString(m) }
func (*GetUserOrdersRequest) ProtoMessage()               {}
func (*GetUserOrdersRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *GetUserOrdersRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

type RemoveOrderRequest struct {
	OrderID string `protobuf:"bytes,1,opt,name=orderID" json:"orderID"`
	UserID  string `protobuf:"bytes,2,opt,name=userID" json:"userID"`
}

func (m *RemoveOrderRequest) Reset()                    { *m = RemoveOrderRequest{} }
func (m *RemoveOrderRequest) String() string            { return proto.CompactTextString(m) }
func (*RemoveOrderRequest) ProtoMessage()               {}
func (*RemoveOrderRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *RemoveOrderRequest) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
}

func (m *RemoveOrderRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

type Order struct {
	OrderID                  string     `protobuf:"bytes,1,opt,name=orderID" json:"orderID"`
	ParentOrderID            string     `protobuf:"bytes,2,opt,name=parentOrderID" json:"parentOrderID"`
	PlanID                   string     `protobuf:"bytes,3,opt,name=planID" json:"planID"`
	PlanDepth                uint32     `protobuf:"varint,4,opt,name=planDepth" json:"planDepth"`
	OrderTemplateID          string     `protobuf:"bytes,5,opt,name=orderTemplateID" json:"orderTemplateID"`
	AccountID                string     `protobuf:"bytes,6,opt,name=accountID" json:"accountID"`
	AccountType              string     `protobuf:"bytes,28,opt,name=accountType" json:"accountType"`
	KeyPublic                string     `protobuf:"bytes,7,opt,name=keyPublic" json:"keyPublic"`
	KeySecret                string     `protobuf:"bytes,8,opt,name=keySecret" json:"keySecret"`
	KeyDescription           string     `protobuf:"bytes,9,opt,name=keyDescription" json:"keyDescription"`
	OrderPriority            uint32     `protobuf:"varint,10,opt,name=orderPriority" json:"orderPriority"`
	OrderType                string     `protobuf:"bytes,11,opt,name=orderType" json:"orderType"`
	Side                     string     `protobuf:"bytes,12,opt,name=side" json:"side"`
	LimitPrice               float64    `protobuf:"fixed64,13,opt,name=limitPrice" json:"limitPrice"`
	Exchange                 string     `protobuf:"bytes,14,opt,name=exchange" json:"exchange"`
	ExchangeMarketName       string     `protobuf:"bytes,15,opt,name=exchangeMarketName" json:"exchangeMarketName"`
	MarketName               string     `protobuf:"bytes,16,opt,name=marketName" json:"marketName"`
	InitialCurrencySymbol    string     `protobuf:"bytes,17,opt,name=initialCurrencySymbol" json:"initialCurrencySymbol"`
	InitialCurrencyBalance   float64    `protobuf:"fixed64,18,opt,name=initialCurrencyBalance" json:"initialCurrencyBalance"`
	InitialCurrencyTraded    float64    `protobuf:"fixed64,19,opt,name=initialCurrencyTraded" json:"initialCurrencyTraded"`
	InitialCurrencyRemainder float64    `protobuf:"fixed64,20,opt,name=initialCurrencyRemainder" json:"initialCurrencyRemainder"`
	FinalCurrencySymbol      string     `protobuf:"bytes,21,opt,name=finalCurrencySymbol" json:"finalCurrencySymbol"`
	FinalCurrencyBalance     float64    `protobuf:"fixed64,22,opt,name=finalCurrencyBalance" json:"finalCurrencyBalance"`
	Status                   string     `protobuf:"bytes,23,opt,name=status" json:"status"`
	Grupo                    string     `protobuf:"bytes,24,opt,name=grupo" json:"grupo"`
	CreatedOn                string     `protobuf:"bytes,25,opt,name=createdOn" json:"createdOn"`
	UpdatedOn                string     `protobuf:"bytes,26,opt,name=updatedOn" json:"updatedOn"`
	Triggers                 []*Trigger `protobuf:"bytes,27,rep,name=triggers" json:"triggers"`
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

func (m *Order) GetParentOrderID() string {
	if m != nil {
		return m.ParentOrderID
	}
	return ""
}

func (m *Order) GetPlanID() string {
	if m != nil {
		return m.PlanID
	}
	return ""
}

func (m *Order) GetPlanDepth() uint32 {
	if m != nil {
		return m.PlanDepth
	}
	return 0
}

func (m *Order) GetOrderTemplateID() string {
	if m != nil {
		return m.OrderTemplateID
	}
	return ""
}

func (m *Order) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *Order) GetAccountType() string {
	if m != nil {
		return m.AccountType
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

func (m *Order) GetKeyDescription() string {
	if m != nil {
		return m.KeyDescription
	}
	return ""
}

func (m *Order) GetOrderPriority() uint32 {
	if m != nil {
		return m.OrderPriority
	}
	return 0
}

func (m *Order) GetOrderType() string {
	if m != nil {
		return m.OrderType
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

func (m *Order) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *Order) GetExchangeMarketName() string {
	if m != nil {
		return m.ExchangeMarketName
	}
	return ""
}

func (m *Order) GetMarketName() string {
	if m != nil {
		return m.MarketName
	}
	return ""
}

func (m *Order) GetInitialCurrencySymbol() string {
	if m != nil {
		return m.InitialCurrencySymbol
	}
	return ""
}

func (m *Order) GetInitialCurrencyBalance() float64 {
	if m != nil {
		return m.InitialCurrencyBalance
	}
	return 0
}

func (m *Order) GetInitialCurrencyTraded() float64 {
	if m != nil {
		return m.InitialCurrencyTraded
	}
	return 0
}

func (m *Order) GetInitialCurrencyRemainder() float64 {
	if m != nil {
		return m.InitialCurrencyRemainder
	}
	return 0
}

func (m *Order) GetFinalCurrencySymbol() string {
	if m != nil {
		return m.FinalCurrencySymbol
	}
	return ""
}

func (m *Order) GetFinalCurrencyBalance() float64 {
	if m != nil {
		return m.FinalCurrencyBalance
	}
	return 0
}

func (m *Order) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Order) GetGrupo() string {
	if m != nil {
		return m.Grupo
	}
	return ""
}

func (m *Order) GetCreatedOn() string {
	if m != nil {
		return m.CreatedOn
	}
	return ""
}

func (m *Order) GetUpdatedOn() string {
	if m != nil {
		return m.UpdatedOn
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
	TriggerID          string   `protobuf:"bytes,1,opt,name=triggerID" json:"triggerID"`
	TriggerTemplateID  string   `protobuf:"bytes,2,opt,name=triggerTemplateID" json:"triggerTemplateID"`
	OrderID            string   `protobuf:"bytes,3,opt,name=orderID" json:"orderID"`
	Index              uint32   `protobuf:"varint,4,opt,name=index" json:"index"`
	Title              string   `protobuf:"bytes,5,opt,name=title" json:"title"`
	Name               string   `protobuf:"bytes,6,opt,name=name" json:"name"`
	Code               string   `protobuf:"bytes,7,opt,name=code" json:"code"`
	Triggered          bool     `protobuf:"varint,8,opt,name=triggered" json:"triggered"`
	TriggeredPrice     float64  `protobuf:"fixed64,9,opt,name=triggeredPrice" json:"triggeredPrice"`
	TriggeredCondition string   `protobuf:"bytes,10,opt,name=triggeredCondition" json:"triggeredCondition"`
	TriggeredTimestamp string   `protobuf:"bytes,11,opt,name=triggeredTimestamp" json:"triggeredTimestamp"`
	CreatedOn          string   `protobuf:"bytes,12,opt,name=createdOn" json:"createdOn"`
	UpdatedOn          string   `protobuf:"bytes,13,opt,name=updatedOn" json:"updatedOn"`
	Actions            []string `protobuf:"bytes,14,rep,name=actions" json:"actions"`
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

func (m *Trigger) GetTriggerTemplateID() string {
	if m != nil {
		return m.TriggerTemplateID
	}
	return ""
}

func (m *Trigger) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
}

func (m *Trigger) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *Trigger) GetTitle() string {
	if m != nil {
		return m.Title
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

func (m *Trigger) GetTriggeredPrice() float64 {
	if m != nil {
		return m.TriggeredPrice
	}
	return 0
}

func (m *Trigger) GetTriggeredCondition() string {
	if m != nil {
		return m.TriggeredCondition
	}
	return ""
}

func (m *Trigger) GetTriggeredTimestamp() string {
	if m != nil {
		return m.TriggeredTimestamp
	}
	return ""
}

func (m *Trigger) GetCreatedOn() string {
	if m != nil {
		return m.CreatedOn
	}
	return ""
}

func (m *Trigger) GetUpdatedOn() string {
	if m != nil {
		return m.UpdatedOn
	}
	return ""
}

func (m *Trigger) GetActions() []string {
	if m != nil {
		return m.Actions
	}
	return nil
}

// Responses
type UserOrderData struct {
	Order *Order `protobuf:"bytes,1,opt,name=order" json:"order"`
}

func (m *UserOrderData) Reset()                    { *m = UserOrderData{} }
func (m *UserOrderData) String() string            { return proto.CompactTextString(m) }
func (*UserOrderData) ProtoMessage()               {}
func (*UserOrderData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *UserOrderData) GetOrder() *Order {
	if m != nil {
		return m.Order
	}
	return nil
}

type UserOrdersData struct {
	Orders []*Order `protobuf:"bytes,1,rep,name=orders" json:"orders"`
}

func (m *UserOrdersData) Reset()                    { *m = UserOrdersData{} }
func (m *UserOrdersData) String() string            { return proto.CompactTextString(m) }
func (*UserOrdersData) ProtoMessage()               {}
func (*UserOrdersData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *UserOrdersData) GetOrders() []*Order {
	if m != nil {
		return m.Orders
	}
	return nil
}

type OrderResponse struct {
	Status  string         `protobuf:"bytes,1,opt,name=status" json:"status"`
	Message string         `protobuf:"bytes,2,opt,name=message" json:"message"`
	Data    *UserOrderData `protobuf:"bytes,3,opt,name=data" json:"data"`
}

func (m *OrderResponse) Reset()                    { *m = OrderResponse{} }
func (m *OrderResponse) String() string            { return proto.CompactTextString(m) }
func (*OrderResponse) ProtoMessage()               {}
func (*OrderResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *OrderResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *OrderResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *OrderResponse) GetData() *UserOrderData {
	if m != nil {
		return m.Data
	}
	return nil
}

type OrderListResponse struct {
	Status  string          `protobuf:"bytes,1,opt,name=status" json:"status"`
	Message string          `protobuf:"bytes,2,opt,name=message" json:"message"`
	Data    *UserOrdersData `protobuf:"bytes,3,opt,name=data" json:"data"`
}

func (m *OrderListResponse) Reset()                    { *m = OrderListResponse{} }
func (m *OrderListResponse) String() string            { return proto.CompactTextString(m) }
func (*OrderListResponse) ProtoMessage()               {}
func (*OrderListResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *OrderListResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *OrderListResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *OrderListResponse) GetData() *UserOrdersData {
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
func (*OrdersPage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

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

func init() {
	proto.RegisterType((*TriggerRequest)(nil), "order.TriggerRequest")
	proto.RegisterType((*NewOrderRequest)(nil), "order.NewOrderRequest")
	proto.RegisterType((*GetUserOrderRequest)(nil), "order.GetUserOrderRequest")
	proto.RegisterType((*GetUserOrdersRequest)(nil), "order.GetUserOrdersRequest")
	proto.RegisterType((*RemoveOrderRequest)(nil), "order.RemoveOrderRequest")
	proto.RegisterType((*Order)(nil), "order.Order")
	proto.RegisterType((*Trigger)(nil), "order.Trigger")
	proto.RegisterType((*UserOrderData)(nil), "order.UserOrderData")
	proto.RegisterType((*UserOrdersData)(nil), "order.UserOrdersData")
	proto.RegisterType((*OrderResponse)(nil), "order.OrderResponse")
	proto.RegisterType((*OrderListResponse)(nil), "order.OrderListResponse")
	proto.RegisterType((*OrdersPage)(nil), "order.OrdersPage")
}

func init() {
	proto.RegisterFile("github.com/asciiu/gomo/plan-service/proto/order/order.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 984 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x56, 0xdd, 0x6e, 0x63, 0x35,
	0x10, 0x56, 0x9a, 0x34, 0x69, 0x26, 0x4d, 0x4a, 0xdd, 0x1f, 0xcc, 0x52, 0xa1, 0x28, 0x5a, 0xa1,
	0x80, 0xa0, 0x85, 0x2e, 0xda, 0x0b, 0xb8, 0x63, 0x23, 0x56, 0x95, 0x60, 0x5b, 0x9d, 0x96, 0x07,
	0x70, 0xcf, 0x19, 0x52, 0xab, 0x39, 0x3f, 0xd8, 0xce, 0xd2, 0xf0, 0x0a, 0xbc, 0x0d, 0x2f, 0xc0,
	0x3d, 0x4f, 0xb5, 0xf2, 0xd8, 0x27, 0xe7, 0x27, 0xa7, 0xdd, 0x95, 0xf6, 0x26, 0xf1, 0x7c, 0xdf,
	0xd8, 0x3e, 0xe3, 0xf9, 0xc6, 0x63, 0xf8, 0x69, 0x2e, 0xcd, 0xdd, 0xf2, 0xf6, 0x34, 0x4c, 0xe3,
	0x33, 0xa1, 0x43, 0x29, 0x97, 0x67, 0xf3, 0x34, 0x4e, 0xcf, 0xb2, 0x85, 0x48, 0xbe, 0xd5, 0xa8,
	0xde, 0xca, 0x10, 0xcf, 0x32, 0x95, 0x9a, 0xf4, 0x2c, 0x55, 0x11, 0x2a, 0xf7, 0x7b, 0x4a, 0x08,
	0xdb, 0x26, 0x63, 0xf2, 0x7f, 0x0b, 0x46, 0x37, 0x4a, 0xce, 0xe7, 0xa8, 0x02, 0xfc, 0x73, 0x89,
	0xda, 0xb0, 0x13, 0xe8, 0x1b, 0x87, 0x5c, 0xcc, 0x78, 0x6b, 0xdc, 0x9a, 0xf6, 0x83, 0x02, 0x60,
	0x0c, 0x3a, 0x61, 0x1a, 0x21, 0xdf, 0x22, 0x82, 0xc6, 0xec, 0x10, 0xb6, 0x65, 0x12, 0xe1, 0x03,
	0x6f, 0x8f, 0x5b, 0xd3, 0x61, 0xe0, 0x0c, 0xeb, 0x99, 0x88, 0x18, 0x79, 0xc7, 0x79, 0xda, 0xb1,
	0xf5, 0x34, 0xd2, 0x2c, 0x90, 0x6f, 0x13, 0xe8, 0x0c, 0xf6, 0x0d, 0xec, 0xfb, 0x0d, 0x6e, 0x30,
	0xce, 0x16, 0xc2, 0xe0, 0xc5, 0x8c, 0x77, 0xc9, 0x63, 0x93, 0x60, 0x1c, 0x7a, 0x22, 0x34, 0x32,
	0x4d, 0x34, 0xef, 0x8d, 0xdb, 0xd3, 0x7e, 0x90, 0x9b, 0x93, 0x7f, 0xdb, 0xb0, 0xf7, 0x06, 0xff,
	0xba, 0xb4, 0x91, 0xe5, 0xd1, 0x70, 0xe8, 0x51, 0xa4, 0xeb, 0x58, 0x72, 0x93, 0x3d, 0x87, 0x21,
	0x0d, 0xaf, 0x94, 0x4c, 0x95, 0x34, 0x2b, 0x0a, 0x69, 0x18, 0x54, 0x41, 0x7b, 0x1a, 0x04, 0xdc,
	0xac, 0x32, 0xa4, 0xf8, 0xfa, 0x41, 0x01, 0xb0, 0x29, 0xec, 0x39, 0xa3, 0xf8, 0x6e, 0x17, 0x6e,
	0x1d, 0xb6, 0xeb, 0x88, 0x30, 0x4c, 0x97, 0x89, 0xb9, 0x98, 0xf9, 0xe8, 0x0b, 0xc0, 0x7e, 0x4b,
	0x26, 0x14, 0x26, 0xe6, 0xd2, 0x7f, 0xab, 0x8b, 0xbe, 0x0a, 0xda, 0xd3, 0x9b, 0xab, 0x65, 0x96,
	0xf2, 0x9e, 0x3b, 0x3d, 0x32, 0xd8, 0x17, 0x00, 0xb1, 0x50, 0xf7, 0x68, 0xde, 0xd8, 0xd3, 0xde,
	0x21, 0xaa, 0x84, 0xd8, 0x3c, 0x68, 0x19, 0x21, 0xef, 0xbb, 0x3c, 0xd8, 0xb1, 0x9d, 0xb3, 0x90,
	0xb1, 0x34, 0x57, 0x4a, 0x86, 0xc8, 0x61, 0xdc, 0x9a, 0xb6, 0x82, 0x12, 0xc2, 0x5e, 0xc2, 0xb1,
	0x4c, 0xa4, 0x91, 0x62, 0xf1, 0x6a, 0xa9, 0x14, 0x26, 0xe1, 0xea, 0x67, 0xb1, 0x10, 0x49, 0x88,
	0x7c, 0x40, 0xbe, 0x8f, 0xb0, 0xec, 0x7b, 0xd8, 0xf1, 0x09, 0xd3, 0x7c, 0x77, 0xdc, 0x9e, 0x0e,
	0xce, 0x8f, 0x4e, 0x9d, 0xea, 0xaa, 0x22, 0x0b, 0xd6, 0x6e, 0x93, 0xd7, 0x70, 0xf0, 0x1a, 0xcd,
	0xef, 0x1a, 0xd5, 0x07, 0xe6, 0xed, 0x18, 0xba, 0x4b, 0x4d, 0x84, 0xd3, 0xa0, 0xb7, 0x26, 0xa7,
	0x70, 0x58, 0x5e, 0x48, 0xe7, 0x2b, 0x15, 0xfe, 0xad, 0x8a, 0xff, 0x2f, 0xc0, 0x02, 0x8c, 0xd3,
	0xb7, 0xf8, 0x91, 0xfb, 0xfe, 0xb3, 0x03, 0xdb, 0xb4, 0xc4, 0xd3, 0x5a, 0xab, 0xe6, 0x77, 0xab,
	0x29, 0xbf, 0xc7, 0xd0, 0xb5, 0xb5, 0x7b, 0x31, 0xf3, 0x42, 0xf3, 0x96, 0xd5, 0x8e, 0x1d, 0xcd,
	0x30, 0x33, 0x77, 0xa4, 0xaf, 0x61, 0x50, 0x00, 0x4d, 0x1a, 0xdc, 0xfe, 0x00, 0x0d, 0x76, 0xeb,
	0x1a, 0x1c, 0xc3, 0xc0, 0x1b, 0xa4, 0xf5, 0x13, 0xe2, 0xcb, 0x90, 0x9d, 0x7f, 0x8f, 0xab, 0xab,
	0xe5, 0xed, 0x42, 0x86, 0x5e, 0x83, 0x05, 0xe0, 0xd9, 0x6b, 0x0c, 0x15, 0x1a, 0x2f, 0xc3, 0x02,
	0x60, 0x5f, 0xc2, 0xe8, 0x1e, 0x57, 0x33, 0xd4, 0xa1, 0x92, 0x99, 0x2d, 0x57, 0xaf, 0xc7, 0x1a,
	0xba, 0x59, 0x95, 0xf0, 0xde, 0xaa, 0x1c, 0xd4, 0xab, 0x32, 0x57, 0xfc, 0xee, 0xa3, 0x8a, 0x1f,
	0x6e, 0x28, 0xfe, 0x19, 0xec, 0xe0, 0x43, 0x78, 0x27, 0x92, 0x39, 0xf2, 0x11, 0xcd, 0x5b, 0xdb,
	0xec, 0x14, 0x58, 0x3e, 0xfe, 0xad, 0xa8, 0xb4, 0x3d, 0xf2, 0x6a, 0x60, 0x6a, 0x15, 0xf9, 0xc9,
	0x46, 0x45, 0xfe, 0x00, 0x47, 0xb5, 0xfa, 0xb9, 0x5e, 0xc5, 0xb7, 0xe9, 0x82, 0xef, 0x93, 0x6b,
	0x33, 0xf9, 0x44, 0x4d, 0xb2, 0x27, 0x6b, 0x72, 0x73, 0xb7, 0x1b, 0x25, 0x22, 0x8c, 0xf8, 0x01,
	0x4d, 0x6b, 0x26, 0xd9, 0x8f, 0xc0, 0x6b, 0x44, 0x80, 0xb1, 0xb0, 0x37, 0xbb, 0xe2, 0x87, 0x34,
	0xf1, 0x51, 0x9e, 0x7d, 0x07, 0x07, 0x7f, 0xc8, 0x64, 0x23, 0xba, 0x23, 0x8a, 0xae, 0x89, 0x62,
	0xe7, 0x70, 0x58, 0x81, 0xf3, 0xc8, 0x8e, 0x69, 0xa7, 0x46, 0xce, 0x56, 0x8b, 0x36, 0xc2, 0x2c,
	0x35, 0xff, 0xd4, 0x55, 0x8b, 0xb3, 0x8a, 0x5b, 0x92, 0x97, 0x6f, 0xc9, 0x13, 0xe8, 0x87, 0x0a,
	0x85, 0xc1, 0xe8, 0x32, 0xe1, 0x9f, 0x39, 0xc5, 0xac, 0x01, 0xcb, 0x2e, 0xb3, 0xc8, 0xb3, 0xcf,
	0x1c, 0xbb, 0x06, 0xd8, 0xd7, 0xa5, 0x5b, 0xed, 0x73, 0xba, 0xd5, 0x46, 0xb5, 0x5b, 0xad, 0xb8,
	0xce, 0xfe, 0x6b, 0x43, 0xcf, 0xa3, 0xef, 0xe9, 0xa4, 0x8d, 0x5d, 0x6f, 0xeb, 0x89, 0xae, 0x97,
	0xdf, 0x2d, 0xed, 0xea, 0xdd, 0xb2, 0xee, 0xbe, 0x9d, 0x72, 0xf7, 0x6d, 0xee, 0xb4, 0x79, 0x4f,
	0xee, 0x96, 0x7a, 0x72, 0xde, 0xd1, 0x7b, 0xa5, 0x8e, 0x5e, 0x7c, 0x39, 0x46, 0x54, 0xcb, 0x3b,
	0x41, 0x01, 0xd8, 0x5a, 0x5e, 0x1b, 0xae, 0x9e, 0xfa, 0x94, 0xa7, 0x1a, 0x6a, 0xeb, 0x66, 0x8d,
	0xbc, 0x4a, 0x93, 0x48, 0x52, 0xdd, 0x83, 0xab, 0x9b, 0x4d, 0xa6, 0xe2, 0x7f, 0x23, 0x63, 0xd4,
	0x46, 0xc4, 0x99, 0x2f, 0xef, 0x06, 0xa6, 0x9a, 0xd3, 0xdd, 0x27, 0x73, 0x3a, 0xac, 0xe7, 0xb4,
	0xf4, 0x8a, 0x18, 0x55, 0x5f, 0x11, 0x2f, 0x60, 0xb8, 0x6e, 0x22, 0x33, 0x61, 0x04, 0x9b, 0x80,
	0x7b, 0x2c, 0x51, 0x0a, 0x07, 0xe7, 0xbb, 0x3e, 0xf7, 0xae, 0x6d, 0xf8, 0x77, 0xd4, 0x4b, 0x18,
	0x15, 0x9d, 0x87, 0x66, 0x3d, 0x87, 0x2e, 0x51, 0x9a, 0xb7, 0x48, 0x32, 0xd5, 0x69, 0x9e, 0x9b,
	0xdc, 0xc3, 0xd0, 0xb7, 0x1f, 0x9d, 0xa5, 0x89, 0x2e, 0xab, 0xba, 0x55, 0x51, 0x35, 0x87, 0x5e,
	0x8c, 0x5a, 0x8b, 0x79, 0xfe, 0xf4, 0xca, 0x4d, 0x36, 0x85, 0x4e, 0x24, 0x8c, 0x20, 0x59, 0x0c,
	0xce, 0x0f, 0xfd, 0x36, 0x95, 0x10, 0x02, 0xf2, 0x98, 0x64, 0xb0, 0x4f, 0xd0, 0xaf, 0x52, 0x9b,
	0x8f, 0xd8, 0xf0, 0xab, 0xca, 0x86, 0x47, 0xf5, 0x0d, 0x75, 0x69, 0xc7, 0x07, 0x00, 0x87, 0x5d,
	0xd9, 0x89, 0x0c, 0x3a, 0x99, 0x5d, 0xaf, 0x45, 0x42, 0xa5, 0xb1, 0xbd, 0x77, 0xed, 0xff, 0xb5,
	0xfc, 0x1b, 0xfd, 0x03, 0x6c, 0x6d, 0x93, 0x86, 0x53, 0x23, 0x16, 0xf9, 0xbb, 0x92, 0x8c, 0xd2,
	0xc1, 0x76, 0x1e, 0x3f, 0xd8, 0xdb, 0x2e, 0x3d, 0x73, 0x5f, 0xbc, 0x0b, 0x00, 0x00, 0xff, 0xff,
	0x5e, 0x25, 0xf5, 0x4f, 0x25, 0x0b, 0x00, 0x00,
}
