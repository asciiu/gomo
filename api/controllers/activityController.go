package controllers

import (
	"net/http"
	"strconv"

	protoActivity "github.com/asciiu/gomo/activity-bulletin/proto"
	constResponse "github.com/asciiu/gomo/common/constants/response"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

// A ResponseActivityPageSuccess will always contain a status of "successful".
// swagger:model ResponseActivityPageSuccess
type ResponseActivityPageSuccess struct {
	Status string                          `json:"status"`
	Data   *protoActivity.UserActivityPage `json:"data"`
}

// A ResponseActivitySuccess will always contain a status of "successful".
// swagger:model ResponseActivitySuccess
type ResponseActivitySuccess struct {
	Status string                      `json:"status"`
	Data   *protoActivity.ActivityData `json:"data"`
}

// This struct is used in the generated swagger docs,
// and it is not used anywhere.
// swagger:parameters searchActivity
type SearchActivity struct {
	// Required: false
	// In: query
	ObjectType string `json:"objectType"`
	// Required: false
	// In: query
	ObjectID string `json:"objectID"`
	// Required: false
	// In: query
	Page string `json:"page"`
	// Required: false
	// In: query
	PageSize string `json:"pageSize"`
}

type ActivityController struct {
	Bulletin protoActivity.ActivityBulletinClient
}

func NewActivityController(service micro.Service) *ActivityController {
	controller := ActivityController{
		Bulletin: protoActivity.NewActivityBulletinClient("bulletin", service.Client()),
	}

	return &controller
}

// swagger:route GET /activity activity searchActivity
//
// get protoActivity (protected)
//
// Returns a list of activity.
//
// responses:
//  200: ResponseActivityPageSuccess "data" will contain array of protoActivity with "status": "success"
func (controller *ActivityController) HandleListActivity(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["jti"].(string)

	objectID := c.QueryParam("objectID")
	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("pageSize")

	// defaults for page and page size here
	// ignore the errors and assume the values are int
	page, _ := strconv.ParseUint(pageStr, 10, 32)
	pageSize, _ := strconv.ParseUint(pageSizeStr, 10, 32)
	if pageSize == 0 {
		pageSize = 20
	}

	req := protoActivity.ActivityRequest{
		UserID:   userID,
		ObjectID: objectID,
		Page:     uint32(page),
		PageSize: uint32(pageSize),
	}

	r, _ := controller.Bulletin.FindUserActivity(context.Background(), &req)
	if r.Status != constResponse.Success {
		response := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		if r.Status == constResponse.Fail {
			return c.JSON(http.StatusBadRequest, response)
		}
		if r.Status == constResponse.Error {
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	// in case activity is null do this
	if r.Data.Activity == nil {
		r.Data.Activity = make([]*protoActivity.Activity, 0)
	}

	response := &ResponseActivityPageSuccess{
		Status: constResponse.Success,
		Data:   r.Data,
	}

	return c.JSON(http.StatusOK, response)
}

// swagger:parameters UpdateActivity
type UpdateActivity struct {
	// Optional e.g. 2006-01-02T15:04:05Z
	// In: query
	ClickedAt string `json:"clickedAt"`
	// Optional e.g. 2006-01-02T15:04:05Z
	// In: query
	SeenAt string `json:"seenAt"`
}

// swagger:route PUT /activity/:activityID activity UpdateActivity
//
// update activity (protected)
//
// Update activity clickedAt or seenAt. Timestamps must be UTC.
//
// responses:
//  200: responseActivitySuccess "data" will contain array of protoActivity with "status": "success"
func (controller *ActivityController) HandleUpdateActivity(c echo.Context) error {
	//token := c.Get("user").(*jwt.Token)
	//claims := token.Claims.(jwt.MapClaims)
	//userID := claims["jti"].(string)
	activityID := c.Param("activityID")

	// read strategy from post body
	updateActivity := new(UpdateActivity)
	err := c.Bind(&updateActivity)
	if err != nil {
		return fail(c, err.Error())
	}

	req := protoActivity.UpdateActivityRequest{
		ActivityID: activityID,
		ClickedAt:  updateActivity.ClickedAt,
		SeenAt:     updateActivity.SeenAt,
	}

	r, _ := controller.Bulletin.UpdateActivity(context.Background(), &req)
	if r.Status != constResponse.Success {
		response := &ResponseError{
			Status:  r.Status,
			Message: r.Message,
		}

		if r.Status == constResponse.Nonentity {
			return c.JSON(http.StatusNotFound, response)
		}
		if r.Status == constResponse.Fail {
			return c.JSON(http.StatusBadRequest, response)
		}
		if r.Status == constResponse.Error {
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	response := &ResponseActivitySuccess{
		Status: constResponse.Success,
		Data:   r.Data,
	}

	return c.JSON(http.StatusOK, response)
}
