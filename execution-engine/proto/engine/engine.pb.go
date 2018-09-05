// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/asciiu/gomo/execution-engine/proto/engine/engine.proto

/*
Package engine is a generated protocol buffer package.

It is generated from these files:
	github.com/asciiu/gomo/execution-engine/proto/engine/engine.proto

It has these top-level messages:
	ActiveRequest
	KillRequest
	KillUserRequest
	NewPlanRequest
	Order
	Trigger
	Plan
	PlanList
	PlanResponse
*/
package engine

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
type ActiveRequest struct {
	Exchange   string `protobuf:"bytes,1,opt,name=exchange" json:"exchange"`
	MarketName string `protobuf:"bytes,2,opt,name=marketName" json:"marketName"`
}

func (m *ActiveRequest) Reset()                    { *m = ActiveRequest{} }
func (m *ActiveRequest) String() string            { return proto.CompactTextString(m) }
func (*ActiveRequest) ProtoMessage()               {}
func (*ActiveRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ActiveRequest) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *ActiveRequest) GetMarketName() string {
	if m != nil {
		return m.MarketName
	}
	return ""
}

type KillRequest struct {
	PlanID string `protobuf:"bytes,1,opt,name=planID" json:"planID"`
}

func (m *KillRequest) Reset()                    { *m = KillRequest{} }
func (m *KillRequest) String() string            { return proto.CompactTextString(m) }
func (*KillRequest) ProtoMessage()               {}
func (*KillRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *KillRequest) GetPlanID() string {
	if m != nil {
		return m.PlanID
	}
	return ""
}

type KillUserRequest struct {
	UserID string `protobuf:"bytes,1,opt,name=userID" json:"userID"`
}

func (m *KillUserRequest) Reset()                    { *m = KillUserRequest{} }
func (m *KillUserRequest) String() string            { return proto.CompactTextString(m) }
func (*KillUserRequest) ProtoMessage()               {}
func (*KillUserRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *KillUserRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

type NewPlanRequest struct {
	PlanID                string   `protobuf:"bytes,1,opt,name=planID" json:"planID"`
	UserID                string   `protobuf:"bytes,2,opt,name=userID" json:"userID"`
	ActiveCurrencySymbol  string   `protobuf:"bytes,3,opt,name=activeCurrencySymbol" json:"activeCurrencySymbol"`
	ActiveCurrencyBalance float64  `protobuf:"fixed64,4,opt,name=activeCurrencyBalance" json:"activeCurrencyBalance"`
	CloseOnComplete       bool     `protobuf:"varint,5,opt,name=closeOnComplete" json:"closeOnComplete"`
	Orders                []*Order `protobuf:"bytes,6,rep,name=orders" json:"orders"`
}

func (m *NewPlanRequest) Reset()                    { *m = NewPlanRequest{} }
func (m *NewPlanRequest) String() string            { return proto.CompactTextString(m) }
func (*NewPlanRequest) ProtoMessage()               {}
func (*NewPlanRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *NewPlanRequest) GetPlanID() string {
	if m != nil {
		return m.PlanID
	}
	return ""
}

func (m *NewPlanRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *NewPlanRequest) GetActiveCurrencySymbol() string {
	if m != nil {
		return m.ActiveCurrencySymbol
	}
	return ""
}

func (m *NewPlanRequest) GetActiveCurrencyBalance() float64 {
	if m != nil {
		return m.ActiveCurrencyBalance
	}
	return 0
}

func (m *NewPlanRequest) GetCloseOnComplete() bool {
	if m != nil {
		return m.CloseOnComplete
	}
	return false
}

func (m *NewPlanRequest) GetOrders() []*Order {
	if m != nil {
		return m.Orders
	}
	return nil
}

type Order struct {
	OrderID     string     `protobuf:"bytes,1,opt,name=orderID" json:"orderID"`
	Exchange    string     `protobuf:"bytes,2,opt,name=exchange" json:"exchange"`
	MarketName  string     `protobuf:"bytes,3,opt,name=marketName" json:"marketName"`
	Side        string     `protobuf:"bytes,4,opt,name=side" json:"side"`
	LimitPrice  float64    `protobuf:"fixed64,5,opt,name=limitPrice" json:"limitPrice"`
	OrderType   string     `protobuf:"bytes,6,opt,name=orderType" json:"orderType"`
	OrderStatus string     `protobuf:"bytes,7,opt,name=orderStatus" json:"orderStatus"`
	AccountID   string     `protobuf:"bytes,8,opt,name=accountID" json:"accountID"`
	KeyPublic   string     `protobuf:"bytes,9,opt,name=keyPublic" json:"keyPublic"`
	KeySecret   string     `protobuf:"bytes,10,opt,name=keySecret" json:"keySecret"`
	Triggers    []*Trigger `protobuf:"bytes,11,rep,name=triggers" json:"triggers"`
}

func (m *Order) Reset()                    { *m = Order{} }
func (m *Order) String() string            { return proto.CompactTextString(m) }
func (*Order) ProtoMessage()               {}
func (*Order) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

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
	TriggerID string   `protobuf:"bytes,1,opt,name=triggerID" json:"triggerID"`
	OrderID   string   `protobuf:"bytes,2,opt,name=orderID" json:"orderID"`
	Name      string   `protobuf:"bytes,5,opt,name=name" json:"name"`
	Code      string   `protobuf:"bytes,6,opt,name=code" json:"code"`
	Triggered bool     `protobuf:"varint,7,opt,name=triggered" json:"triggered"`
	Actions   []string `protobuf:"bytes,8,rep,name=actions" json:"actions"`
}

