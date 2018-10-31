// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/notification/notification.proto

/*
Package ios is a generated protocol buffer package.

It is generated from these files:
	proto/notification/notification.proto

It has these top-level messages:
	EmailRequest
	ActivityRequest
	RecentActivityRequest
	ActivityCountRequest
	UpdateActivityRequest
	EmailResponse
	Activity
	UserActivityPage
	ActivityPagedResponse
	ActivityData
	ActivityResponse
	ActivityList
	ActivityListResponse
	ActivityCount
	ActivityCountResponse
*/
package ios

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
type EmailRequest struct {
	EmailSender    string `protobuf:"bytes,1,opt,name=emailSender" json:"emailSender"`
	EmailRecipient string `protobuf:"bytes,2,opt,name=emailRecipient" json:"emailRecipient"`
	Subject        string `protobuf:"bytes,3,opt,name=subject" json:"subject"`
	HtmlBody       string `protobuf:"bytes,4,opt,name=htmlBody" json:"htmlBody"`
	TextBody       string `protobuf:"bytes,5,opt,name=textBody" json:"textBody"`
}

func (m *EmailRequest) Reset()                    { *m = EmailRequest{} }
func (m *EmailRequest) String() string            { return proto.CompactTextString(m) }
func (*EmailRequest) ProtoMessage()               {}
func (*EmailRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *EmailRequest) GetEmailSender() string {
	if m != nil {
		return m.EmailSender
	}
	return ""
}

func (m *EmailRequest) GetEmailRecipient() string {
	if m != nil {
		return m.EmailRecipient
	}
	return ""
}

func (m *EmailRequest) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *EmailRequest) GetHtmlBody() string {
	if m != nil {
		return m.HtmlBody
	}
	return ""
}

func (m *EmailRequest) GetTextBody() string {
	if m != nil {
		return m.TextBody
	}
	return ""
}

type ActivityRequest struct {
	UserID   string `protobuf:"bytes,1,opt,name=userID" json:"userID"`
	ObjectID string `protobuf:"bytes,2,opt,name=objectID" json:"objectID"`
	Page     uint32 `protobuf:"varint,3,opt,name=page" json:"page"`
	PageSize uint32 `protobuf:"varint,4,opt,name=pageSize" json:"pageSize"`
}

func (m *ActivityRequest) Reset()                    { *m = ActivityRequest{} }
func (m *ActivityRequest) String() string            { return proto.CompactTextString(m) }
func (*ActivityRequest) ProtoMessage()               {}
func (*ActivityRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ActivityRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *ActivityRequest) GetObjectID() string {
	if m != nil {
		return m.ObjectID
	}
	return ""
}

func (m *ActivityRequest) GetPage() uint32 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *ActivityRequest) GetPageSize() uint32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

type RecentActivityRequest struct {
	ObjectID string `protobuf:"bytes,1,opt,name=objectID" json:"objectID"`
	Count    uint32 `protobuf:"varint,2,opt,name=count" json:"count"`
}

func (m *RecentActivityRequest) Reset()                    { *m = RecentActivityRequest{} }
func (m *RecentActivityRequest) String() string            { return proto.CompactTextString(m) }
func (*RecentActivityRequest) ProtoMessage()               {}
func (*RecentActivityRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *RecentActivityRequest) GetObjectID() string {
	if m != nil {
		return m.ObjectID
	}
	return ""
}

func (m *RecentActivityRequest) GetCount() uint32 {
	if m != nil {
		return m.Count
	}
	return 0
}

type ActivityCountRequest struct {
	ObjectID string `protobuf:"bytes,1,opt,name=objectID" json:"objectID"`
}

func (m *ActivityCountRequest) Reset()                    { *m = ActivityCountRequest{} }
func (m *ActivityCountRequest) String() string            { return proto.CompactTextString(m) }
func (*ActivityCountRequest) ProtoMessage()               {}
func (*ActivityCountRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ActivityCountRequest) GetObjectID() string {
	if m != nil {
		return m.ObjectID
	}
	return ""
}

