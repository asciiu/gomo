package main

import (
	"context"
	"database/sql"
	"log"

	keyRepo "github.com/asciiu/gomo/apikey-service/db/sql"
	keyProto "github.com/asciiu/gomo/apikey-service/proto/apikey"
)

type KeyVerifiedListener struct {
	DB *sql.DB
}

func (listener *KeyVerifiedListener) Process(ctx context.Context, key *keyProto.ApiKey) error {
	log.Println("received verified key ", key.ApiKeyId)

	_, error := keyRepo.UpdateApiKeyStatus(listener.DB, key)
	return error
}
