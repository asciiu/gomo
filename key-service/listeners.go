package main

import (
	"context"
	"database/sql"
	"log"

	repo "github.com/asciiu/gomo/key-service/db/sql"
	kp "github.com/asciiu/gomo/key-service/proto/key"
	notifications "github.com/asciiu/gomo/notification-service/proto"
	micro "github.com/micro/go-micro"
)

type KeyVerifiedListener struct {
	DB        *sql.DB
	NotifyPub micro.Publisher
}

func (listener *KeyVerifiedListener) Process(ctx context.Context, key *kp.Key) error {
	log.Println("received verified key ", key.KeyID)

	_, error := repo.UpdateKeyStatus(listener.DB, key)

	notification := notifications.Notification{
		UserID:      key.UserID,
		Description: "your key has been verified",
	}

	// publish verify key event
	if err := listener.NotifyPub.Publish(context.Background(), &notification); err != nil {
		log.Println("could not publish verified key event: ", err)
	}

	return error
}