type UpdateActivityRequest struct {
	ActivityID string `protobuf:"bytes,1,opt,name=activityID" json:"activityID"`
	SeenAt     string `protobuf:"bytes,2,opt,name=seenAt" json:"seenAt"`
	ClickedAt  string `protobuf:"bytes,3,opt,name=clickedAt" json:"clickedAt"`
}

func (m *UpdateActivityRequest) Reset()                    { *m = UpdateActivityRequest{} }
func (m *UpdateActivityRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdateActivityRequest) ProtoMessage()               {}
func (*UpdateActivityRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *UpdateActivityRequest) GetActivityID() string {
	if m != nil {
		return m.ActivityID
	}
	return ""
}

func (m *UpdateActivityRequest) GetSeenAt() string {
	if m != nil {
		return m.SeenAt
	}
	return ""
}

func (m *UpdateActivityRequest) GetClickedAt() string {
	if m != nil {
		return m.ClickedAt
	}
	return ""
}

// Responses
type EmailResponse struct {
	Status  string `protobuf:"bytes,1,opt,name=status" json:"status"`
	Message string `protobuf:"bytes,2,opt,name=message" json:"message"`
}

func (m *EmailResponse) Reset()                    { *m = EmailResponse{} }
func (m *EmailResponse) String() string            { return proto.CompactTextString(m) }
func (*EmailResponse) ProtoMessage()               {}
func (*EmailResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *EmailResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *EmailResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type Activity struct {
	ActivityID  string `protobuf:"bytes,1,opt,name=activityID" json:"activityID"`
	UserID      string `protobuf:"bytes,3,opt,name=userID" json:"userID"`
	Type        string `protobuf:"bytes,2,opt,name=type" json:"type"`
	ObjectID    string `protobuf:"bytes,4,opt,name=objectID" json:"objectID"`
	Title       string `protobuf:"bytes,5,opt,name=title" json:"title"`
	Subtitle    string `protobuf:"bytes,6,opt,name=subtitle" json:"subtitle"`
	Description string `protobuf:"bytes,7,opt,name=description" json:"description"`
	Details     string `protobuf:"bytes,8,opt,name=details" json:"details"`
	Timestamp   string `protobuf:"bytes,9,opt,name=timestamp" json:"timestamp"`
	ClickedAt   string `protobuf:"bytes,10,opt,name=clickedAt" json:"clickedAt"`
	SeenAt      string `protobuf:"bytes,11,opt,name=seenAt" json:"seenAt"`
}

func (m *Activity) Reset()                    { *m = Activity{} }
func (m *Activity) String() string            { return proto.CompactTextString(m) }
func (*Activity) ProtoMessage()               {}
func (*Activity) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *Activity) GetActivityID() string {
	if m != nil {
		return m.ActivityID
	}
	return ""
}

func (m *Activity) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *Activity) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Activity) GetObjectID() string {
	if m != nil {
		return m.ObjectID
	}
	return ""
}

func (m *Activity) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Activity) GetSubtitle() string {
	if m != nil {
		return m.Subtitle
	}
	return ""
}

func (m *Activity) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Activity) GetDetails() string {
	if m != nil {
		return m.Details
	}
	return ""
}

func (m *Activity) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

func (m *Activity) GetClickedAt() string {
	if m != nil {
		return m.ClickedAt
	}
	return ""
}

func (m *Activity) GetSeenAt() string {
	if m != nil {
		return m.SeenAt
	}
	return ""
}

type UserActivityPage struct {
	Page     uint32      `protobuf:"varint,1,opt,name=page" json:"page"`
	PageSize uint32      `protobuf:"varint,2,opt,name=pageSize" json:"pageSize"`
	Total    uint32      `protobuf:"varint,3,opt,name=total" json:"total"`
	Activity []*Activity `protobuf:"bytes,4,rep,name=activity" json:"activity"`
}

