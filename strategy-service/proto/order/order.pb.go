// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/order/order.proto

/*
Package fomo_orders is a generated protocol buffer package.

It is generated from these files:
	proto/order/order.proto

It has these top-level messages:
	OrdersRequest
	OrderRequest
	GetUserOrderRequest
	GetUserOrdersRequest
	RemoveOrderRequest
	Order
	UserOrderData
	UserOrdersData
	OrderResponse
	OrderListResponse
*/
package fomo_orders

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
type OrdersRequest struct {
	Orders []*OrderRequest `protobuf:"bytes,1,rep,name=orders" json:"orders"`
}

func (m *OrdersRequest) Reset()                    { *m = OrdersRequest{} }
func (m *OrdersRequest) String() string            { return proto.CompactTextString(m) }
func (*OrdersRequest) ProtoMessage()               {}
func (*OrdersRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *OrdersRequest) GetOrders() []*OrderRequest {
	if m != nil {
		return m.Orders
	}
	return nil
}

type OrderRequest struct {
	OrderID          string  `protobuf:"bytes,1,opt,name=orderID" json:"orderID"`
	UserID           string  `protobuf:"bytes,2,opt,name=userID" json:"userID"`
	KeyID            string  `protobuf:"bytes,3,opt,name=keyID" json:"keyID"`
	Exchange         string  `protobuf:"bytes,4,opt,name=exchange" json:"exchange"`
	MarketName       string  `protobuf:"bytes,5,opt,name=marketName" json:"marketName"`
	Side             string  `protobuf:"bytes,6,opt,name=side" json:"side"`
	OrderType        string  `protobuf:"bytes,7,opt,name=orderType" json:"orderType"`
	BaseQuantity     float64 `protobuf:"fixed64,8,opt,name=baseQuantity" json:"baseQuantity"`
	BasePercent      float64 `protobuf:"fixed64,9,opt,name=basePercent" json:"basePercent"`
	CurrencyQuantity float64 `protobuf:"fixed64,10,opt,name=currencyQuantity" json:"currencyQuantity"`
	CurrencyPercent  float64 `protobuf:"fixed64,11,opt,name=currencyPercent" json:"currencyPercent"`
	Conditions       string  `protobuf:"bytes,12,opt,name=conditions" json:"conditions"`
	Price            float64 `protobuf:"fixed64,13,opt,name=price" json:"price"`
	ParentOrderID    string  `protobuf:"bytes,14,opt,name=parentOrderID" json:"parentOrderID"`
	Active           bool    `protobuf:"varint,15,opt,name=active" json:"active"`
}

func (m *OrderRequest) Reset()                    { *m = OrderRequest{} }
func (m *OrderRequest) String() string            { return proto.CompactTextString(m) }
func (*OrderRequest) ProtoMessage()               {}
func (*OrderRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *OrderRequest) GetOrderID() string {
	if m != nil {
		return m.OrderID
	}
	return ""
}

func (m *OrderRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *OrderRequest) GetKeyID() string {
	if m != nil {
		return m.KeyID
	}
	return ""
}

func (m *OrderRequest) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *OrderRequest) GetMarketName() string {
	if m != nil {
		return m.MarketName
	}
	return ""
}

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

func (m *OrderRequest) GetBaseQuantity() float64 {
	if m != nil {
		return m.BaseQuantity
	}
	return 0
}

func (m *OrderRequest) GetBasePercent() float64 {
	if m != nil {
		return m.BasePercent
	}
	return 0
}

func (m *OrderRequest) GetCurrencyQuantity() float64 {
	if m != nil {
		return m.CurrencyQuantity
	}
	return 0
}

func (m *OrderRequest) GetCurrencyPercent() float64 {
	if m != nil {
		return m.CurrencyPercent
	}
	return 0
}

func (m *OrderRequest) GetConditions() string {
	if m != nil {
		return m.Conditions
	}
	return ""
}

func (m *OrderRequest) GetPrice() float64 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *OrderRequest) GetParentOrderID() string {
	if m != nil {
		return m.ParentOrderID
	}
	return ""
}

