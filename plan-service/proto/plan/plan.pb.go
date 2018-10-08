// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/asciiu/gomo/plan-service/proto/plan/plan.proto

/*
Package plan is a generated protocol buffer package.

It is generated from these files:
	github.com/asciiu/gomo/plan-service/proto/plan/plan.proto

It has these top-level messages:
	NewPlanRequest
	UpdatePlanRequest
	GetUserPlanRequest
	GetUserPlansRequest
	DeletePlanRequest
	Plan
	PlanData
	PlanResponse
	PlansPageResponse
	PlansPage
*/
package plan

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import order "github.com/asciiu/gomo/plan-service/proto/order"

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
type NewPlanRequest struct {
	UserID                  string                   `protobuf:"bytes,1,opt,name=userID" json:"userID"`
	UserCurrencySymbol      string                   `protobuf:"bytes,2,opt,name=userCurrencySymbol" json:"userCurrencySymbol"`
	InitialTimestamp        string                   `protobuf:"bytes,3,opt,name=initialTimestamp" json:"initialTimestamp"`
	Title                   string                   `protobuf:"bytes,4,opt,name=title" json:"title"`
	PlanTemplateID          string                   `protobuf:"bytes,5,opt,name=planTemplateID" json:"planTemplateID"`
	Status                  string                   `protobuf:"bytes,6,opt,name=status" json:"status"`
	CloseOnComplete         bool                     `protobuf:"varint,7,opt,name=closeOnComplete" json:"closeOnComplete"`
	CommittedCurrencySymbol string                   `protobuf:"bytes,8,opt,name=committedCurrencySymbol" json:"committedCurrencySymbol"`
	CommittedCurrencyAmount float64                  `protobuf:"fixed64,9,opt,name=committedCurrencyAmount" json:"committedCurrencyAmount"`
	Orders                  []*order.NewOrderRequest `protobuf:"bytes,10,rep,name=orders" json:"orders"`
}

func (m *NewPlanRequest) Reset()                    { *m = NewPlanRequest{} }
func (m *NewPlanRequest) String() string            { return proto.CompactTextString(m) }
func (*NewPlanRequest) ProtoMessage()               {}
func (*NewPlanRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *NewPlanRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *NewPlanRequest) GetUserCurrencySymbol() string {
	if m != nil {
		return m.UserCurrencySymbol
	}
	return ""
}

func (m *NewPlanRequest) GetInitialTimestamp() string {
	if m != nil {
		return m.InitialTimestamp
	}
	return ""
}

func (m *NewPlanRequest) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *NewPlanRequest) GetPlanTemplateID() string {
	if m != nil {
		return m.PlanTemplateID
	}
	return ""
}

func (m *NewPlanRequest) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *NewPlanRequest) GetCloseOnComplete() bool {
	if m != nil {
		return m.CloseOnComplete
	}
	return false
}

func (m *NewPlanRequest) GetCommittedCurrencySymbol() string {
	if m != nil {
		return m.CommittedCurrencySymbol
	}
	return ""
}

func (m *NewPlanRequest) GetCommittedCurrencyAmount() float64 {
	if m != nil {
		return m.CommittedCurrencyAmount
	}
	return 0
}

func (m *NewPlanRequest) GetOrders() []*order.NewOrderRequest {
	if m != nil {
		return m.Orders
	}
	return nil
}

type UpdatePlanRequest struct {
	PlanID                  string                   `protobuf:"bytes,1,opt,name=planID" json:"planID"`
	UserID                  string                   `protobuf:"bytes,2,opt,name=userID" json:"userID"`
	UserCurrencySymbol      string                   `protobuf:"bytes,3,opt,name=userCurrencySymbol" json:"userCurrencySymbol"`
	InitialTimestamp        string                   `protobuf:"bytes,4,opt,name=initialTimestamp" json:"initialTimestamp"`
	Title                   string                   `protobuf:"bytes,5,opt,name=title" json:"title"`
	PlanTemplateID          string                   `protobuf:"bytes,6,opt,name=planTemplateID" json:"planTemplateID"`
	Status                  string                   `protobuf:"bytes,7,opt,name=status" json:"status"`
	CloseOnComplete         bool                     `protobuf:"varint,8,opt,name=closeOnComplete" json:"closeOnComplete"`
	CommittedCurrencySymbol string                   `protobuf:"bytes,9,opt,name=committedCurrencySymbol" json:"committedCurrencySymbol"`
	CommittedCurrencyAmount float64                  `protobuf:"fixed64,10,opt,name=committedCurrencyAmount" json:"committedCurrencyAmount"`
	Orders                  []*order.NewOrderRequest `protobuf:"bytes,11,rep,name=orders" json:"orders"`
}

