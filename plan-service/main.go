package main

import (
	"fmt"
	"log"
	"os"

	bp "github.com/asciiu/gomo/balance-service/proto/balance"
	msg "github.com/asciiu/gomo/common/constants/messages"
	"github.com/asciiu/gomo/common/db"
	keys "github.com/asciiu/gomo/key-service/proto/key"
	protoPlan "github.com/asciiu/gomo/plan-service/proto/plan"
	micro "github.com/micro/go-micro"
	k8s "github.com/micro/kubernetes/go/micro"
)

func main() {
	dbURL := fmt.Sprintf("%s", os.Getenv("DB_URL"))

	// Create a new service. Include some options here.
	srv := k8s.NewService(
		// This name must match the package name given in your protobuf definition
		micro.Name("fomo.plans"),
		micro.Version("latest"),
	)

	// Init will parse the command line flags.
	srv.Init()

	gomoDB, err := db.NewDB(dbURL)

	if err != nil {
		log.Fatalf(err.Error())
	}

	planService := PlanService{
		DB:        gomoDB,
		Client:    bp.NewBalanceServiceClient("balances", srv.Client()),
		KeyClient: keys.NewKeyServiceClient("keys", srv.Client()),
		NewPlan:   micro.NewPublisher(msg.TopicNewOrder, srv.Client()),
	}

	filledReceiver := OrderFilledReceiver{
		DB:        gomoDB,
		Service:   &planService,
		NotifyPub: micro.NewPublisher(msg.TopicNotification, srv.Client()),
	}

	engineReceiver := EngineStartReceiver{
		Service: &planService,
	}

	micro.RegisterSubscriber(msg.TopicOrderFilled, srv.Server(), &filledReceiver)
	micro.RegisterSubscriber(msg.TopicEngineStart, srv.Server(), &engineReceiver)

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	protoPlan.RegisterPlanServiceHandler(srv.Server(), &planService)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
