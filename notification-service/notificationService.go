package main

import (
	"context"
	"database/sql"

	repo "github.com/asciiu/gomo/notification-service/db/sql"
	protoNotification "github.com/asciiu/gomo/notification-service/proto"
)

type NotificationService struct {
	DB *sql.DB
}

// GetUserNotifications returns error to conform to protobuf def, but the error will always be returned as nil.
// Can't return an error with a response object - response object is returned as nil when error is non nil.
// Therefore, return error in response object.
func (service *NotificationService) GetUserNotificationsByType(ctx context.Context, req *protoNotification.GetNotifcationsByType, res *protoNotification.NotificationPagedResponse) error {

	var pagedResult *protoNotification.UserNotificationsPage
	var err error
	if req.NotificationType == "" {
		pagedResult, err = repo.FindNotifications(service.DB, req)
	} else {
		pagedResult, err = repo.FindNotificationsByType(service.DB, req)
	}

	switch {
	case err == nil:
		res.Status = "success"
		res.Data = pagedResult
	default:
		res.Status = "error"
		res.Message = err.Error()
	}
	return nil
}
