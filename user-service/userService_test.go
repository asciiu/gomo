package main

import (
	"context"
	"testing"

	"github.com/asciiu/gomo/common/db"
	pb "github.com/asciiu/gomo/user-service/proto/user"
)

func setupService() *UserService {
	dbUrl := "postgres://postgres@localhost:5432/gomo_dev?&sslmode=disable"
	db, _ := db.NewDB(dbUrl)
	service := UserService{db}
	return &service
}

func TestInsertUser(t *testing.T) {
	service := setupService()

	request := pb.CreateUserRequest{
		First:    "test",
		Last:     "last",
		Email:    "email@email",
		Password: "password",
	}

	response := pb.UserResponse{
		Data: &pb.UserData{
			&pb.User{},
		},
	}

	service.CreateUser(context.Background(), &request, &response)

	if response.Status != "success" {
		t.Errorf(response.Message)
	}
	if response.Data.User.Email != request.Email {
		t.Errorf("emails do not")
	}

	requestDelete := pb.DeleteUserRequest{
		UserId: response.Data.User.UserId,
		Hard:   true,
	}

	responseDel := pb.Response{}
	service.DeleteUser(context.Background(), &requestDelete, &responseDel)

	if responseDel.Status != "success" {
		t.Errorf(responseDel.Message)
	}
}

func TestChangePassword(t *testing.T) {
	service := setupService()

	request := pb.ChangePasswordRequest{
		UserId:      "7cccddd9-0ee3-4832-a22c-ece20b3084bf",
		OldPassword: "old",
		NewPassword: "new",
	}

	response := pb.Response{}
	service.ChangePassword(context.Background(), &request, &response)
}
