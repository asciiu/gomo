package test

import (
	"context"
	"database/sql"

	protoAnalytics "github.com/asciiu/gomo/analytics-service/proto/analytics"
	"github.com/micro/go-micro/client"
)

// Test clients of the Key service should use this client interface.
type mockAnalyticsService struct {
	db *sql.DB
}

func (m *mockAnalyticsService) ConvertCurrency(ctx context.Context, req *protoAnalytics.ConversionRequest, opts ...client.CallOption) (*protoAnalytics.ConversionResponse, error) {
	value := 0.0
	switch {
	case req.From == "BTC" && req.To == "USDT":
		value = 100.0
	case req.From == "USDT" && req.To == "USDT":
		value = 100.0
	}

	return &protoAnalytics.ConversionResponse{
		Status: "success",
		Data: &protoAnalytics.ConversionAmount{
			ConvertedAmount: value,
		},
	}, nil
}

func (m *mockAnalyticsService) GetMarketInfo(ctx context.Context, req *protoAnalytics.SearchMarketsRequest, opts ...client.CallOption) (*protoAnalytics.MarketsResponse, error) {
	return &protoAnalytics.MarketsResponse{
		Status: "success",
	}, nil
}

func MockAnalyticsServiceClient(db *sql.DB) protoAnalytics.AnalyticsServiceClient {
	return &mockAnalyticsService{db}
}
