package test

import (
	"context"

	protoBalance "github.com/asciiu/gomo/binance-service/proto/balance"
	protoBinance "github.com/asciiu/gomo/binance-service/proto/binance"
	"github.com/micro/go-micro/client"
)

// Test clients of the Key service should use this client interface.
type mockBinanceService struct {
	count uint32
}

func (m *mockBinanceService) GetBalances(ctx context.Context, in *protoBalance.BalanceRequest, opts ...client.CallOption) (*protoBalance.BalancesResponse, error) {
	call1 := []*protoBalance.Balance{
		&protoBalance.Balance{
			CurrencySymbol: "BTC",
			Free:           1.0,
			Locked:         1.0,
		},
		&protoBalance.Balance{
			CurrencySymbol: "USDT",
			Free:           100.00,
			Locked:         20.00,
		},
	}

	call2 := []*protoBalance.Balance{
		&protoBalance.Balance{
			CurrencySymbol: "BTC",
			Free:           0.5,
			Locked:         1.5,
		},
		&protoBalance.Balance{
			CurrencySymbol: "USDT",
			Free:           0.00,
			Locked:         0.00,
		},
	}

	calls := [][]*protoBalance.Balance{call1, call2}
	balances := calls[m.count%2]
	m.count++

	return &protoBalance.BalancesResponse{
		Status: "success",
		Data: &protoBalance.BalanceList{
			Balances: balances,
		},
	}, nil
}

func (m *mockBinanceService) GetMarketRestrictions(ctx context.Context, in *protoBinance.MarketRestrictionRequest, opts ...client.CallOption) (*protoBinance.MarketRestrictionResponse, error) {
	return &protoBinance.MarketRestrictionResponse{
		Status: "success",
		Data: &protoBinance.RestrictionData{
			Restrictions: &protoBinance.MarketRestriction{
				MinTradeSize:    1,
				MaxTradeSize:    10,
				TradeSizeStep:   0.1,
				MinMarketPrice:  1,
				MaxMarketPrice:  100,
				MarketPriceStep: 1,
			},
		},
	}, nil
}

func (m *mockBinanceService) GetCandles(ctx context.Context, in *protoBinance.MarketRequest, opts ...client.CallOption) (*protoBinance.CandlesResponse, error) {
	return &protoBinance.CandlesResponse{
		Status: "success",
	}, nil
}

func MockBinanceServiceClient() protoBinance.BinanceServiceClient {
	return new(mockBinanceService)
}
