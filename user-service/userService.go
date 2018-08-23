package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	constRes "github.com/asciiu/gomo/common/constants/response"
	repoUser "github.com/asciiu/gomo/user-service/db/sql"
	"github.com/asciiu/gomo/user-service/models"
	protoUser "github.com/asciiu/gomo/user-service/proto/user"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB *sql.DB
}

// CreateUser returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object always.
func (service *UserService) CreateUser(ctx context.Context, req *protoUser.CreateUserRequest, res *protoUser.UserResponse) error {
	user := models.NewUser(req.First, req.Last, req.Email, req.Password)
	_, error := repoUser.InsertUser(service.DB, user)

	switch {
	case error == nil:
		res.Status = constRes.Success
		res.Data = &protoUser.UserData{
			User: &protoUser.User{
				UserID: user.ID,
				First:  user.First,
				Last:   user.Last,
				Email:  user.Email,
			},
		}
		return nil

	case strings.Contains(error.Error(), "violates unique constraint \"users_email_key\""):
		res.Status = constRes.Fail
		res.Message = "email already exists"
		return nil

	default:
		res.Status = constRes.Error
		res.Message = error.Error()
		return nil
	}
}

// DeleteUser returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object always.
func (service *UserService) DeleteUser(ctx context.Context, req *protoUser.DeleteUserRequest, res *protoUser.Response) error {
	var err error
	if req.Hard {
		err = repoUser.DeleteUserHard(service.DB, req.UserID)
	} else {
		err = repoUser.DeleteUserSoft(service.DB, req.UserID)
	}

	if err == nil {
		res.Status = constRes.Success
	} else {
		res.Status = constRes.Error
		res.Message = err.Error()
	}
	return nil
}

// ChangePassword returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object always.
// Changes the user's password. Password is updated when the request's
// old password matches the current user's password hash.
func (service *UserService) ChangePassword(ctx context.Context, req *protoUser.ChangePasswordRequest, res *protoUser.Response) error {
	user, error := repoUser.FindUserByID(service.DB, req.UserID)

	switch {
	case error == nil:
		if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)) == nil {

			err := repoUser.UpdateUserPassword(service.DB, req.UserID, models.HashAndSalt([]byte(req.NewPassword)))
			if err != nil {
				res.Status = constRes.Error
				res.Message = err.Error()
			} else {
				res.Status = constRes.Success
			}

		} else {
			res.Status = constRes.Fail
			res.Message = "current password mismatch"
		}

	case strings.Contains(error.Error(), "no rows in result set"):
		res.Status = constRes.Fail
		res.Message = fmt.Sprintf("user id not found: %s", req.UserID)

	default:
		res.Status = constRes.Error
		res.Message = error.Error()
	}

	return nil
}

// GetUserInfo returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object always.
func (service *UserService) GetUserInfo(ctx context.Context, req *protoUser.GetUserInfoRequest, res *protoUser.UserResponse) error {
	user, error := repoUser.FindUserByID(service.DB, req.UserID)
	if error != nil {
		res.Status = constRes.Error
		res.Message = error.Error()
	} else if error == nil {
		res.Status = constRes.Success
		res.Data = &protoUser.UserData{
			User: &protoUser.User{
				UserID: user.ID,
				First:  user.First,
				Last:   user.Last,
				Email:  user.Email,
			},
		}
	}

	return nil
}

// UpdateUser returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object always.
func (service *UserService) UpdateUser(ctx context.Context, req *protoUser.UpdateUserRequest, res *protoUser.UserResponse) error {
	user, err := repoUser.FindUserByID(service.DB, req.UserID)
	switch {
	case err == nil:
		user.Email = req.Email
		user.First = req.First
		user.Last = req.Last

		user, err = repoUser.UpdateUserInfo(service.DB, user)
		if err != nil {
			res.Status = constRes.Error
			res.Message = err.Error()
		} else {
			res.Status = constRes.Success
			res.Data = &protoUser.UserData{
				User: &protoUser.User{
					UserID: user.ID,
					First:  user.First,
					Last:   user.Last,
					Email:  user.Email,
				},
			}
		}

	case err == sql.ErrNoRows:
		res.Status = constRes.Fail
		res.Message = "user does not exist by that id"

	default:
		res.Status = constRes.Error
		res.Message = err.Error()
	}

	return nil
}
