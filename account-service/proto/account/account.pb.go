// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/asciiu/gomo/account-service/proto/account/account.proto

/*
Package account is a generated protocol buffer package.

It is generated from these files:
	github.com/asciiu/gomo/account-service/proto/account/account.proto

It has these top-level messages:
	NewAccountRequest
	UpdateAccountRequest
	AccountRequest
	AccountsRequest
	Account
	UserAccount
	AccountResponse
	UserAccounts
	AccountsResponse
*/
package account

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import bal "github.com/asciiu/gomo/account-service/proto/balance"

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
type NewAccountRequest struct {
	UserID      string                   `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	Exchange    string                   `protobuf:"bytes,2,opt,name=exchange" json:"exchange,omitempty"`
	KeyPublic   string                   `protobuf:"bytes,3,opt,name=keyPublic" json:"keyPublic,omitempty"`
	KeySecret   string                   `protobuf:"bytes,4,opt,name=keySecret" json:"keySecret,omitempty"`
	Description string                   `protobuf:"bytes,5,opt,name=description" json:"description,omitempty"`
	Balances    []*bal.NewBalanceRequest `protobuf:"bytes,6,rep,name=balances" json:"balances,omitempty"`
}

func (m *NewAccountRequest) Reset()                    { *m = NewAccountRequest{} }
func (m *NewAccountRequest) String() string            { return proto.CompactTextString(m) }
func (*NewAccountRequest) ProtoMessage()               {}
func (*NewAccountRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *NewAccountRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *NewAccountRequest) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *NewAccountRequest) GetKeyPublic() string {
	if m != nil {
		return m.KeyPublic
	}
	return ""
}

func (m *NewAccountRequest) GetKeySecret() string {
	if m != nil {
		return m.KeySecret
	}
	return ""
}

func (m *NewAccountRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *NewAccountRequest) GetBalances() []*bal.NewBalanceRequest {
	if m != nil {
		return m.Balances
	}
	return nil
}

type UpdateAccountRequest struct {
	AccountID   string                   `protobuf:"bytes,1,opt,name=accountID" json:"accountID,omitempty"`
	UserID      string                   `protobuf:"bytes,2,opt,name=userID" json:"userID,omitempty"`
	KeyPublic   string                   `protobuf:"bytes,3,opt,name=keyPublic" json:"keyPublic,omitempty"`
	KeySecret   string                   `protobuf:"bytes,4,opt,name=keySecret" json:"keySecret,omitempty"`
	Description string                   `protobuf:"bytes,5,opt,name=description" json:"description,omitempty"`
	Balances    []*bal.NewBalanceRequest `protobuf:"bytes,6,rep,name=balances" json:"balances,omitempty"`
}

func (m *UpdateAccountRequest) Reset()                    { *m = UpdateAccountRequest{} }
func (m *UpdateAccountRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdateAccountRequest) ProtoMessage()               {}
func (*UpdateAccountRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *UpdateAccountRequest) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *UpdateAccountRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *UpdateAccountRequest) GetKeyPublic() string {
	if m != nil {
		return m.KeyPublic
	}
	return ""
}

func (m *UpdateAccountRequest) GetKeySecret() string {
	if m != nil {
		return m.KeySecret
	}
	return ""
}

func (m *UpdateAccountRequest) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *UpdateAccountRequest) GetBalances() []*bal.NewBalanceRequest {
	if m != nil {
		return m.Balances
	}
	return nil
}

type AccountRequest struct {
	AccountID string `protobuf:"bytes,1,opt,name=accountID" json:"accountID,omitempty"`
	UserID    string `protobuf:"bytes,2,opt,name=userID" json:"userID,omitempty"`
}

func (m *AccountRequest) Reset()                    { *m = AccountRequest{} }
func (m *AccountRequest) String() string            { return proto.CompactTextString(m) }
func (*AccountRequest) ProtoMessage()               {}
func (*AccountRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *AccountRequest) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *AccountRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

type AccountsRequest struct {
	UserID string `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
}

func (m *AccountsRequest) Reset()                    { *m = AccountsRequest{} }
func (m *AccountsRequest) String() string            { return proto.CompactTextString(m) }
func (*AccountsRequest) ProtoMessage()               {}
func (*AccountsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *AccountsRequest) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

