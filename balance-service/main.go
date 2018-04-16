package main

import (
	"fmt"
	"log"
	"os"

	"github.com/asciiu/gomo/common/db"
	msg "github.com/asciiu/gomo/common/messages"
	micro "github.com/micro/go-micro"
)

func main() {
	srv := micro.NewService(
		micro.Name("go.micro.srv.balance"),
	)

	srv.Init()

	dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	gomoDB, err := db.NewDB(dbUrl)

	if err != nil {
		log.Fatalf(err.Error())
	}

	// subscribe to new key topic with a key validator
	micro.RegisterSubscriber(msg.TopicBalanceUpdate, srv.Server(), &BalancerUpdater{gomoDB, srv})
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
