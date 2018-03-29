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
	gsql "github.com/asciiu/gomo/common/db/sql"

	apiModels "github.com/asciiu/gomo/api/models"
	models "github.com/asciiu/gomo/common/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

const refreshDuration = 720 * time.Hour
const jwtDuration = 5 * time.Minute

type AuthController struct {
	DB *sql.DB
}

type JwtClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type SignupRequest struct {
	First    string `json:"first"`
	Last     string `json:"last"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResponseSuccess struct {
	Status string    `json:"status"`
	Data   *UserData `json:"data"`
}

type UserData struct {
	User *models.UserInfo `json:"user"`
}

type ResponseError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func createJwtToken(userId string, duration time.Duration) (string, error) {
	claims := jwt.StandardClaims{
		Id:        userId,
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
	accessToken, err := createJwtToken(refreshToken.UserId, jwtDuration)
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

// Handles a login request.
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
			accessToken, err := createJwtToken(user.Id, 5*time.Minute)
			if err != nil {
				response := &ResponseError{
					Status:  "error",
					Message: err.Error(),
				}
				return c.JSON(http.StatusInternalServerError, response)
			}

			// issue a refresh token if remember is true
			if loginRequest.Remember {
				refreshToken := apiModels.NewRefreshToken(user.Id)
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

			response := &ResponseSuccess{
				Status: "success",
				Data:   &UserData{user.Info()},
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

// Handles logout cleanup
func (controller *AuthController) HandleLogout(c echo.Context) error {
	selectAuth := c.Request().Header.Get("Refresh")
	if selectAuth != "" {
		sa := strings.Split(selectAuth, ":")

		if len(sa) != 2 {
			response := &ResponseError{
				Status:  "fail",
				Message: "refresh token invalid",
			}
			return c.JSON(http.StatusOK, response)
		}

		asql.DeleteRefreshTokenBySelector(controller.DB, sa[0])
	}

	response := &ResponseSuccess{
		Status: "success",
	}
	return c.JSON(http.StatusOK, response)
}

// Handles a new signup request
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

	user := models.NewUser(signupRequest.First, signupRequest.Last, signupRequest.Email, signupRequest.Password)
	_, error := gsql.InsertUser(controller.DB, user)
	if error != nil {
		response := &ResponseError{
			Status:  "fail",
			Message: error.Error(),
		}

		return c.JSON(http.StatusConflict, response)
	}

	response := &ResponseSuccess{
		Status: "success",
		Data:   &UserData{user.Info()},
	}

	return c.JSON(http.StatusOK, response)
}
