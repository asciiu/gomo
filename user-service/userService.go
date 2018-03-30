package main

import (
	"context"
	"database/sql"
	"fmt"

	pb "github.com/asciiu/gomo/user-service/proto/user"
)

type UserService struct {
	DB *sql.DB
}

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *UserService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest, res *pb.Response) error {
	fmt.Println(req)

	return nil
}

func (s *UserService) GetUserInfo(ctx context.Context, req *pb.UserInfoRequest, res *pb.Response) error {
	return nil
}
