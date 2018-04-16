package main

import (
	"context"
	"database/sql"
	"fmt"

	bp "github.com/asciiu/gomo/balance-service/proto/balance"
	micro "github.com/micro/go-micro"
)

type BalancerUpdater struct {
	DB    *sql.DB
	Micro micro.Service
}

func (service *BalancerUpdater) Process(ctx context.Context, balances *bp.AccountBalances) error {
	fmt.Println(balances)
	return nil
}
