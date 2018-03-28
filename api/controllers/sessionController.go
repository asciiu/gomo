package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/antonlindstrom/pgstore"
	"github.com/labstack/echo"
)

type SessionController struct {
	DB *sql.DB
}

func (controller *SessionController) Session(c echo.Context) error {
	dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))

	store, err := pgstore.NewPGStore(dbUrl, []byte("secret-key"))
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer store.Close()

	// Run a background goroutine to clean up expired sessions from the database.
	defer store.StopCleanup(store.Cleanup(time.Minute * 5))

	// Get a session.
	session, err := store.Get(c.Request(), "session-key")
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Add a value.
	session.Values["foo"] = "bar"

	// Save.
	if err = session.Save(c.Request(), c.Response().Writer); err != nil {
		log.Fatalf("Error saving session: %v", err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"sessionKey": session.ID,
	})
}
