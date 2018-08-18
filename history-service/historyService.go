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

func (service *HistoryService) Archive(ctx context.Context, history *protoHistory.History) error {

	_, error := repoHistory.InsertHistory(service.db, history)
	if error != nil {
		log.Println("could not insert new notification ", error)
	}

	log.Println("notification ", history.Description)

	getRequest := protoDevice.GetUserDevicesRequest{
		UserID: history.UserID,
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
		Message:  "fomo",
		Badge:    1,
		Category: "test",
		Sound:    "1",
		Topic:    service.topic,
		Alert: &protoGorush.Alert{
			Title:    history.Title,
			Body:     history.Description,
			Subtitle: history.Subtitle,
			LocKey:   "Test loc key",
			LocArgs:  []string{"test", "test"},
		},
	})

	if err != nil {
		log.Printf("could not send: %v\n", err)
	} else {
		log.Printf("Success: %t\n", r.Success)
	}

	return nil
}

func (service *HistoryService) FindUserHistory(ctx context.Context, req *protoHistory.HistoryRequest, res *protoHistory.HistoryPagedResponse) error {
	var pagedResult *protoHistory.UserHistoryPage
	var err error

	if req.ObjectID != "" {
		// history associated with object ID only
		pagedResult, err = repoHistory.FindObjectHistory(service.db, req)
	} else {
		// all user history
		pagedResult, err = repoHistory.FindUserHistory(service.db, req.UserID, req.Page, req.PageSize)
	}

	if err == nil {
		res.Status = "success"
		res.Data = pagedResult
	} else {
		res.Status = "error"
		res.Message = err.Error()
	}
	return nil
}

func (service *HistoryService) FindMostRecentHistory(ctx context.Context, req *protoHistory.RecentHistoryRequest, res *protoHistory.HistoryListResponse) error {
	history, err := repoHistory.FindRecentObjectHistory(service.db, req)
	if err == nil {
		res.Status = "success"
		res.Data = &protoHistory.HistoryList{
			History: history,
		}
	} else {
		res.Status = "error"
		res.Message = err.Error()
	}
	return nil
}

func (service *HistoryService) FindHistoryCount(ctx context.Context, req *protoHistory.HistoryCountRequest, res *protoHistory.HistoryCountResponse) error {
	return nil
}

func (service *HistoryService) UpdateHistory(ctx context.Context, req *protoHistory.UpdateHistoryRequest, res *protoHistory.HistoryResponse) error {
	return nil
}