func (m *UserActivityPage) Reset()                    { *m = UserActivityPage{} }
func (m *UserActivityPage) String() string            { return proto.CompactTextString(m) }
func (*UserActivityPage) ProtoMessage()               {}
func (*UserActivityPage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *UserActivityPage) GetPage() uint32 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *UserActivityPage) GetPageSize() uint32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *UserActivityPage) GetTotal() uint32 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *UserActivityPage) GetActivity() []*Activity {
	if m != nil {
		return m.Activity
	}
	return nil
}

type ActivityPagedResponse struct {
	Status  string            `protobuf:"bytes,1,opt,name=status" json:"status"`
	Message string            `protobuf:"bytes,2,opt,name=message" json:"message"`
	Data    *UserActivityPage `protobuf:"bytes,3,opt,name=data" json:"data"`
}

func (m *ActivityPagedResponse) Reset()                    { *m = ActivityPagedResponse{} }
func (m *ActivityPagedResponse) String() string            { return proto.CompactTextString(m) }
func (*ActivityPagedResponse) ProtoMessage()               {}
func (*ActivityPagedResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *ActivityPagedResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *ActivityPagedResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *ActivityPagedResponse) GetData() *UserActivityPage {
	if m != nil {
		return m.Data
	}
	return nil
}

type ActivityData struct {
	Activity *Activity `protobuf:"bytes,1,opt,name=activity" json:"activity"`
}

func (m *ActivityData) Reset()                    { *m = ActivityData{} }
func (m *ActivityData) String() string            { return proto.CompactTextString(m) }
func (*ActivityData) ProtoMessage()               {}
func (*ActivityData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *ActivityData) GetActivity() *Activity {
	if m != nil {
		return m.Activity
	}
	return nil
}

type ActivityResponse struct {
	Status  string        `protobuf:"bytes,1,opt,name=status" json:"status"`
	Message string        `protobuf:"bytes,2,opt,name=message" json:"message"`
	Data    *ActivityData `protobuf:"bytes,3,opt,name=data" json:"data"`
}

func (m *ActivityResponse) Reset()                    { *m = ActivityResponse{} }
func (m *ActivityResponse) String() string            { return proto.CompactTextString(m) }
func (*ActivityResponse) ProtoMessage()               {}
func (*ActivityResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *ActivityResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *ActivityResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *ActivityResponse) GetData() *ActivityData {
	if m != nil {
		return m.Data
	}
	return nil
}

type ActivityList struct {
	Activity []*Activity `protobuf:"bytes,1,rep,name=activity" json:"activity"`
}

func (m *ActivityList) Reset()                    { *m = ActivityList{} }
func (m *ActivityList) String() string            { return proto.CompactTextString(m) }
func (*ActivityList) ProtoMessage()               {}
func (*ActivityList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *ActivityList) GetActivity() []*Activity {
	if m != nil {
		return m.Activity
	}
	return nil
}

type ActivityListResponse struct {
	Status  string        `protobuf:"bytes,1,opt,name=status" json:"status"`
	Message string        `protobuf:"bytes,2,opt,name=message" json:"message"`
	Data    *ActivityList `protobuf:"bytes,3,opt,name=data" json:"data"`
}

func (m *ActivityListResponse) Reset()                    { *m = ActivityListResponse{} }
func (m *ActivityListResponse) String() string            { return proto.CompactTextString(m) }
func (*ActivityListResponse) ProtoMessage()               {}
func (*ActivityListResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *ActivityListResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *ActivityListResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *ActivityListResponse) GetData() *ActivityList {
	if m != nil {
		return m.Data
	}
	return nil
}

type ActivityCount struct {
	Count uint32 `protobuf:"varint,1,opt,name=count" json:"count"`
}

