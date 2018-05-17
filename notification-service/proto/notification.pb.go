// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/notification/notification.proto

/*
Package go_micro_srv_notification is a generated protocol buffer package.

It is generated from these files:
	proto/notification/notification.proto

It has these top-level messages:
	GetUserNotifcationsRequest
	Notification
	UserNotifications
	NotificationListResponse
*/
package go_micro_srv_notification

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
type GetUserNotifcationsRequest struct {
	UserID string `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
}

func (m *GetUserNotifcationsRequest) Reset()                    { *m = GetUserNotifcationsRequest{} }
func (m *GetUserNotifcationsRequest) String() string            { return proto.CompactTextString(m) }
func (*GetUserNotifcationsRequest) ProtoMessage()               {}
func (*GetUserNotifcationsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *GetUserNotifcationsRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

// Responses
type Notification struct {
	NotificationID string `protobuf:"bytes,1,opt,name=notificationID" json:"notificationID,omitempty"`
	UserID         string `protobuf:"bytes,2,opt,name=userID" json:"userID,omitempty"`
	Description    string `protobuf:"bytes,3,opt,name=description" json:"description,omitempty"`
}

func (m *Notification) Reset()                    { *m = Notification{} }
func (m *Notification) String() string            { return proto.CompactTextString(m) }
func (*Notification) ProtoMessage()               {}
func (*Notification) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Notification) GetNotificationID() string {
	if m != nil {
		return m.NotificationID
	}
	return ""
}

func (m *Notification) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *Notification) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

type UserNotifications struct {
	Notifications []*Notification `protobuf:"bytes,1,rep,name=notifications" json:"notifications,omitempty"`
}

func (m *UserNotifications) Reset()                    { *m = UserNotifications{} }
func (m *UserNotifications) String() string            { return proto.CompactTextString(m) }
func (*UserNotifications) ProtoMessage()               {}
func (*UserNotifications) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *UserNotifications) GetNotifications() []*Notification {
	if m != nil {
		return m.Notifications
	}
	return nil
}

type NotificationListResponse struct {
	Status  string             `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	Message string             `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	Data    *UserNotifications `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *NotificationListResponse) Reset()                    { *m = NotificationListResponse{} }
func (m *NotificationListResponse) String() string            { return proto.CompactTextString(m) }
func (*NotificationListResponse) ProtoMessage()               {}
func (*NotificationListResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *NotificationListResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *NotificationListResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *NotificationListResponse) GetData() *UserNotifications {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*GetUserNotifcationsRequest)(nil), "go.micro.srv.notification.GetUserNotifcationsRequest")
	proto.RegisterType((*Notification)(nil), "go.micro.srv.notification.Notification")
	proto.RegisterType((*UserNotifications)(nil), "go.micro.srv.notification.UserNotifications")
	proto.RegisterType((*NotificationListResponse)(nil), "go.micro.srv.notification.NotificationListResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for NotificationService service

type NotificationServiceClient interface {
	GetUserNotifications(ctx context.Context, in *GetUserNotifcationsRequest, opts ...client.CallOption) (*NotificationListResponse, error)
}

type notificationServiceClient struct {
	c           client.Client
	serviceName string
}

func NewNotificationServiceClient(serviceName string, c client.Client) NotificationServiceClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "go.micro.srv.notification"
	}
	return &notificationServiceClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *notificationServiceClient) GetUserNotifications(ctx context.Context, in *GetUserNotifcationsRequest, opts ...client.CallOption) (*NotificationListResponse, error) {
	req := c.c.NewRequest(c.serviceName, "NotificationService.GetUserNotifications", in)
	out := new(NotificationListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for NotificationService service

type NotificationServiceHandler interface {
	GetUserNotifications(context.Context, *GetUserNotifcationsRequest, *NotificationListResponse) error
}

func RegisterNotificationServiceHandler(s server.Server, hdlr NotificationServiceHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&NotificationService{hdlr}, opts...))
}

type NotificationService struct {
	NotificationServiceHandler
}

