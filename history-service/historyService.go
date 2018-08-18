package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	protoGorush "github.com/appleboy/gorush/rpc/proto"
	protoDevice "github.com/asciiu/gomo/device-service/proto/device"
	repoHistory "github.com/asciiu/gomo/history-service/db/sql"
	protoHistory "github.com/asciiu/gomo/history-service/proto"
	notification "github.com/asciiu/gomo/notification-service/proto"
	micro "github.com/micro/go-micro"
	"google.golang.org/grpc"
)

type HistoryService struct {
	db      *sql.DB
	devices protoDevice.DeviceServiceClient
	client  protoGorush.GorushClient
	topic   string
}

func NewHistoryService(db *sql.DB, service micro.Service) *HistoryService {
	address := fmt.Sprintf("%s", os.Getenv("GORUSH_ADDRESS"))
	topic := fmt.Sprintf("%s", os.Getenv("APNS_TOPIC"))

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := protoGorush.NewGorushClient(conn)

	hs := HistoryService{
		db:      db,
		client:  client,
		topic:   topic,
		devices: protoDevice.NewDeviceServiceClient("devices", service.Client()),
	}

	return &hs
}

func (service *HistoryService) Archive(ctx context.Context, note *notification.Notification) error {

	_, error := repoHistory.InsertNotification(service.db, note)
	if error != nil {
		log.Println("could not insert new notification ", error)
	}

	log.Println("notification ", note.Description)

	getRequest := protoDevice.GetUserDevicesRequest{
		UserID: note.UserID,
	}

	// get device tokens from DB for user ID
	ds, _ := service.devices.GetUserDevices(context.Background(), &getRequest)

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
	// TODO fill in the rest
	r, err := service.client.Send(context.Background(), &protoGorush.NotificationRequest{
		Platform: 1,
		Tokens:   iosTokens,
		Message:  "test message",
		Badge:    1,
		Category: "test",
		Sound:    "1",
		Topic:    service.topic,
		Alert: &protoGorush.Alert{
			Title:    note.Title,
			Body:     note.Description,
			Subtitle: note.Subtitle,
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

func (service *HistoryService) FindUserHistory(ctx context.Context, req *protoHistory.HistoryRequest, res *protoHistory.HistoryPagedResponse) error {

	// var pagedResult *protoHistory.UserNotificationsPage
	// var err error
	// if req.NotificationType == "" {
	// 	pagedResult, err = repoHistory.FindNotifications(service.DB, req)
	// } else {
	// 	pagedResult, err = repoHistory.FindNotificationsByType(service.DB, req)
	// }

	// switch {
	// case err == nil:
	// 	res.Status = "success"
	// 	res.Data = pagedResult
	// default:
	// 	res.Status = "error"
	// 	res.Message = err.Error()
	// }
	return nil

}

func (service *HistoryService) FindMostRecentHistory(ctx context.Context, req *protoHistory.RecentHistoryRequest, res *protoHistory.HistoryListResponse) error {
	return nil
}

func (service *HistoryService) FindHistoryCount(ctx context.Context, req *protoHistory.HistoryCountRequest, res *protoHistory.HistoryCountResponse) error {
	return nil
}

func (service *HistoryService) UpdateHistory(ctx context.Context, req *protoHistory.UpdateHistoryRequest, res *protoHistory.HistoryResponse) error {
	return nil
}
