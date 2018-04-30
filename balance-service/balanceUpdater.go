package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	balRepo "github.com/asciiu/gomo/balance-service/db/sql"
	bp "github.com/asciiu/gomo/balance-service/proto/balance"
)

type BalanceUpdateListener struct {
	DB      *sql.DB
	Service *BalanceService
}

func (service *BalanceUpdateListener) Process(ctx context.Context, balances *bp.AccountBalances) error {

	count, err := balRepo.UpsertBalances(service.DB, balances)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("inserted ", count)
	return nil
}
