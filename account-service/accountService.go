package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	protoAccount "github.com/asciiu/gomo/account-service/proto/account"
	protoActivity "github.com/asciiu/gomo/activity-bulletin/proto"
	constExch "github.com/asciiu/gomo/common/constants/exchange"
	constRes "github.com/asciiu/gomo/common/constants/response"
	repoKey "github.com/asciiu/gomo/key-service/db/sql"
	protoKey "github.com/asciiu/gomo/key-service/proto/key"
	micro "github.com/micro/go-micro"
)

type AccountService struct {
	DB        *sql.DB
	KeyPub    micro.Publisher
	NotifyPub micro.Publisher
}

// IMPORTANT! When adding support for new exchanges we must add there names here!
// TODO: read these from a config or from the DB
var supported = [...]string{constExch.Binance}

func supportedExchange(a string) bool {
	for _, b := range supported {
		if b == a {
			return true
		}
	}
	return false
}

func (service *AccountService) HandleVerifiedKey(ctx context.Context, key *protoKey.Key) error {
	log.Println("received verified key ", key.KeyID)

	_, error := repoKey.UpdateKeyStatus(service.DB, key)

	description := fmt.Sprintf("%s key verified", key.Exchange)
	notification := protoActivity.Activity{
		UserID:      key.UserID,
		Type:        "key",
		ObjectID:    key.KeyID,
		Title:       "Exchange Setup",
		Description: description,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	// publish verify key event
	if err := service.NotifyPub.Publish(context.Background(), &notification); err != nil {
		log.Println("could not publish verified key event: ", err)
	}

	return error
}

func (service *AccountService) AddAccount(ctx context.Context, req *protoAccount.NewAccountRequest, res *protoAccount.AccountResponse) error {
	// supported exchange keys check
	if !supportedExchange(req.Exchange) {
		res.Status = constRes.Fail
		res.Message = fmt.Sprintf("%s is not a supported exchange", req.Exchange)
		return nil
	}
	return nil
}
func (service *AccountService) DeleteAccount(ctx context.Context, req *protoAccount.AccountRequest, res *protoAccount.AccountResponse) error {
	return nil
}
func (service *AccountService) GetAccounts(ctx context.Context, req *protoAccount.GetAccountsRequest, res *protoAccount.AccountsResponse) error {
	return nil
}
func (service *AccountService) GetAccount(ctx context.Context, req *protoAccount.AccountRequest, res *protoAccount.AccountResponse) error {
	return nil
}
func (service *AccountService) UpdateAccount(ctx context.Context, req *protoAccount.UpdateAccountRequest, res *protoAccount.AccountResponse) error {
	return nil
}
