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
	"github.com/asciiu/gomo/common/db"
	"github.com/asciiu/gomo/common/util"
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
	service := new(BinanceService)

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := repoUser.InsertUser(db, user)
	checkErr(error)

	return service, db, user
}

func TestExchangeInfo(t *testing.T) {
	bexinfo := NewBinanceExchangeInfo()
	marketName := "ADA-BTC"
	lotSize := bexinfo.LotSize(marketName)
	priceFilter := bexinfo.PriceFilter(marketName)
	minNotional := bexinfo.MinNotional(marketName)
	icebergParts := bexinfo.IcebergParts(marketName)
	maxAlg := bexinfo.MaxAlgoOrders(marketName)

	fmt.Println(priceFilter)
	fmt.Println(lotSize)
	fmt.Println(minNotional)
	fmt.Println(icebergParts)
	fmt.Println(maxAlg)
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
	//orders, _ := b.AllOrders(binance.AllOrdersRequest{
	//	Symbol:     marketName,
	//	OrderID:    59016988,
	//	Limit:      1,
	//	RecvWindow: time.Duration(2) * time.Second,
	//	Timestamp:  time.Now(),
	//})

	trades, _ := b.MyTrades(binance.MyTradesRequest{
		Symbol: marketName,
		//FromID:     18887917,
		Limit:      200,
		RecvWindow: time.Duration(2) * time.Second,
		Timestamp:  time.Now(),
	})

	//if err != nil {
	//	fmt.Println(err)
	//}

	//executedOrder, _ := b.QueryOrder(binance.QueryOrderRequest{
	//	Symbol:            marketName,
	//	OrderID:           59016988,
	//	OrigClientOrderID: "1d2b57a8-a695-11e8-98d0-529269fb1459",
	//	RecvWindow:        time.Duration(2) * time.Second,
	//	Timestamp:         time.Now(),
	//})

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