// Responses
type Account struct {
	AccountID   string         `protobuf:"bytes,1,opt,name=accountID" json:"accountID,omitempty"`
	UserID      string         `protobuf:"bytes,2,opt,name=userID" json:"userID,omitempty"`
	Exchange    string         `protobuf:"bytes,3,opt,name=exchange" json:"exchange,omitempty"`
	KeyPublic   string         `protobuf:"bytes,4,opt,name=keyPublic" json:"keyPublic,omitempty"`
	KeySecret   string         `protobuf:"bytes,5,opt,name=keySecret" json:"keySecret,omitempty"`
	Title       string         `protobuf:"bytes,6,opt,name=title" json:"title,omitempty"`
	Description string         `protobuf:"bytes,7,opt,name=description" json:"description,omitempty"`
	Status      string         `protobuf:"bytes,8,opt,name=status" json:"status,omitempty"`
	CreatedOn   string         `protobuf:"bytes,9,opt,name=createdOn" json:"createdOn,omitempty"`
	UpdatedOn   string         `protobuf:"bytes,10,opt,name=updatedOn" json:"updatedOn,omitempty"`
	Balances    []*bal.Balance `protobuf:"bytes,11,rep,name=balances" json:"balances,omitempty"`
}

func (m *Account) Reset()                    { *m = Account{} }
func (m *Account) String() string            { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()               {}
func (*Account) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Account) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *Account) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *Account) GetExchange() string {
	if m != nil {
		return m.Exchange
	}
	return ""
}

func (m *Account) GetKeyPublic() string {
	if m != nil {
		return m.KeyPublic
	}
	return ""
}

func (m *Account) GetKeySecret() string {
	if m != nil {
		return m.KeySecret
	}
	return ""
}

func (m *Account) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Account) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Account) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Account) GetCreatedOn() string {
	if m != nil {
		return m.CreatedOn
	}
	return ""
}

func (m *Account) GetUpdatedOn() string {
	if m != nil {
		return m.UpdatedOn
	}
	return ""
}

func (m *Account) GetBalances() []*bal.Balance {
	if m != nil {
		return m.Balances
	}
	return nil
}

type UserAccount struct {
	Account *Account `protobuf:"bytes,1,opt,name=account" json:"account,omitempty"`
}

func (m *UserAccount) Reset()                    { *m = UserAccount{} }
func (m *UserAccount) String() string            { return proto.CompactTextString(m) }
func (*UserAccount) ProtoMessage()               {}
func (*UserAccount) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *UserAccount) GetAccount() *Account {
	if m != nil {
		return m.Account
	}
	return nil
}

type AccountResponse struct {
	Status  string       `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	Message string       `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	Data    *UserAccount `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *AccountResponse) Reset()                    { *m = AccountResponse{} }
func (m *AccountResponse) String() string            { return proto.CompactTextString(m) }
func (*AccountResponse) ProtoMessage()               {}
func (*AccountResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *AccountResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *AccountResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *AccountResponse) GetData() *UserAccount {
	if m != nil {
		return m.Data
	}
	return nil
}

type UserAccounts struct {
	Accounts []*Account `protobuf:"bytes,1,rep,name=accounts" json:"accounts,omitempty"`
}

func (m *UserAccounts) Reset()                    { *m = UserAccounts{} }
func (m *UserAccounts) String() string            { return proto.CompactTextString(m) }
func (*UserAccounts) ProtoMessage()               {}
func (*UserAccounts) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *UserAccounts) GetAccounts() []*Account {
	if m != nil {
		return m.Accounts
	}
	return nil
}

type AccountsResponse struct {
	Status  string        `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	Message string        `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
	Data    *UserAccounts `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
}

