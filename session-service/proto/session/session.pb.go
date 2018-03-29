// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/session/session.proto

/*
Package go_micro_srv_session is a generated protocol buffer package.

It is generated from these files:
	proto/session/session.proto

It has these top-level messages:
	SessionRequest
	SessionResponse
	Response
	GetRequest
*/
package go_micro_srv_session

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

type SessionRequest struct {
	Id     string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	UserId string `protobuf:"bytes,2,opt,name=userId" json:"userId,omitempty"`
}

func (m *SessionRequest) Reset()                    { *m = SessionRequest{} }
func (m *SessionRequest) String() string            { return proto.CompactTextString(m) }
func (*SessionRequest) ProtoMessage()               {}
func (*SessionRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *SessionRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *SessionRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

type SessionResponse struct {
	SessionId string `protobuf:"bytes,1,opt,name=sessionId" json:"sessionId,omitempty"`
}

func (m *SessionResponse) Reset()                    { *m = SessionResponse{} }
func (m *SessionResponse) String() string            { return proto.CompactTextString(m) }
func (*SessionResponse) ProtoMessage()               {}
func (*SessionResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SessionResponse) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

type Response struct {
	Success  bool             `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
	Response *SessionResponse `protobuf:"bytes,2,opt,name=response" json:"response,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Response) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *Response) GetResponse() *SessionResponse {
	if m != nil {
		return m.Response
	}
	return nil
}

// Created a blank get request
type GetRequest struct {
	SessionId  string `protobuf:"bytes,1,opt,name=sessionId" json:"sessionId,omitempty"`
	SessionKey string `protobuf:"bytes,2,opt,name=sessionKey" json:"sessionKey,omitempty"`
}

func (m *GetRequest) Reset()                    { *m = GetRequest{} }
func (m *GetRequest) String() string            { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()               {}
func (*GetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *GetRequest) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

func (m *GetRequest) GetSessionKey() string {
	if m != nil {
		return m.SessionKey
	}
	return ""
}

func init() {
	proto.RegisterType((*SessionRequest)(nil), "go.micro.srv.session.SessionRequest")
	proto.RegisterType((*SessionResponse)(nil), "go.micro.srv.session.SessionResponse")
	proto.RegisterType((*Response)(nil), "go.micro.srv.session.Response")
	proto.RegisterType((*GetRequest)(nil), "go.micro.srv.session.GetRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for SessionService service

type SessionServiceClient interface {
	CreateSession(ctx context.Context, in *SessionRequest, opts ...client.CallOption) (*Response, error)
	GetSession(ctx context.Context, in *GetRequest, opts ...client.CallOption) (*Response, error)
	DeleteSession(ctx context.Context, in *GetRequest, opts ...client.CallOption) (*Response, error)
}

type sessionServiceClient struct {
	c           client.Client
	serviceName string
}

func NewSessionServiceClient(serviceName string, c client.Client) SessionServiceClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "go.micro.srv.session"
	}
	return &sessionServiceClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *sessionServiceClient) CreateSession(ctx context.Context, in *SessionRequest, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.serviceName, "SessionService.CreateSession", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionServiceClient) GetSession(ctx context.Context, in *GetRequest, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.serviceName, "SessionService.GetSession", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sessionServiceClient) DeleteSession(ctx context.Context, in *GetRequest, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.serviceName, "SessionService.DeleteSession", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for SessionService service

type SessionServiceHandler interface {
	CreateSession(context.Context, *SessionRequest, *Response) error
	GetSession(context.Context, *GetRequest, *Response) error
	DeleteSession(context.Context, *GetRequest, *Response) error
}

func RegisterSessionServiceHandler(s server.Server, hdlr SessionServiceHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&SessionService{hdlr}, opts...))
}

type SessionService struct {
	SessionServiceHandler
}

func (h *SessionService) CreateSession(ctx context.Context, in *SessionRequest, out *Response) error {
	return h.SessionServiceHandler.CreateSession(ctx, in, out)
}

func (h *SessionService) GetSession(ctx context.Context, in *GetRequest, out *Response) error {
	return h.SessionServiceHandler.GetSession(ctx, in, out)
}

func (h *SessionService) DeleteSession(ctx context.Context, in *GetRequest, out *Response) error {
	return h.SessionServiceHandler.DeleteSession(ctx, in, out)
}

func init() { proto.RegisterFile("proto/session/session.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 266 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x92, 0x41, 0x4b, 0x03, 0x31,
	0x10, 0x85, 0xed, 0x1e, 0xea, 0x76, 0xa4, 0x15, 0x06, 0x91, 0x45, 0xa5, 0x94, 0xa0, 0xe0, 0x29,
	0x85, 0x7a, 0xf1, 0x2a, 0x0a, 0x52, 0xbd, 0xc8, 0xee, 0xc1, 0xb3, 0x66, 0x87, 0x12, 0xd0, 0xa6,
	0x66, 0xb2, 0x0b, 0xfe, 0x0b, 0x7f, 0xb2, 0x10, 0x93, 0xac, 0xc8, 0xa2, 0x1e, 0x7a, 0xda, 0x9d,
	0x97, 0x37, 0x6f, 0xbe, 0x09, 0x81, 0xe3, 0x8d, 0x35, 0xce, 0xcc, 0x99, 0x98, 0xb5, 0x59, 0xc7,
	0xaf, 0xf4, 0x2a, 0x1e, 0xac, 0x8c, 0x7c, 0xd5, 0xca, 0x1a, 0xc9, 0xb6, 0x95, 0xe1, 0x4c, 0x5c,
	0xc2, 0xa4, 0xfa, 0xfa, 0x2d, 0xe9, 0xad, 0x21, 0x76, 0x38, 0x81, 0x4c, 0xd7, 0xc5, 0x60, 0x36,
	0x38, 0x1f, 0x95, 0x99, 0xae, 0xf1, 0x10, 0x86, 0x0d, 0x93, 0x5d, 0xd6, 0x45, 0xe6, 0xb5, 0x50,
	0x89, 0x39, 0xec, 0xa7, 0x4e, 0xde, 0x98, 0x35, 0x13, 0x9e, 0xc0, 0x28, 0xe4, 0x2e, 0x63, 0x42,
	0x27, 0x88, 0x15, 0xe4, 0xc9, 0x59, 0xc0, 0x2e, 0x37, 0x4a, 0x11, 0xb3, 0xf7, 0xe5, 0x65, 0x2c,
	0xf1, 0x0a, 0x72, 0x1b, 0x5c, 0x7e, 0xe0, 0xde, 0xe2, 0x4c, 0xf6, 0x91, 0xcb, 0x1f, 0xc3, 0xcb,
	0xd4, 0x26, 0xee, 0x00, 0x6e, 0xc9, 0xc5, 0x7d, 0x7e, 0x85, 0xc2, 0x29, 0x40, 0x28, 0xee, 0xe9,
	0x3d, 0x6c, 0xf8, 0x4d, 0x59, 0x7c, 0x64, 0xe9, 0x82, 0x2a, 0xb2, 0xad, 0x56, 0x84, 0x8f, 0x30,
	0xbe, 0xb6, 0xf4, 0xe4, 0x28, 0xe8, 0x78, 0xfa, 0x07, 0xa0, 0xe7, 0x38, 0x9a, 0xf6, 0xbb, 0x22,
	0xbf, 0xd8, 0xc1, 0x07, 0xcf, 0x1d, 0x53, 0x67, 0xfd, 0xfe, 0x6e, 0xb3, 0x7f, 0x24, 0x56, 0x30,
	0xbe, 0xa1, 0x17, 0xea, 0x50, 0xb7, 0x10, 0xfa, 0x3c, 0xf4, 0xef, 0xe9, 0xe2, 0x33, 0x00, 0x00,
	0xff, 0xff, 0x79, 0x2d, 0xea, 0x72, 0x6e, 0x02, 0x00, 0x00,
}