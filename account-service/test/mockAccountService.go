package test

import (
	"context"
	"database/sql"

	repoAccount "github.com/asciiu/gomo/account-service/db/sql"
	protoAccount "github.com/asciiu/gomo/account-service/proto/account"
	protoBalance "github.com/asciiu/gomo/account-service/proto/balance"
	"github.com/micro/go-micro/client"
)

// Test clients of the Key service should use this client interface.
type mockAccountService struct {
	db *sql.DB
}

func (m *mockAccountService) AddAccount(ctx context.Context, req *protoAccount.NewAccountRequest, opts ...client.CallOption) (*protoAccount.AccountResponse, error) {
	return &protoAccount.AccountResponse{
		Status: "success",
	}, nil
}

func (m *mockAccountService) AddBalance(ctx context.Context, req *protoBalance.NewBalanceRequest, opts ...client.CallOption) (*protoBalance.BalanceResponse, error) {
	return &protoBalance.BalanceResponse{
		Status: "success",
	}, nil
}

func (m *mockAccountService) AdjustAvailableBalance(ctx context.Context, req *protoBalance.AdjustBalanceRequest, opts ...client.CallOption) (*protoBalance.BalanceResponse, error) {
	return &protoBalance.BalanceResponse{
		Status: "success",
	}, nil
}

func (m *mockAccountService) AdjustLockedBalance(ctx context.Context, req *protoBalance.AdjustBalanceRequest, opts ...client.CallOption) (*protoBalance.BalanceResponse, error) {
	return &protoBalance.BalanceResponse{
		Status: "success",
	}, nil
}

func (m *mockAccountService) AdjustBalance(ctx context.Context, req *protoBalance.AdjustBalanceRequest, opts ...client.CallOption) (*protoBalance.BalanceResponse, error) {
	return &protoBalance.BalanceResponse{
		Status: "success",
	}, nil
}

func (m *mockAccountService) DeleteAccount(ctx context.Context, req *protoAccount.AccountRequest, opts ...client.CallOption) (*protoAccount.AccountResponse, error) {
	return &protoAccount.AccountResponse{
		Status: "success",
	}, nil
}

func (m *mockAccountService) GetAccounts(ctx context.Context, req *protoAccount.AccountsRequest, opts ...client.CallOption) (*protoAccount.AccountsResponse, error) {
	return &protoAccount.AccountsResponse{
		Status: "success",
	}, nil
}

func (m *mockAccountService) GetAccount(ctx context.Context, req *protoAccount.AccountRequest, opts ...client.CallOption) (*protoAccount.AccountResponse, error) {
	return &protoAccount.AccountResponse{
		Status: "success",
	}, nil
}

func (m *mockAccountService) GetAccountKeys(ctx context.Context, req *protoAccount.GetAccountKeysRequest, opts ...client.CallOption) (*protoAccount.AccountKeysResponse, error) {
	keys, _ := repoAccount.FindAccountKeys(m.db, req.AccountIDs)

	return &protoAccount.AccountKeysResponse{
		Status: "success",
		Data: &protoAccount.KeysList{
			Keys: keys,
		},
	}, nil
}

func (m *mockAccountService) GetAccountBalance(ctx context.Context, req *protoBalance.BalanceRequest, opts ...client.CallOption) (*protoBalance.BalanceResponse, error) {
	balance, _ := repoAccount.FindAccountBalance(m.db, req.UserID, req.AccountID, req.CurrencySymbol)

	return &protoBalance.BalanceResponse{
		Status: "success",
		Data: &protoBalance.BalanceData{
			Balance: balance,
		},
	}, nil
}

func (m *mockAccountService) ResyncAccounts(ctx context.Context, req *protoAccount.AccountsRequest, opts ...client.CallOption) (*protoAccount.AccountsResponse, error) {
	return &protoAccount.AccountsResponse{
		Status: "success",
	}, nil
}

func (m *mockAccountService) UpdateAccount(ctx context.Context, req *protoAccount.UpdateAccountRequest, opts ...client.CallOption) (*protoAccount.AccountResponse, error) {
	return &protoAccount.AccountResponse{
		Status: "success",
	}, nil
}

func (m *mockAccountService) ValidateAccountBalance(ctx context.Context, req *protoBalance.ValidateBalanceRequest, opts ...client.CallOption) (*protoBalance.ValidateBalanceResponse, error) {
	balance, _ := repoAccount.FindAccountBalance(m.db, req.UserID, req.AccountID, req.CurrencySymbol)

	valid := balance.Available >= req.RequestedAmount

	return &protoBalance.ValidateBalanceResponse{
		Status: "success",
		Data:   valid,
	}, nil
}

func MockAccountServiceClient(db *sql.DB) protoAccount.AccountServiceClient {
	return &mockAccountService{db}
}
