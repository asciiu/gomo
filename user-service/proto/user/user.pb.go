// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/user/user.proto

/*
Package go_micro_srv_user is a generated protocol buffer package.

It is generated from these files:
	proto/user/user.proto

It has these top-level messages:
	ChangePasswordRequest
	Response
	UserInfoRequest
*/
package go_micro_srv_user

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

type ChangePasswordRequest struct {
	Id     string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	UserId string `protobuf:"bytes,2,opt,name=userId" json:"userId,omitempty"`
}

func (m *ChangePasswordRequest) Reset()                    { *m = ChangePasswordRequest{} }
func (m *ChangePasswordRequest) String() string            { return proto.CompactTextString(m) }
func (*ChangePasswordRequest) ProtoMessage()               {}
func (*ChangePasswordRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ChangePasswordRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ChangePasswordRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

type Response struct {
	Success bool `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Response) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

// Created a blank get request
type UserInfoRequest struct {
	UserId string `protobuf:"bytes,1,opt,name=userId" json:"userId,omitempty"`
}

func (m *UserInfoRequest) Reset()                    { *m = UserInfoRequest{} }
func (m *UserInfoRequest) String() string            { return proto.CompactTextString(m) }
func (*UserInfoRequest) ProtoMessage()               {}
func (*UserInfoRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *UserInfoRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func init() {
	proto.RegisterType((*ChangePasswordRequest)(nil), "go.micro.srv.user.ChangePasswordRequest")
	proto.RegisterType((*Response)(nil), "go.micro.srv.user.Response")
	proto.RegisterType((*UserInfoRequest)(nil), "go.micro.srv.user.UserInfoRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for UserService service

type UserServiceClient interface {
	GetUserInfo(ctx context.Context, in *UserInfoRequest, opts ...client.CallOption) (*Response, error)
	ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...client.CallOption) (*Response, error)
}

type userServiceClient struct {
	c           client.Client
	serviceName string
}

func NewUserServiceClient(serviceName string, c client.Client) UserServiceClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "go.micro.srv.user"
	}
	return &userServiceClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *userServiceClient) GetUserInfo(ctx context.Context, in *UserInfoRequest, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.serviceName, "UserService.GetUserInfo", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.serviceName, "UserService.ChangePassword", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserService service

type UserServiceHandler interface {
	GetUserInfo(context.Context, *UserInfoRequest, *Response) error
	ChangePassword(context.Context, *ChangePasswordRequest, *Response) error
}

func RegisterUserServiceHandler(s server.Server, hdlr UserServiceHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&UserService{hdlr}, opts...))
}

type UserService struct {
	UserServiceHandler
}

func (h *UserService) GetUserInfo(ctx context.Context, in *UserInfoRequest, out *Response) error {
	return h.UserServiceHandler.GetUserInfo(ctx, in, out)
}

func (h *UserService) ChangePassword(ctx context.Context, in *ChangePasswordRequest, out *Response) error {
	return h.UserServiceHandler.ChangePassword(ctx, in, out)
}

func init() { proto.RegisterFile("proto/user/user.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 219 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2d, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x2f, 0x2d, 0x4e, 0x2d, 0x02, 0x13, 0x7a, 0x60, 0xbe, 0x90, 0x60, 0x7a, 0xbe, 0x5e,
	0x6e, 0x66, 0x72, 0x51, 0xbe, 0x5e, 0x71, 0x51, 0x99, 0x1e, 0x48, 0x42, 0xc9, 0x9e, 0x4b, 0xd4,
	0x39, 0x23, 0x31, 0x2f, 0x3d, 0x35, 0x20, 0xb1, 0xb8, 0xb8, 0x3c, 0xbf, 0x28, 0x25, 0x28, 0xb5,
	0xb0, 0x34, 0xb5, 0xb8, 0x44, 0x88, 0x8f, 0x8b, 0x29, 0x33, 0x45, 0x82, 0x51, 0x81, 0x51, 0x83,
	0x33, 0x88, 0x29, 0x33, 0x45, 0x48, 0x8c, 0x8b, 0x0d, 0xa4, 0xc1, 0x33, 0x45, 0x82, 0x09, 0x2c,
	0x06, 0xe5, 0x29, 0xa9, 0x70, 0x71, 0x04, 0xa5, 0x16, 0x17, 0xe4, 0xe7, 0x15, 0xa7, 0x0a, 0x49,
	0x70, 0xb1, 0x17, 0x97, 0x26, 0x27, 0xa7, 0x16, 0x17, 0x83, 0x35, 0x72, 0x04, 0xc1, 0xb8, 0x4a,
	0x9a, 0x5c, 0xfc, 0xa1, 0x20, 0xf5, 0x79, 0x69, 0xf9, 0x30, 0x0b, 0x10, 0x06, 0x32, 0x22, 0x1b,
	0x68, 0xb4, 0x8b, 0x91, 0x8b, 0x1b, 0xa4, 0x36, 0x38, 0xb5, 0xa8, 0x2c, 0x33, 0x39, 0x55, 0x28,
	0x80, 0x8b, 0xdb, 0x3d, 0xb5, 0x04, 0xa6, 0x5b, 0x48, 0x49, 0x0f, 0xc3, 0x13, 0x7a, 0x68, 0x46,
	0x4b, 0x49, 0x63, 0x51, 0x03, 0x73, 0xa4, 0x12, 0x83, 0x50, 0x24, 0x17, 0x1f, 0xaa, 0x9f, 0x85,
	0x34, 0xb0, 0x68, 0xc0, 0x1a, 0x2c, 0x04, 0x8c, 0x4e, 0x62, 0x03, 0x07, 0xb4, 0x31, 0x20, 0x00,
	0x00, 0xff, 0xff, 0x33, 0xf9, 0x26, 0x6f, 0x81, 0x01, 0x00, 0x00,
}
