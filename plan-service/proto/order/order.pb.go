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
	// 974 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x56, 0xdd, 0x6e, 0x63, 0x35,
	0x10, 0x56, 0x9a, 0x34, 0x69, 0xa6, 0x4d, 0x4a, 0xdd, 0x1f, 0xcc, 0x52, 0xa1, 0x2a, 0x5a, 0xa1,
	0x80, 0xa0, 0x85, 0x2e, 0xda, 0x0b, 0xb8, 0x63, 0x23, 0x56, 0x95, 0x60, 0x5b, 0x9d, 0x96, 0x07,
	0x70, 0xcf, 0x19, 0x52, 0xab, 0x39, 0x3f, 0xd8, 0xce, 0xd2, 0xf0, 0x58, 0xbc, 0x00, 0xf7, 0xbc,
	0x09, 0x6f, 0xb1, 0xf2, 0xd8, 0xe7, 0x37, 0xa7, 0xdd, 0x95, 0xf6, 0x26, 0xf1, 0x7c, 0x33, 0xb6,
	0xcf, 0x78, 0xbe, 0xcf, 0x63, 0xf8, 0x69, 0x2e, 0xcd, 0xdd, 0xf2, 0xf6, 0x34, 0x4c, 0xe3, 0x33,
	0xa1, 0x43, 0x29, 0x97, 0x67, 0xf3, 0x34, 0x4e, 0xcf, 0xb2, 0x85, 0x48, 0xbe, 0xd5, 0xa8, 0xde,
	0xca, 0x10, 0xcf, 0x32, 0x95, 0x9a, 0xf4, 0x2c, 0x55, 0x11, 0x2a, 0xf7, 0x7b, 0x4a, 0x08, 0xdb,
	0x24, 0x63, 0xf2, 0x5f, 0x07, 0xc6, 0x37, 0x4a, 0xce, 0xe7, 0xa8, 0x02, 0xfc, 0x73, 0x89, 0xda,
	0xb0, 0x63, 0x18, 0x1a, 0x87, 0x5c, 0xcc, 0x78, 0xe7, 0xa4, 0x33, 0x1d, 0x06, 0x25, 0xc0, 0x18,
	0xf4, 0xc2, 0x34, 0x42, 0xbe, 0x41, 0x0e, 0x1a, 0xb3, 0x03, 0xd8, 0x94, 0x49, 0x84, 0x0f, 0xbc,
	0x7b, 0xd2, 0x99, 0x8e, 0x02, 0x67, 0xd8, 0xc8, 0x44, 0xc4, 0xc8, 0x7b, 0x2e, 0xd2, 0x8e, 0x6d,
	0xa4, 0x91, 0x66, 0x81, 0x7c, 0x93, 0x40, 0x67, 0xb0, 0x6f, 0x60, 0xcf, 0x6f, 0x70, 0x83, 0x71,
	0xb6, 0x10, 0x06, 0x2f, 0x66, 0xbc, 0x4f, 0x11, 0xeb, 0x0e, 0xc6, 0x61, 0x20, 0x42, 0x23, 0xd3,
	0x44, 0xf3, 0xc1, 0x49, 0x77, 0x3a, 0x0c, 0x72, 0x73, 0xf2, 0x4f, 0x17, 0x76, 0xdf, 0xe0, 0x5f,
	0x97, 0x36, 0xb3, 0x3c, 0x1b, 0x0e, 0x03, 0xca, 0xb4, 0xc8, 0x25, 0x37, 0xd9, 0x73, 0x18, 0xd1,
	0xf0, 0x4a, 0xc9, 0x54, 0x49, 0xb3, 0xa2, 0x94, 0x46, 0x41, 0x1d, 0xb4, 0xa7, 0x41, 0xc0, 0xcd,
	0x2a, 0x43, 0xca, 0x6f, 0x18, 0x94, 0x00, 0x9b, 0xc2, 0xae, 0x33, 0xca, 0xef, 0x76, 0xe9, 0x36,
	0x61, 0xbb, 0x8e, 0x08, 0xc3, 0x74, 0x99, 0x98, 0x8b, 0x99, 0xcf, 0xbe, 0x04, 0xec, 0xb7, 0x64,
	0x42, 0x61, 0x62, 0x2e, 0xfd, 0xb7, 0xba, 0xec, 0xeb, 0xa0, 0x3d, 0xbd, 0xb9, 0x5a, 0x66, 0x29,
	0x1f, 0xb8, 0xd3, 0x23, 0x83, 0x7d, 0x01, 0x10, 0x0b, 0x75, 0x8f, 0xe6, 0x8d, 0x3d, 0xed, 0x2d,
	0x72, 0x55, 0x10, 0x5b, 0x07, 0x2d, 0x23, 0xe4, 0x43, 0x57, 0x07, 0x3b, 0xb6, 0x73, 0x16, 0x32,
	0x96, 0xe6, 0x4a, 0xc9, 0x10, 0x39, 0x9c, 0x74, 0xa6, 0x9d, 0xa0, 0x82, 0xb0, 0x97, 0x70, 0x24,
	0x13, 0x69, 0xa4, 0x58, 0xbc, 0x5a, 0x2a, 0x85, 0x49, 0xb8, 0xfa, 0x59, 0x2c, 0x44, 0x12, 0x22,
	0xdf, 0xa6, 0xd8, 0x47, 0xbc, 0xec, 0x7b, 0xd8, 0xf2, 0x05, 0xd3, 0x7c, 0xe7, 0xa4, 0x3b, 0xdd,
	0x3e, 0x3f, 0x3c, 0x75, 0xac, 0xab, 0x93, 0x2c, 0x28, 0xc2, 0x26, 0xaf, 0x61, 0xff, 0x35, 0x9a,
	0xdf, 0x35, 0xaa, 0x0f, 0xac, 0xdb, 0x11, 0xf4, 0x97, 0x9a, 0x1c, 0x8e, 0x83, 0xde, 0x9a, 0x9c,
	0xc2, 0x41, 0x75, 0x21, 0x9d, 0xaf, 0x54, 0xc6, 0x77, 0x6a, 0xf1, 0xbf, 0x00, 0x0b, 0x30, 0x4e,
	0xdf, 0xe2, 0x47, 0xee, 0xfb, 0xff, 0x00, 0x36, 0x69, 0x89, 0xa7, 0xb9, 0x56, 0xaf, 0xef, 0x46,
	0x5b, 0x7d, 0x8f, 0xa0, 0x6f, 0xb5, 0x7b, 0x31, 0xf3, 0x44, 0xf3, 0x96, 0xe5, 0x8e, 0x1d, 0xcd,
	0x30, 0x33, 0x77, 0xc4, 0xaf, 0x51, 0x50, 0x02, 0x6d, 0x1c, 0xdc, 0xfc, 0x00, 0x0e, 0xf6, 0x9b,
	0x1c, 0x3c, 0x86, 0xe1, 0x3d, 0xae, 0xae, 0x96, 0xb7, 0x0b, 0x19, 0x7a, 0x86, 0x95, 0x80, 0xf7,
	0x5e, 0x63, 0xa8, 0xd0, 0x78, 0x92, 0x95, 0x00, 0xfb, 0x12, 0xc6, 0xf7, 0xb8, 0x9a, 0xa1, 0x0e,
	0x95, 0xcc, 0xac, 0x18, 0x3d, 0xdb, 0x1a, 0xe8, 0xba, 0xe6, 0xe0, 0xbd, 0x9a, 0xdb, 0x6e, 0x6a,
	0x2e, 0xe7, 0xf3, 0xce, 0xa3, 0x7c, 0x1e, 0xad, 0xf1, 0xf9, 0x19, 0x6c, 0xe1, 0x43, 0x78, 0x27,
	0x92, 0x39, 0xf2, 0x31, 0xcd, 0x2b, 0x6c, 0x76, 0x0a, 0x2c, 0x1f, 0xff, 0x56, 0xea, 0x68, 0x97,
	0xa2, 0x5a, 0x3c, 0x0d, 0xbd, 0x7d, 0xb2, 0xa6, 0xb7, 0x1f, 0xe0, 0xb0, 0xa1, 0x8e, 0xeb, 0x55,
	0x7c, 0x9b, 0x2e, 0xf8, 0x1e, 0x85, 0xb6, 0x3b, 0x9f, 0x50, 0x1c, 0x7b, 0x52, 0x71, 0xeb, 0xbb,
	0xdd, 0x28, 0x11, 0x61, 0xc4, 0xf7, 0x69, 0x5a, 0xbb, 0x93, 0xfd, 0x08, 0xbc, 0xe1, 0x08, 0x30,
	0x16, 0xf6, 0xde, 0x56, 0xfc, 0x80, 0x26, 0x3e, 0xea, 0x67, 0xdf, 0xc1, 0xfe, 0x1f, 0x32, 0x59,
	0xcb, 0xee, 0x90, 0xb2, 0x6b, 0x73, 0xb1, 0x73, 0x38, 0xa8, 0xc1, 0x79, 0x66, 0x47, 0xb4, 0x53,
	0xab, 0xcf, 0x6a, 0x41, 0x1b, 0x61, 0x96, 0x9a, 0x7f, 0xea, 0xb4, 0xe0, 0xac, 0xf2, 0x0e, 0xe4,
	0xd5, 0x3b, 0xf0, 0x18, 0x86, 0xa1, 0x42, 0x61, 0x30, 0xba, 0x4c, 0xf8, 0x67, 0x8e, 0x31, 0x05,
	0x60, 0xbd, 0xcb, 0x2c, 0xf2, 0xde, 0x67, 0xce, 0x5b, 0x00, 0xec, 0xeb, 0xca, 0x9d, 0xf5, 0x39,
	0xdd, 0x59, 0xe3, 0xc6, 0x9d, 0x55, 0x5e, 0x56, 0xff, 0x76, 0x61, 0xe0, 0xd1, 0xf7, 0xf4, 0xc9,
	0xd6, 0x9e, 0xb6, 0xf1, 0x44, 0x4f, 0xcb, 0x6f, 0x8e, 0x6e, 0xfd, 0xe6, 0x28, 0x7a, 0x6b, 0xaf,
	0xda, 0x5b, 0xdb, 0xfb, 0x68, 0xde, 0x71, 0xfb, 0x95, 0x8e, 0x9b, 0xf7, 0xeb, 0x41, 0xa5, 0x5f,
	0x97, 0x5f, 0x8e, 0x11, 0x69, 0x79, 0x2b, 0x28, 0x01, 0xab, 0xe5, 0xc2, 0x70, 0x7a, 0x1a, 0x52,
	0x9d, 0x1a, 0xa8, 0xd5, 0x4d, 0x81, 0xbc, 0x4a, 0x93, 0x48, 0x92, 0xee, 0xc1, 0xe9, 0x66, 0xdd,
	0x53, 0x8b, 0xbf, 0x91, 0x31, 0x6a, 0x23, 0xe2, 0xcc, 0xcb, 0xbb, 0xc5, 0x53, 0xaf, 0xe9, 0xce,
	0x93, 0x35, 0x1d, 0x35, 0x6b, 0x5a, 0x79, 0x23, 0x8c, 0xeb, 0x6f, 0x84, 0x17, 0x30, 0x2a, 0x5a,
	0xc4, 0x4c, 0x18, 0xc1, 0x26, 0xe0, 0x9e, 0x42, 0x54, 0xc2, 0xed, 0xf3, 0x1d, 0x5f, 0x7b, 0xd7,
	0x14, 0xfc, 0x2b, 0xe9, 0x25, 0x8c, 0xcb, 0xbe, 0x42, 0xb3, 0x9e, 0x43, 0x9f, 0x5c, 0x9a, 0x77,
	0x88, 0x32, 0xf5, 0x69, 0xde, 0x37, 0xb9, 0x87, 0x91, 0x6f, 0x2e, 0x3a, 0x4b, 0x13, 0x5d, 0x65,
	0x75, 0xa7, 0xc6, 0x6a, 0x0e, 0x83, 0x18, 0xb5, 0x16, 0xf3, 0xfc, 0x61, 0x95, 0x9b, 0x6c, 0x0a,
	0xbd, 0x48, 0x18, 0x41, 0xb4, 0xd8, 0x3e, 0x3f, 0xf0, 0xdb, 0xd4, 0x52, 0x08, 0x28, 0x62, 0x92,
	0xc1, 0x1e, 0x41, 0xbf, 0x4a, 0x6d, 0x3e, 0x62, 0xc3, 0xaf, 0x6a, 0x1b, 0x1e, 0x36, 0x37, 0xd4,
	0x95, 0x1d, 0x1f, 0x00, 0x1c, 0x76, 0x65, 0x27, 0x32, 0xe8, 0x65, 0x76, 0xbd, 0x0e, 0x11, 0x95,
	0xc6, 0xf6, 0xde, 0xb5, 0xff, 0xd7, 0xf2, 0x6f, 0xf4, 0xcf, 0xab, 0xc2, 0x26, 0x0e, 0xa7, 0x46,
	0x2c, 0xf2, 0x57, 0x23, 0x19, 0x95, 0x83, 0xed, 0x3d, 0x7e, 0xb0, 0xb7, 0x7d, 0x7a, 0xc4, 0xbe,
	0x78, 0x17, 0x00, 0x00, 0xff, 0xff, 0x51, 0x13, 0xa6, 0xf5, 0x03, 0x0b, 0x00, 0x00,
}
