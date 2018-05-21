package controllers

import (
	"net/http"
	"strconv"

	notifications "github.com/asciiu/gomo/notification-service/proto"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

// A ResponseNotificationSuccess will always contain a status of "successful".
// swagger:model responseNotificationSuccess
type ResponseNotificationSuccess struct {
	Status string                               `json:"status"`
	Data   *notifications.UserNotificationsPage `json:"data"`
}

// This struct is used in the generated swagger docs,
// and it is not used anywhere.
// swagger:parameters searchNotifications
type SearchType struct {
	// Required: false
	// In: query
	Type string `json:"type"`
	// Required: false
	// In: query
	Page string `json:"page"`
	// Required: false
	// In: query
	PageSize string `json:"pageSize"`
}

type NotificationController struct {
	Notifications notifications.NotificationServiceClient
}

func NewNotificationController() *NotificationController {
	service := micro.NewService(micro.Name("notification.client"))
	service.Init()

	controller := NotificationController{
		Notifications: notifications.NewNotificationServiceClient("go.srv.notification-service", service.Client()),
	}

	return &controller
}

// swagger:route GET /notifications notifications searchNotifications
//
// get notifications (protected)
//
// Returns a list of notifications.
//
// responses:
//  200: responseNotificationSuccess "data" will contain array of notifications with "status": "success"
func (controller *NotificationController) HandleListNotifications(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)

	noteType := c.QueryParam("type")
	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("pageSize")

	// defaults for page and page size here
	// ignore the errors and assume the values are int
	page, _ := strconv.ParseUint(pageStr, 10, 32)
	pageSize, _ := strconv.ParseUint(pageSizeStr, 10, 32)
	if pageSize == 0 {
		pageSize = 20
	}

	req := notifications.GetNotifcationsByType{
		UserID:           userID,
		NotificationType: noteType,
		Page:             uint32(page),
		PageSize:         uint32(pageSize),
	}

	r, _ := controller.Notifications.GetUserNotificationsByType(context.Background(), &req)
	if r.Status != "success" {
		response := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		if r.Status == "fail" {
			return c.JSON(http.StatusBadRequest, response)
		}
		if r.Status == "error" {
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	response := &ResponseNotificationSuccess{
		Status: "success",
		Data:   r.Data,
	}

	return c.JSON(http.StatusOK, response)
}
