package main

import (
	"fmt"
	"log"
	"os"

	msg "github.com/asciiu/gomo/apikey-service/models"
	"github.com/asciiu/gomo/common/db"
	micro "github.com/micro/go-micro"
)

func main() {
	srv := micro.NewService(
		micro.Name("go.micro.srv.binance"),
	)

	srv.Init()

	dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	gomoDB, err := db.NewDB(dbUrl)

	if err != nil {
		log.Fatalf(err.Error())
	}

	// subscribe to new key topic with a key validator
	micro.RegisterSubscriber(msg.TopicNewKey, srv.Server(), &KeyValidator{gomoDB})
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
