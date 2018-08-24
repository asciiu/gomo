package main

import (
	"fmt"
	"log"
	"os"

	protoBalance "github.com/asciiu/gomo/balance-service/proto/balance"
	constMessage "github.com/asciiu/gomo/common/constants/message"
	"github.com/asciiu/gomo/common/db"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	k8s "github.com/micro/kubernetes/go/micro"
)

func main() {
	srv := k8s.NewService(
		micro.Name("fomo.balances"),
	)

	srv.Init()

	dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	gomoDB, err := db.NewDB(dbUrl)
	if err != nil {
		log.Fatalf(err.Error())
	}

	balanceService := BalanceService{gomoDB}

	// subscribe to new key topic with a key validator
	micro.RegisterSubscriber(constMessage.TopicBalanceUpdate, srv.Server(), balanceService.HandleBalances, server.SubscriberQueue("update.balances"))

	protoBalance.RegisterBalanceServiceHandler(srv.Server(), &balanceService)
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