func (m *ActivityCount) Reset()                    { *m = ActivityCount{} }
func (m *ActivityCount) String() string            { return proto.CompactTextString(m) }
func (*ActivityCount) ProtoMessage()               {}
func (*ActivityCount) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *ActivityCount) GetCount() uint32 {
	if m != nil {
		return m.Count
	}
	return 0
}

type ActivityCountResponse struct {
	Status  string         `protobuf:"bytes,1,opt,name=status" json:"status"`
	Message string         `protobuf:"bytes,2,opt,name=message" json:"message"`
	Data    *ActivityCount `protobuf:"bytes,3,opt,name=data" json:"data"`
}

func (m *ActivityCountResponse) Reset()                    { *m = ActivityCountResponse{} }
func (m *ActivityCountResponse) String() string            { return proto.CompactTextString(m) }
func (*ActivityCountResponse) ProtoMessage()               {}
func (*ActivityCountResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *ActivityCountResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *ActivityCountResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *ActivityCountResponse) GetData() *ActivityCount {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*EmailRequest)(nil), "ios.EmailRequest")
	proto.RegisterType((*ActivityRequest)(nil), "ios.ActivityRequest")
	proto.RegisterType((*RecentActivityRequest)(nil), "ios.RecentActivityRequest")
	proto.RegisterType((*ActivityCountRequest)(nil), "ios.ActivityCountRequest")
	proto.RegisterType((*UpdateActivityRequest)(nil), "ios.UpdateActivityRequest")
	proto.RegisterType((*EmailResponse)(nil), "ios.EmailResponse")
	proto.RegisterType((*Activity)(nil), "ios.Activity")
	proto.RegisterType((*UserActivityPage)(nil), "ios.UserActivityPage")
	proto.RegisterType((*ActivityPagedResponse)(nil), "ios.ActivityPagedResponse")
	proto.RegisterType((*ActivityData)(nil), "ios.ActivityData")
	proto.RegisterType((*ActivityResponse)(nil), "ios.ActivityResponse")
	proto.RegisterType((*ActivityList)(nil), "ios.ActivityList")
	proto.RegisterType((*ActivityListResponse)(nil), "ios.ActivityListResponse")
	proto.RegisterType((*ActivityCount)(nil), "ios.ActivityCount")
	proto.RegisterType((*ActivityCountResponse)(nil), "ios.ActivityCountResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Notification service

type NotificationClient interface {
	FindUserActivity(ctx context.Context, in *ActivityRequest, opts ...client.CallOption) (*ActivityPagedResponse, error)
	FindMostRecentActivity(ctx context.Context, in *RecentActivityRequest, opts ...client.CallOption) (*ActivityListResponse, error)
	FindActivityCount(ctx context.Context, in *ActivityCountRequest, opts ...client.CallOption) (*ActivityCountResponse, error)
	UpdateActivity(ctx context.Context, in *UpdateActivityRequest, opts ...client.CallOption) (*ActivityResponse, error)
	SendEmail(ctx context.Context, in *EmailRequest, opts ...client.CallOption) (*EmailResponse, error)
}

type notificationClient struct {
	c           client.Client
	serviceName string
}

func NewNotificationClient(serviceName string, c client.Client) NotificationClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "ios"
	}
	return &notificationClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *notificationClient) FindUserActivity(ctx context.Context, in *ActivityRequest, opts ...client.CallOption) (*ActivityPagedResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Notification.FindUserActivity", in)
	out := new(ActivityPagedResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationClient) FindMostRecentActivity(ctx context.Context, in *RecentActivityRequest, opts ...client.CallOption) (*ActivityListResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Notification.FindMostRecentActivity", in)
	out := new(ActivityListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationClient) FindActivityCount(ctx context.Context, in *ActivityCountRequest, opts ...client.CallOption) (*ActivityCountResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Notification.FindActivityCount", in)
	out := new(ActivityCountResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationClient) UpdateActivity(ctx context.Context, in *UpdateActivityRequest, opts ...client.CallOption) (*ActivityResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Notification.UpdateActivity", in)
	out := new(ActivityResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationClient) SendEmail(ctx context.Context, in *EmailRequest, opts ...client.CallOption) (*EmailResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Notification.SendEmail", in)
	out := new(EmailResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Notification service

type NotificationHandler interface {
	FindUserActivity(context.Context, *ActivityRequest, *ActivityPagedResponse) error
	FindMostRecentActivity(context.Context, *RecentActivityRequest, *ActivityListResponse) error
	FindActivityCount(context.Context, *ActivityCountRequest, *ActivityCountResponse) error
	UpdateActivity(context.Context, *UpdateActivityRequest, *ActivityResponse) error
	SendEmail(context.Context, *EmailRequest, *EmailResponse) error
}

func RegisterNotificationHandler(s server.Server, hdlr NotificationHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&Notification{hdlr}, opts...))
}

type Notification struct {
	NotificationHandler
}

func (h *Notification) FindUserActivity(ctx context.Context, in *ActivityRequest, out *ActivityPagedResponse) error {
	return h.NotificationHandler.FindUserActivity(ctx, in, out)
}

func (h *Notification) FindMostRecentActivity(ctx context.Context, in *RecentActivityRequest, out *ActivityListResponse) error {
	return h.NotificationHandler.FindMostRecentActivity(ctx, in, out)
}

func (h *Notification) FindActivityCount(ctx context.Context, in *ActivityCountRequest, out *ActivityCountResponse) error {
	return h.NotificationHandler.FindActivityCount(ctx, in, out)
}

func (h *Notification) UpdateActivity(ctx context.Context, in *UpdateActivityRequest, out *ActivityResponse) error {
	return h.NotificationHandler.UpdateActivity(ctx, in, out)
}

func (h *Notification) SendEmail(ctx context.Context, in *EmailRequest, out *EmailResponse) error {
	return h.NotificationHandler.SendEmail(ctx, in, out)
}

func init() { proto.RegisterFile("proto/notification/notification.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 700 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x55, 0xdd, 0x4e, 0xdb, 0x4a,
	0x10, 0xc6, 0x49, 0x80, 0x64, 0x42, 0x38, 0xb0, 0x22, 0xc8, 0x44, 0x47, 0x47, 0xc8, 0x12, 0x08,
	0x6e, 0x38, 0x52, 0xce, 0xb9, 0xe9, 0x65, 0x5a, 0x40, 0x42, 0xa2, 0x55, 0x6b, 0xc4, 0x03, 0x2c,
	0xf6, 0xb6, 0xdd, 0x92, 0x78, 0x4d, 0x76, 0x5c, 0x95, 0xde, 0xf6, 0x6d, 0x7a, 0xdd, 0x47, 0xeb,
	0x03, 0x54, 0xfb, 0xe7, 0xec, 0x9a, 0x50, 0x21, 0xa4, 0x5e, 0x25, 0xdf, 0xcc, 0xce, 0xcc, 0xb7,
	0xdf, 0x8c, 0x67, 0xe1, 0xa0, 0x9c, 0x0b, 0x14, 0xff, 0x16, 0x02, 0xf9, 0x7b, 0x9e, 0x51, 0xe4,
	0xa2, 0x08, 0xc0, 0x89, 0xf6, 0x93, 0x36, 0x17, 0x32, 0xf9, 0x1e, 0xc1, 0xc6, 0xd9, 0x8c, 0xf2,
	0x69, 0xca, 0xee, 0x2a, 0x26, 0x91, 0xec, 0x43, 0x9f, 0x29, 0x7c, 0xc5, 0x8a, 0x9c, 0xcd, 0xe3,
	0x68, 0x3f, 0x3a, 0xea, 0xa5, 0xbe, 0x89, 0x1c, 0xc2, 0x26, 0x33, 0x11, 0x19, 0x2f, 0x39, 0x2b,
	0x30, 0x6e, 0xe9, 0x43, 0x0d, 0x2b, 0x89, 0x61, 0x5d, 0x56, 0x37, 0x9f, 0x58, 0x86, 0x71, 0x5b,
	0x1f, 0x70, 0x90, 0x8c, 0xa0, 0xfb, 0x11, 0x67, 0xd3, 0x97, 0x22, 0xbf, 0x8f, 0x3b, 0xda, 0x55,
	0x63, 0xe5, 0x43, 0xf6, 0x05, 0xb5, 0x6f, 0xd5, 0xf8, 0x1c, 0x4e, 0x2a, 0xf8, 0x6b, 0x92, 0x21,
	0xff, 0xcc, 0xf1, 0xde, 0xd1, 0xdd, 0x85, 0xb5, 0x4a, 0xb2, 0xf9, 0xc5, 0xa9, 0x65, 0x6a, 0x91,
	0x4a, 0x23, 0x74, 0xb1, 0x8b, 0x53, 0x4b, 0xaf, 0xc6, 0x84, 0x40, 0xa7, 0xa4, 0x1f, 0x98, 0x66,
	0x35, 0x48, 0xf5, 0x7f, 0x75, 0x5e, 0xfd, 0x5e, 0xf1, 0xaf, 0x4c, 0x53, 0x1a, 0xa4, 0x35, 0x4e,
	0x2e, 0x60, 0x98, 0xb2, 0x8c, 0x15, 0xd8, 0x2c, 0xee, 0x17, 0x89, 0x1a, 0x45, 0x76, 0x60, 0x35,
	0x13, 0x95, 0x15, 0x67, 0x90, 0x1a, 0x90, 0x8c, 0x61, 0xc7, 0x25, 0x79, 0xa5, 0x0c, 0x4f, 0xc8,
	0x94, 0xcc, 0x60, 0x78, 0x5d, 0xe6, 0x14, 0x59, 0xb3, 0xfc, 0x3f, 0x00, 0xd4, 0x9a, 0xea, 0x30,
	0xcf, 0xa2, 0xb4, 0x91, 0x8c, 0x15, 0x13, 0xd7, 0x20, 0x8b, 0xc8, 0xdf, 0xd0, 0xcb, 0xa6, 0x3c,
	0xbb, 0x65, 0xf9, 0xc4, 0xb5, 0x66, 0x61, 0x48, 0x26, 0x30, 0xb0, 0x03, 0x21, 0x4b, 0x51, 0x48,
	0xa6, 0xd3, 0x20, 0xc5, 0x4a, 0x3a, 0x89, 0x0d, 0x52, 0xfd, 0x9d, 0x31, 0x29, 0x95, 0x92, 0x26,
	0xbf, 0x83, 0xc9, 0x8f, 0x16, 0x74, 0x1d, 0xd9, 0xa7, 0xb0, 0xb4, 0x1d, 0x6c, 0x07, 0x1d, 0x24,
	0xd0, 0xc1, 0xfb, 0xd2, 0xe5, 0xd6, 0xff, 0x03, 0x99, 0x3a, 0x0f, 0x05, 0x47, 0x8e, 0x53, 0x66,
	0xa7, 0xc6, 0x00, 0x15, 0x21, 0xab, 0x1b, 0xe3, 0x58, 0x33, 0x11, 0x0e, 0xab, 0x51, 0xcf, 0x99,
	0xcc, 0xe6, 0xbc, 0x54, 0x5f, 0x45, 0xbc, 0x6e, 0x46, 0xdd, 0x33, 0xa9, 0x2b, 0xe6, 0x0c, 0x29,
	0x9f, 0xca, 0xb8, 0x6b, 0xae, 0x68, 0xa1, 0xd2, 0x10, 0xf9, 0x8c, 0x49, 0xa4, 0xb3, 0x32, 0xee,
	0x19, 0x0d, 0x6b, 0x43, 0xa8, 0x30, 0x34, 0x14, 0xf6, 0xfa, 0xd2, 0xf7, 0xfb, 0x92, 0x7c, 0x8b,
	0x60, 0xeb, 0x5a, 0xb2, 0xb9, 0x93, 0xee, 0xad, 0x1a, 0x4c, 0x37, 0xac, 0xd1, 0x23, 0xc3, 0xda,
	0x0a, 0x87, 0x55, 0xcb, 0x20, 0x90, 0x4e, 0xed, 0x74, 0x1b, 0x40, 0x8e, 0xa1, 0xeb, 0x24, 0x8f,
	0x3b, 0xfb, 0xed, 0xa3, 0xfe, 0x78, 0x70, 0xc2, 0x85, 0x3c, 0xa9, 0x47, 0xaa, 0x76, 0x27, 0x08,
	0x43, 0x9f, 0x40, 0xfe, 0xfc, 0x39, 0x20, 0xc7, 0xd0, 0xc9, 0x29, 0x52, 0x4d, 0xa5, 0x3f, 0x1e,
	0xea, 0x8a, 0xcd, 0x0b, 0xa6, 0xfa, 0x48, 0xf2, 0x02, 0x36, 0x9c, 0xf5, 0x94, 0x22, 0x0d, 0x08,
	0x47, 0x3a, 0xfc, 0x51, 0xc2, 0xb7, 0xb0, 0xb5, 0xf8, 0x32, 0x9e, 0xcd, 0xf5, 0x20, 0xe0, 0xba,
	0x1d, 0x14, 0x53, 0x8c, 0x1e, 0xf2, 0xbc, 0xe4, 0x12, 0x1b, 0x3c, 0x7f, 0x2b, 0xac, 0x58, 0x7c,
	0xfb, 0x2a, 0xf4, 0x0f, 0x71, 0xd5, 0xa9, 0x0d, 0xd7, 0x03, 0x18, 0x04, 0xcb, 0x66, 0xb1, 0x93,
	0x22, 0x7f, 0x27, 0xdd, 0x2d, 0x1a, 0x6e, 0x77, 0xd2, 0xb3, 0x89, 0x1d, 0x06, 0xc4, 0x48, 0x40,
	0xcc, 0xe4, 0xd6, 0xfe, 0xf1, 0xcf, 0x16, 0x6c, 0xbc, 0xf1, 0x5e, 0x24, 0x72, 0x0e, 0x5b, 0xe7,
	0xbc, 0xc8, 0xfd, 0xe1, 0x20, 0x3b, 0xa1, 0x90, 0x66, 0xe9, 0x8d, 0x46, 0x81, 0x35, 0x98, 0xd0,
	0x64, 0x85, 0xbc, 0x83, 0x5d, 0x95, 0xe7, 0xb5, 0x50, 0xfa, 0xfa, 0x2b, 0x9b, 0x98, 0xb8, 0xa5,
	0x7b, 0x7c, 0xb4, 0xf7, 0x50, 0xc1, 0x45, 0xca, 0x4b, 0xd8, 0x56, 0x29, 0x43, 0x25, 0xf7, 0x96,
	0x5c, 0x6d, 0x29, 0xc1, 0x40, 0xd1, 0x64, 0x85, 0x9c, 0xc1, 0x66, 0xb8, 0xcc, 0x2d, 0xb1, 0xa5,
	0x1b, 0x7e, 0x34, 0x6c, 0x48, 0x50, 0xa7, 0xf9, 0x1f, 0x7a, 0xea, 0x35, 0xd6, 0x8b, 0x9a, 0x98,
	0x01, 0xf0, 0x5f, 0xf1, 0x11, 0xf1, 0x4d, 0x2e, 0xea, 0x66, 0x4d, 0x3f, 0xfc, 0xff, 0xfd, 0x0a,
	0x00, 0x00, 0xff, 0xff, 0x2b, 0x1a, 0x74, 0x93, 0x21, 0x08, 0x00, 0x00,
}
