package main

import (
	"fmt"
	"log"
	"os"

	constMessage "github.com/asciiu/gomo/common/constants/message"
	"github.com/asciiu/gomo/common/db"
	protoKey "github.com/asciiu/gomo/key-service/proto/key"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	k8s "github.com/micro/kubernetes/go/micro"
)

func main() {
	srv := k8s.NewService(
		micro.Name("keys"),
	)

	// Init will parse the command line flags.
	srv.Init()

	dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	gomoDB, err := db.NewDB(dbUrl)
	if err != nil {
		log.Fatalf(err.Error())
	}

	accountService := AccountService{
		DB:        gomoDB,
		KeyPub:    micro.NewPublisher(constMessage.TopicNewKey, srv.Client()),
		NotifyPub: micro.NewPublisher(constMessage.TopicNotification, srv.Client()),
	}

	protoKey.RegisterKeyServiceHandler(srv.Server(), &accountService)

	// handles key verified events
	micro.RegisterSubscriber(constMessage.TopicKeyVerified, srv.Server(), accountService.HandleVerifiedKey, server.SubscriberQueue("verified.key"))

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
