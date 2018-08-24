package test

import (
	"context"
	"database/sql"

	protoAccount "github.com/asciiu/gomo/account-service/proto/account"
	"github.com/micro/go-micro/client"
)

// Test clients of the Key service should use this client interface.
type mockAccountService struct {
	db *sql.DB
}

func (m *mockAccountService) GetKeys(ctx context.Context, req *protoKey.GetKeysRequest, opts ...client.CallOption) (*protoKey.KeyListResponse, error) {
	//keys, _ := repoKey.FindKeys(m.db, req)

	//return &protoKey.KeyListResponse{
	//	Status: "success",
	//	Data: &protoKey.UserKeysData{
	//		Keys: keys,
	//	},
	//}, nil

}

func (m *mockAccountService) AddKey(ctx context.Context, in *protoKey.KeyRequest, opts ...client.CallOption) (*protoKey.KeyResponse, error) {
	return &protoAccount.AccountResponse{}, nil
}
func (m *mockAccountService) GetUserKey(ctx context.Context, in *protoKey.GetUserKeyRequest, opts ...client.CallOption) (*protoKey.KeyResponse, error) {
	return &protoKey.AccountResponse{}, nil
}
func (m *mockAccountService) GetUserKeys(ctx context.Context, in *protoKey.GetUserKeysRequest, opts ...client.CallOption) (*protoKey.KeyListResponse, error) {
	return &protoKey.KeyListResponse{}, nil
}
func (m *mockAccountService) RemoveKey(ctx context.Context, in *protoKey.RemoveKeyRequest, opts ...client.CallOption) (*protoKey.KeyResponse, error) {
	return &protoKey.KeyResponse{}, nil
}
func (m *mockAccountService) UpdateKeyDescription(ctx context.Context, in *protoKey.KeyRequest, opts ...client.CallOption) (*protoKey.KeyResponse, error) {
	return &protoKey.KeyResponse{}, nil
}

func MockAccountServiceClient(db *sql.DB) protoKey.KeyServiceClient {
	return &mockAccountService{db}
}
