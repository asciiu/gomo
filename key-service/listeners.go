package main

import (
	"context"
	"database/sql"
	"log"

	repo "github.com/asciiu/gomo/key-service/db/sql"
	kp "github.com/asciiu/gomo/key-service/proto/key"
)

type KeyVerifiedListener struct {
	DB *sql.DB
}

func (listener *KeyVerifiedListener) Process(ctx context.Context, key *kp.Key) error {
	log.Println("received verified key ", key.KeyID)

	_, error := repo.UpdateKeyStatus(listener.DB, key)
	return error
}
