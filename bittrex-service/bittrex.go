package main

import (
	"context"
	"fmt"

	kp "github.com/asciiu/gomo/apikey-service/proto/apikey"
	micro "github.com/micro/go-micro"
)

type KeyValidator struct {
	KeyVerifiedPub micro.Publisher
	BalancePub     micro.Publisher
}

func (service *KeyValidator) Process(ctx context.Context, key *kp.ApiKey) error {
	if key.Exchange != "Bittrex" {
		return nil
	}

	fmt.Println(key)

	return nil
}
