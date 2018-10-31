package main

import (
	"fmt"
	"log"
	"os"

	constMessage "github.com/asciiu/gomo/common/constants/message"
	"github.com/asciiu/gomo/common/db"
	protoNotification "github.com/asciiu/gomo/notification-service/proto/notification"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	k8s "github.com/micro/kubernetes/go/micro"
)

func main() {
	srv := k8s.NewService(
		micro.Name("fomo.notifications"),
	)

	// Init will parse the command line flags.
	srv.Init()

	dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	gomoDB, err := db.NewDB(dbUrl)
	if err != nil {
		log.Fatalf(err.Error())
	}

	note := NewNotificationService(gomoDB, srv)
	protoNotification.RegisterNotificationHandler(srv.Server(), note)

	// handles key verified events
	micro.RegisterSubscriber(constMessage.TopicNotification, srv.Server(), note.LogActivity, server.SubscriberQueue("notifications"))

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
