package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

type MainRoutes struct {
	DB *sql.DB
}

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

func (routes *MainRoutes) Login(c echo.Context) error {
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

func (routes *MainRoutes) Signup(c echo.Context) error {
	// panic on error
	u1 := uuid.Must(uuid.NewV4())

	stmt, err := routes.DB.Prepare("INSERT INTO users(id, first_name, last_name, email, password, salt) VALUES($1,$2,$3,$4,$5,$6)")
	if err != nil {
		log.Print("HERE")
		log.Fatal(err)
	}
	res, err := stmt.Exec(u1, "test", "name", "test@email", "password", "salt")
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("affected = %d\n", rowCnt)

	return c.String(http.StatusOK, "registered")
}
