package main

import (
	"fmt"
	"log"
	"os"

	protoActivity "github.com/asciiu/gomo/activity-bulletin/proto"
	constMessage "github.com/asciiu/gomo/common/constants/messages"
	"github.com/asciiu/gomo/common/db"
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

	bulletin := NewBulletin(gomoDB, srv)

	protoActivity.RegisterActivityBulletinHandler(srv.Server(), bulletin)

	// handles key verified events
	micro.RegisterSubscriber(constMessage.TopicNotification, srv.Server(), bulletin.LogActivity, server.SubscriberQueue("archive"))

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