func (h *NotificationService) GetUserNotifications(ctx context.Context, in *GetUserNotifcationsRequest, out *NotificationListResponse) error {
	return h.NotificationServiceHandler.GetUserNotifications(ctx, in, out)
}

func init() { proto.RegisterFile("proto/notification/notification.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 291 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0x41, 0x4b, 0xf3, 0x40,
	0x10, 0xfd, 0xf6, 0xab, 0x54, 0x9c, 0xa8, 0xe0, 0x2a, 0x12, 0x7b, 0x0a, 0x01, 0xb5, 0x07, 0x59,
	0x21, 0xd5, 0xbb, 0x07, 0x41, 0x04, 0xed, 0x21, 0xe2, 0x0f, 0xd8, 0xa6, 0x63, 0xd9, 0x43, 0xb3,
	0x71, 0x67, 0xd3, 0x5f, 0xe0, 0xd9, 0x5f, 0xe0, 0x8f, 0x95, 0x6c, 0xd3, 0x30, 0x55, 0x5a, 0x3c,
	0xbe, 0x37, 0xb3, 0x6f, 0xe6, 0xbd, 0x59, 0x38, 0xaf, 0x9c, 0xf5, 0xf6, 0xba, 0xb4, 0xde, 0xbc,
	0x99, 0x42, 0x7b, 0x63, 0xcb, 0x35, 0xa0, 0x42, 0x5d, 0x9e, 0xcd, 0xac, 0x9a, 0x9b, 0xc2, 0x59,
	0x45, 0x6e, 0xa1, 0x78, 0x43, 0x7a, 0x03, 0x83, 0x07, 0xf4, 0xaf, 0x84, 0x6e, 0xdc, 0xd0, 0x4b,
	0x96, 0x72, 0x7c, 0xaf, 0x91, 0xbc, 0x3c, 0x85, 0x7e, 0x4d, 0xe8, 0x1e, 0xef, 0x63, 0x91, 0x88,
	0xe1, 0x5e, 0xde, 0xa2, 0xb4, 0x82, 0xfd, 0x31, 0x53, 0x91, 0x17, 0x70, 0xc8, 0x55, 0xbb, 0xfe,
	0x1f, 0x2c, 0xd3, 0xfb, 0xcf, 0xf5, 0x64, 0x02, 0xd1, 0x14, 0xa9, 0x70, 0xa6, 0x6a, 0x1a, 0xe3,
	0x5e, 0x28, 0x72, 0x2a, 0x9d, 0xc0, 0x51, 0xb7, 0x64, 0xab, 0x47, 0xf2, 0x19, 0x0e, 0xf8, 0x00,
	0x8a, 0x45, 0xd2, 0x1b, 0x46, 0xd9, 0xa5, 0xda, 0xe8, 0x57, 0x71, 0x81, 0x7c, 0xfd, 0x75, 0xfa,
	0x29, 0x20, 0xe6, 0xf5, 0x27, 0x43, 0x3e, 0x47, 0xaa, 0x6c, 0x49, 0xd8, 0xac, 0x4e, 0x5e, 0xfb,
	0x9a, 0x56, 0x51, 0x2c, 0x91, 0x8c, 0x61, 0x77, 0x8e, 0x44, 0x7a, 0x86, 0xad, 0xa7, 0x15, 0x94,
	0x77, 0xb0, 0x33, 0xd5, 0x5e, 0x07, 0x37, 0x51, 0x76, 0xb5, 0x65, 0xa9, 0x5f, 0xce, 0xf2, 0xf0,
	0x32, 0xfb, 0x12, 0x70, 0xcc, 0xf9, 0x17, 0x74, 0x0b, 0x53, 0xa0, 0xfc, 0x10, 0x70, 0xc2, 0xaf,
	0xd6, 0x05, 0x72, 0xbb, 0x65, 0xc8, 0xe6, 0x33, 0x0f, 0x46, 0x7f, 0x0c, 0x8c, 0x07, 0x92, 0xfe,
	0x9b, 0xf4, 0xc3, 0xef, 0x1a, 0x7d, 0x07, 0x00, 0x00, 0xff, 0xff, 0xf3, 0x24, 0x59, 0x91, 0x86,
	0x02, 0x00, 0x00,
}