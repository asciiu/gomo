// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/asciiu/gomo/account-service/proto/balance/balance.proto

/*
Package balance is a generated protocol buffer package.

It is generated from these files:
	github.com/asciiu/gomo/account-service/proto/balance/balance.proto

It has these top-level messages:
	NewBalanceRequest
	BalanceRequest
	UpdateBalanceRequest
	ValidateBalanceRequest
	Balance
	BalanceData
	BalanceResponse
	ValidateBalanceResponse
*/
package balance

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

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
type NewBalanceRequest struct {
	AccountID      string  `protobuf:"bytes,1,opt,name=accountID" json:"accountID"`
	CurrencySymbol string  `protobuf:"bytes,2,opt,name=currencySymbol" json:"currencySymbol"`
	Available      float64 `protobuf:"fixed64,3,opt,name=available" json:"available"`
}

func (m *NewBalanceRequest) Reset()                    { *m = NewBalanceRequest{} }
func (m *NewBalanceRequest) String() string            { return proto.CompactTextString(m) }
func (*NewBalanceRequest) ProtoMessage()               {}
func (*NewBalanceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *NewBalanceRequest) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *NewBalanceRequest) GetCurrencySymbol() string {
	if m != nil {
		return m.CurrencySymbol
	}
	return ""
}

func (m *NewBalanceRequest) GetAvailable() float64 {
	if m != nil {
		return m.Available
	}
	return 0
}

type BalanceRequest struct {
	AccountID      string `protobuf:"bytes,1,opt,name=accountID" json:"accountID"`
	CurrencySymbol string `protobuf:"bytes,2,opt,name=currencySymbol" json:"currencySymbol"`
}

func (m *BalanceRequest) Reset()                    { *m = BalanceRequest{} }
func (m *BalanceRequest) String() string            { return proto.CompactTextString(m) }
func (*BalanceRequest) ProtoMessage()               {}
func (*BalanceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *BalanceRequest) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *BalanceRequest) GetCurrencySymbol() string {
	if m != nil {
		return m.CurrencySymbol
	}
	return ""
}

type UpdateBalanceRequest struct {
	AccountID         string  `protobuf:"bytes,1,opt,name=accountID" json:"accountID"`
	CurrencySymbol    string  `protobuf:"bytes,2,opt,name=currencySymbol" json:"currencySymbol"`
	Available         float64 `protobuf:"fixed64,3,opt,name=available" json:"available"`
	Locked            float64 `protobuf:"fixed64,4,opt,name=locked" json:"locked"`
	ExchangeTotal     float64 `protobuf:"fixed64,5,opt,name=exchangeTotal" json:"exchangeTotal"`
	ExchangeAvailable float64 `protobuf:"fixed64,6,opt,name=exchangeAvailable" json:"exchangeAvailable"`
	ExchangeLocked    float64 `protobuf:"fixed64,7,opt,name=exchangeLocked" json:"exchangeLocked"`
}

func (m *UpdateBalanceRequest) Reset()                    { *m = UpdateBalanceRequest{} }
func (m *UpdateBalanceRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdateBalanceRequest) ProtoMessage()               {}
func (*UpdateBalanceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *UpdateBalanceRequest) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *UpdateBalanceRequest) GetCurrencySymbol() string {
	if m != nil {
		return m.CurrencySymbol
	}
	return ""
}

func (m *UpdateBalanceRequest) GetAvailable() float64 {
	if m != nil {
		return m.Available
	}
	return 0
}

func (m *UpdateBalanceRequest) GetLocked() float64 {
	if m != nil {
		return m.Locked
	}
	return 0
}

func (m *UpdateBalanceRequest) GetExchangeTotal() float64 {
	if m != nil {
		return m.ExchangeTotal
	}
	return 0
}

func (m *UpdateBalanceRequest) GetExchangeAvailable() float64 {
	if m != nil {
		return m.ExchangeAvailable
	}
	return 0
}

func (m *UpdateBalanceRequest) GetExchangeLocked() float64 {
	if m != nil {
		return m.ExchangeLocked
	}
	return 0
}

