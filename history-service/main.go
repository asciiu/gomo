package main

import (
	"fmt"
	"log"
	"os"

	constMessage "github.com/asciiu/gomo/common/constants/messages"
	"github.com/asciiu/gomo/common/db"
	protoHistory "github.com/asciiu/gomo/history-service/proto"
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

	historyService := NewHistoryService(gomoDB, srv)

	protoHistory.RegisterHistoryServiceHandler(srv.Server(), historyService)

	// handles key verified events
	micro.RegisterSubscriber(constMessage.TopicNotification, srv.Server(), historyService.Archive, server.SubscriberQueue("archive"))

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
