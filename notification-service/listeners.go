package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/appleboy/gorush/rpc/proto"
	devices "github.com/asciiu/gomo/device-service/proto/device"
	notification "github.com/asciiu/gomo/notification-service/proto"
	micro "github.com/micro/go-micro"
	"google.golang.org/grpc"
)

type NotificationListener struct {
	db      *sql.DB
	devices devices.DeviceServiceClient
	client  proto.GorushClient
	topic   string
}

func NewNotificationListener(db *sql.DB, service micro.Service) *NotificationListener {
	address := fmt.Sprintf("%s", os.Getenv("GORUSH_ADDRESS"))
	topic := fmt.Sprintf("%s", os.Getenv("APNS_TOPIC"))

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := proto.NewGorushClient(conn)

	listener := NotificationListener{
		db:      db,
		client:  client,
		topic:   topic,
		devices: devices.NewDeviceServiceClient("go.srv.device-service", service.Client()),
	}

	return &listener
}

func (listener *NotificationListener) ProcessNotification(ctx context.Context, note *notification.Notification) error {
	// get device tokens from DB for user ID
	log.Println("notification ", note.Description)

	getRequest := devices.GetUserDevicesRequest{
		UserID: note.UserID,
	}

	ds, _ := listener.devices.GetUserDevices(context.Background(), &getRequest)

	if ds.Status != "success" {
		log.Println("error from GetUserDevices ", ds.Message)
		return errors.New(ds.Message)
	}

	iosTokens := make([]string, 0)

	for _, thing := range ds.Data.Devices {
		deviceType := thing.DeviceType
		deviceToken := thing.DeviceToken
		if deviceType == "ios" {
			iosTokens = append(iosTokens, deviceToken)
		}
	}

	// loop over device tokens and send
	r, err := listener.client.Send(context.Background(), &proto.NotificationRequest{
		Platform: 1,
		Tokens:   iosTokens,
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
