package main

import (
	"context"
	"database/sql"
	"fmt"

	balProto "github.com/asciiu/gomo/balance-service/proto/balance"
	micro "github.com/micro/go-micro"
)

type BalancerUpdater struct {
	DB    *sql.DB
	Micro micro.Service
}

func (service *BalancerUpdater) Process(ctx context.Context, balances *balProto.AccountBalances) error {
	fmt.Println(balances)
	return nil
}
