// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/apikey/apikey.proto

/*
Package go_micro_srv_apikey is a generated protocol buffer package.

It is generated from these files:
	proto/apikey/apikey.proto

It has these top-level messages:
	ApiKeyRequest
	GetUserApiKeyRequest
	GetUserApiKeysRequest
	RemoveApiKeyRequest
	ApiKey
	UserApiKeyData
	UserApiKeysData
	ApiKeyResponse
	ApiKeyListResponse
*/
package go_micro_srv_apikey

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
type ApiKeyRequest struct {
	ApiKeyId    string `protobuf:"bytes,1,opt,name=apiKeyId" json:"apiKeyId,omitempty"`
	UserId      string `protobuf:"bytes,2,opt,name=userId" json:"userId,omitempty"`
	Exchange    string `protobuf:"bytes,3,opt,name=exchange" json:"exchange,omitempty"`
	Key         string `protobuf:"bytes,4,opt,name=key" json:"key,omitempty"`
	Secret      string `protobuf:"bytes,5,opt,name=secret" json:"secret,omitempty"`
	Description string `protobuf:"bytes,6,opt,name=description" json:"description,omitempty"`
	Status      string `protobuf:"bytes,7,opt,name=status" json:"status,omitempty"`
}

func (m *ApiKeyRequest) Reset()                    { *m = ApiKeyRequest{} }
func (m *ApiKeyRequest) String() string            { return proto.CompactTextString(m) }
func (*ApiKeyRequest) ProtoMessage()               {}
func (*ApiKeyRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ApiKeyRequest) GetApiKeyId() string {
	if m != nil {
		return m.ApiKeyId
	}
	return ""
}

func (m *ApiKeyRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *ApiKeyRequest) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *ApiKeyRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *ApiKeyRequest) GetSecret() string {
	if m != nil {
		return m.Secret
	}
	return ""
}

func (m *ApiKeyRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *ApiKeyRequest) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type GetUserApiKeyRequest struct {
	UserId   string `protobuf:"bytes,1,opt,name=userId" json:"userId,omitempty"`
	ApiKeyId string `protobuf:"bytes,2,opt,name=apiKeyId" json:"apiKeyId,omitempty"`
}

func (m *GetUserApiKeyRequest) Reset()                    { *m = GetUserApiKeyRequest{} }
func (m *GetUserApiKeyRequest) String() string            { return proto.CompactTextString(m) }
func (*GetUserApiKeyRequest) ProtoMessage()               {}
func (*GetUserApiKeyRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *GetUserApiKeyRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *GetUserApiKeyRequest) GetApiKeyId() string {
	if m != nil {
		return m.ApiKeyId
	}
	return ""
}

type GetUserApiKeysRequest struct {
	UserId string `protobuf:"bytes,1,opt,name=userId" json:"userId,omitempty"`
}

func (m *GetUserApiKeysRequest) Reset()                    { *m = GetUserApiKeysRequest{} }
func (m *GetUserApiKeysRequest) String() string            { return proto.CompactTextString(m) }
func (*GetUserApiKeysRequest) ProtoMessage()               {}
func (*GetUserApiKeysRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *GetUserApiKeysRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

type RemoveApiKeyRequest struct {
	UserId   string `protobuf:"bytes,1,opt,name=userId" json:"userId,omitempty"`
	ApiKeyId string `protobuf:"bytes,2,opt,name=apiKeyId" json:"apiKeyId,omitempty"`
}

func (m *RemoveApiKeyRequest) Reset()                    { *m = RemoveApiKeyRequest{} }
func (m *RemoveApiKeyRequest) String() string            { return proto.CompactTextString(m) }
func (*RemoveApiKeyRequest) ProtoMessage()               {}
func (*RemoveApiKeyRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *RemoveApiKeyRequest) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *RemoveApiKeyRequest) GetApiKeyId() string {
	if m != nil {
		return m.ApiKeyId
	}
	return ""
}