func (m *Trigger) Reset()                    { *m = Trigger{} }
func (m *Trigger) String() string            { return proto.CompactTextString(m) }
func (*Trigger) ProtoMessage()               {}
func (*Trigger) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

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

// Responses
type Plan struct {
	PlanID string `protobuf:"bytes,1,opt,name=planID" json:"planID"`
}

func (m *Plan) Reset()                    { *m = Plan{} }
func (m *Plan) String() string            { return proto.CompactTextString(m) }
func (*Plan) ProtoMessage()               {}
func (*Plan) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *Plan) GetPlanID() string {
	if m != nil {
		return m.PlanID
	}
	return ""
}

type PlanList struct {
	Plans []*Plan `protobuf:"bytes,1,rep,name=plans" json:"plans"`
}

func (m *PlanList) Reset()                    { *m = PlanList{} }
func (m *PlanList) String() string            { return proto.CompactTextString(m) }
func (*PlanList) ProtoMessage()               {}
func (*PlanList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *PlanList) GetPlans() []*Plan {
	if m != nil {
		return m.Plans
	}
	return nil
}

type PlanResponse struct {
	Status  string    `protobuf:"bytes,1,opt,name=status" json:"status"`
	Message string    `protobuf:"bytes,2,opt,name=message" json:"message"`
	Data    *PlanList `protobuf:"bytes,3,opt,name=data" json:"data"`
}