func (m *OrderRequest) GetActive() bool {
	if m != nil {
		return m.Active
	}
	return false
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

// Responses
type Order struct {
	OrderID            string  `protobuf:"bytes,1,opt,name=orderID" json:"orderID"`
	UserID             string  `protobuf:"bytes,2,opt,name=userID" json:"userID"`
	KeyID              string  `protobuf:"bytes,3,opt,name=keyID" json:"keyID"`
	Key                string  `protobuf:"bytes,4,opt,name=key" json:"key"`
	Secret             string  `protobuf:"bytes,5,opt,name=secret" json:"secret"`
	Exchange           string  `protobuf:"bytes,6,opt,name=exchange" json:"exchange"`
	ExchangeOrderID    string  `protobuf:"bytes,7,opt,name=exchangeOrderID" json:"exchangeOrderID"`
	ExchangeMarketName string  `protobuf:"bytes,8,opt,name=exchangeMarketName" json:"exchangeMarketName"`
	MarketName         string  `protobuf:"bytes,9,opt,name=marketName" json:"marketName"`
	Side               string  `protobuf:"bytes,10,opt,name=side" json:"side"`
	OrderType          string  `protobuf:"bytes,11,opt,name=orderType" json:"orderType"`
	Price              float64 `protobuf:"fixed64,13,opt,name=price" json:"price"`
	BaseQuantity       float64 `protobuf:"fixed64,14,opt,name=baseQuantity" json:"baseQuantity"`
	BasePercent        float64 `protobuf:"fixed64,15,opt,name=basePercent" json:"basePercent"`
	CurrencyQuantity   float64 `protobuf:"fixed64,16,opt,name=currencyQuantity" json:"currencyQuantity"`
	CurrencyPercent    float64 `protobuf:"fixed64,17,opt,name=currencyPercent" json:"currencyPercent"`
	Status             string  `protobuf:"bytes,18,opt,name=status" json:"status"`
	ChainStatus        string  `protobuf:"bytes,19,opt,name=chainStatus" json:"chainStatus"`
	Conditions         string  `protobuf:"bytes,20,opt,name=conditions" json:"conditions"`
	Condition          string  `protobuf:"bytes,21,opt,name=condition" json:"condition"`
	ParentOrderID      string  `protobuf:"bytes,22,opt,name=parentOrderID" json:"parentOrderID"`
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

func (m *Order) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *Order) GetKeyID() string {
	if m != nil {
		return m.KeyID
	}
	return ""
}

func (m *Order) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Order) GetSecret() string {
	if m != nil {
		return m.Secret
	}
	return ""
}

