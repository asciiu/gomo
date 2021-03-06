package main

import (
	"fmt"
	"log"
	"os"

	protoAccount "github.com/asciiu/gomo/account-service/proto/account"
	protoAnalytics "github.com/asciiu/gomo/analytics-service/proto/analytics"
	constMessage "github.com/asciiu/gomo/common/constants/message"
	"github.com/asciiu/gomo/common/db"
	protoEngine "github.com/asciiu/gomo/execution-engine/proto/engine"
	protoPlan "github.com/asciiu/gomo/plan-service/proto/plan"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
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
		DB:              gomoDB,
		AccountClient:   protoAccount.NewAccountServiceClient("accounts", srv.Client()),
		AnalyticsClient: protoAnalytics.NewAnalyticsServiceClient("analytics", srv.Client()),
		EngineClient:    protoEngine.NewExecutionEngineClient("engine", srv.Client()),
		NotifyPub:       micro.NewPublisher(constMessage.TopicNotification, srv.Client()),
	}

	micro.RegisterSubscriber(constMessage.TopicCompletedOrder, srv.Server(), planService.HandleCompletedOrder, server.SubscriberQueue("complete.order"))
	micro.RegisterSubscriber(constMessage.TopicEngineStart, srv.Server(), planService.HandleStartEngine, server.SubscriberQueue("pop.engine"))
	micro.RegisterSubscriber(constMessage.TopicAccountDeleted, srv.Server(), planService.HandleAccountDeleted, server.SubscriberQueue("delete.account"))

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	protoPlan.RegisterPlanServiceHandler(srv.Server(), &planService)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
