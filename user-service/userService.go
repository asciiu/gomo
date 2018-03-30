package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	userRepo "github.com/asciiu/gomo/user-service/db/sql"
	"github.com/asciiu/gomo/user-service/models"
	pb "github.com/asciiu/gomo/user-service/proto/user"
	"golang.org/x/crypto/bcrypt"
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

// Changes the user's password. Password is updated when the request's
// old password matches the current user's password hash.
func (service *UserService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest, res *pb.Response) error {
	user, error := userRepo.FindUserById(service.DB, req.UserId)

	switch {
	case error == nil:
		if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)) == nil {

			err := userRepo.UpdateUserPassword(service.DB, req.UserId, models.HashAndSalt([]byte(req.NewPassword)))
			if err != nil {
				res.Status = "error"
				res.Message = err.Error()
			} else {
				res.Status = "success"
			}

		} else {
			res.Status = "fail"
			res.Message = "current password mismatch"
		}

	case strings.Contains(error.Error(), "no rows in result set"):
		res.Status = "fail"
		res.Message = fmt.Sprintf("user id not found: %s", req.UserId)

	default:
		res.Status = "error"
		res.Message = error.Error()
	}

	return error
}

func (service *UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest, res *pb.UserResponse) error {
	user, error := userRepo.FindUserById(service.DB, req.UserId)
	if error != nil {
		res.Status = "fail"
		res.Message = error.Error()
	} else if error == nil {
		res.Status = "success"
		res.Data.User.UserId = user.Id
		res.Data.User.First = user.First
		res.Data.User.Last = user.Last
		res.Data.User.Email = user.Email
	}

	return error
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest, res *pb.UserResponse) error {
	fmt.Println(req)
	return nil
}
