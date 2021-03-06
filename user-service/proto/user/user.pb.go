// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/user/user.proto

/*
Package user is a generated protocol buffer package.

It is generated from these files:
	proto/user/user.proto

It has these top-level messages:
	ChangePasswordRequest
	CreateUserRequest
	DeleteUserRequest
	GetUserInfoRequest
	UpdateUserRequest
	Response
	User
	UserData
	UserResponse
*/
package user

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
type ChangePasswordRequest struct {
	UserID      string `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	NewPassword string `protobuf:"bytes,2,opt,name=newPassword" json:"newPassword,omitempty"`
	OldPassword string `protobuf:"bytes,3,opt,name=oldPassword" json:"oldPassword,omitempty"`
}

func (m *ChangePasswordRequest) Reset()                    { *m = ChangePasswordRequest{} }
func (m *ChangePasswordRequest) String() string            { return proto.CompactTextString(m) }
func (*ChangePasswordRequest) ProtoMessage()               {}
func (*ChangePasswordRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ChangePasswordRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *ChangePasswordRequest) GetNewPassword() string {
	if m != nil {
		return m.NewPassword
	}
	return ""
}

func (m *ChangePasswordRequest) GetOldPassword() string {
	if m != nil {
		return m.OldPassword
	}
	return ""
}

type CreateUserRequest struct {
	First    string `protobuf:"bytes,1,opt,name=first" json:"first,omitempty"`
	Last     string `protobuf:"bytes,2,opt,name=last" json:"last,omitempty"`
	Email    string `protobuf:"bytes,3,opt,name=email" json:"email,omitempty"`
	Password string `protobuf:"bytes,4,opt,name=password" json:"password,omitempty"`
}

func (m *CreateUserRequest) Reset()                    { *m = CreateUserRequest{} }
func (m *CreateUserRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateUserRequest) ProtoMessage()               {}
func (*CreateUserRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *CreateUserRequest) GetFirst() string {
	if m != nil {
		return m.First
	}
	return ""
}

func (m *CreateUserRequest) GetLast() string {
	if m != nil {
		return m.Last
	}
	return ""
}

func (m *CreateUserRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *CreateUserRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type DeleteUserRequest struct {
	UserID string `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	Hard   bool   `protobuf:"varint,2,opt,name=hard" json:"hard,omitempty"`
}

func (m *DeleteUserRequest) Reset()                    { *m = DeleteUserRequest{} }
func (m *DeleteUserRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteUserRequest) ProtoMessage()               {}
func (*DeleteUserRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *DeleteUserRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *DeleteUserRequest) GetHard() bool {
	if m != nil {
		return m.Hard
	}
	return false
}

type GetUserInfoRequest struct {
	UserID string `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
}

func (m *GetUserInfoRequest) Reset()                    { *m = GetUserInfoRequest{} }
func (m *GetUserInfoRequest) String() string            { return proto.CompactTextString(m) }
func (*GetUserInfoRequest) ProtoMessage()               {}
func (*GetUserInfoRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *GetUserInfoRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

type UpdateUserRequest struct {
	UserID string `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	First  string `protobuf:"bytes,2,opt,name=first" json:"first,omitempty"`
	Last   string `protobuf:"bytes,3,opt,name=last" json:"last,omitempty"`
	Email  string `protobuf:"bytes,4,opt,name=email" json:"email,omitempty"`
}

func (m *UpdateUserRequest) Reset()                    { *m = UpdateUserRequest{} }
func (m *UpdateUserRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdateUserRequest) ProtoMessage()               {}
func (*UpdateUserRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *UpdateUserRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *UpdateUserRequest) GetFirst() string {
	if m != nil {
		return m.First
	}
	return ""
}

func (m *UpdateUserRequest) GetLast() string {
	if m != nil {
		return m.Last
	}
	return ""
}

func (m *UpdateUserRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

// Responses
type Response struct {
	Status  string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Response) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Response) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type User struct {
	UserID string `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	First  string `protobuf:"bytes,2,opt,name=first" json:"first,omitempty"`
	Last   string `protobuf:"bytes,3,opt,name=last" json:"last,omitempty"`
	Email  string `protobuf:"bytes,4,opt,name=email" json:"email,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *User) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *User) GetFirst() string {
	if m != nil {
		return m.First
	}
	return ""
}

func (m *User) GetLast() string {
	if m != nil {
		return m.Last
	}
	return ""
}

func (m *User) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

type UserData struct {
	User *User `protobuf:"bytes,1,opt,name=user" json:"user,omitempty"`
}

func (m *UserData) Reset()                    { *m = UserData{} }
func (m *UserData) String() string            { return proto.CompactTextString(m) }
func (*UserData) ProtoMessage()               {}
func (*UserData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *UserData) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

