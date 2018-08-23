package main

import (
	"context"
	"strings"
	"testing"

	"github.com/asciiu/gomo/common/db"
	protoUser "github.com/asciiu/gomo/user-service/proto/user"
)

func setupService() *UserService {
	dbUrl := "postgres://postgres@localhost:5432/gomo_dev?&sslmode=disable"
	db, _ := db.NewDB(dbUrl)
	service := UserService{db}
	return &service
}

func TestInsertUser(t *testing.T) {
	service := setupService()

	request := protoUser.CreateUserRequest{
		First:    "test",
		Last:     "last",
		Email:    "email@email",
		Password: "password",
	}

	response := protoUser.UserResponse{}

	service.CreateUser(context.Background(), &request, &response)

	if response.Status != "success" {
		t.Errorf(response.Message)
	}
	if response.Data.User.Email != request.Email {
		t.Errorf("emails do not")
	}

	requestDelete := protoUser.DeleteUserRequest{
		UserID: response.Data.User.UserID,
		Hard:   true,
	}

	responseDel := protoUser.Response{}
	service.DeleteUser(context.Background(), &requestDelete, &responseDel)

	if responseDel.Status != "success" {
		t.Errorf(responseDel.Message)
	}
}

func TestChangePassword(t *testing.T) {
	service := setupService()
	request := protoUser.CreateUserRequest{
		First:    "test",
		Last:     "last",
		Email:    "email@email",
		Password: "password",
	}

	response := protoUser.UserResponse{}

	service.CreateUser(context.Background(), &request, &response)

	if response.Status != "success" {
		t.Errorf(response.Message)
	}

	invalidChangeReq := protoUser.ChangePasswordRequest{
		UserID:      response.Data.User.UserID,
		OldPassword: "pass",
		NewPassword: "new",
	}

	response2 := protoUser.Response{}
	service.ChangePassword(context.Background(), &invalidChangeReq, &response2)
	if !strings.Contains(response2.Message, "current password mismatch") {
		t.Errorf(response.Message)
	}

	validChangeReq := protoUser.ChangePasswordRequest{
		UserID:      response.Data.User.UserID,
		OldPassword: "password",
		NewPassword: "new",
	}

	response3 := protoUser.Response{}
	eor := service.ChangePassword(context.Background(), &validChangeReq, &response3)
	if response3.Status != "success" {
		t.Errorf(eor.Error())
	}

	requestDelete := protoUser.DeleteUserRequest{
		UserID: response.Data.User.UserID,
		Hard:   true,
	}

	responseDel := protoUser.Response{}
	service.DeleteUser(context.Background(), &requestDelete, &responseDel)
}

func TestGetUserInfo(t *testing.T) {
	service := setupService()
	request := protoUser.CreateUserRequest{
		First:    "Bobbie",
		Last:     "McGee",
		Email:    "bobbie@luv",
		Password: "password",
	}

	response := protoUser.UserResponse{}

	service.CreateUser(context.Background(), &request, &response)

	if response.Status != "success" {
		t.Errorf(response.Message)
	}

	getRequest := protoUser.GetUserInfoRequest{
		UserID: response.Data.User.UserID,
	}
	service.GetUserInfo(context.Background(), &getRequest, &response)

	if response.Status != "success" {
		t.Errorf(response.Message)
	}

	requestDelete := protoUser.DeleteUserRequest{
		UserID: response.Data.User.UserID,
		Hard:   true,
	}

	responseDel := protoUser.Response{}
	service.DeleteUser(context.Background(), &requestDelete, &responseDel)
}

func TestUpdateUser(t *testing.T) {
	service := setupService()
	request := protoUser.CreateUserRequest{
		First:    "Bobbie",
		Last:     "McGee",
		Email:    "bobbie@luv",
		Password: "password",
	}

	response := protoUser.UserResponse{}

	service.CreateUser(context.Background(), &request, &response)

	updateRequest := protoUser.UpdateUserRequest{
		UserID: response.Data.User.UserID,
		First:  "Bobby",
		Last:   "McLovin",
		Email:  "bobby@mcLovin",
	}

	service.UpdateUser(context.Background(), &updateRequest, &response)

	if response.Status != "success" {
		t.Errorf(response.Message)
	}

	if response.Data.User.First != updateRequest.First {
		t.Errorf("first not updated")
	}
	if response.Data.User.Last != updateRequest.Last {
		t.Errorf("last not updated")
	}
	if response.Data.User.Email != updateRequest.Email {
		t.Errorf("email not updated")
	}

	requestDelete := protoUser.DeleteUserRequest{
		UserID: response.Data.User.UserID,
		Hard:   true,
	}

	responseDel := protoUser.Response{}
	service.DeleteUser(context.Background(), &requestDelete, &responseDel)
}
