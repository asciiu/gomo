package sql_test

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/asciiu/gomo/common/db"
	repoNote "github.com/asciiu/gomo/notification-service/db/sql"
	protoNote "github.com/asciiu/gomo/notification-service/proto"
	repoUser "github.com/asciiu/gomo/user-service/db/sql"
	user "github.com/asciiu/gomo/user-service/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func setupService() (*sql.DB, *user.User) {
	dbUrl := "postgres://postgres@localhost:5432/gomo_test?&sslmode=disable"
	db, _ := db.NewDB(dbUrl)

	user := user.NewUser("first", "last", "test@email", "hash")
	_, err := repoUser.InsertUser(db, user)
	checkErr(err)

	return db, user
}

func TestInsertDevice(t *testing.T) {
	db, user := setupService()

	n := protoNote.Notification{
		NotificationID:   uuid.New().String(),
		NotificationType: "test",
		UserID:           user.ID,
		ObjectID:         uuid.New().String(),
		Title:            "Test",
		Subtitle:         "123",
		Description:      "this is a test!",
		Timestamp:        string(pq.FormatTimestamp(time.Now().UTC())),
	}

	note, err := repoNote.InsertNotification(db, &n)
	assert.Equal(t, nil, err, "error for insert notification should be nil")
	assert.Equal(t, "Test", note.Title, "titles did not match")

	repoUser.DeleteUserHard(db, user.ID)
}
