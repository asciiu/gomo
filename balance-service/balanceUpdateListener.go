package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	repo "github.com/asciiu/gomo/balance-service/db/sql"
	balances "github.com/asciiu/gomo/balance-service/proto/balance"
)

type BalanceUpdateListener struct {
	DB      *sql.DB
	Service *BalanceService
}

func (service *BalanceUpdateListener) Process(ctx context.Context, balances *balances.AccountBalances) error {

	count, err := repo.UpsertBalances(service.DB, balances)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("inserted ", count)
	return nil
}
