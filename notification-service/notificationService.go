package main

import (
	"context"
	"database/sql"

	notifications "github.com/asciiu/gomo/notification-service/proto"
)

type NotificationService struct {
	DB *sql.DB
}

// GetUserNotifications returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *NotificationService) GetUserNotifications(ctx context.Context, req *notifications.GetUserNotifcationsRequest, res *notifications.NotificationListResponse) error {

	//keys, error := repo.FindKeysByUserID(service.DB, req)

	//switch {
	//case error == nil:
	//	res.Status = "success"
	//	res.Data = &kp.UserKeysData{
	//		Keys: keys,
	//	}
	//default:
	//	res.Status = "error"
	//	res.Message = error.Error()
	//}
	return nil
}
