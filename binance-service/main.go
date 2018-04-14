package main

import (
	"fmt"
	"log"

	//"github.com/go-kit/kit/log"

	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/cmd"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

func main() {
	cmd.Init()
	if err := broker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}
	if err := broker.Connect(); err != nil {
		log.Fatalf("Broker Connect error: %v", err)
	}

	_, err := broker.Subscribe("new.key", func(p broker.Publication) error {
		VerifyKey(p.Message().Body)
		//fmt.Println("[sub] received message:", string(p.Message().Body), "header", p.Message().Header)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan struct{})
	<-forever
}
