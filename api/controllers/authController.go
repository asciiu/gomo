package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	gsql "github.com/asciiu/gomo/common/db/sql"
	"github.com/asciiu/gomo/common/models"
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
}

type SignupRequest struct {
	First    string `json:"first"`
	Last     string `json:"last"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func createJwtToken(user *models.User) (string, error) {
	claims := jwt.StandardClaims{
		Id:        user.Id,
		ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
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
		log.Printf("failed reading the request %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	user, err := gsql.FindUser(controller.DB, loginRequest.Email)
	switch {
	case err == sql.ErrNoRows:
		return echo.ErrUnauthorized
	case err != nil:
		log.Fatal(err)
		return echo.ErrUnauthorized
	default:
		if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password)) == nil {
			// Generate encoded token and send it as response.
			// TODO read secret from env var
			token, err := createJwtToken(user)
			if err != nil {
				return err
			}
			return c.JSON(http.StatusOK, map[string]string{
				"token": token,
			})
		}
	}

	return echo.ErrUnauthorized
}

type JSendUserRegisteredSuccess struct {
	Status string       `json:"status"`
	Data   *models.User `json:"data"`
}

func (controller *AuthController) Signup(c echo.Context) error {
	signupRequest := SignupRequest{}

	err := json.NewDecoder(c.Request().Body).Decode(&signupRequest)
	if err != nil {
		log.Printf("failed reading the signup request %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	user := models.NewUser(signupRequest.First, signupRequest.Last, signupRequest.Email, signupRequest.Password)
	_, error := gsql.InsertUser(controller.DB, user)
	if error != nil {
		log.Printf("failed reading the request %s", error)
	}

	response := &JSendUserRegisteredSuccess{
		Status: "success",
		Data:   user,
	}

	return c.JSON(http.StatusOK, response)
}
