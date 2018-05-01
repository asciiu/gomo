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
		res.Data = &pb.UserData{
			User: &pb.User{
				UserID: user.ID,
				First:  user.First,
				Last:   user.Last,
				Email:  user.Email,
			},
		}
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
		err = userRepo.DeleteUserHard(service.DB, req.UserID)
	} else {
		err = userRepo.DeleteUserSoft(service.DB, req.UserID)
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
	user, error := userRepo.FindUserByID(service.DB, req.UserID)

	switch {
	case error == nil:
		if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)) == nil {

			err := userRepo.UpdateUserPassword(service.DB, req.UserID, models.HashAndSalt([]byte(req.NewPassword)))
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
		res.Message = fmt.Sprintf("user id not found: %s", req.UserID)

	default:
		res.Status = "error"
		res.Message = error.Error()
	}

	return error
}

func (service *UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest, res *pb.UserResponse) error {
	user, error := userRepo.FindUserByID(service.DB, req.UserID)
	if error != nil {
		res.Status = "fail"
		res.Message = error.Error()
	} else if error == nil {
		res.Status = "success"
		res.Data = &pb.UserData{
			User: &pb.User{
				UserID: user.ID,
				First:  user.First,
				Last:   user.Last,
				Email:  user.Email,
			},
		}
	}

	return error
}

func (service *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest, res *pb.UserResponse) error {
	user, error := userRepo.FindUserByID(service.DB, req.UserID)
	switch {
	case error == nil:
		user.Email = req.Email
		user.First = req.First
		user.Last = req.Last

		user, error = userRepo.UpdateUserInfo(service.DB, user)
		if error != nil {
			res.Status = "error"
			res.Message = error.Error()
		} else {
			res.Status = "success"
			res.Data = &pb.UserData{
				User: &pb.User{
					UserID: user.ID,
					First:  user.First,
					Last:   user.Last,
					Email:  user.Email,
				},
			}
		}

	case strings.Contains(error.Error(), "no rows in result set"):
		res.Status = "fail"
		res.Message = "user does not exist by that id"

	default:
		res.Status = "error"
		res.Message = error.Error()
	}

	return error
}
