package main

import (
	"fmt"
	"log"
	"os"

	"github.com/asciiu/gomo/common/db"
	pb "github.com/asciiu/gomo/device-service/proto/device"
	micro "github.com/micro/go-micro"
)

func NewDeviceService(name, dbUrl string) micro.Service {
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

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	pb.RegisterDeviceServiceHandler(srv.Server(), &DeviceService{gomoDB})

	return srv
}

func main() {
	dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	srv := NewDeviceService("go.srv.device-service", dbUrl)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}