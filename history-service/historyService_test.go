package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/asciiu/gomo/common/db"
	repoHistory "github.com/asciiu/gomo/history-service/db/sql"
	protoHistory "github.com/asciiu/gomo/history-service/proto"
	repoUser "github.com/asciiu/gomo/user-service/db/sql"
	user "github.com/asciiu/gomo/user-service/models"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func setupService() (*HistoryService, *user.User) {
	dbUrl := "postgres://postgres@localhost:5432/gomo_test?&sslmode=disable"
	db, _ := db.NewDB(dbUrl)

	hs := HistoryService{
		db: db,
		//devices: protoDevice.NewDeviceServiceClient("devices", service.Client()),
	}

	user := user.NewUser("first", "last", "test@email", "hash")
	repoUser.InsertUser(db, user)

	return &hs, user
}

// You shouldn't be able to insert a plan with no orders. A new plan requires at least a single order.
func TestHistory(t *testing.T) {
	service, user := setupService()

	defer service.db.Close()

	note := protoHistory.History{
		UserID:      user.ID,
		Type:        "order",
		ObjectID:    "bf24b117-1c0f-4c4f-82bc-7586c99b8d40",
		Title:       "Test",
		Subtitle:    "test",
		Description: "this is a test",
		Timestamp:   string(pq.FormatTimestamp(time.Now().UTC())),
	}
	repoHistory.InsertHistory(service.db, &note)

	req := protoHistory.HistoryRequest{
		UserID:   user.ID,
		ObjectID: "bf24b117-1c0f-4c4f-82bc-7586c99b8d40",
		Page:     0,
		PageSize: 10,
	}
	res := protoHistory.HistoryPagedResponse{}
	service.FindUserHistory(context.Background(), &req, &res)

	assert.Equal(t, "success", res.Status, fmt.Sprintf("%s", res.Message))
	assert.Equal(t, uint32(1), res.Data.Total, "must be 1 history entry")

	repoUser.DeleteUserHard(service.db, user.ID)
}
