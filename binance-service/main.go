package main

import (
	"log"

	msg "github.com/asciiu/gomo/apikey-service/models"
	micro "github.com/micro/go-micro"
)

func main() {
	srv := micro.NewService(
		micro.Name("go.micro.srv.binance"),
	)

	srv.Init()

	// subscribe to new key topic with a key validator
	micro.RegisterSubscriber(msg.TopicNewKey, srv.Server(), new(KeyValidator))
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
