package main

import (
	"fmt"
	"log"
	"os"

	bp "github.com/asciiu/gomo/balance-service/proto/balance"
	"github.com/asciiu/gomo/common/db"
	op "github.com/asciiu/gomo/order-service/proto/order"

	micro "github.com/micro/go-micro"
)

func NewOrderService(name, dbUrl string) micro.Service {
	// Create a new service. Include some options here.
	srv := micro.NewService(
		// This name must match the package name given in your protobuf definition
		micro.Name(name),
		micro.Version("latest"),
	)

	// Init will parse the command line flags.
	srv.Init()

	gomoDB, err := db.NewDB(dbUrl)

	// TODO read secret from env var
	//dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf(err.Error())
	}

	orderService := OrderService{
		DB:     gomoDB,
		Client: bp.NewBalanceServiceClient("go.micro.srv.balance", srv.Client()),
	}

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	op.RegisterOrderServiceHandler(srv.Server(), &orderService)

	return srv
}

func main() {
	dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	srv := NewOrderService("go.srv.order-service", dbUrl)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
