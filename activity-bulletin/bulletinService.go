package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	protoGorush "github.com/appleboy/gorush/rpc/proto"
	repoActivity "github.com/asciiu/gomo/activity-bulletin/db/sql"
	protoActivity "github.com/asciiu/gomo/activity-bulletin/proto"
	constResponse "github.com/asciiu/gomo/common/constants/response"
	protoDevice "github.com/asciiu/gomo/device-service/proto/device"
	"github.com/google/uuid"
	micro "github.com/micro/go-micro"
	"google.golang.org/grpc"
)

type Bulletin struct {
	db      *sql.DB
	devices protoDevice.DeviceServiceClient
	client  protoGorush.GorushClient
	topic   string
}

func NewBulletin(db *sql.DB, service micro.Service) *Bulletin {
	address := fmt.Sprintf("%s", os.Getenv("GORUSH_ADDRESS"))
	topic := fmt.Sprintf("%s", os.Getenv("APNS_TOPIC"))

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := protoGorush.NewGorushClient(conn)

	hs := Bulletin{
		db:      db,
		client:  client,
		topic:   topic,
		devices: protoDevice.NewDeviceServiceClient("devices", service.Client()),
	}

	return &hs
}

func (service *Bulletin) LogActivity(ctx context.Context, history *protoActivity.Activity) error {

	_, error := repoActivity.InsertActivity(service.db, history)
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

func (service *Bulletin) FindUserActivity(ctx context.Context, req *protoActivity.ActivityRequest, res *protoActivity.ActivityPagedResponse) error {
	var pagedResult *protoActivity.UserActivityPage
	var err error

	if req.ObjectID != "" {
		if _, err := uuid.Parse(req.ObjectID); err != nil {
			res.Status = constResponse.Fail
			res.Message = fmt.Sprintf("object %s not found", req.ObjectID)
			return nil
		}
		// history associated with object ID only
		pagedResult, err = repoActivity.FindObjectActivity(service.db, req)
	} else {
		// all user history
		pagedResult, err = repoActivity.FindUserActivity(service.db, req.UserID, req.Page, req.PageSize)
	}

	if err == nil {
		res.Status = constResponse.Success
		res.Data = pagedResult
	} else {
		res.Status = constResponse.Error
		res.Message = err.Error()
	}
	return nil
}

func (service *Bulletin) FindMostRecentActivity(ctx context.Context, req *protoActivity.RecentActivityRequest, res *protoActivity.ActivityListResponse) error {
	history, err := repoActivity.FindRecentObjectActivity(service.db, req)
	if err == nil {
		res.Status = constResponse.Success
		res.Data = &protoActivity.ActivityList{
			Activity: history,
		}
	} else {
		res.Status = constResponse.Error
		res.Message = err.Error()
	}
	return nil
}

func (service *Bulletin) FindActivityCount(ctx context.Context, req *protoActivity.ActivityCountRequest, res *protoActivity.ActivityCountResponse) error {
	count := repoActivity.FindObjectActivityCount(service.db, req.ObjectID)
	res.Status = constResponse.Success
	res.Data = &protoActivity.ActivityCount{
		Count: count,
	}
	return nil
}

func (service *Bulletin) UpdateActivity(ctx context.Context, req *protoActivity.UpdateActivityRequest, res *protoActivity.ActivityResponse) error {

	var history *protoActivity.Activity
	var err error

	if req.ClickedAt != "" {
		history, err = repoActivity.UpdateActivityClickedAt(service.db, req.ActivityID, req.ClickedAt)
		if err != nil {
			res.Status = constResponse.Error
			res.Message = err.Error()
			return nil
		}
	}
	if req.SeenAt != "" {
		history, err = repoActivity.UpdateActivitySeenAt(service.db, req.ActivityID, req.SeenAt)
		if err != nil {
			res.Status = constResponse.Error
			res.Message = err.Error()
			return nil
		}
	}

	res.Status = constResponse.Success
	res.Data = &protoActivity.ActivityData{
		Activity: history,
	}

	return nil
}
