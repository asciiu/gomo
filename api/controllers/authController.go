package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	asql "github.com/asciiu/gomo/api/db/sql"
	balances "github.com/asciiu/gomo/balance-service/proto/balance"
	devices "github.com/asciiu/gomo/device-service/proto/device"
	keys "github.com/asciiu/gomo/key-service/proto/key"
	gsql "github.com/asciiu/gomo/user-service/db/sql"
	users "github.com/asciiu/gomo/user-service/proto/user"
	micro "github.com/micro/go-micro"
	k8s "github.com/micro/kubernetes/go/micro"

	apiModels "github.com/asciiu/gomo/api/models"
	models "github.com/asciiu/gomo/user-service/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

//const refreshDuration = 720 * time.Hour
//const jwtDuration = 1440 * time.Minute
const refreshDuration = 30 * time.Minute
const jwtDuration = 30 * time.Minute

type AuthController struct {
	DB       *sql.DB
	Users    users.UserServiceClient
	Balances balances.BalanceServiceClient
	Keys     keys.KeyServiceClient
	Devices  devices.DeviceServiceClient
}

type JwtClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// swagger:parameters login
type LoginRequest struct {
	// Required. Backend code does not check email atm.
	// in: body
	Email string `json:"email"`
	// Required. Backend code does not have any password requirements atm.
	// in: body
	Password string `json:"password"`
	// Optional. Return refresh token in response
	// in: body
	Remember bool `json:"remember"`
}

// swagger:parameters signup
type SignupRequest struct {
	// Optional.
	// in: body
	First string `json:"first"`
	// Optional.
	// in: body
	Last string `json:"last"`
	// Required. Must be unique! We need to validate these coming in.
	// in: body
	Email string `json:"email"`
	// Required. We need password requirements.
	// in: body
	Password string `json:"password"`
}

// A ResponseSuccess will always contain a status of "successful".
// swagger:model responseSuccess
type ResponseSuccess struct {
	Status string    `json:"status"`
	Data   *UserData `json:"data"`
}

type Device struct {
	DeviceID         string `json:"deviceID"`
	DeviceType       string `json:"deviceType"`
	ExternalDeviceID string `json:"externalDeviceID"`
	DeviceToken      string `json:"deviceToken"`
}

type UserData struct {
	User    *models.UserInfo `json:"user"`
	Devices []*Device        `json:"devices"`
}

// A ResponseSuccess will always contain a status of "successful".
// This response may or may not include data encapsulating the user information.
// swagger:model responseError
type ResponseError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewAuthController(db *sql.DB) *AuthController {

	service := k8s.NewService(micro.Name("user.client"))
	service.Init()

	controller := AuthController{
		DB:       db,
		Users:    users.NewUserServiceClient("users", service.Client()),
		Balances: balances.NewBalanceServiceClient("balances", service.Client()),
		Keys:     keys.NewKeyServiceClient("keys", service.Client()),
		Devices:  devices.NewDeviceServiceClient("devices", service.Client()),
	}

	return &controller
}

func createJwtToken(userID string, duration time.Duration) (string, error) {
	claims := jwt.StandardClaims{
		Id:        userID,
		ExpiresAt: time.Now().Add(duration).Unix(),
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// Generate encoded token and send it as response.
	token, err := rawToken.SignedString([]byte(os.Getenv("GOMO_JWT")))
	if err != nil {
		return "", err
	}

	return token, nil
}

// Renews the refresh token and the access token in the reponse headers.
func renewTokens(c echo.Context, refreshToken *apiModels.RefreshToken) {
	// renew access
	accessToken, err := createJwtToken(refreshToken.UserID, jwtDuration)
	if err != nil {
		log.Fatal(err)
	}

	// renew the refresh token
	expiresOn := time.Now().Add(refreshDuration)
	selectAuth := refreshToken.Renew(expiresOn)

	c.Response().Header().Set("set-authorization", accessToken)
	c.Response().Header().Set("set-refresh", selectAuth)
}

// A custom middleware function to check the refresh token on each request.
func (controller *AuthController) RefreshAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		if auth == "" {
			return next(c)
		}

		tokenString := strings.Split(auth, " ")[1]

		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(os.Getenv("GOMO_JWT")), nil
		})

		if err != nil {

			selectAuth := c.Request().Header.Get("Refresh")
			if selectAuth != "" {
				sa := strings.Split(selectAuth, ":")

				if len(sa) != 2 {
					return next(c)
				}

				selector := sa[0]
				authenticator := sa[1]

				refreshToken, err := asql.FindRefreshToken(controller.DB, selector)
				if err != nil {
					return next(c)
				}

				if refreshToken.Valid(authenticator) {
					// renew access
					renewTokens(c, refreshToken)
					_, err3 := asql.UpdateRefreshToken(controller.DB, refreshToken)

					if err3 != nil {
						log.Fatal(err3)
					}
				}

				if refreshToken.ExpiresOn.Before(time.Now()) {
					asql.DeleteRefreshToken(controller.DB, refreshToken)
				}
			}
		}

		return next(c)
	}
}