type UserResponse struct {
	Status  string    `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	Message string    `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	Data    *UserData `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *UserResponse) Reset()                    { *m = UserResponse{} }
func (m *UserResponse) String() string            { return proto.CompactTextString(m) }
func (*UserResponse) ProtoMessage()               {}
func (*UserResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *UserResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *UserResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *UserResponse) GetData() *UserData {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*ChangePasswordRequest)(nil), "ChangePasswordRequest")
	proto.RegisterType((*CreateUserRequest)(nil), "CreateUserRequest")
	proto.RegisterType((*DeleteUserRequest)(nil), "DeleteUserRequest")
	proto.RegisterType((*GetUserInfoRequest)(nil), "GetUserInfoRequest")
	proto.RegisterType((*UpdateUserRequest)(nil), "UpdateUserRequest")
	proto.RegisterType((*Response)(nil), "Response")
	proto.RegisterType((*User)(nil), "User")
	proto.RegisterType((*UserData)(nil), "UserData")
	proto.RegisterType((*UserResponse)(nil), "UserResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for UserService service

type UserServiceClient interface {
	ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...client.CallOption) (*Response, error)
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...client.CallOption) (*UserResponse, error)
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...client.CallOption) (*Response, error)
	GetUserInfo(ctx context.Context, in *GetUserInfoRequest, opts ...client.CallOption) (*UserResponse, error)
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...client.CallOption) (*UserResponse, error)
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
		serviceName = "userservice"
	}
	return &userServiceClient{
		c:           c,
		serviceName: serviceName,
	}
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

func (c *userServiceClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...client.CallOption) (*UserResponse, error) {
	req := c.c.NewRequest(c.serviceName, "UserService.CreateUser", in)
	out := new(UserResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.serviceName, "UserService.DeleteUser", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUserInfo(ctx context.Context, in *GetUserInfoRequest, opts ...client.CallOption) (*UserResponse, error) {
	req := c.c.NewRequest(c.serviceName, "UserService.GetUserInfo", in)
	out := new(UserResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...client.CallOption) (*UserResponse, error) {
	req := c.c.NewRequest(c.serviceName, "UserService.UpdateUser", in)
	out := new(UserResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UserService service

type UserServiceHandler interface {
	ChangePassword(context.Context, *ChangePasswordRequest, *Response) error
	CreateUser(context.Context, *CreateUserRequest, *UserResponse) error
	DeleteUser(context.Context, *DeleteUserRequest, *Response) error
	GetUserInfo(context.Context, *GetUserInfoRequest, *UserResponse) error
	UpdateUser(context.Context, *UpdateUserRequest, *UserResponse) error
}

func RegisterUserServiceHandler(s server.Server, hdlr UserServiceHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&UserService{hdlr}, opts...))
}

type UserService struct {
	UserServiceHandler
}

func (h *UserService) ChangePassword(ctx context.Context, in *ChangePasswordRequest, out *Response) error {
	return h.UserServiceHandler.ChangePassword(ctx, in, out)
}

func (h *UserService) CreateUser(ctx context.Context, in *CreateUserRequest, out *UserResponse) error {
	return h.UserServiceHandler.CreateUser(ctx, in, out)
}

func (h *UserService) DeleteUser(ctx context.Context, in *DeleteUserRequest, out *Response) error {
	return h.UserServiceHandler.DeleteUser(ctx, in, out)
}

func (h *UserService) GetUserInfo(ctx context.Context, in *GetUserInfoRequest, out *UserResponse) error {
	return h.UserServiceHandler.GetUserInfo(ctx, in, out)
}

func (h *UserService) UpdateUser(ctx context.Context, in *UpdateUserRequest, out *UserResponse) error {
	return h.UserServiceHandler.UpdateUser(ctx, in, out)
}

func init() { proto.RegisterFile("proto/user/user.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 405 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x93, 0x5f, 0x4b, 0xeb, 0x40,
	0x10, 0xc5, 0xfb, 0x27, 0xed, 0x4d, 0x27, 0xf7, 0x5e, 0xe8, 0xdc, 0xdb, 0x52, 0x0b, 0x42, 0x59,
	0x10, 0x7c, 0xd0, 0x95, 0xb6, 0xf8, 0x26, 0xf8, 0xd0, 0x82, 0xf4, 0x4d, 0x22, 0x7d, 0x55, 0x56,
	0x33, 0xfd, 0x83, 0x69, 0x12, 0xb3, 0x5b, 0xfb, 0x2d, 0xfc, 0xcc, 0xb2, 0x9b, 0xa4, 0x4d, 0x49,
	0x45, 0x11, 0x7c, 0x09, 0x39, 0x93, 0xc9, 0xfc, 0xe6, 0x2c, 0x67, 0xa1, 0x15, 0xc5, 0xa1, 0x0a,
	0x2f, 0xd6, 0x92, 0x62, 0xf3, 0xe0, 0x46, 0x33, 0x09, 0xad, 0xd1, 0x42, 0x04, 0x73, 0xba, 0x15,
	0x52, 0x6e, 0xc2, 0xd8, 0x73, 0xe9, 0x65, 0x4d, 0x52, 0x61, 0x1b, 0xea, 0xba, 0x6d, 0x32, 0xee,
	0x94, 0x7b, 0xe5, 0xd3, 0x86, 0x9b, 0x2a, 0xec, 0x81, 0x13, 0xd0, 0x26, 0xeb, 0xee, 0x54, 0xcc,
	0xc7, 0x7c, 0x49, 0x77, 0x84, 0xbe, 0xb7, 0xed, 0xa8, 0x26, 0x1d, 0xb9, 0x12, 0x0b, 0xa1, 0x39,
	0x8a, 0x49, 0x28, 0x9a, 0x4a, 0x8a, 0x33, 0xe0, 0x7f, 0xa8, 0xcd, 0x96, 0xb1, 0x54, 0x29, 0x2f,
	0x11, 0x88, 0x60, 0xf9, 0x42, 0xaa, 0x94, 0x63, 0xde, 0x75, 0x27, 0xad, 0xc4, 0xd2, 0x4f, 0x47,
	0x27, 0x02, 0xbb, 0x60, 0x47, 0x19, 0xd3, 0x32, 0x1f, 0xb6, 0x9a, 0x5d, 0x43, 0x73, 0x4c, 0x3e,
	0xed, 0x03, 0x3f, 0x72, 0x88, 0x60, 0x2d, 0x44, 0x6a, 0xcd, 0x76, 0xcd, 0x3b, 0x3b, 0x03, 0xbc,
	0x21, 0xa5, 0xff, 0x9e, 0x04, 0xb3, 0xf0, 0x93, 0x09, 0xec, 0x19, 0x9a, 0xd3, 0xc8, 0x13, 0x5f,
	0xc3, 0x6d, 0x7d, 0x57, 0x0e, 0xf9, 0xae, 0x1e, 0xf2, 0x6d, 0xe5, 0x7c, 0xb3, 0x2b, 0xb0, 0x5d,
	0x92, 0x51, 0x18, 0x48, 0xd2, 0x0c, 0xa9, 0x84, 0x5a, 0xcb, 0x8c, 0x91, 0x28, 0xec, 0xc0, 0xaf,
	0x15, 0x49, 0x29, 0xe6, 0x94, 0x52, 0x32, 0xc9, 0xee, 0xc1, 0xd2, 0x4b, 0xfe, 0xd8, 0x76, 0x27,
	0x60, 0xeb, 0xf9, 0x63, 0xa1, 0x04, 0x1e, 0x81, 0xa5, 0xa7, 0x1a, 0x82, 0x33, 0xa8, 0x71, 0x73,
	0x3a, 0xa6, 0xc4, 0x1e, 0xe0, 0x77, 0x72, 0x56, 0xdf, 0x35, 0x82, 0xc7, 0x60, 0x79, 0x42, 0x09,
	0xb3, 0x92, 0x33, 0x68, 0xf0, 0x8c, 0xea, 0x9a, 0xf2, 0xe0, 0xad, 0x02, 0x8e, 0x2e, 0xdd, 0x51,
	0xfc, 0xba, 0x7c, 0x22, 0xbc, 0x84, 0xbf, 0xfb, 0xb9, 0xc7, 0x36, 0x3f, 0x78, 0x11, 0xba, 0x0d,
	0x9e, 0x6d, 0xc5, 0x4a, 0xd8, 0x07, 0xd8, 0x25, 0x17, 0x91, 0x17, 0x62, 0xdc, 0xfd, 0xc3, 0xf3,
	0x46, 0x58, 0x09, 0xcf, 0x01, 0x76, 0xd9, 0x43, 0xe4, 0x85, 0x20, 0xee, 0x13, 0x86, 0xe0, 0xe4,
	0x92, 0x86, 0xff, 0x78, 0x31, 0x77, 0x45, 0x46, 0x1f, 0x60, 0x17, 0x38, 0x44, 0x5e, 0x48, 0x5f,
	0xe1, 0x97, 0xc7, 0xba, 0xb9, 0xff, 0xc3, 0xf7, 0x00, 0x00, 0x00, 0xff, 0xff, 0x71, 0xd4, 0xb5,
	0x1d, 0x18, 0x04, 0x00, 0x00,
}
