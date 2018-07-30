package main

import (
	"fmt"
	"log"
	"os"

	msg "github.com/asciiu/gomo/common/constants/messages"
	"github.com/asciiu/gomo/common/db"
	kp "github.com/asciiu/gomo/key-service/proto/key"
	micro "github.com/micro/go-micro"
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

	keyService := KeyService{
		DB:     gomoDB,
		KeyPub: micro.NewPublisher(msg.TopicNewKey, srv.Client()),
	}

	kp.RegisterKeyServiceHandler(srv.Server(), &keyService)

	listener1 := KeyVerifiedListener{
		DB:        gomoDB,
		NotifyPub: micro.NewPublisher(msg.TopicNotification, srv.Client()),
	}

	// handles key verified events
	micro.RegisterSubscriber(msg.TopicKeyVerified, srv.Server(), &listener1)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
