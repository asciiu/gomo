package main

import (
	"log"

	msg "github.com/asciiu/gomo/common/messages"
	micro "github.com/micro/go-micro"
)

func main() {
	srv := micro.NewService(
		micro.Name("go.micro.srv.bittrex"),
	)

	srv.Init()

	verifiedPub := micro.NewPublisher(msg.TopicKeyVerified, srv.Client())
	balancePub := micro.NewPublisher(msg.TopicBalanceUpdate, srv.Client())
	keyValidator := KeyValidator{
		KeyVerifiedPub: verifiedPub,
		BalancePub:     balancePub,
	}

	// subscribe to new key topic with a key validator
	micro.RegisterSubscriber(msg.TopicNewKey, srv.Server(), &keyValidator)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
