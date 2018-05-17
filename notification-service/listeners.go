package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/appleboy/gorush/rpc/proto"
	notification "github.com/asciiu/gomo/notification-service/proto"
	"google.golang.org/grpc"
)

type NotificationListener struct {
	DB     *sql.DB
	client proto.GorushClient
	topic  string
}

func NewNotificationListener(db *sql.DB) *NotificationListener {
	address := fmt.Sprintf("%s", os.Getenv("GORUSH_ADDRESS"))
	topic := fmt.Sprintf("%s", os.Getenv("APNS_TOPIC"))

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := proto.NewGorushClient(conn)

	listener := NotificationListener{
		DB:     db,
		client: client,
		topic:  topic,
	}

	return &listener
}

func (listener *NotificationListener) Process(ctx context.Context, note *notification.Notification) error {
	// get device tokens from DB for user ID
	log.Println("notification ", note.Description)

	// loop over device tokens and send
	r, err := listener.client.Send(context.Background(), &proto.NotificationRequest{
		Platform: 1,
		Tokens:   []string{"d98b367c3cdb9d2c2a52a4fe0cc40f95c693ac0d87e7c4fb41988afc3d46111c"},
		Message:  "test message",
		Badge:    1,
		Category: "test",
		Sound:    "test",
		Topic:    listener.topic,
		Alert: &proto.Alert{
			Title:    "Test Title",
			Body:     note.Description,
			Subtitle: "Test Alert Sub Title",
			LocKey:   "Test loc key",
			LocArgs:  []string{"test", "test"},
		},
	})
	if err != nil {
		log.Printf("could not send: %v\n", err)
	}
	log.Printf("Success: %t\n", r.Success)
	log.Printf("Count: %d\n", r.Counts)

	return nil
}
