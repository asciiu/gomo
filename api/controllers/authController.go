package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	asql "github.com/asciiu/gomo/api/db/sql"
	gsql "github.com/asciiu/gomo/common/db/sql"

	apiModels "github.com/asciiu/gomo/api/models"
	models "github.com/asciiu/gomo/common/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

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
	Status string           `json:"status"`
	Data   *models.UserInfo `json:"data"`
}

type ResponseError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func createJwtToken(user *models.User, duration time.Duration) (string, error) {
	claims := jwt.StandardClaims{
		Id:        user.Id,
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

func (controller *AuthController) Login(c echo.Context) error {
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
			accessToken, err := createJwtToken(user, 5*time.Minute)
			if err != nil {
				response := &ResponseError{
					Status:  "error",
					Message: err.Error(),
				}
				return c.JSON(http.StatusInternalServerError, response)
			}

			// issue a refresh token if remember is true
			if loginRequest.Remember {
				expiresOn := time.Now().Add(720 * time.Hour)
				selectAuth := apiModels.NewSelectorAuth()
				token := apiModels.NewRefreshToken(user.Id, selectAuth, expiresOn)

				_, err3 := asql.InsertRefreshToken(controller.DB, token)

				if err3 != nil {
					response := &ResponseError{
						Status:  "error",
						Message: err.Error(),
					}
					return c.JSON(http.StatusInternalServerError, response)
				}

				return c.JSON(http.StatusOK, map[string]string{
					"access":  accessToken,
					"refresh": selectAuth,
				})

			} else {
				return c.JSON(http.StatusOK, map[string]string{
					"access": accessToken,
				})
			}
		}
	}

	response := &ResponseError{
		Status:  "fail",
		Message: "password/login incorrect",
	}
	return c.JSON(http.StatusUnauthorized, response)
}

func (controller *AuthController) Signup(c echo.Context) error {
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
		Data:   user.Info(),
	}

	return c.JSON(http.StatusOK, response)
}