// Responses
type ApiKey struct {
	ApiKeyId    string `protobuf:"bytes,1,opt,name=apiKeyId" json:"apiKeyId,omitempty"`
	UserId      string `protobuf:"bytes,2,opt,name=userId" json:"userId,omitempty"`
	Exchange    string `protobuf:"bytes,3,opt,name=exchange" json:"exchange,omitempty"`
	Key         string `protobuf:"bytes,4,opt,name=key" json:"key,omitempty"`
	Secret      string `protobuf:"bytes,5,opt,name=secret" json:"secret,omitempty"`
	Status      string `protobuf:"bytes,6,opt,name=status" json:"status,omitempty"`
	Description string `protobuf:"bytes,7,opt,name=description" json:"description,omitempty"`
}

func (m *ApiKey) Reset()                    { *m = ApiKey{} }
func (m *ApiKey) String() string            { return proto.CompactTextString(m) }
func (*ApiKey) ProtoMessage()               {}
func (*ApiKey) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *ApiKey) GetApiKeyId() string {
	if m != nil {
		return m.ApiKeyId
	}
	return ""
}

func (m *ApiKey) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *ApiKey) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *ApiKey) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *ApiKey) GetSecret() string {
	if m != nil {
		return m.Secret
	}
	return ""
}

func (m *ApiKey) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *ApiKey) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

type UserApiKeyData struct {
	ApiKey *ApiKey `protobuf:"bytes,1,opt,name=apiKey" json:"apiKey,omitempty"`
}

func (m *UserApiKeyData) Reset()                    { *m = UserApiKeyData{} }
func (m *UserApiKeyData) String() string            { return proto.CompactTextString(m) }
func (*UserApiKeyData) ProtoMessage()               {}
func (*UserApiKeyData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *UserApiKeyData) GetApiKey() *ApiKey {
	if m != nil {
		return m.ApiKey
	}
	return nil
}

type UserApiKeysData struct {
	ApiKey []*ApiKey `protobuf:"bytes,1,rep,name=apiKey" json:"apiKey,omitempty"`
}

func (m *UserApiKeysData) Reset()                    { *m = UserApiKeysData{} }
func (m *UserApiKeysData) String() string            { return proto.CompactTextString(m) }
func (*UserApiKeysData) ProtoMessage()               {}
func (*UserApiKeysData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *UserApiKeysData) GetApiKey() []*ApiKey {
	if m != nil {
		return m.ApiKey
	}
	return nil
}

type ApiKeyResponse struct {
	Status  string          `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	Message string          `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	Data    *UserApiKeyData `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *ApiKeyResponse) Reset()                    { *m = ApiKeyResponse{} }
func (m *ApiKeyResponse) String() string            { return proto.CompactTextString(m) }
func (*ApiKeyResponse) ProtoMessage()               {}
func (*ApiKeyResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *ApiKeyResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *ApiKeyResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *ApiKeyResponse) GetData() *UserApiKeyData {
	if m != nil {
		return m.Data
	}
	return nil
}

type ApiKeyListResponse struct {
	Status  string           `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	Message string           `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	Data    *UserApiKeysData `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *ApiKeyListResponse) Reset()                    { *m = ApiKeyListResponse{} }
