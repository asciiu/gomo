package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

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

	description := fmt.Sprintf("%s key verified", key.Exchange)
	notification := notifications.Notification{
		UserID:           key.UserID,
		NotificationType: "key",
		ObjectID:         key.KeyID,
		Title:            "Exchange Setup",
		Description:      description,
		Timestamp:        time.Now().UTC().Format(time.RFC3339),
	}

	// publish verify key event
	if err := listener.NotifyPub.Publish(context.Background(), &notification); err != nil {
		log.Println("could not publish verified key event: ", err)
	}

	return error
}
