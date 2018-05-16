package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/appleboy/gorush/rpc/proto"
	"github.com/asciiu/gomo/common/db"
	msg "github.com/asciiu/gomo/common/messages"
	notification "github.com/asciiu/gomo/notification-service/proto/notification"
	micro "github.com/micro/go-micro"
	"google.golang.org/grpc"
)

const (
	address = "localhost:9000"
)

func TestPush() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewGorushClient(conn)

	r, err := c.Send(context.Background(), &proto.NotificationRequest{
		Platform: 2,
		Tokens:   []string{"1234567890"},
		Message:  "test message",
		Badge:    1,
		Category: "test",
		Sound:    "test",
		Alert: &proto.Alert{
			Title:    "Test Title",
			Body:     "Test Alert Body",
			Subtitle: "Test Alert Sub Title",
			LocKey:   "Test loc key",
			LocArgs:  []string{"test", "test"},
		},
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Success: %t\n", r.Success)
	log.Printf("Count: %d\n", r.Counts)
}

func main() {
	srv := micro.NewService(
		micro.Name("go.srv.notification-service"),
	)

	// Init will parse the command line flags.
	srv.Init()

	dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	gomoDB, err := db.NewDB(dbUrl)
	if err != nil {
		log.Fatalf(err.Error())
	}

	notificationService := NotificationService{gomoDB}

	notification.RegisterNotificationServiceHandler(srv.Server(), &notificationService)

	listener1 := NotificationListener{gomoDB}
	// handles key verified events
	micro.RegisterSubscriber(msg.TopicNotification, srv.Server(), &listener1)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