func (m *PlanResponse) Reset()                    { *m = PlanResponse{} }
func (m *PlanResponse) String() string            { return proto.CompactTextString(m) }
func (*PlanResponse) ProtoMessage()               {}
func (*PlanResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

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

func (m *PlanResponse) GetData() *PlanList {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*ActiveRequest)(nil), "engine.ActiveRequest")
	proto.RegisterType((*KillRequest)(nil), "engine.KillRequest")
	proto.RegisterType((*KillUserRequest)(nil), "engine.KillUserRequest")
	proto.RegisterType((*NewPlanRequest)(nil), "engine.NewPlanRequest")
	proto.RegisterType((*Order)(nil), "engine.Order")
	proto.RegisterType((*Trigger)(nil), "engine.Trigger")
	proto.RegisterType((*Plan)(nil), "engine.Plan")
	proto.RegisterType((*PlanList)(nil), "engine.PlanList")
	proto.RegisterType((*PlanResponse)(nil), "engine.PlanResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for ExecutionEngine service

type ExecutionEngineClient interface {
	AddPlan(ctx context.Context, in *NewPlanRequest, opts ...client.CallOption) (*PlanResponse, error)
	GetActivePlans(ctx context.Context, in *ActiveRequest, opts ...client.CallOption) (*PlanResponse, error)
	KillPlan(ctx context.Context, in *KillRequest, opts ...client.CallOption) (*PlanResponse, error)
	KillUserPlans(ctx context.Context, in *KillUserRequest, opts ...client.CallOption) (*PlanResponse, error)
}

type executionEngineClient struct {
	c           client.Client
	serviceName string
}

func NewExecutionEngineClient(serviceName string, c client.Client) ExecutionEngineClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "engine"
	}
	return &executionEngineClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *executionEngineClient) AddPlan(ctx context.Context, in *NewPlanRequest, opts ...client.CallOption) (*PlanResponse, error) {
	req := c.c.NewRequest(c.serviceName, "ExecutionEngine.AddPlan", in)
	out := new(PlanResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executionEngineClient) GetActivePlans(ctx context.Context, in *ActiveRequest, opts ...client.CallOption) (*PlanResponse, error) {
	req := c.c.NewRequest(c.serviceName, "ExecutionEngine.GetActivePlans", in)
	out := new(PlanResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executionEngineClient) KillPlan(ctx context.Context, in *KillRequest, opts ...client.CallOption) (*PlanResponse, error) {
	req := c.c.NewRequest(c.serviceName, "ExecutionEngine.KillPlan", in)
	out := new(PlanResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executionEngineClient) KillUserPlans(ctx context.Context, in *KillUserRequest, opts ...client.CallOption) (*PlanResponse, error) {
	req := c.c.NewRequest(c.serviceName, "ExecutionEngine.KillUserPlans", in)
	out := new(PlanResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ExecutionEngine service

type ExecutionEngineHandler interface {
	AddPlan(context.Context, *NewPlanRequest, *PlanResponse) error
	GetActivePlans(context.Context, *ActiveRequest, *PlanResponse) error
	KillPlan(context.Context, *KillRequest, *PlanResponse) error
	KillUserPlans(context.Context, *KillUserRequest, *PlanResponse) error
}

func RegisterExecutionEngineHandler(s server.Server, hdlr ExecutionEngineHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&ExecutionEngine{hdlr}, opts...))
}

type ExecutionEngine struct {
	ExecutionEngineHandler
}

func (h *ExecutionEngine) AddPlan(ctx context.Context, in *NewPlanRequest, out *PlanResponse) error {
	return h.ExecutionEngineHandler.AddPlan(ctx, in, out)
}

func (h *ExecutionEngine) GetActivePlans(ctx context.Context, in *ActiveRequest, out *PlanResponse) error {
	return h.ExecutionEngineHandler.GetActivePlans(ctx, in, out)
}

func (h *ExecutionEngine) KillPlan(ctx context.Context, in *KillRequest, out *PlanResponse) error {
	return h.ExecutionEngineHandler.KillPlan(ctx, in, out)
}

func (h *ExecutionEngine) KillUserPlans(ctx context.Context, in *KillUserRequest, out *PlanResponse) error {
	return h.ExecutionEngineHandler.KillUserPlans(ctx, in, out)
}

func init() {
	proto.RegisterFile("github.com/asciiu/gomo/execution-engine/proto/engine/engine.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 650 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x54, 0x6d, 0x6e, 0xdb, 0x38,
	0x10, 0x5d, 0xf9, 0x53, 0x1e, 0x27, 0xf1, 0x82, 0x9b, 0xec, 0x0a, 0xc1, 0x22, 0x30, 0x84, 0x06,
	0x70, 0x51, 0xd4, 0x06, 0xdc, 0x02, 0x45, 0x7f, 0xb5, 0x69, 0x12, 0x14, 0x41, 0x8a, 0x24, 0x50,
	0xd2, 0x03, 0xd0, 0xd4, 0xd4, 0x21, 0x22, 0x91, 0xae, 0x48, 0xb5, 0xf1, 0x0d, 0x7a, 0x8d, 0xde,
	0xaa, 0x77, 0xe9, 0x9f, 0x82, 0xa4, 0x24, 0xcb, 0x41, 0xdc, 0xfe, 0x92, 0xe6, 0xbd, 0x99, 0x47,
	0xce, 0x23, 0x39, 0x70, 0x34, 0xe7, 0xfa, 0x36, 0x9f, 0x8d, 0x99, 0x4c, 0x27, 0x54, 0x31, 0xce,
	0xf3, 0xc9, 0x5c, 0xa6, 0x72, 0x82, 0xf7, 0xc8, 0x72, 0xcd, 0xa5, 0x78, 0x8e, 0x62, 0xce, 0x05,
	0x4e, 0x16, 0x99, 0xd4, 0x72, 0x52, 0x04, 0xee, 0x33, 0xb6, 0x18, 0xe9, 0xb8, 0x28, 0x3c, 0x87,
	0xed, 0x23, 0xa6, 0xf9, 0x17, 0x8c, 0xf0, 0x73, 0x8e, 0x4a, 0x93, 0x7d, 0xf0, 0xf1, 0x9e, 0xdd,
	0x52, 0x31, 0xc7, 0xc0, 0x1b, 0x7a, 0xa3, 0x5e, 0x54, 0xc5, 0xe4, 0x00, 0x20, 0xa5, 0xd9, 0x1d,
	0xea, 0x0b, 0x9a, 0x62, 0xd0, 0xb0, 0x6c, 0x0d, 0x09, 0x0f, 0xa1, 0x7f, 0xce, 0x93, 0xa4, 0x94,
	0xfa, 0x17, 0x3a, 0x8b, 0x84, 0x8a, 0xb3, 0x93, 0x42, 0xa8, 0x88, 0xc2, 0xa7, 0x30, 0x30, 0x69,
	0x1f, 0x15, 0x66, 0xb5, 0xd4, 0x5c, 0x61, 0xb6, 0x4a, 0x75, 0x51, 0xf8, 0xd3, 0x83, 0x9d, 0x0b,
	0xfc, 0x7a, 0x95, 0x50, 0xf1, 0x07, 0xd5, 0x9a, 0x44, 0xa3, 0x2e, 0x41, 0xa6, 0xb0, 0x4b, 0x6d,
	0x87, 0xc7, 0x79, 0x96, 0xa1, 0x60, 0xcb, 0xeb, 0x65, 0x3a, 0x93, 0x49, 0xd0, 0xb4, 0x59, 0x8f,
	0x72, 0xe4, 0x25, 0xec, 0xad, 0xe3, 0xef, 0x68, 0x42, 0x05, 0xc3, 0xa0, 0x35, 0xf4, 0x46, 0x5e,
	0xf4, 0x38, 0x49, 0x46, 0x30, 0x60, 0x89, 0x54, 0x78, 0x29, 0x8e, 0x65, 0xba, 0x48, 0x50, 0x63,
	0xd0, 0x1e, 0x7a, 0x23, 0x3f, 0x7a, 0x08, 0x93, 0x43, 0xe8, 0xc8, 0x2c, 0xc6, 0x4c, 0x05, 0x9d,
	0x61, 0x73, 0xd4, 0x9f, 0x6e, 0x8f, 0x8b, 0xc3, 0xb9, 0x34, 0x68, 0x54, 0x90, 0xe1, 0x8f, 0x06,
	0xb4, 0x2d, 0x42, 0x02, 0xe8, 0x5a, 0xac, 0xea, 0xba, 0x0c, 0xd7, 0xce, 0xab, 0xf1, 0xdb, 0xf3,
	0x6a, 0x3e, 0x3c, 0x2f, 0x42, 0xa0, 0xa5, 0x78, 0xec, 0xba, 0xea, 0x45, 0xf6, 0xdf, 0xd4, 0x24,
	0x3c, 0xe5, 0xfa, 0x2a, 0xe3, 0xcc, 0xed, 0xdf, 0x8b, 0x6a, 0x08, 0xf9, 0x1f, 0x7a, 0x76, 0xe9,
	0x9b, 0xe5, 0x02, 0x83, 0x8e, 0x2d, 0x5c, 0x01, 0x64, 0x08, 0x7d, 0x1b, 0x5c, 0x6b, 0xaa, 0x73,
	0x15, 0x74, 0x2d, 0x5f, 0x87, 0x4c, 0x3d, 0x65, 0x4c, 0xe6, 0x42, 0x9f, 0x9d, 0x04, 0xbe, 0xab,
	0xaf, 0x00, 0xc3, 0xde, 0xe1, 0xf2, 0x2a, 0x9f, 0x25, 0x9c, 0x05, 0x3d, 0xc7, 0x56, 0x40, 0xc1,
	0x5e, 0x23, 0xcb, 0x50, 0x07, 0x50, 0xb1, 0x0e, 0x20, 0xcf, 0xc0, 0xd7, 0x19, 0x9f, 0xcf, 0x8d,
	0xad, 0x7d, 0x6b, 0xeb, 0xa0, 0xb4, 0xf5, 0xc6, 0xe1, 0x51, 0x95, 0x10, 0x7e, 0xf7, 0xa0, 0x5b,
	0xa0, 0x46, 0xb6, 0xc0, 0x2b, 0x7b, 0x57, 0x40, 0xdd, 0xfa, 0xc6, 0xba, 0xf5, 0x04, 0x5a, 0xc2,
	0x18, 0xdb, 0x76, 0xf6, 0x89, 0xc2, 0x52, 0x26, 0xe3, 0xd2, 0x19, 0xfb, 0x5f, 0xd3, 0xc7, 0xd8,
	0x5a, 0xe2, 0x47, 0x2b, 0xc0, 0xe8, 0x9b, 0xeb, 0x24, 0x85, 0x0a, 0xfc, 0x61, 0xd3, 0xe8, 0x17,
	0x61, 0x78, 0x00, 0x2d, 0x73, 0xf1, 0x37, 0xbe, 0xa3, 0x31, 0xf8, 0x86, 0xff, 0xc0, 0x95, 0x26,
	0x21, 0xb4, 0x0d, 0xaa, 0x02, 0xcf, 0x76, 0xbe, 0x55, 0x76, 0x6e, 0x5f, 0x8e, 0xa3, 0xc2, 0x4f,
	0xb0, 0xe5, 0x1e, 0x92, 0x5a, 0x48, 0xa1, 0xd0, 0xe8, 0x2a, 0x77, 0x4e, 0x85, 0xae, 0x8b, 0xcc,
	0x8e, 0x52, 0x54, 0x8a, 0x56, 0x37, 0xaa, 0x0c, 0xc9, 0x13, 0x68, 0xc5, 0x54, 0x53, 0x7b, 0x95,
	0xfa, 0xd3, 0xbf, 0xeb, 0x8b, 0x98, 0x5d, 0x44, 0x96, 0x9d, 0x7e, 0x6b, 0xc0, 0xe0, 0xb4, 0x1c,
	0x45, 0xa7, 0x36, 0x85, 0xbc, 0x86, 0xee, 0x51, 0x1c, 0xbb, 0x76, 0xca, 0xb2, 0xf5, 0x87, 0xbd,
	0xbf, 0xbb, 0xb6, 0xe7, 0x62, 0x93, 0xe1, 0x5f, 0xe4, 0x0d, 0xec, 0xbc, 0x47, 0xed, 0xa6, 0x94,
	0xa1, 0x14, 0xd9, 0x2b, 0x33, 0xd7, 0x46, 0xd7, 0x46, 0x81, 0x57, 0xe0, 0x9b, 0x79, 0x63, 0x17,
	0xff, 0xa7, 0xcc, 0xa9, 0x0d, 0xaa, 0x8d, 0x85, 0x6f, 0x61, 0xbb, 0x1c, 0x54, 0x6e, 0xe1, 0xff,
	0xea, 0xd5, 0xb5, 0xf9, 0xb5, 0x49, 0x61, 0xd6, 0xb1, 0xd3, 0xf6, 0xc5, 0xaf, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x3b, 0xef, 0x58, 0x56, 0xb2, 0x05, 0x00, 0x00,
}