func (m *ApiKeyListResponse) String() string            { return proto.CompactTextString(m) }
func (*ApiKeyListResponse) ProtoMessage()               {}
func (*ApiKeyListResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *ApiKeyListResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *ApiKeyListResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *ApiKeyListResponse) GetData() *UserApiKeysData {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*ApiKeyRequest)(nil), "go.micro.srv.apikey.ApiKeyRequest")
	proto.RegisterType((*GetUserApiKeyRequest)(nil), "go.micro.srv.apikey.GetUserApiKeyRequest")
	proto.RegisterType((*GetUserApiKeysRequest)(nil), "go.micro.srv.apikey.GetUserApiKeysRequest")
	proto.RegisterType((*RemoveApiKeyRequest)(nil), "go.micro.srv.apikey.RemoveApiKeyRequest")
	proto.RegisterType((*ApiKey)(nil), "go.micro.srv.apikey.ApiKey")
	proto.RegisterType((*UserApiKeyData)(nil), "go.micro.srv.apikey.UserApiKeyData")
	proto.RegisterType((*UserApiKeysData)(nil), "go.micro.srv.apikey.UserApiKeysData")
	proto.RegisterType((*ApiKeyResponse)(nil), "go.micro.srv.apikey.ApiKeyResponse")
	proto.RegisterType((*ApiKeyListResponse)(nil), "go.micro.srv.apikey.ApiKeyListResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for ApiKeyService service

type ApiKeyServiceClient interface {
	AddApiKey(ctx context.Context, in *ApiKeyRequest, opts ...client.CallOption) (*ApiKeyResponse, error)
	GetUserApiKey(ctx context.Context, in *GetUserApiKeyRequest, opts ...client.CallOption) (*ApiKeyResponse, error)
	GetUserApiKeys(ctx context.Context, in *GetUserApiKeysRequest, opts ...client.CallOption) (*ApiKeyListResponse, error)
	RemoveApiKey(ctx context.Context, in *RemoveApiKeyRequest, opts ...client.CallOption) (*ApiKeyResponse, error)
	UpdateApiKey(ctx context.Context, in *ApiKeyRequest, opts ...client.CallOption) (*ApiKeyResponse, error)
}

type apiKeyServiceClient struct {
	c           client.Client
	serviceName string
}

func NewApiKeyServiceClient(serviceName string, c client.Client) ApiKeyServiceClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "go.micro.srv.apikey"
	}
	return &apiKeyServiceClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *apiKeyServiceClient) AddApiKey(ctx context.Context, in *ApiKeyRequest, opts ...client.CallOption) (*ApiKeyResponse, error) {
	req := c.c.NewRequest(c.serviceName, "ApiKeyService.AddApiKey", in)
	out := new(ApiKeyResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiKeyServiceClient) GetUserApiKey(ctx context.Context, in *GetUserApiKeyRequest, opts ...client.CallOption) (*ApiKeyResponse, error) {
	req := c.c.NewRequest(c.serviceName, "ApiKeyService.GetUserApiKey", in)
	out := new(ApiKeyResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiKeyServiceClient) GetUserApiKeys(ctx context.Context, in *GetUserApiKeysRequest, opts ...client.CallOption) (*ApiKeyListResponse, error) {
	req := c.c.NewRequest(c.serviceName, "ApiKeyService.GetUserApiKeys", in)
	out := new(ApiKeyListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiKeyServiceClient) RemoveApiKey(ctx context.Context, in *RemoveApiKeyRequest, opts ...client.CallOption) (*ApiKeyResponse, error) {
	req := c.c.NewRequest(c.serviceName, "ApiKeyService.RemoveApiKey", in)
	out := new(ApiKeyResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiKeyServiceClient) UpdateApiKey(ctx context.Context, in *ApiKeyRequest, opts ...client.CallOption) (*ApiKeyResponse, error) {
	req := c.c.NewRequest(c.serviceName, "ApiKeyService.UpdateApiKey", in)
	out := new(ApiKeyResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ApiKeyService service

type ApiKeyServiceHandler interface {
	AddApiKey(context.Context, *ApiKeyRequest, *ApiKeyResponse) error
	GetUserApiKey(context.Context, *GetUserApiKeyRequest, *ApiKeyResponse) error
	GetUserApiKeys(context.Context, *GetUserApiKeysRequest, *ApiKeyListResponse) error
	RemoveApiKey(context.Context, *RemoveApiKeyRequest, *ApiKeyResponse) error
	UpdateApiKey(context.Context, *ApiKeyRequest, *ApiKeyResponse) error
}

func RegisterApiKeyServiceHandler(s server.Server, hdlr ApiKeyServiceHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&ApiKeyService{hdlr}, opts...))
}

type ApiKeyService struct {
	ApiKeyServiceHandler
}

func (h *ApiKeyService) AddApiKey(ctx context.Context, in *ApiKeyRequest, out *ApiKeyResponse) error {
	return h.ApiKeyServiceHandler.AddApiKey(ctx, in, out)
}

func (h *ApiKeyService) GetUserApiKey(ctx context.Context, in *GetUserApiKeyRequest, out *ApiKeyResponse) error {
	return h.ApiKeyServiceHandler.GetUserApiKey(ctx, in, out)
}

func (h *ApiKeyService) GetUserApiKeys(ctx context.Context, in *GetUserApiKeysRequest, out *ApiKeyListResponse) error {
	return h.ApiKeyServiceHandler.GetUserApiKeys(ctx, in, out)
}

func (h *ApiKeyService) RemoveApiKey(ctx context.Context, in *RemoveApiKeyRequest, out *ApiKeyResponse) error {
	return h.ApiKeyServiceHandler.RemoveApiKey(ctx, in, out)
}

func (h *ApiKeyService) UpdateApiKey(ctx context.Context, in *ApiKeyRequest, out *ApiKeyResponse) error {
	return h.ApiKeyServiceHandler.UpdateApiKey(ctx, in, out)
}

func init() { proto.RegisterFile("proto/apikey/apikey.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 463 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x54, 0xcf, 0xaf, 0xd2, 0x40,
	0x10, 0x7e, 0x7d, 0x7d, 0x16, 0xdf, 0xf0, 0x1e, 0x9a, 0x7d, 0x6a, 0x56, 0xbc, 0x90, 0xc5, 0x44,
	0xf4, 0x50, 0x12, 0x38, 0xe8, 0x95, 0xc4, 0x1f, 0x41, 0x3d, 0xd5, 0x60, 0xe2, 0xc9, 0xac, 0xed,
	0xa4, 0x36, 0x04, 0x5a, 0x77, 0x17, 0x22, 0xf1, 0xe2, 0x1f, 0xe6, 0xcd, 0x3f, 0xc8, 0x7f, 0xc1,
	0xb0, 0xdb, 0xd2, 0x2e, 0x34, 0x40, 0x82, 0x89, 0x27, 0x98, 0xd9, 0x9d, 0x99, 0xef, 0xfb, 0x66,
	0xbf, 0xc2, 0xc3, 0x4c, 0xa4, 0x2a, 0xed, 0xf3, 0x2c, 0x99, 0xe2, 0x2a, 0xff, 0xf1, 0x75, 0x8e,
	0xdc, 0xc4, 0xa9, 0x3f, 0x4b, 0x42, 0x91, 0xfa, 0x52, 0x2c, 0x7d, 0x73, 0xc4, 0x7e, 0x3b, 0x70,
	0x3d, 0xca, 0x92, 0x77, 0xb8, 0x0a, 0xf0, 0xdb, 0x02, 0xa5, 0x22, 0x6d, 0xb8, 0xcd, 0x75, 0x62,
	0x1c, 0x51, 0xa7, 0xe3, 0xf4, 0x2e, 0x83, 0x4d, 0x4c, 0x1e, 0x80, 0xb7, 0x90, 0x28, 0xc6, 0x11,
	0x3d, 0xd7, 0x27, 0x79, 0xb4, 0xae, 0xc1, 0xef, 0xe1, 0x57, 0x3e, 0x8f, 0x91, 0xba, 0xa6, 0xa6,
	0x88, 0xc9, 0x5d, 0x70, 0xa7, 0xb8, 0xa2, 0x17, 0x3a, 0xbd, 0xfe, 0xbb, 0xee, 0x22, 0x31, 0x14,
	0xa8, 0xe8, 0x2d, 0xd3, 0xc5, 0x44, 0xa4, 0x03, 0xcd, 0x08, 0x65, 0x28, 0x92, 0x4c, 0x25, 0xe9,
	0x9c, 0x7a, 0xfa, 0xb0, 0x9a, 0xd2, 0x95, 0x8a, 0xab, 0x85, 0xa4, 0x8d, 0xbc, 0x52, 0x47, 0xec,
	0x2d, 0xdc, 0x7b, 0x83, 0x6a, 0x22, 0x51, 0xd8, 0x5c, 0x4a, 0xbc, 0xce, 0x36, 0xde, 0x0d, 0xc7,
	0x73, 0x9b, 0x23, 0xeb, 0xc3, 0x7d, 0xab, 0x97, 0x3c, 0xd0, 0x8c, 0x8d, 0xe1, 0x26, 0xc0, 0x59,
	0xba, 0xc4, 0xd3, 0x67, 0xff, 0x72, 0xc0, 0x33, 0x5d, 0xfe, 0xe3, 0x1a, 0x4a, 0x91, 0xbd, 0xaa,
	0xc8, 0xdb, 0xeb, 0x69, 0xec, 0xac, 0x87, 0xbd, 0x82, 0x56, 0xa9, 0xdb, 0x4b, 0xae, 0x38, 0x19,
	0x82, 0x67, 0x50, 0x6b, 0x0e, 0xcd, 0xc1, 0x23, 0xbf, 0xe6, 0x11, 0xfa, 0xb9, 0x70, 0xf9, 0x55,
	0xf6, 0x1a, 0xee, 0x54, 0xe4, 0xdf, 0xe9, 0xe3, 0x1e, 0xdb, 0xe7, 0x07, 0xb4, 0x8a, 0x95, 0xc8,
	0x2c, 0x9d, 0x4b, 0xac, 0x50, 0x73, 0x2c, 0x6a, 0x14, 0x1a, 0x33, 0x94, 0x92, 0xc7, 0x98, 0x2b,
	0x5a, 0x84, 0xe4, 0x39, 0x5c, 0x44, 0x5c, 0x71, 0x2d, 0x67, 0x73, 0xd0, 0xad, 0x1d, 0x6b, 0x73,
	0x0e, 0x74, 0x01, 0xfb, 0xe9, 0x00, 0x31, 0xc9, 0xf7, 0x89, 0x54, 0x27, 0x20, 0x78, 0x61, 0x21,
	0x78, 0x7c, 0x00, 0x81, 0x2c, 0x21, 0x0c, 0xfe, 0xb8, 0x85, 0xb7, 0x3f, 0xa0, 0x58, 0x26, 0x21,
	0x92, 0x8f, 0x70, 0x39, 0x8a, 0xa2, 0xfc, 0x85, 0xb1, 0x7d, 0x1a, 0x9a, 0x47, 0xdc, 0xee, 0xee,
	0xbd, 0x63, 0x38, 0xb1, 0x33, 0xc2, 0xe1, 0xda, 0xf2, 0x0c, 0x79, 0x5a, 0x5b, 0x57, 0xe7, 0xd1,
	0x63, 0x47, 0xc4, 0xd0, 0xb2, 0x6d, 0x49, 0x9e, 0x1d, 0x9e, 0x51, 0x78, 0xb7, 0xfd, 0x64, 0xcf,
	0x90, 0xea, 0x7e, 0xd8, 0x19, 0xf9, 0x0c, 0x57, 0x55, 0x3b, 0x93, 0x5e, 0x6d, 0x69, 0x8d, 0xe3,
	0x8f, 0x65, 0xf2, 0x09, 0xae, 0x26, 0x59, 0xc4, 0x15, 0xfe, 0xf3, 0x3d, 0x7c, 0xf1, 0xf4, 0x97,
	0x7e, 0xf8, 0x37, 0x00, 0x00, 0xff, 0xff, 0x4e, 0x71, 0x8c, 0xe9, 0x06, 0x06, 0x00, 0x00,
}