// swagger:route POST /login authentication login
//
// user login (open)
//
// The login endpoint returns an authorization token in the "set-authorization" response header.
// You may also receive an optional refresh token if you include 'remember': true in your post request.
// Both tokens will be returned in the reponse headers as "set-refresh" and "set-authorization".
//
// responses:
//  200: responseSuccess "data" will be non null with "status": "success"
//  400: responseError email and password are not found in request with "status": "fail"
//  401: responseError unauthorized user because of incorrect password with "status": "fail"
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *AuthController) HandleLogin(c echo.Context) error {
	loginRequest := LoginRequest{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&loginRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "fail",
			Message: "malformed json request for 'email' and 'password'",
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	// lookup user by email
	user, err := gsql.FindUser(controller.DB, loginRequest.Email)
	switch {
	case err == sql.ErrNoRows:
		response := &ResponseError{
			Status:  "fail",
			Message: "password/login incorrect",
		}
		// no user by this email send unauthorized response
		return c.JSON(http.StatusUnauthorized, response)

	case err != nil:
		log.Fatal(err)
		response := &ResponseError{
			Status:  "error",
			Message: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, response)

	default:
		if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password)) == nil {

			accessToken, err := createJwtToken(user.ID, jwtDuration)
			if err != nil {
				response := &ResponseError{
					Status:  "error",
					Message: err.Error(),
				}
				return c.JSON(http.StatusInternalServerError, response)
			}

			// issue a refresh token if remember is true
			if loginRequest.Remember {
				refreshToken := apiModels.NewRefreshToken(user.ID)
				renewTokens(c, refreshToken)

				_, err3 := asql.InsertRefreshToken(controller.DB, refreshToken)

				if err3 != nil {
					response := &ResponseError{
						Status:  "error",
						Message: err.Error(),
					}
					return c.JSON(http.StatusInternalServerError, response)
				}
			} else {
				c.Response().Header().Set("set-authorization", accessToken)
			}

			// TODO refactor with device controller implementtion
			// get user devices here
			getRequest := devices.GetUserDevicesRequest{
				UserID: user.ID,
			}

			r, _ := controller.Devices.GetUserDevices(context.Background(), &getRequest)
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

			devices := make([]*Device, 0)
			for _, d := range r.Data.Devices {
				// api removes the secret
				device := Device{
					DeviceID:         d.DeviceID,
					DeviceType:       d.DeviceType,
					ExternalDeviceID: d.ExternalDeviceID,
					DeviceToken:      d.DeviceToken,
				}
				devices = append(devices, &device)
			}

			response := &ResponseSuccess{
				Status: "success",
				Data: &UserData{
					User:    user.Info(),
					Devices: devices,
				},
			}

			return c.JSON(http.StatusOK, response)
		}
	}

	response := &ResponseError{
		Status:  "fail",
		Message: "password/login incorrect",
	}
	return c.JSON(http.StatusUnauthorized, response)
}

// swagger:route GET /logout authentication logout
//
// logout user (protected)
//
// If a valid refresh token is found in the request headers, the server
// will attempt to remove the refresh token from the database.
//
//	Responses:
//	  200: responseSuccess data will be null with status "success"
//	  400: responseError you either sent in no refresh token or the refresh token in the header is not valid with status: "fail"
func (controller *AuthController) HandleLogout(c echo.Context) error {
	selectAuth := c.Request().Header.Get("Refresh")
	if selectAuth != "" {
		sa := strings.Split(selectAuth, ":")

		if len(sa) != 2 {
			response := &ResponseError{
				Status:  "fail",
				Message: "refresh token invalid",
			}
			return c.JSON(http.StatusBadRequest, response)
		}

		asql.DeleteRefreshTokenBySelector(controller.DB, sa[0])
	}

	response := &ResponseSuccess{
		Status: "success",
	}
	return c.JSON(http.StatusOK, response)
}

// swagger:route POST /signup authentication signup
//
// user registration (open)
//
// Registers a new user. Expects email to be unique. Duplicate email will result
// in a bad request.
//
// responses:
//  200: responseSuccess "data" will be non null with "status": "success"
//  400: responseError message should relay information with regard to bad request with "status": "fail"
//  410: responseError the user-service is not reachable. The user-service is a microservice that runs independantly from the api. When we take it offline you will receive this error.
//  500: responseError the message will state what the internal server error was with "status": "error"
func (controller *AuthController) HandleSignup(c echo.Context) error {
	signupRequest := SignupRequest{}

	err := json.NewDecoder(c.Request().Body).Decode(&signupRequest)
	if err != nil {
		response := &ResponseError{
			Status:  "fail",
			Message: err.Error(),
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	if signupRequest.Email == "" || signupRequest.Password == "" {
		response := &ResponseError{
			Status:  "fail",
			Message: "email and password are required",
		}

		return c.JSON(http.StatusBadRequest, response)
	}

	createRequest := users.CreateUserRequest{
		First:    signupRequest.First,
		Last:     signupRequest.Last,
		Email:    signupRequest.Email,
		Password: signupRequest.Password,
	}

	r, e := controller.Users.CreateUser(context.Background(), &createRequest)
	fmt.Printf("error was %+v\n", e)
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

	response := &ResponseSuccess{
		Status: "success",
		Data: &UserData{
			User: &models.UserInfo{
				UserID: r.Data.User.UserID,
				First:  r.Data.User.First,
				Last:   r.Data.User.Last,
				Email:  r.Data.User.Email,
			},
		},
	}

	return c.JSON(http.StatusOK, response)
}