func (m *AccountsResponse) Reset()                    { *m = AccountsResponse{} }
func (m *AccountsResponse) String() string            { return proto.CompactTextString(m) }
func (*AccountsResponse) ProtoMessage()               {}
func (*AccountsResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *AccountsResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *AccountsResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *AccountsResponse) GetData() *UserAccounts {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*NewAccountRequest)(nil), "account.NewAccountRequest")
	proto.RegisterType((*UpdateAccountRequest)(nil), "account.UpdateAccountRequest")
	proto.RegisterType((*AccountRequest)(nil), "account.AccountRequest")
	proto.RegisterType((*AccountsRequest)(nil), "account.AccountsRequest")
	proto.RegisterType((*Account)(nil), "account.Account")
	proto.RegisterType((*UserAccount)(nil), "account.UserAccount")
	proto.RegisterType((*AccountResponse)(nil), "account.AccountResponse")
	proto.RegisterType((*UserAccounts)(nil), "account.UserAccounts")
	proto.RegisterType((*AccountsResponse)(nil), "account.AccountsResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for AccountService service

type AccountServiceClient interface {
	AddAccount(ctx context.Context, in *NewAccountRequest, opts ...client.CallOption) (*AccountResponse, error)
	AddAccountBalance(ctx context.Context, in *bal.NewBalanceRequest, opts ...client.CallOption) (*bal.BalanceResponse, error)
	DeleteAccount(ctx context.Context, in *AccountRequest, opts ...client.CallOption) (*AccountResponse, error)
	GetAccounts(ctx context.Context, in *AccountsRequest, opts ...client.CallOption) (*AccountsResponse, error)
	GetAccount(ctx context.Context, in *AccountRequest, opts ...client.CallOption) (*AccountResponse, error)
	GetAccountBalance(ctx context.Context, in *bal.BalanceRequest, opts ...client.CallOption) (*bal.BalanceResponse, error)
	ResyncAccounts(ctx context.Context, in *AccountsRequest, opts ...client.CallOption) (*AccountsResponse, error)
	UpdateAccount(ctx context.Context, in *UpdateAccountRequest, opts ...client.CallOption) (*AccountResponse, error)
	UpdateAccountBalance(ctx context.Context, in *bal.UpdateBalanceRequest, opts ...client.CallOption) (*bal.BalanceResponse, error)
	ValidateAccountBalance(ctx context.Context, in *bal.ValidateBalanceRequest, opts ...client.CallOption) (*bal.ValidateBalanceResponse, error)
}

type accountServiceClient struct {
	c           client.Client
	serviceName string
}

func NewAccountServiceClient(serviceName string, c client.Client) AccountServiceClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "account"
	}
	return &accountServiceClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *accountServiceClient) AddAccount(ctx context.Context, in *NewAccountRequest, opts ...client.CallOption) (*AccountResponse, error) {
	req := c.c.NewRequest(c.serviceName, "AccountService.AddAccount", in)
	out := new(AccountResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) AddAccountBalance(ctx context.Context, in *bal.NewBalanceRequest, opts ...client.CallOption) (*bal.BalanceResponse, error) {
	req := c.c.NewRequest(c.serviceName, "AccountService.AddAccountBalance", in)
	out := new(bal.BalanceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) DeleteAccount(ctx context.Context, in *AccountRequest, opts ...client.CallOption) (*AccountResponse, error) {
	req := c.c.NewRequest(c.serviceName, "AccountService.DeleteAccount", in)
	out := new(AccountResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetAccounts(ctx context.Context, in *AccountsRequest, opts ...client.CallOption) (*AccountsResponse, error) {
	req := c.c.NewRequest(c.serviceName, "AccountService.GetAccounts", in)
	out := new(AccountsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetAccount(ctx context.Context, in *AccountRequest, opts ...client.CallOption) (*AccountResponse, error) {
	req := c.c.NewRequest(c.serviceName, "AccountService.GetAccount", in)
	out := new(AccountResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetAccountBalance(ctx context.Context, in *bal.BalanceRequest, opts ...client.CallOption) (*bal.BalanceResponse, error) {
	req := c.c.NewRequest(c.serviceName, "AccountService.GetAccountBalance", in)
	out := new(bal.BalanceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) ResyncAccounts(ctx context.Context, in *AccountsRequest, opts ...client.CallOption) (*AccountsResponse, error) {
	req := c.c.NewRequest(c.serviceName, "AccountService.ResyncAccounts", in)
	out := new(AccountsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) UpdateAccount(ctx context.Context, in *UpdateAccountRequest, opts ...client.CallOption) (*AccountResponse, error) {
	req := c.c.NewRequest(c.serviceName, "AccountService.UpdateAccount", in)
	out := new(AccountResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) UpdateAccountBalance(ctx context.Context, in *bal.UpdateBalanceRequest, opts ...client.CallOption) (*bal.BalanceResponse, error) {
	req := c.c.NewRequest(c.serviceName, "AccountService.UpdateAccountBalance", in)
	out := new(bal.BalanceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) ValidateAccountBalance(ctx context.Context, in *bal.ValidateBalanceRequest, opts ...client.CallOption) (*bal.ValidateBalanceResponse, error) {
	req := c.c.NewRequest(c.serviceName, "AccountService.ValidateAccountBalance", in)
	out := new(bal.ValidateBalanceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AccountService service

type AccountServiceHandler interface {
	AddAccount(context.Context, *NewAccountRequest, *AccountResponse) error
	AddAccountBalance(context.Context, *bal.NewBalanceRequest, *bal.BalanceResponse) error
	DeleteAccount(context.Context, *AccountRequest, *AccountResponse) error
	GetAccounts(context.Context, *AccountsRequest, *AccountsResponse) error
	GetAccount(context.Context, *AccountRequest, *AccountResponse) error
	GetAccountBalance(context.Context, *bal.BalanceRequest, *bal.BalanceResponse) error
	ResyncAccounts(context.Context, *AccountsRequest, *AccountsResponse) error
	UpdateAccount(context.Context, *UpdateAccountRequest, *AccountResponse) error
	UpdateAccountBalance(context.Context, *bal.UpdateBalanceRequest, *bal.BalanceResponse) error
	ValidateAccountBalance(context.Context, *bal.ValidateBalanceRequest, *bal.ValidateBalanceResponse) error
}

func RegisterAccountServiceHandler(s server.Server, hdlr AccountServiceHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&AccountService{hdlr}, opts...))
}

type AccountService struct {
	AccountServiceHandler
}

func (h *AccountService) AddAccount(ctx context.Context, in *NewAccountRequest, out *AccountResponse) error {
	return h.AccountServiceHandler.AddAccount(ctx, in, out)
}

func (h *AccountService) AddAccountBalance(ctx context.Context, in *bal.NewBalanceRequest, out *bal.BalanceResponse) error {
	return h.AccountServiceHandler.AddAccountBalance(ctx, in, out)
}

func (h *AccountService) DeleteAccount(ctx context.Context, in *AccountRequest, out *AccountResponse) error {
	return h.AccountServiceHandler.DeleteAccount(ctx, in, out)
}

func (h *AccountService) GetAccounts(ctx context.Context, in *AccountsRequest, out *AccountsResponse) error {
	return h.AccountServiceHandler.GetAccounts(ctx, in, out)
}

func (h *AccountService) GetAccount(ctx context.Context, in *AccountRequest, out *AccountResponse) error {
	return h.AccountServiceHandler.GetAccount(ctx, in, out)
}

func (h *AccountService) GetAccountBalance(ctx context.Context, in *bal.BalanceRequest, out *bal.BalanceResponse) error {
	return h.AccountServiceHandler.GetAccountBalance(ctx, in, out)
}

func (h *AccountService) ResyncAccounts(ctx context.Context, in *AccountsRequest, out *AccountsResponse) error {
	return h.AccountServiceHandler.ResyncAccounts(ctx, in, out)
}

func (h *AccountService) UpdateAccount(ctx context.Context, in *UpdateAccountRequest, out *AccountResponse) error {
	return h.AccountServiceHandler.UpdateAccount(ctx, in, out)
}

func (h *AccountService) UpdateAccountBalance(ctx context.Context, in *bal.UpdateBalanceRequest, out *bal.BalanceResponse) error {
	return h.AccountServiceHandler.UpdateAccountBalance(ctx, in, out)
}

func (h *AccountService) ValidateAccountBalance(ctx context.Context, in *bal.ValidateBalanceRequest, out *bal.ValidateBalanceResponse) error {
	return h.AccountServiceHandler.ValidateAccountBalance(ctx, in, out)
}

func init() {
	proto.RegisterFile("github.com/asciiu/gomo/account-service/proto/account/account.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 631 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x55, 0x4d, 0x6f, 0xd3, 0x4c,
	0x10, 0x7e, 0xd3, 0x34, 0x5f, 0xe3, 0xb6, 0x6f, 0xbb, 0x84, 0xe0, 0x9a, 0x22, 0x55, 0x3e, 0xa5,
	0x08, 0x12, 0x29, 0x9c, 0x90, 0x38, 0xd0, 0x12, 0x51, 0x95, 0x43, 0x41, 0xae, 0xca, 0x7d, 0xb3,
	0x1e, 0xb5, 0x16, 0x8e, 0x1d, 0xbc, 0x6b, 0x4a, 0x7f, 0x24, 0x3f, 0x80, 0x0b, 0x17, 0xfe, 0x08,
	0xca, 0xae, 0x77, 0x6d, 0x27, 0x26, 0x85, 0xc0, 0x81, 0x93, 0x35, 0xf3, 0x8c, 0x67, 0x9f, 0x67,
	0x3e, 0x76, 0xe1, 0xe4, 0x2a, 0x10, 0xd7, 0xe9, 0x64, 0xc0, 0xe2, 0xe9, 0x90, 0x72, 0x16, 0x04,
	0xe9, 0xf0, 0x2a, 0x9e, 0xc6, 0x43, 0xca, 0x58, 0x9c, 0x46, 0xe2, 0x29, 0xc7, 0xe4, 0x53, 0xc0,
	0x70, 0x38, 0x4b, 0x62, 0x61, 0xbc, 0xfa, 0x3b, 0x90, 0x5e, 0xd2, 0xca, 0x4c, 0xe7, 0xf7, 0x92,
	0x4d, 0x68, 0x48, 0x23, 0x86, 0xfa, 0xab, 0x92, 0xb9, 0x5f, 0x6b, 0xb0, 0x77, 0x8e, 0x37, 0xc7,
	0xea, 0x17, 0x0f, 0x3f, 0xa6, 0xc8, 0x05, 0xe9, 0x41, 0x33, 0xe5, 0x98, 0x9c, 0x8d, 0xed, 0xda,
	0x61, 0xad, 0xdf, 0xf1, 0x32, 0x8b, 0x38, 0xd0, 0xc6, 0xcf, 0xec, 0x9a, 0x46, 0x57, 0x68, 0x6f,
	0x48, 0xc4, 0xd8, 0xe4, 0x00, 0x3a, 0x1f, 0xf0, 0xf6, 0x5d, 0x3a, 0x09, 0x03, 0x66, 0xd7, 0x25,
	0x98, 0x3b, 0x32, 0xf4, 0x02, 0x59, 0x82, 0xc2, 0xde, 0x34, 0xa8, 0x72, 0x90, 0x43, 0xb0, 0x7c,
	0xe4, 0x2c, 0x09, 0x66, 0x22, 0x88, 0x23, 0xbb, 0x21, 0xf1, 0xa2, 0x8b, 0x8c, 0xa0, 0x9d, 0x11,
	0xe7, 0x76, 0xf3, 0xb0, 0xde, 0xb7, 0x46, 0xbd, 0xc1, 0x84, 0x86, 0x83, 0x73, 0xbc, 0x39, 0x51,
	0xfe, 0x8c, 0xbb, 0x67, 0xe2, 0xdc, 0x6f, 0x35, 0xe8, 0x5e, 0xce, 0x7c, 0x2a, 0x70, 0x41, 0xde,
	0x01, 0x74, 0xb2, 0x1a, 0x19, 0x85, 0xb9, 0xa3, 0x20, 0x7e, 0xa3, 0x24, 0xfe, 0xdf, 0x13, 0xf8,
	0x1a, 0x76, 0xfe, 0x86, 0x32, 0xf7, 0x08, 0xfe, 0xcf, 0xf2, 0xf0, 0x3b, 0x26, 0xc0, 0xfd, 0xb2,
	0x01, 0xad, 0x2c, 0x76, 0xcd, 0x32, 0x16, 0x67, 0xa8, 0xbe, 0x6a, 0x86, 0x36, 0x57, 0x96, 0xb8,
	0xb1, 0x58, 0xe2, 0x2e, 0x34, 0x44, 0x20, 0x42, 0xb4, 0x9b, 0x12, 0x51, 0xc6, 0x62, 0xe1, 0x5b,
	0xcb, 0x85, 0xef, 0x41, 0x93, 0x0b, 0x2a, 0x52, 0x6e, 0xb7, 0x15, 0x4f, 0x65, 0xcd, 0x4f, 0x63,
	0x09, 0x52, 0x81, 0xfe, 0xdb, 0xc8, 0xee, 0xa8, 0xd3, 0x8c, 0x63, 0x8e, 0xa6, 0x72, 0xb4, 0xe6,
	0x28, 0x28, 0xd4, 0x38, 0x48, 0xbf, 0xd0, 0x4c, 0x4b, 0x36, 0x73, 0x4b, 0x36, 0x53, 0x77, 0x32,
	0x6f, 0xe1, 0x73, 0xb0, 0x2e, 0x39, 0x26, 0xba, 0xa4, 0x8f, 0x41, 0x6f, 0xb7, 0x2c, 0xa8, 0x35,
	0xda, 0x1d, 0xe8, 0xe5, 0xd7, 0x9d, 0xd6, 0x01, 0xee, 0xd4, 0x74, 0xcd, 0x43, 0x3e, 0x8b, 0x23,
	0x8e, 0x05, 0x2d, 0xb5, 0x92, 0x16, 0x1b, 0x5a, 0x53, 0xe4, 0x9c, 0x9a, 0xb5, 0xd5, 0x26, 0xe9,
	0xc3, 0xa6, 0x4f, 0x05, 0x95, 0x9d, 0xb0, 0x46, 0x5d, 0x73, 0x5a, 0x81, 0x94, 0x27, 0x23, 0xdc,
	0x17, 0xb0, 0x55, 0x70, 0x72, 0xf2, 0x04, 0xda, 0x59, 0xf0, 0xfc, 0xb4, 0x7a, 0x25, 0x57, 0x13,
	0xe1, 0xc6, 0xb0, 0x9b, 0x8f, 0xd8, 0xda, 0x6c, 0x8f, 0x4a, 0x6c, 0xef, 0x57, 0xb1, 0xe5, 0x8a,
	0xee, 0xe8, 0x7b, 0xc3, 0x2c, 0xc7, 0x85, 0xba, 0x07, 0xc9, 0x18, 0xe0, 0xd8, 0xf7, 0x75, 0xa9,
	0x1d, 0xf3, 0xf7, 0xd2, 0xfd, 0xe7, 0xd8, 0x4b, 0x4a, 0x32, 0xce, 0xee, 0x7f, 0xe4, 0x15, 0xec,
	0xe5, 0x59, 0xb2, 0x86, 0x92, 0x9f, 0xec, 0xaa, 0xd3, 0x2d, 0xb5, 0x3d, 0x4f, 0x32, 0x86, 0xed,
	0x31, 0x86, 0x68, 0x6e, 0x26, 0xf2, 0x60, 0xf9, 0xc4, 0xbb, 0xa9, 0x8c, 0xc1, 0x3a, 0x45, 0x61,
	0x3a, 0xb2, 0x14, 0xaa, 0xb7, 0xd9, 0xd9, 0xaf, 0x40, 0x4c, 0x96, 0x63, 0x80, 0x3c, 0xcb, 0x7a,
	0x44, 0x5e, 0xc2, 0x5e, 0x9e, 0x42, 0xd7, 0xe4, 0x5e, 0x59, 0xfb, 0xea, 0x82, 0x9c, 0xc2, 0x8e,
	0x87, 0xfc, 0x36, 0x62, 0x7f, 0xaa, 0xe6, 0x0d, 0x6c, 0x97, 0xee, 0x7c, 0xf2, 0x28, 0x9f, 0x92,
	0x8a, 0xb7, 0x60, 0xa5, 0xac, 0xb3, 0x85, 0xf7, 0x43, 0x2b, 0xdb, 0x97, 0x22, 0x14, 0xf4, 0x8b,
	0xfa, 0x2e, 0xa1, 0xf7, 0x9e, 0x86, 0x41, 0x45, 0xb2, 0x87, 0xf2, 0x0f, 0x0d, 0x2e, 0xa4, 0x3b,
	0xa8, 0x06, 0x75, 0xda, 0x49, 0x53, 0xbe, 0xe2, 0xcf, 0x7e, 0x04, 0x00, 0x00, 0xff, 0xff, 0xc8,
	0x1e, 0x6d, 0xa6, 0x58, 0x08, 0x00, 0x00,
}
