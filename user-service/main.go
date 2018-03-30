package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/asciiu/gomo/common/db"
	pb "github.com/asciiu/gomo/user-service/proto/user"
	micro "github.com/micro/go-micro"
)

type service struct {
	DB *sql.DB
}

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *service) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest, res *pb.Response) error {
	return nil
}

func (s *service) GetUserInfo(ctx context.Context, req *pb.UserInfoRequest, res *pb.Response) error {
	return nil
}

func main() {
	// Create a new service. Include some options here.
	srv := micro.NewService(
		// This name must match the package name given in your protobuf definition
		micro.Name("go.micro.srv.user"),
		micro.Version("latest"),
	)
	// Init will parse the command line flags.
	srv.Init()

	dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	fmt.Println(dbUrl)

	gomoDB, err := db.NewDB(dbUrl)

	// TODO read secret from env var
	//dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	pb.RegisterUserServiceHandler(srv.Server(), &service{gomoDB})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
