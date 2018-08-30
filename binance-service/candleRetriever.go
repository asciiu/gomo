package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	binance "github.com/asciiu/go-binance"
	"github.com/asciiu/gomo/common/constants/exchange"

	evt "github.com/asciiu/gomo/common/proto/events"
	"github.com/go-kit/kit/log"
)

type CandleRetriever struct {
}

func (service *CandleRetriever) ProcessGetCandle(ctx context.Context, request *evt.CandleDataRequest) error {
	if request.Exchange != exchange.Binance {
		return nil
	}

	go func() {
		symbol := strings.Replace(request.MarketName, "-", "", 0)

		//url := fmt.Sprintf("https://api.binance.com/api/v1/klines?symbol=%s&interval=%s", symbol, request.Interval)

		//	req, _ := http.NewRequest("GET", url, nil)

		//	res, _ := http.DefaultClient.Do(req)

		//	defer res.Body.Close()

		//	body, _ := ioutil.ReadAll(res.Body)

		//	fmt.Printf("%s: %s\n", request.Exchange, request.MarketName)
		//	fmt.Println(string(body))
		var logger log.Logger
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

		hmacSigner := &binance.HmacSigner{
			Key: []byte(""),
		}

		binanceService := binance.NewAPIService(
			"https://www.binance.com",
			"",
			hmacSigner,
			logger,
			ctx,
		)
		b := binance.NewBinance(binanceService)
		kl, err := b.Klines(binance.KlinesRequest{
			Symbol:   symbol,
			Interval: binance.Minute,
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", kl)

		//accountBalances := bp.AccountBalances{
		//	Balances: balances,
		//}

		// publish account balances
		//if err := service.BalancePub.Publish(context.Background(), &accountBalances); err != nil {
		//	logger.Log("could not publish account balances event: ", err)
		//}
	}()
	return nil
}
