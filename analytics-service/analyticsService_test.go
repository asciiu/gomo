package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	repoAnalytics "github.com/asciiu/gomo/analytics-service/db/sql"
	protoAnalytics "github.com/asciiu/gomo/analytics-service/proto/analytics"
	"github.com/asciiu/gomo/common/db"
	commonUtil "github.com/asciiu/gomo/common/util"
	repoUser "github.com/asciiu/gomo/user-service/db/sql"
	user "github.com/asciiu/gomo/user-service/models"
	"github.com/stretchr/testify/assert"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func setupService() (*AnalyticsService, *user.User) {
	dbUrl := "postgres://postgres@localhost:5432/gomo_test?&sslmode=disable"
	db, _ := db.NewDB(dbUrl)

	analyticsService := AnalyticsService{
		DB: db,
	}

	user := user.NewUser("first", "last", "test@email", "hash")
	repoUser.InsertUser(db, user)

	return &analyticsService, user
}

func TestDirectConversion(t *testing.T) {
	service, user := setupService()

	defer service.DB.Close()

	rates := []*protoAnalytics.MarketPrice{
		&protoAnalytics.MarketPrice{
			Exchange:      "test",
			MarketName:    "BTC-USDT",
			ClosedAtTime:  "2018-08-27 06:00:00",
			ClosedAtPrice: 6716.08,
		},
		&protoAnalytics.MarketPrice{
			Exchange:      "test",
			MarketName:    "ADA-USDT",
			ClosedAtTime:  "2018-08-27 06:00:00",
			ClosedAtPrice: 0.09585,
		},
		&protoAnalytics.MarketPrice{
			Exchange:      "test",
			MarketName:    "ADA-BTC",
			ClosedAtTime:  "2018-08-27 06:00:00",
			ClosedAtPrice: 0.00001425,
		},
	}
	repoAnalytics.InsertPrices(service.DB, rates)

	req := protoAnalytics.ConversionRequest{
		Exchange:    "test",
		From:        "ADA",
		FromAmount:  100,
		To:          "USDT",
		AtTimestamp: "2018-08-27T06:02:35.168652Z",
	}
	res := protoAnalytics.ConversionResponse{}
	service.ConvertCurrency(context.Background(), &req, &res)

	assert.Equal(t, "success", res.Status, fmt.Sprintf("return status of inserting plan should be success got: %s", res.Message))
	assert.Equal(t, 0.09585*100, res.Data.ConvertedAmount, "converted amount is not correct")

	repoAnalytics.DeleteExchangeRates(service.DB)
	repoUser.DeleteUserHard(service.DB, user.ID)
}

func TestIndirectConversion(t *testing.T) {
	service, user := setupService()

	defer service.DB.Close()

	rates := []*protoAnalytics.MarketPrice{
		&protoAnalytics.MarketPrice{
			Exchange:      "test",
			MarketName:    "ADA-BNB",
			ClosedAtTime:  "2018-08-27 06:00:00",
			ClosedAtPrice: 0.00911,
		},
		&protoAnalytics.MarketPrice{
			Exchange:      "test",
			MarketName:    "BTC-USDT",
			ClosedAtTime:  "2018-08-27 06:00:00",
			ClosedAtPrice: 6716.08,
		},
		&protoAnalytics.MarketPrice{
			Exchange:      "test",
			MarketName:    "ADA-BTC",
			ClosedAtTime:  "2018-08-27 06:00:00",
			ClosedAtPrice: 0.00001425,
		},
	}
	err := repoAnalytics.InsertPrices(service.DB, rates)
	assert.Equal(t, nil, err, "nope")

	req := protoAnalytics.ConversionRequest{
		Exchange:    "test",
		From:        "ADA",
		FromAmount:  100,
		To:          "USDT",
		AtTimestamp: "2018-08-27T06:02:35.168652Z",
	}
	res := protoAnalytics.ConversionResponse{}
	service.ConvertCurrency(context.Background(), &req, &res)

	assert.Equal(t, "success", res.Status, fmt.Sprintf("return status of inserting plan should be success got: %s", res.Message))
	assert.Equal(t, 9.5704, commonUtil.ToFixed(res.Data.ConvertedAmount, 4), "converted amount is not correct")

	repoAnalytics.DeleteExchangeRates(service.DB)
	repoUser.DeleteUserHard(service.DB, user.ID)
}