type ValidateBalanceRequest struct {
	AccountID      string  `protobuf:"bytes,1,opt,name=accountID" json:"accountID"`
	CurrencySymbol string  `protobuf:"bytes,2,opt,name=currencySymbol" json:"currencySymbol"`
	Available      float64 `protobuf:"fixed64,3,opt,name=available" json:"available"`
}

func (m *ValidateBalanceRequest) Reset()                    { *m = ValidateBalanceRequest{} }
func (m *ValidateBalanceRequest) String() string            { return proto.CompactTextString(m) }
func (*ValidateBalanceRequest) ProtoMessage()               {}
func (*ValidateBalanceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ValidateBalanceRequest) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *ValidateBalanceRequest) GetCurrencySymbol() string {
	if m != nil {
		return m.CurrencySymbol
	}
	return ""
}

func (m *ValidateBalanceRequest) GetAvailable() float64 {
	if m != nil {
		return m.Available
	}
	return 0
}

// Responses
type Balance struct {
	UserID            string  `protobuf:"bytes,1,opt,name=userID" json:"userID"`
	AccountID         string  `protobuf:"bytes,2,opt,name=accountID" json:"accountID"`
	CurrencySymbol    string  `protobuf:"bytes,3,opt,name=currencySymbol" json:"currencySymbol"`
	Available         float64 `protobuf:"fixed64,4,opt,name=available" json:"available"`
	Locked            float64 `protobuf:"fixed64,5,opt,name=locked" json:"locked"`
	ExchangeTotal     float64 `protobuf:"fixed64,6,opt,name=exchangeTotal" json:"exchangeTotal"`
	ExchangeAvailable float64 `protobuf:"fixed64,7,opt,name=exchangeAvailable" json:"exchangeAvailable"`
	ExchangeLocked    float64 `protobuf:"fixed64,8,opt,name=exchangeLocked" json:"exchangeLocked"`
	CreatedOn         string  `protobuf:"bytes,9,opt,name=created_on,json=createdOn" json:"created_on"`
	UpdatedOn         string  `protobuf:"bytes,10,opt,name=updated_on,json=updatedOn" json:"updated_on"`
}

func (m *Balance) Reset()                    { *m = Balance{} }
func (m *Balance) String() string            { return proto.CompactTextString(m) }
func (*Balance) ProtoMessage()               {}
func (*Balance) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *Balance) GetUserID() string {
	if m != nil {
		return m.UserID
	}
	return ""
}

func (m *Balance) GetAccountID() string {
	if m != nil {
		return m.AccountID
	}
	return ""
}

func (m *Balance) GetCurrencySymbol() string {
	if m != nil {
		return m.CurrencySymbol
	}
	return ""
}

func (m *Balance) GetAvailable() float64 {
	if m != nil {
		return m.Available
	}
	return 0
}

func (m *Balance) GetLocked() float64 {
	if m != nil {
		return m.Locked
	}
	return 0
}

func (m *Balance) GetExchangeTotal() float64 {
	if m != nil {
		return m.ExchangeTotal
	}
	return 0
}

func (m *Balance) GetExchangeAvailable() float64 {
	if m != nil {
		return m.ExchangeAvailable
	}
	return 0
}

func (m *Balance) GetExchangeLocked() float64 {
	if m != nil {
		return m.ExchangeLocked
	}
	return 0
}

func (m *Balance) GetCreatedOn() string {
	if m != nil {
		return m.CreatedOn
	}
	return ""
}

func (m *Balance) GetUpdatedOn() string {
	if m != nil {
		return m.UpdatedOn
	}
	return ""
}

type BalanceData struct {
	Balance *Balance `protobuf:"bytes,1,opt,name=balance" json:"balance"`
}

func (m *BalanceData) Reset()                    { *m = BalanceData{} }
func (m *BalanceData) String() string            { return proto.CompactTextString(m) }
func (*BalanceData) ProtoMessage()               {}
func (*BalanceData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *BalanceData) GetBalance() *Balance {
	if m != nil {
		return m.Balance
	}
	return nil
}