func (m *UpdatePlanRequest) Reset()                    { *m = UpdatePlanRequest{} }
func (m *UpdatePlanRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdatePlanRequest) ProtoMessage()               {}
func (*UpdatePlanRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

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

func (m *UpdatePlanRequest) GetUserCurrencySymbol() string {
	if m != nil {
		return m.UserCurrencySymbol
	}
	return ""
}

func (m *UpdatePlanRequest) GetInitialTimestamp() string {
	if m != nil {
		return m.InitialTimestamp
	}
	return ""
}

func (m *UpdatePlanRequest) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *UpdatePlanRequest) GetPlanTemplateID() string {
	if m != nil {
		return m.PlanTemplateID
	}
	return ""
}

func (m *UpdatePlanRequest) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *UpdatePlanRequest) GetCloseOnComplete() bool {
	if m != nil {
		return m.CloseOnComplete
	}
	return false
}

func (m *UpdatePlanRequest) GetCommittedCurrencySymbol() string {
	if m != nil {
		return m.CommittedCurrencySymbol
	}
	return ""
}

func (m *UpdatePlanRequest) GetCommittedCurrencyAmount() float64 {
	if m != nil {
		return m.CommittedCurrencyAmount
	}
	return 0
}

func (m *UpdatePlanRequest) GetOrders() []*order.NewOrderRequest {
	if m != nil {
		return m.Orders
	}
	return nil
}

type GetUserPlanRequest struct {
	PlanID     string `protobuf:"bytes,1,opt,name=planID" json:"planID"`
	UserID     string `protobuf:"bytes,2,opt,name=userID" json:"userID"`
	PlanDepth  uint32 `protobuf:"varint,3,opt,name=planDepth" json:"planDepth"`
	PlanLength uint32 `protobuf:"varint,4,opt,name=planLength" json:"planLength"`
}

