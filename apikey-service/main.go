package main

import (
	"fmt"
	"log"
	"os"

	pb "github.com/asciiu/gomo/apikey-service/proto/apikey"
	"github.com/asciiu/gomo/common/db"
	micro "github.com/micro/go-micro"
)

func NewKeyService(name, dbUrl string) micro.Service {
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
	pb.RegisterApiKeyServiceHandler(srv.Server(), &KeyService{gomoDB})

	return srv
}

func main() {
	dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	srv := NewKeyService("go.srv.apikey-service", dbUrl)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