type BalanceResponse struct {
	Status  string       `protobuf:"bytes,1,opt,name=status" json:"status"`
	Message string       `protobuf:"bytes,2,opt,name=message" json:"message"`
	Data    *BalanceData `protobuf:"bytes,3,opt,name=data" json:"data"`
}

func (m *BalanceResponse) Reset()                    { *m = BalanceResponse{} }
func (m *BalanceResponse) String() string            { return proto.CompactTextString(m) }
func (*BalanceResponse) ProtoMessage()               {}
func (*BalanceResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *BalanceResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *BalanceResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *BalanceResponse) GetData() *BalanceData {
	if m != nil {
		return m.Data
	}
	return nil
}

type ValidateBalanceResponse struct {
	Status  string `protobuf:"bytes,1,opt,name=status" json:"status"`
	Message string `protobuf:"bytes,2,opt,name=message" json:"message"`
}

func (m *ValidateBalanceResponse) Reset()                    { *m = ValidateBalanceResponse{} }
func (m *ValidateBalanceResponse) String() string            { return proto.CompactTextString(m) }
func (*ValidateBalanceResponse) ProtoMessage()               {}
func (*ValidateBalanceResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *ValidateBalanceResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *ValidateBalanceResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*NewBalanceRequest)(nil), "balance.NewBalanceRequest")
	proto.RegisterType((*BalanceRequest)(nil), "balance.BalanceRequest")
	proto.RegisterType((*UpdateBalanceRequest)(nil), "balance.UpdateBalanceRequest")
	proto.RegisterType((*ValidateBalanceRequest)(nil), "balance.ValidateBalanceRequest")
	proto.RegisterType((*Balance)(nil), "balance.Balance")
	proto.RegisterType((*BalanceData)(nil), "balance.BalanceData")
	proto.RegisterType((*BalanceResponse)(nil), "balance.BalanceResponse")
	proto.RegisterType((*ValidateBalanceResponse)(nil), "balance.ValidateBalanceResponse")
}

