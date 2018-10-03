package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	binance "github.com/asciiu/go-binance"
	protoBalance "github.com/asciiu/gomo/binance-service/proto/balance"
	protoBinance "github.com/asciiu/gomo/binance-service/proto/binance"
	constExch "github.com/asciiu/gomo/common/constants/exchange"
	"github.com/asciiu/gomo/common/db"
	protoEvt "github.com/asciiu/gomo/common/proto/events"
	"github.com/asciiu/gomo/common/util"
	constPlan "github.com/asciiu/gomo/plan-service/constants"
	repoUser "github.com/asciiu/gomo/user-service/db/sql"
	user "github.com/asciiu/gomo/user-service/models"
	kitLog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func setupService() (*BinanceService, *sql.DB, *user.User) {
	dbUrl := "postgres://postgres@localhost:5432/gomo_test?&sslmode=disable"
	db, _ := db.NewDB(dbUrl)
	service := BinanceService{
		Info: NewBinanceExchangeInfo(),
	}

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := repoUser.InsertUser(db, user)
	checkErr(error)

	return &service, db, user
}

func TestExchangeInfo(t *testing.T) {
	bexinfo := NewBinanceExchangeInfo()
	marketName := "TRX-BTC"
	lotSize := bexinfo.LotSize(marketName)
	priceFilter := bexinfo.PriceFilter(marketName)
	minNotional := bexinfo.MinNotional(marketName)
	icebergParts := bexinfo.IcebergParts(marketName)
	maxAlg := bexinfo.MaxAlgoOrders(marketName)

	fmt.Println("price filter: ", priceFilter)
	fmt.Println("lot size: ", lotSize.StepSize)
	fmt.Println("min notional: ", minNotional)
	fmt.Println("iceberg parts: ", icebergParts)
	fmt.Println("max alg: ", maxAlg)
}

func TestOrderQuery(t *testing.T) {
	var logger kitLog.Logger
	logger = kitLog.NewLogfmtLogger(os.Stdout)
	logger = kitLog.With(logger, "time", kitLog.DefaultTimestampUTC, "caller", kitLog.DefaultCaller)

	hmacSigner := &binance.HmacSigner{
		Key: []byte(util.Rot32768("secret")),
	}

	binanceService := binance.NewAPIService(
		"https://www.binance.com",
		"public",
		hmacSigner,
		logger,
		context.Background(),
	)
	b := binance.NewBinance(binanceService)
	marketName := "ADABTC"

	trades, _ := b.MyTrades(binance.MyTradesRequest{
		Symbol:     marketName,
		Limit:      200,
		RecvWindow: time.Duration(2) * time.Second,
		Timestamp:  time.Now(),
	})

	for _, t := range trades {
		fmt.Printf("%+v\n", t)
	}
}

func TestInvalidKey1(t *testing.T) {
	service, db, user := setupService()

	defer db.Close()

	request := protoBalance.BalanceRequest{
		UserID:    user.ID,
		KeyPublic: "public",
		KeySecret: "secret",
	}

	response := protoBalance.BalancesResponse{}
	service.GetBalances(context.Background(), &request, &response)

	assert.Equal(t, "fail", response.Status, response.Message)

	repoUser.DeleteUserHard(db, user.ID)
}

func TestInvalidKey2(t *testing.T) {
	service, db, user := setupService()

	defer db.Close()

	request := protoBalance.BalanceRequest{
		UserID:    user.ID,
		KeyPublic: "Sn54bfgy5FILCvhtXSAlqqPhCgF74VLDlLYpJFNyVYeDMRiFCAo6g0F96CPb6xml",
		KeySecret: "AWxEQXvLQuyx218tZeeEHEWbfvdVXZ0zKjQgYEM3aDutkVIxQmtUeJWQVfHkPT1I",
	}

	response := protoBalance.BalancesResponse{}
	service.GetBalances(context.Background(), &request, &response)

	assert.Equal(t, "fail", response.Status, response.Message)

	repoUser.DeleteUserHard(db, user.ID)
}

func TestInvalidKey3(t *testing.T) {
	service, db, user := setupService()

	defer db.Close()

	request := protoBalance.BalanceRequest{
		UserID:    user.ID,
		KeyPublic: "O5oYc5b2TFSdcdWFqjQz8DnvVExeJUFeiGshmSVFet8WLHFVk3Iy1sQ5c",
		KeySecret: "cudE4yw1fQxrk5BfPXb4X5lKLZC3ypmIIWNOzUvE8e9p8sX40CnACo24nxID",
	}

	response := protoBalance.BalancesResponse{}
	service.GetBalances(context.Background(), &request, &response)

	assert.Equal(t, "fail", response.Status, response.Message)

	repoUser.DeleteUserHard(db, user.ID)
}

func TestGetCandle(t *testing.T) {
	service, db, user := setupService()

	defer db.Close()

	request := protoBinance.MarketRequest{
		MarketName: "ADA-BTC",
	}

	response := protoBinance.CandlesResponse{}
	service.GetCandles(context.Background(), &request, &response)

	assert.Equal(t, "success", response.Status, response.Message)
	assert.Equal(t, 500, len(response.Data.Candles), "should be 500 candle data")

	repoUser.DeleteUserHard(db, user.ID)
}

// Convert TriggeredOrderEvent -> binance.NewOrderRequest
func TestEventConversion(t *testing.T) {
	service, db, user := setupService()

	defer db.Close()

	event := protoEvt.TriggeredOrderEvent{
		Exchange:              constExch.Binance,
		OrderID:               "7010c9ef-33f1-f704-dcf3-abf086b75ffc",
		PlanID:                "f989125e-a40c-4cd1-8c94-9342f2dba92e",
		UserID:                "189170e9-8729-40ae-bcf2-afebd8aa69f1",
		AccountID:             "ae7dd5f8-9e97-4265-b873-2e23128ee176",
		ActiveCurrencySymbol:  "USDT",
		ActiveCurrencyBalance: 2000,
		Quantity:              0.3052004632943033,
		Side:                  constPlan.Buy,
		OrderType:             constPlan.MarketOrder,
		TriggerID:             "f276387d-6d14-4b22-bb11-72c04f95a17b",
		TriggeredPrice:        6553.07,
		TriggeredCondition:    "6553.07000000",
		MarketName:            "BTC-USDT",
	}

	c, b := service.FormatTriggerEvent(&event)

	assert.Equal(t, c.MarketName, event.MarketName, "market name no match")
	assert.Equal(t, c.InitialCurrencyBalance, event.ActiveCurrencyBalance, "initial balance no match")
	assert.Equal(t, c.FinalCurrencySymbol, "BTC", "final asset should be BTC")

	// 2000 / 6553.07 -> rounded to nearest lot size for BTC-USDT
	// at time of this test the lotsize for BTCUSDT was 0.000001
	// qty should be computed as 0.3052 but the alg results in 0.30519999
	// when multiplying by the stepsize due to rounding error
	assert.Equal(t, b.Quantity, 0.3052, "final quantity does not match")

	repoUser.DeleteUserHard(db, user.ID)
}