func (m *GetUserPlanRequest) Reset()                    { *m = GetUserPlanRequest{} }
func (m *GetUserPlanRequest) String() string            { return proto.CompactTextString(m) }
func (*GetUserPlanRequest) ProtoMessage()               {}
func (*GetUserPlanRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

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

func (m *GetUserPlanRequest) GetPlanDepth() uint32 {
	if m != nil {
		return m.PlanDepth
	}
	return 0
}

func (m *GetUserPlanRequest) GetPlanLength() uint32 {
	if m != nil {
		return m.PlanLength
	}
	return 0
}

type GetUserPlansRequest struct {
	UserID   string `protobuf:"bytes,1,opt,name=userID" json:"userID"`
	Status   string `protobuf:"bytes,2,opt,name=status" json:"status"`
	Page     uint32 `protobuf:"varint,3,opt,name=page" json:"page"`
	PageSize uint32 `protobuf:"varint,4,opt,name=pageSize" json:"pageSize"`
}

func (m *GetUserPlansRequest) Reset()                    { *m = GetUserPlansRequest{} }
func (m *GetUserPlansRequest) String() string            { return proto.CompactTextString(m) }
func (*GetUserPlansRequest) ProtoMessage()               {}
func (*GetUserPlansRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *GetUserPlansRequest) GetUserID() string {
	if m != nil {
		return m.UserID
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
func (*DeletePlanRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

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

// Responses
type Plan struct {
	PlanID                    string         `protobuf:"bytes,1,opt,name=planID" json:"planID"`
	PlanTemplateID            string         `protobuf:"bytes,2,opt,name=planTemplateID" json:"planTemplateID"`
	UserID                    string         `protobuf:"bytes,3,opt,name=userID" json:"userID"`
	UserPlanNumber            uint64         `protobuf:"varint,4,opt,name=userPlanNumber" json:"userPlanNumber"`
	Exchange                  string         `protobuf:"bytes,5,opt,name=exchange" json:"exchange"`
	Title                     string         `protobuf:"bytes,6,opt,name=title" json:"title"`
	TotalDepth                uint32         `protobuf:"varint,7,opt,name=totalDepth" json:"totalDepth"`
	UserCurrencySymbol        string         `protobuf:"bytes,8,opt,name=userCurrencySymbol" json:"userCurrencySymbol"`
	UserCurrencyBalanceAtInit float64        `protobuf:"fixed64,9,opt,name=userCurrencyBalanceAtInit" json:"userCurrencyBalanceAtInit"`
	CommittedCurrencySymbol   string         `protobuf:"bytes,10,opt,name=committedCurrencySymbol" json:"committedCurrencySymbol"`
	CommittedCurrencyAmount   float64        `protobuf:"fixed64,11,opt,name=committedCurrencyAmount" json:"committedCurrencyAmount"`
	InitialCurrencySymbol     string         `protobuf:"bytes,14,opt,name=initialCurrencySymbol" json:"initialCurrencySymbol"`
	InitialCurrencyBalance    float64        `protobuf:"fixed64,15,opt,name=initialCurrencyBalance" json:"initialCurrencyBalance"`
	InitialTimestamp          string         `protobuf:"bytes,16,opt,name=initialTimestamp" json:"initialTimestamp"`
	LastExecutedPlanDepth     uint32         `protobuf:"varint,17,opt,name=lastExecutedPlanDepth" json:"lastExecutedPlanDepth"`
	LastExecutedOrderID       string         `protobuf:"bytes,18,opt,name=lastExecutedOrderID" json:"lastExecutedOrderID"`
	Status                    string         `protobuf:"bytes,19,opt,name=status" json:"status"`
	CloseOnComplete           bool           `protobuf:"varint,20,opt,name=closeOnComplete" json:"closeOnComplete"`
	CreatedOn                 string         `protobuf:"bytes,21,opt,name=createdOn" json:"createdOn"`
	UpdatedOn                 string         `protobuf:"bytes,22,opt,name=updatedOn" json:"updatedOn"`
	Orders                    []*order.Order `protobuf:"bytes,23,rep,name=orders" json:"orders"`
}

func (m *Plan) Reset()                    { *m = Plan{} }
func (m *Plan) String() string            { return proto.CompactTextString(m) }
func (*Plan) ProtoMessage()               {}
func (*Plan) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

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

func (m *Plan) GetUserPlanNumber() uint64 {
	if m != nil {
		return m.UserPlanNumber
	}
	return 0
}

func (m *Plan) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *Plan) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Plan) GetTotalDepth() uint32 {
	if m != nil {
		return m.TotalDepth
	}
	return 0
}

func (m *Plan) GetUserCurrencySymbol() string {
	if m != nil {
		return m.UserCurrencySymbol
	}
	return ""
}

func (m *Plan) GetUserCurrencyBalanceAtInit() float64 {
	if m != nil {
		return m.UserCurrencyBalanceAtInit
	}
	return 0
}

func (m *Plan) GetCommittedCurrencySymbol() string {
	if m != nil {
		return m.CommittedCurrencySymbol
	}
	return ""
}

func (m *Plan) GetCommittedCurrencyAmount() float64 {
	if m != nil {
		return m.CommittedCurrencyAmount
	}
	return 0
}

func (m *Plan) GetInitialCurrencySymbol() string {
	if m != nil {
		return m.InitialCurrencySymbol
	}
	return ""
}

func (m *Plan) GetInitialCurrencyBalance() float64 {
	if m != nil {
		return m.InitialCurrencyBalance
	}
	return 0
}

func (m *Plan) GetInitialTimestamp() string {
	if m != nil {
		return m.InitialTimestamp
	}
	return ""
}

func (m *Plan) GetLastExecutedPlanDepth() uint32 {
	if m != nil {
		return m.LastExecutedPlanDepth
	}
	return 0
}

func (m *Plan) GetLastExecutedOrderID() string {
	if m != nil {
		return m.LastExecutedOrderID
	}
	return ""
}

func (m *Plan) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Plan) GetCloseOnComplete() bool {
	if m != nil {
		return m.CloseOnComplete
	}
	return false
}

func (m *Plan) GetCreatedOn() string {
	if m != nil {
		return m.CreatedOn
	}
	return ""
}

func (m *Plan) GetUpdatedOn() string {
	if m != nil {
		return m.UpdatedOn
	}
	return ""
}

func (m *Plan) GetOrders() []*order.Order {
	if m != nil {
		return m.Orders
	}
	return nil
}

type PlanData struct {
	Plan *Plan `protobuf:"bytes,1,opt,name=plan" json:"plan"`
}

func (m *PlanData) Reset()                    { *m = PlanData{} }
func (m *PlanData) String() string            { return proto.CompactTextString(m) }
func (*PlanData) ProtoMessage()               {}
func (*PlanData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

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
func (*PlanResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

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

type PlansPageResponse struct {
	Status  string     `protobuf:"bytes,1,opt,name=status" json:"status"`
	Message string     `protobuf:"bytes,2,opt,name=message" json:"message"`
	Data    *PlansPage `protobuf:"bytes,3,opt,name=data" json:"data"`
}

func (m *PlansPageResponse) Reset()                    { *m = PlansPageResponse{} }
func (m *PlansPageResponse) String() string            { return proto.CompactTextString(m) }
func (*PlansPageResponse) ProtoMessage()               {}
func (*PlansPageResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

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

type PlansPage struct {
	Page     uint32  `protobuf:"varint,1,opt,name=page" json:"page"`
	PageSize uint32  `protobuf:"varint,2,opt,name=pageSize" json:"pageSize"`
	Total    uint32  `protobuf:"varint,3,opt,name=total" json:"total"`
	Plans    []*Plan `protobuf:"bytes,4,rep,name=plans" json:"plans"`
}

func (m *PlansPage) Reset()                    { *m = PlansPage{} }
func (m *PlansPage) String() string            { return proto.CompactTextString(m) }
func (*PlansPage) ProtoMessage()               {}
func (*PlansPage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

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
	proto.RegisterType((*NewPlanRequest)(nil), "plan.NewPlanRequest")
	proto.RegisterType((*UpdatePlanRequest)(nil), "plan.UpdatePlanRequest")
	proto.RegisterType((*GetUserPlanRequest)(nil), "plan.GetUserPlanRequest")
	proto.RegisterType((*GetUserPlansRequest)(nil), "plan.GetUserPlansRequest")
	proto.RegisterType((*DeletePlanRequest)(nil), "plan.DeletePlanRequest")
	proto.RegisterType((*Plan)(nil), "plan.Plan")
	proto.RegisterType((*PlanData)(nil), "plan.PlanData")
	proto.RegisterType((*PlanResponse)(nil), "plan.PlanResponse")
	proto.RegisterType((*PlansPageResponse)(nil), "plan.PlansPageResponse")
	proto.RegisterType((*PlansPage)(nil), "plan.PlansPage")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for PlanService service

type PlanServiceClient interface {
	NewPlan(ctx context.Context, in *NewPlanRequest, opts ...client.CallOption) (*PlanResponse, error)
	GetUserPlan(ctx context.Context, in *GetUserPlanRequest, opts ...client.CallOption) (*PlanResponse, error)
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
		serviceName = "plan"
	}
	return &planServiceClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *planServiceClient) NewPlan(ctx context.Context, in *NewPlanRequest, opts ...client.CallOption) (*PlanResponse, error) {
	req := c.c.NewRequest(c.serviceName, "PlanService.NewPlan", in)
	out := new(PlanResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *planServiceClient) GetUserPlan(ctx context.Context, in *GetUserPlanRequest, opts ...client.CallOption) (*PlanResponse, error) {
	req := c.c.NewRequest(c.serviceName, "PlanService.GetUserPlan", in)
	out := new(PlanResponse)
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
	NewPlan(context.Context, *NewPlanRequest, *PlanResponse) error
	GetUserPlan(context.Context, *GetUserPlanRequest, *PlanResponse) error
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

func (h *PlanService) NewPlan(ctx context.Context, in *NewPlanRequest, out *PlanResponse) error {
	return h.PlanServiceHandler.NewPlan(ctx, in, out)
}

func (h *PlanService) GetUserPlan(ctx context.Context, in *GetUserPlanRequest, out *PlanResponse) error {
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

func init() {
	proto.RegisterFile("github.com/asciiu/gomo/plan-service/proto/plan/plan.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 883 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x56, 0xdb, 0x6e, 0xdb, 0x46,
	0x10, 0x2d, 0x25, 0x5a, 0x97, 0x91, 0x23, 0xc7, 0xeb, 0x1b, 0x63, 0x04, 0x86, 0xc0, 0x16, 0x81,
	0x10, 0xa0, 0x72, 0xe1, 0xa6, 0x45, 0x0b, 0xb7, 0x0f, 0xa9, 0x55, 0x14, 0x06, 0x0a, 0xc7, 0xa0,
	0x93, 0x0f, 0x58, 0x51, 0x03, 0x89, 0x05, 0x6f, 0xe5, 0x2e, 0x9b, 0xa4, 0x8f, 0xfd, 0x9a, 0xfe,
	0x47, 0x7e, 0xa4, 0x9f, 0x52, 0xec, 0x2c, 0x25, 0x52, 0x14, 0x69, 0x29, 0xf6, 0x8b, 0xc4, 0xb9,
	0xec, 0x0e, 0x67, 0x0f, 0xcf, 0xd9, 0x81, 0x1f, 0x67, 0x9e, 0x9c, 0xa7, 0x93, 0x91, 0x1b, 0x05,
	0xe7, 0x5c, 0xb8, 0x9e, 0x97, 0x9e, 0xcf, 0xa2, 0x20, 0x3a, 0x8f, 0x7d, 0x1e, 0x7e, 0x2d, 0x30,
	0xf9, 0xcb, 0x73, 0xf1, 0x3c, 0x4e, 0x22, 0xa9, 0x5d, 0xf4, 0x33, 0x22, 0x9b, 0x99, 0xea, 0xf9,
	0xf4, 0x72, 0xfb, 0x0d, 0xa2, 0x64, 0x8a, 0x89, 0xfe, 0xd5, 0x5b, 0xd8, 0xff, 0x36, 0xa1, 0x7f,
	0x83, 0xef, 0x6f, 0x7d, 0x1e, 0x3a, 0xf8, 0x67, 0x8a, 0x42, 0xb2, 0x63, 0x68, 0xa5, 0x02, 0x93,
	0xeb, 0xb1, 0x65, 0x0c, 0x8c, 0x61, 0xd7, 0xc9, 0x2c, 0x36, 0x02, 0xa6, 0x9e, 0xae, 0xd2, 0x24,
	0xc1, 0xd0, 0xfd, 0x78, 0xf7, 0x31, 0x98, 0x44, 0xbe, 0xd5, 0xa0, 0x9c, 0x8a, 0x08, 0x7b, 0x09,
	0x4f, 0xbd, 0xd0, 0x93, 0x1e, 0xf7, 0xdf, 0x7a, 0x01, 0x0a, 0xc9, 0x83, 0xd8, 0x6a, 0x52, 0xf6,
	0x9a, 0x9f, 0x1d, 0xc2, 0x8e, 0xf4, 0xa4, 0x8f, 0x96, 0x49, 0x09, 0xda, 0x60, 0x2f, 0xa0, 0xaf,
	0x9a, 0x78, 0x8b, 0x41, 0xec, 0x73, 0x89, 0xd7, 0x63, 0x6b, 0x87, 0xc2, 0x25, 0xaf, 0x7a, 0x63,
	0x21, 0xb9, 0x4c, 0x85, 0xd5, 0xd2, 0x6f, 0xac, 0x2d, 0x36, 0x84, 0x3d, 0xd7, 0x8f, 0x04, 0xbe,
	0x09, 0xaf, 0xa2, 0x20, 0xf6, 0x51, 0xa2, 0xd5, 0x1e, 0x18, 0xc3, 0x8e, 0x53, 0x76, 0xb3, 0x1f,
	0xe0, 0xc4, 0x8d, 0x82, 0xc0, 0x93, 0x12, 0xa7, 0xa5, 0x06, 0x3b, 0xb4, 0x65, 0x5d, 0xb8, 0x72,
	0xe5, 0xeb, 0x20, 0x4a, 0x43, 0x69, 0x75, 0x07, 0xc6, 0xd0, 0x70, 0xea, 0xc2, 0x6c, 0x04, 0x2d,
	0x42, 0x42, 0x58, 0x30, 0x68, 0x0e, 0x7b, 0x17, 0xc7, 0x23, 0x0d, 0xcc, 0x0d, 0xbe, 0x7f, 0xa3,
	0x1e, 0x32, 0x3c, 0x9c, 0x2c, 0xcb, 0xfe, 0xd4, 0x84, 0xfd, 0x77, 0xf1, 0x94, 0x4b, 0x2c, 0xa1,
	0xa5, 0x4e, 0x23, 0x47, 0x4b, 0x5b, 0x05, 0x14, 0x1b, 0x5b, 0xa0, 0xd8, 0xfc, 0x2c, 0x14, 0xcd,
	0x4d, 0x28, 0xee, 0xdc, 0x8f, 0x62, 0x6b, 0x03, 0x8a, 0xed, 0x4d, 0x28, 0x76, 0x3e, 0x1b, 0xc5,
	0xee, 0x83, 0x51, 0x84, 0x6d, 0x51, 0xec, 0x6d, 0x85, 0xe2, 0x3f, 0x06, 0xb0, 0xdf, 0x50, 0xbe,
	0x13, 0x98, 0x3c, 0x06, 0xc6, 0xe7, 0xd0, 0x55, 0x19, 0x63, 0x8c, 0xe5, 0x9c, 0xd0, 0x7b, 0xe2,
	0xe4, 0x0e, 0x76, 0x06, 0xa0, 0x8c, 0xdf, 0x31, 0x9c, 0xc9, 0x39, 0xc1, 0xf5, 0xc4, 0x29, 0x78,
	0xec, 0x14, 0x0e, 0x0a, 0xef, 0x20, 0x36, 0x31, 0x3f, 0x47, 0xa6, 0xb1, 0x82, 0x0c, 0x03, 0x33,
	0xe6, 0x33, 0xcc, 0xea, 0xd3, 0x33, 0x3b, 0x85, 0x8e, 0xfa, 0xbf, 0xf3, 0xfe, 0xc6, 0xac, 0xf0,
	0xd2, 0xb6, 0xaf, 0x60, 0x7f, 0x8c, 0x0a, 0xa9, 0x47, 0x74, 0x6e, 0xff, 0xd7, 0x02, 0x53, 0xad,
	0xaf, 0x5d, 0xb8, 0xfe, 0xbd, 0x35, 0xea, 0xbe, 0xb7, 0xac, 0x40, 0x73, 0xa5, 0xdb, 0x17, 0xd0,
	0x4f, 0xb3, 0x93, 0xb9, 0x49, 0x83, 0x09, 0x26, 0xd4, 0x87, 0xe9, 0x94, 0xbc, 0xaa, 0x53, 0xfc,
	0xe0, 0xce, 0x79, 0x38, 0x5b, 0x7c, 0xf0, 0x4b, 0x3b, 0x67, 0x42, 0xab, 0xc8, 0x84, 0x33, 0x00,
	0x19, 0x49, 0xee, 0x6b, 0xd4, 0xda, 0x1a, 0x96, 0xdc, 0x53, 0xc3, 0xcd, 0x4e, 0x2d, 0x37, 0x7f,
	0x82, 0x67, 0x45, 0xef, 0x2f, 0xdc, 0xe7, 0xa1, 0x8b, 0xaf, 0xe5, 0x75, 0xe8, 0x2d, 0xd4, 0xa7,
	0x3e, 0xe1, 0x3e, 0xb6, 0xc0, 0x83, 0xd9, 0xd2, 0xbb, 0x9f, 0x2d, 0xaf, 0xe0, 0x28, 0x53, 0x8d,
	0x52, 0xc5, 0x3e, 0x55, 0xac, 0x0e, 0xb2, 0xef, 0xe1, 0xb8, 0x14, 0xc8, 0x3a, 0xb1, 0xf6, 0xa8,
	0x5c, 0x4d, 0xb4, 0x52, 0xbb, 0x9e, 0xd6, 0x68, 0xd7, 0x2b, 0x38, 0xf2, 0xb9, 0x90, 0xbf, 0x7e,
	0x40, 0x37, 0x95, 0x38, 0xbd, 0x5d, 0x92, 0x6b, 0x9f, 0x60, 0xaa, 0x0e, 0xb2, 0x6f, 0xe0, 0xa0,
	0x18, 0x20, 0xc6, 0x5f, 0x8f, 0x2d, 0x46, 0x45, 0xaa, 0x42, 0x05, 0x2e, 0x1d, 0x6c, 0x52, 0xb9,
	0xc3, 0x6a, 0x95, 0x7b, 0x0e, 0x5d, 0x37, 0x41, 0xae, 0xf6, 0x0c, 0xad, 0x23, 0xda, 0x24, 0x77,
	0xa8, 0x68, 0x4a, 0x97, 0x84, 0x8a, 0x1e, 0xeb, 0xe8, 0xd2, 0xc1, 0xbe, 0x5a, 0xaa, 0xd5, 0x09,
	0xa9, 0xd5, 0x6e, 0xa6, 0x56, 0x5a, 0xaa, 0x16, 0x1a, 0xf5, 0x12, 0x3a, 0xd4, 0x22, 0x97, 0x9c,
	0x9d, 0x01, 0x4d, 0x19, 0xc4, 0xb1, 0xde, 0x05, 0x8c, 0x68, 0xfc, 0x20, 0xfe, 0x92, 0xdf, 0x9e,
	0xc2, 0xae, 0x66, 0xb3, 0x88, 0xa3, 0x50, 0x60, 0xa1, 0x3f, 0x63, 0xa5, 0x3f, 0x0b, 0xda, 0x01,
	0x0a, 0xa1, 0xe4, 0x42, 0xd3, 0x71, 0x61, 0x32, 0x1b, 0xcc, 0x29, 0x97, 0x9c, 0x58, 0xd8, 0xbb,
	0xe8, 0xe7, 0x15, 0x54, 0x7d, 0x87, 0x62, 0xf6, 0x1f, 0xb0, 0x4f, 0x4a, 0x75, 0xcb, 0x67, 0xf8,
	0x88, 0x52, 0x5f, 0xae, 0x94, 0xda, 0xcb, 0x4b, 0xe9, 0x8d, 0x75, 0x2d, 0x01, 0xdd, 0xa5, 0x6b,
	0x29, 0x71, 0x46, 0x8d, 0xc4, 0x35, 0x56, 0x25, 0x8e, 0x88, 0xaf, 0x08, 0x9d, 0x69, 0xa2, 0x36,
	0xd8, 0x00, 0x76, 0x54, 0x29, 0x61, 0x99, 0x74, 0xea, 0xc5, 0x53, 0xd4, 0x81, 0x8b, 0x4f, 0x0d,
	0xe8, 0x29, 0xfb, 0x4e, 0xcf, 0x6b, 0xec, 0x3b, 0x68, 0x67, 0x63, 0x19, 0x3b, 0xd4, 0xd9, 0xab,
	0x53, 0xda, 0x29, 0x2b, 0xec, 0x91, 0x1d, 0x88, 0xfd, 0x05, 0xfb, 0x19, 0x7a, 0x05, 0x61, 0x67,
	0x96, 0x4e, 0x5a, 0xbf, 0x6f, 0x6a, 0x96, 0x8f, 0x61, 0xb7, 0x78, 0x2f, 0xb0, 0x67, 0x6b, 0xeb,
	0x17, 0x77, 0xc5, 0xe9, 0x49, 0xf9, 0xf0, 0xf2, 0x5d, 0x2e, 0x01, 0x72, 0x99, 0x67, 0x59, 0xe2,
	0x9a, 0xf0, 0xd7, 0xbc, 0xc2, 0x25, 0x40, 0x3e, 0xe4, 0x2c, 0x16, 0xaf, 0x8d, 0x3d, 0xd5, 0x8b,
	0x27, 0x2d, 0x1a, 0x6a, 0xbf, 0xfd, 0x3f, 0x00, 0x00, 0xff, 0xff, 0xab, 0x17, 0xf5, 0xd9, 0x54,
	0x0b, 0x00, 0x00,
}