func init() {
	proto.RegisterFile("github.com/asciiu/gomo/account-service/proto/balance/balance.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 423 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x54, 0xd1, 0x6a, 0xd4, 0x40,
	0x14, 0x25, 0xe9, 0x36, 0x71, 0xef, 0x62, 0xb5, 0x43, 0xa9, 0x79, 0x50, 0x28, 0x41, 0x64, 0x11,
	0xdd, 0x40, 0x7d, 0xf2, 0xd1, 0xd2, 0x17, 0x51, 0x2c, 0x44, 0xed, 0xab, 0xdc, 0x4c, 0x2e, 0x69,
	0x30, 0x99, 0x59, 0x33, 0x33, 0xad, 0x05, 0xbf, 0xc1, 0xdf, 0xf3, 0x77, 0x24, 0x93, 0x99, 0x94,
	0x6e, 0xd5, 0x5d, 0x04, 0xed, 0x53, 0x38, 0xe7, 0x5c, 0xe6, 0xdc, 0xb9, 0x67, 0x6e, 0xe0, 0xa8,
	0xaa, 0xf5, 0x99, 0x29, 0x16, 0x5c, 0xb6, 0x19, 0x2a, 0x5e, 0xd7, 0x26, 0xab, 0x64, 0x2b, 0x33,
	0xe4, 0x5c, 0x1a, 0xa1, 0x9f, 0x2b, 0xea, 0xce, 0x6b, 0x4e, 0xd9, 0xb2, 0x93, 0x5a, 0x66, 0x05,
	0x36, 0x28, 0x38, 0xf9, 0xef, 0xc2, 0xb2, 0x2c, 0x76, 0x30, 0xbd, 0x80, 0xdd, 0x77, 0x74, 0x71,
	0x34, 0xa0, 0x9c, 0xbe, 0x18, 0x52, 0x9a, 0x3d, 0x84, 0xa9, 0x3b, 0xec, 0xf5, 0x71, 0x12, 0x1c,
	0x04, 0xf3, 0x69, 0x7e, 0x45, 0xb0, 0x27, 0xb0, 0xc3, 0x4d, 0xd7, 0x91, 0xe0, 0x97, 0xef, 0x2f,
	0xdb, 0x42, 0x36, 0x49, 0x68, 0x4b, 0x56, 0x58, 0x7b, 0xca, 0x39, 0xd6, 0x0d, 0x16, 0x0d, 0x25,
	0x5b, 0x07, 0xc1, 0x3c, 0xc8, 0xaf, 0x88, 0xf4, 0x14, 0x76, 0xfe, 0x85, 0x6b, 0xfa, 0x3d, 0x84,
	0xbd, 0x8f, 0xcb, 0x12, 0x35, 0xfd, 0xff, 0x4b, 0xb1, 0x7d, 0x88, 0x1a, 0xc9, 0x3f, 0x53, 0x99,
	0x4c, 0xac, 0xe4, 0x10, 0x7b, 0x0c, 0x77, 0xe9, 0x2b, 0x3f, 0x43, 0x51, 0xd1, 0x07, 0xa9, 0xb1,
	0x49, 0xb6, 0xad, 0x7c, 0x9d, 0x64, 0xcf, 0x60, 0xd7, 0x13, 0xaf, 0x46, 0x8f, 0xc8, 0x56, 0xde,
	0x14, 0xfa, 0x8e, 0x3d, 0xf9, 0x76, 0xf0, 0x8c, 0x6d, 0xe9, 0x0a, 0x9b, 0x7e, 0x83, 0xfd, 0x53,
	0x6c, 0xea, 0xdb, 0x99, 0x48, 0xfa, 0x23, 0x84, 0xd8, 0xd9, 0xf6, 0xd3, 0x31, 0x8a, 0xba, 0xd1,
	0xcc, 0xa1, 0xeb, 0x7d, 0x84, 0xeb, 0xfb, 0xd8, 0x5a, 0xdf, 0xc7, 0xe4, 0xf7, 0xc9, 0x6c, 0xff,
	0x39, 0x99, 0x68, 0xe3, 0x64, 0xe2, 0xcd, 0x93, 0xb9, 0xf3, 0xab, 0x64, 0xd8, 0x23, 0x00, 0xde,
	0x11, 0x6a, 0x2a, 0x3f, 0x49, 0x91, 0x4c, 0x87, 0x8b, 0x3b, 0xe6, 0x44, 0xf4, 0xb2, 0xb1, 0x0f,
	0xd9, 0xca, 0x30, 0xc8, 0x8e, 0x39, 0x11, 0xe9, 0x4b, 0x98, 0xb9, 0xc1, 0x1e, 0xa3, 0x46, 0xf6,
	0x14, 0xfc, 0x4e, 0xdb, 0xe9, 0xce, 0x0e, 0xef, 0x2f, 0xfc, 0xca, 0xfb, 0xd8, 0xc7, 0xa5, 0x6f,
	0xe1, 0xde, 0xf8, 0x14, 0xd4, 0x52, 0x0a, 0x65, 0xe7, 0xa3, 0x34, 0x6a, 0xa3, 0x7c, 0x36, 0x03,
	0x62, 0x09, 0xc4, 0x2d, 0x29, 0x85, 0x15, 0xb9, 0x64, 0x3c, 0x64, 0x73, 0x98, 0x94, 0xa8, 0xd1,
	0xa6, 0x31, 0x3b, 0xdc, 0x5b, 0x75, 0xeb, 0x9b, 0xca, 0x6d, 0x45, 0xfa, 0x06, 0x1e, 0xdc, 0x78,
	0x81, 0x7f, 0x6b, 0x5b, 0x44, 0xf6, 0x07, 0xf6, 0xe2, 0x67, 0x00, 0x00, 0x00, 0xff, 0xff, 0x55,
	0x4b, 0xe9, 0x08, 0x06, 0x05, 0x00, 0x00,
}
