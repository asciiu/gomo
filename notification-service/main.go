package main

import (
	"fmt"
	"log"
	"os"

	"github.com/asciiu/gomo/common/db"
	msg "github.com/asciiu/gomo/common/messages"
	notification "github.com/asciiu/gomo/notification-service/proto"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
)

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

	listener1 := NewNotificationListener(gomoDB, srv)

	// handles key verified events
	micro.RegisterSubscriber(msg.TopicNotification, srv.Server(), listener1.ProcessNotification, server.SubscriberQueue("queue.pubsub"))

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
