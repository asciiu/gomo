package main

import (
	"context"
	"database/sql"
	"log"

	notification "github.com/asciiu/gomo/notification-service/proto"
)

type NotificationListener struct {
	DB *sql.DB
}

func (listener *NotificationListener) Process(ctx context.Context, note *notification.Notification) error {
	log.Println("notification ", note.Description)

	//_, error := repo.UpdateKeyStatus(listener.DB, key)
	return nil
}
