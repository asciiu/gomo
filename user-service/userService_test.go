package main

import (
	"context"
	"strings"
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

	invalidChangeReq := pb.ChangePasswordRequest{
		UserId:      response.Data.User.UserId,
		OldPassword: "pass",
		NewPassword: "new",
	}

	response2 := pb.Response{}
	service.ChangePassword(context.Background(), &invalidChangeReq, &response2)
	if !strings.Contains(response2.Message, "current password mismatch") {
		t.Errorf(response.Message)
	}

	validChangeReq := pb.ChangePasswordRequest{
		UserId:      response.Data.User.UserId,
		OldPassword: "password",
		NewPassword: "new",
	}

	response3 := pb.Response{}
	eor := service.ChangePassword(context.Background(), &validChangeReq, &response3)
	if response3.Status != "success" {
		t.Errorf(eor.Error())
	}

	requestDelete := pb.DeleteUserRequest{
		UserId: response.Data.User.UserId,
		Hard:   true,
	}

	responseDel := pb.Response{}
	service.DeleteUser(context.Background(), &requestDelete, &responseDel)
}