func (m *Order) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *Order) GetExchangeOrderID() string {
	if m != nil {
		return m.ExchangeOrderID
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

func (m *Order) GetSide() string {
	if m != nil {
		return m.Side
	}
	return ""
}

func (m *Order) GetOrderType() string {
	if m != nil {
		return m.OrderType
	}
	return ""
}

func (m *Order) GetPrice() float64 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *Order) GetBaseQuantity() float64 {
	if m != nil {
		return m.BaseQuantity
	}
	return 0
}

func (m *Order) GetBasePercent() float64 {
	if m != nil {
		return m.BasePercent
	}
	return 0
}

func (m *Order) GetCurrencyQuantity() float64 {
	if m != nil {
		return m.CurrencyQuantity
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

func (m *Order) GetChainStatus() string {
	if m != nil {
		return m.ChainStatus
	}
	return ""
}

func (m *Order) GetConditions() string {
	if m != nil {
		return m.Conditions
	}
	return ""
}

func (m *Order) GetCondition() string {
	if m != nil {
		return m.Condition
	}
	return ""
}

func (m *Order) GetParentOrderID() string {
	if m != nil {
		return m.ParentOrderID
	}
	return ""
}

type UserOrderData struct {
	Order *Order `protobuf:"bytes,1,opt,name=order" json:"order"`
}

func (m *UserOrderData) Reset()                    { *m = UserOrderData{} }
func (m *UserOrderData) String() string            { return proto.CompactTextString(m) }
func (*UserOrderData) ProtoMessage()               {}
func (*UserOrderData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

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
func (*UserOrdersData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

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
func (*OrderResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

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
func (*OrderListResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

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

func init() {
	proto.RegisterType((*OrdersRequest)(nil), "fomo.orders.OrdersRequest")
	proto.RegisterType((*OrderRequest)(nil), "fomo.orders.OrderRequest")
	proto.RegisterType((*GetUserOrderRequest)(nil), "fomo.orders.GetUserOrderRequest")
	proto.RegisterType((*GetUserOrdersRequest)(nil), "fomo.orders.GetUserOrdersRequest")
	proto.RegisterType((*RemoveOrderRequest)(nil), "fomo.orders.RemoveOrderRequest")
	proto.RegisterType((*Order)(nil), "fomo.orders.Order")
	proto.RegisterType((*UserOrderData)(nil), "fomo.orders.UserOrderData")
	proto.RegisterType((*UserOrdersData)(nil), "fomo.orders.UserOrdersData")
	proto.RegisterType((*OrderResponse)(nil), "fomo.orders.OrderResponse")
	proto.RegisterType((*OrderListResponse)(nil), "fomo.orders.OrderListResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for OrderService service

type OrderServiceClient interface {
	AddOrders(ctx context.Context, in *OrdersRequest, opts ...client.CallOption) (*OrderListResponse, error)
	GetUserOrder(ctx context.Context, in *GetUserOrderRequest, opts ...client.CallOption) (*OrderResponse, error)
	GetUserOrders(ctx context.Context, in *GetUserOrdersRequest, opts ...client.CallOption) (*OrderListResponse, error)
	RemoveOrder(ctx context.Context, in *RemoveOrderRequest, opts ...client.CallOption) (*OrderResponse, error)
	UpdateOrder(ctx context.Context, in *OrderRequest, opts ...client.CallOption) (*OrderResponse, error)
}

type orderServiceClient struct {
	c           client.Client
	serviceName string
}

func NewOrderServiceClient(serviceName string, c client.Client) OrderServiceClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "fomo.orders"
	}
	return &orderServiceClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *orderServiceClient) AddOrders(ctx context.Context, in *OrdersRequest, opts ...client.CallOption) (*OrderListResponse, error) {
	req := c.c.NewRequest(c.serviceName, "OrderService.AddOrders", in)
	out := new(OrderListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) GetUserOrder(ctx context.Context, in *GetUserOrderRequest, opts ...client.CallOption) (*OrderResponse, error) {
	req := c.c.NewRequest(c.serviceName, "OrderService.GetUserOrder", in)
	out := new(OrderResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) GetUserOrders(ctx context.Context, in *GetUserOrdersRequest, opts ...client.CallOption) (*OrderListResponse, error) {
	req := c.c.NewRequest(c.serviceName, "OrderService.GetUserOrders", in)
	out := new(OrderListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) RemoveOrder(ctx context.Context, in *RemoveOrderRequest, opts ...client.CallOption) (*OrderResponse, error) {
	req := c.c.NewRequest(c.serviceName, "OrderService.RemoveOrder", in)
	out := new(OrderResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) UpdateOrder(ctx context.Context, in *OrderRequest, opts ...client.CallOption) (*OrderResponse, error) {
	req := c.c.NewRequest(c.serviceName, "OrderService.UpdateOrder", in)
	out := new(OrderResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for OrderService service

type OrderServiceHandler interface {
	AddOrders(context.Context, *OrdersRequest, *OrderListResponse) error
	GetUserOrder(context.Context, *GetUserOrderRequest, *OrderResponse) error
	GetUserOrders(context.Context, *GetUserOrdersRequest, *OrderListResponse) error
	RemoveOrder(context.Context, *RemoveOrderRequest, *OrderResponse) error
	UpdateOrder(context.Context, *OrderRequest, *OrderResponse) error
}

func RegisterOrderServiceHandler(s server.Server, hdlr OrderServiceHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&OrderService{hdlr}, opts...))
}

type OrderService struct {
	OrderServiceHandler
}

func (h *OrderService) AddOrders(ctx context.Context, in *OrdersRequest, out *OrderListResponse) error {
	return h.OrderServiceHandler.AddOrders(ctx, in, out)
}

func (h *OrderService) GetUserOrder(ctx context.Context, in *GetUserOrderRequest, out *OrderResponse) error {
	return h.OrderServiceHandler.GetUserOrder(ctx, in, out)
}

func (h *OrderService) GetUserOrders(ctx context.Context, in *GetUserOrdersRequest, out *OrderListResponse) error {
	return h.OrderServiceHandler.GetUserOrders(ctx, in, out)
}

func (h *OrderService) RemoveOrder(ctx context.Context, in *RemoveOrderRequest, out *OrderResponse) error {
	return h.OrderServiceHandler.RemoveOrder(ctx, in, out)
}

func (h *OrderService) UpdateOrder(ctx context.Context, in *OrderRequest, out *OrderResponse) error {
	return h.OrderServiceHandler.UpdateOrder(ctx, in, out)
}

func init() { proto.RegisterFile("proto/order/order.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 698 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0xdd, 0x4e, 0xdb, 0x4a,
	0x10, 0x3e, 0x3e, 0x21, 0x21, 0x19, 0x27, 0x04, 0x06, 0x0e, 0x67, 0x4f, 0x0e, 0xa2, 0xa9, 0xd5,
	0x8b, 0x88, 0x8b, 0xa0, 0xa6, 0x57, 0x95, 0x7a, 0xd3, 0x0a, 0x81, 0x90, 0x28, 0xb4, 0x06, 0x1e,
	0x60, 0xb1, 0xa7, 0x60, 0xa1, 0xd8, 0xc1, 0xbb, 0x89, 0x9a, 0xc7, 0xe8, 0x45, 0xdf, 0xac, 0x0f,
	0x54, 0x79, 0xfc, 0x13, 0xdb, 0x71, 0x0a, 0x15, 0xbd, 0x89, 0x3c, 0xdf, 0xf7, 0xcd, 0xec, 0x7a,
	0x77, 0xbe, 0x89, 0xe1, 0xdf, 0x49, 0x18, 0xe8, 0xe0, 0x30, 0x08, 0x5d, 0x0a, 0xe3, 0xdf, 0x21,
	0x23, 0x68, 0x7e, 0x09, 0xc6, 0xc1, 0x90, 0x11, 0x65, 0x7d, 0x80, 0xce, 0x05, 0x3f, 0xd9, 0xf4,
	0x30, 0x25, 0xa5, 0xf1, 0x35, 0x34, 0x62, 0x4a, 0x18, 0xfd, 0xda, 0xc0, 0x1c, 0xfd, 0x37, 0xcc,
	0xc9, 0x87, 0xac, 0x4d, 0xa4, 0x76, 0x22, 0xb4, 0x7e, 0xd4, 0xa0, 0x9d, 0x27, 0x50, 0xc0, 0x3a,
	0x53, 0xa7, 0x47, 0xc2, 0xe8, 0x1b, 0x83, 0x96, 0x9d, 0x86, 0xb8, 0x0b, 0x8d, 0xa9, 0x62, 0xe2,
	0x6f, 0x26, 0x92, 0x08, 0x77, 0xa0, 0x7e, 0x4f, 0xf3, 0xd3, 0x23, 0x51, 0x63, 0x38, 0x0e, 0xb0,
	0x07, 0x4d, 0xfa, 0xea, 0xdc, 0x49, 0xff, 0x96, 0xc4, 0x1a, 0x13, 0x59, 0x8c, 0xfb, 0x00, 0x63,
	0x19, 0xde, 0x93, 0x3e, 0x97, 0x63, 0x12, 0x75, 0x66, 0x73, 0x08, 0x22, 0xac, 0x29, 0xcf, 0x25,
	0xd1, 0x60, 0x86, 0x9f, 0x71, 0x0f, 0x5a, 0xbc, 0x91, 0xab, 0xf9, 0x84, 0xc4, 0x3a, 0x13, 0x0b,
	0x00, 0x2d, 0x68, 0xdf, 0x48, 0x45, 0x9f, 0xa7, 0xd2, 0xd7, 0x9e, 0x9e, 0x8b, 0x66, 0xdf, 0x18,
	0x18, 0x76, 0x01, 0xc3, 0x3e, 0x98, 0x51, 0xfc, 0x89, 0x42, 0x87, 0x7c, 0x2d, 0x5a, 0x2c, 0xc9,
	0x43, 0x78, 0x00, 0x9b, 0xce, 0x34, 0x0c, 0xc9, 0x77, 0xe6, 0x59, 0x25, 0x60, 0xd9, 0x12, 0x8e,
	0x03, 0xe8, 0xa6, 0x58, 0x5a, 0xd1, 0x64, 0x69, 0x19, 0x8e, 0xde, 0xd6, 0x09, 0x7c, 0xd7, 0xd3,
	0x5e, 0xe0, 0x2b, 0xd1, 0x8e, 0xdf, 0x76, 0x81, 0x44, 0xe7, 0x37, 0x09, 0x3d, 0x87, 0x44, 0x87,
	0xf3, 0xe3, 0x00, 0x5f, 0x41, 0x67, 0x22, 0x43, 0xf2, 0xf5, 0x45, 0x72, 0x1b, 0x1b, 0x9c, 0x58,
	0x04, 0xa3, 0x3b, 0x91, 0x8e, 0xf6, 0x66, 0x24, 0xba, 0x7d, 0x63, 0xd0, 0xb4, 0x93, 0xc8, 0x3a,
	0x81, 0xed, 0x13, 0xd2, 0xd7, 0x8a, 0xc2, 0xe7, 0x5d, 0xae, 0x35, 0x84, 0x9d, 0x7c, 0xa1, 0xac,
	0xd5, 0x16, 0x7a, 0xa3, 0xa0, 0x3f, 0x06, 0xb4, 0x69, 0x1c, 0xcc, 0xe8, 0x99, 0xeb, 0x7e, 0xab,
	0x43, 0x9d, 0x4b, 0xfc, 0xb1, 0x86, 0xdc, 0x84, 0xda, 0x3d, 0xcd, 0x93, 0x5e, 0x8c, 0x1e, 0xa3,
	0x7c, 0x45, 0x4e, 0x48, 0x3a, 0x69, 0xc1, 0x24, 0x2a, 0xb4, 0x6e, 0xa3, 0xd4, 0xba, 0x03, 0xe8,
	0xa6, 0xcf, 0xe9, 0xc5, 0xc4, 0xcd, 0x58, 0x86, 0x71, 0x08, 0x98, 0x42, 0x1f, 0x17, 0xcd, 0xde,
	0x64, 0x71, 0x05, 0x53, 0x32, 0x45, 0x6b, 0xa5, 0x29, 0x60, 0x95, 0x29, 0xcc, 0xb2, 0x29, 0xaa,
	0x1b, 0xab, 0x6c, 0x95, 0x8d, 0xc7, 0xad, 0xd2, 0x7d, 0x9a, 0x55, 0x36, 0x9f, 0x6e, 0x95, 0xad,
	0x6a, 0xab, 0x44, 0x37, 0xa2, 0xa5, 0x9e, 0x2a, 0x81, 0xc9, 0x8d, 0x70, 0x14, 0xed, 0xc7, 0xb9,
	0x93, 0x9e, 0x7f, 0x19, 0x93, 0xdb, 0x4c, 0xe6, 0xa1, 0x92, 0xc9, 0x76, 0x96, 0x4c, 0xb6, 0x07,
	0xad, 0x2c, 0x12, 0xff, 0xc4, 0x27, 0x95, 0x01, 0xcb, 0x66, 0xdb, 0xad, 0x30, 0x9b, 0xf5, 0x16,
	0x3a, 0x99, 0x11, 0x8e, 0xa4, 0x96, 0x38, 0x80, 0x3a, 0x9f, 0x36, 0x37, 0xa6, 0x39, 0xc2, 0x8a,
	0x71, 0x1b, 0x0b, 0xac, 0x77, 0xb0, 0xb1, 0xf0, 0x10, 0xe7, 0x1e, 0x94, 0x66, 0x75, 0x55, 0x72,
	0x3a, 0xa4, 0x1f, 0x92, 0x41, 0x6f, 0x93, 0x9a, 0x04, 0xbe, 0xa2, 0xdc, 0x39, 0x19, 0x85, 0x73,
	0x12, 0xb0, 0x3e, 0x26, 0xa5, 0xe4, 0x2d, 0x25, 0x96, 0x48, 0x43, 0x1c, 0xc2, 0x9a, 0x2b, 0xb5,
	0x64, 0x4b, 0x98, 0xa3, 0x5e, 0x61, 0xb1, 0xc2, 0x4b, 0xd9, 0xac, 0xb3, 0x66, 0xb0, 0xc5, 0xd0,
	0x99, 0xa7, 0xf4, 0x33, 0x96, 0x3d, 0x2c, 0x2c, 0xfb, 0x7f, 0xf5, 0xb2, 0x6a, 0xb1, 0xee, 0xe8,
	0x7b, 0xfa, 0x7f, 0x74, 0x49, 0xe1, 0x2c, 0x6a, 0xd7, 0x53, 0x68, 0xbd, 0x77, 0xdd, 0x58, 0x87,
	0xbd, 0xe5, 0x43, 0x4a, 0x27, 0x52, 0x6f, 0x7f, 0x99, 0xcb, 0x6f, 0xde, 0xfa, 0x0b, 0xcf, 0xa1,
	0x9d, 0x9f, 0x65, 0xd8, 0x2f, 0x64, 0x54, 0xcc, 0xcb, 0x5e, 0xaf, 0xea, 0x0f, 0x34, 0xab, 0x77,
	0x05, 0x9d, 0xc2, 0x6c, 0xc4, 0x97, 0x2b, 0x0b, 0xfe, 0xc6, 0x2e, 0xcf, 0xc0, 0xcc, 0x4d, 0x50,
	0x7c, 0x51, 0x48, 0x58, 0x9e, 0xad, 0x8f, 0xec, 0xf1, 0x18, 0xcc, 0xeb, 0x89, 0x2b, 0x75, 0x52,
	0x6d, 0xf5, 0x17, 0xc1, 0xaf, 0xeb, 0xdc, 0x34, 0xf8, 0xfb, 0xe3, 0xcd, 0xcf, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x87, 0xe0, 0x60, 0x67, 0x9a, 0x08, 0x00, 0x00,
}