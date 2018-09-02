package test

import (
	"context"

	protoBalance "github.com/asciiu/gomo/binance-service/proto/balance"
	protoBinance "github.com/asciiu/gomo/binance-service/proto/binance"
	"github.com/micro/go-micro/client"
)

// Test clients of the Key service should use this client interface.
type mockBinanceService struct{}

func (m *mockBinanceService) GetBalances(ctx context.Context, in *protoBalance.BalanceRequest, opts ...client.CallOption) (*protoBalance.BalancesResponse, error) {
	return &protoBalance.BalancesResponse{
		Status: "success",
		Data: &protoBalance.BalanceList{
			Balances: []*protoBalance.Balance{},
		},
	}, nil
}

func MockBinanceServiceClient() protoBinance.BinanceServiceClient {
	return new(mockBinanceService)
}
