package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	userRepo "github.com/asciiu/gomo/user-service/db/sql"
	"github.com/asciiu/gomo/user-service/models"
	pb "github.com/asciiu/gomo/user-service/proto/user"
)

type UserService struct {
	DB *sql.DB
}

func (service *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest, res *pb.UserResponse) error {
	user := models.NewUser(req.First, req.Last, req.Email, req.Password)
	_, error := userRepo.InsertUser(service.DB, user)

	switch {
	case error == nil:
		res.Status = "success"
		res.Data.User.UserId = user.Id
		res.Data.User.First = user.First
		res.Data.User.Last = user.Last
		res.Data.User.Email = user.Email
		return nil

	case strings.Contains(error.Error(), "violates unique constraint \"users_email_key\""):
		res.Status = "fail"
		res.Message = "email already exists"
		return error

	default:
		res.Status = "error"
		res.Message = error.Error()
		return error
	}
}

// Deletes a user
func (service *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest, res *pb.Response) error {
	var err error
	if req.Hard {
		err = userRepo.DeleteUserHard(service.DB, req.UserId)
	} else {
		err = userRepo.DeleteUserSoft(service.DB, req.UserId)
	}

	if err == nil {
		res.Status = "success"
	} else {
		res.Status = "error"
		res.Message = err.Error()
	}
	return err
}

func (s *UserService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest, res *pb.Response) error {
	_, error := userRepo.FindUserById(s.DB, req.UserId)
	if error != nil {
		log.Println(error)
	}

	//fmt.Println(user)

	//fmt.Printf("%s %s %s", req.UserId, req.OldPassword, req.NewPassword)

	return nil
}

func (s *UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest, res *pb.UserResponse) error {
	fmt.Println(req)
	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest, res *pb.UserResponse) error {
	fmt.Println(req)
	return nil
}
