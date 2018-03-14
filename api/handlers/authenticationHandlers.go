package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func createJwtToken() (string, error) {
	claims := JwtClaims{
		"jack",
		jwt.StandardClaims{
			Id:        "userId",
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// Generate encoded token and send it as response.
	// TODO read secret from env var
	token, err := rawToken.SignedString([]byte("cuddlegang"))
	if err != nil {
		return "", err
	}

	return token, nil
}

func Login(c echo.Context) error {
	loginRequest := LoginRequest{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&loginRequest)
	if err != nil {
		log.Printf("failed reading the request %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	if loginRequest.Username == "jon" && loginRequest.Password == "shhh!" {
		// Generate encoded token and send it as response.
		// TODO read secret from env var
		token, err := createJwtToken()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": token,
		})
	}

	return echo.ErrUnauthorized
}
