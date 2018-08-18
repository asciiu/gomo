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
	user, err := repoUser.InsertUser(db, user)
	if err != nil {
		fmt.Println(err)
	}

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

func TestRecentHistory(t *testing.T) {
	service, user := setupService()

	defer service.db.Close()

	note1 := protoHistory.History{
		UserID:      user.ID,
		Type:        "order",
		ObjectID:    "bf24b117-1c0f-4c4f-82bc-7586c99b8d40",
		Title:       "Test",
		Subtitle:    "test",
		Description: "this is a test",
		Timestamp:   "2018-08-18 05:34:27.462218561Z",
	}
	repoHistory.InsertHistory(service.db, &note1)
	note2 := protoHistory.History{
		UserID:      user.ID,
		Type:        "order",
		ObjectID:    "bf24b117-1c0f-4c4f-82bc-7586c99b8d40",
		Title:       "Test",
		Subtitle:    "test",
		Description: "this is a test",
		Timestamp:   "2018-08-18 05:44:00.000000000Z",
	}
	repoHistory.InsertHistory(service.db, &note2)
	note3 := protoHistory.History{
		UserID:      user.ID,
		Type:        "order",
		ObjectID:    "bf24b117-1c0f-4c4f-82bc-7586c99b8d40",
		Title:       "Test",
		Subtitle:    "test",
		Description: "this is a test",
		Timestamp:   "2018-08-18 05:54:00.000000000Z",
	}
	repoHistory.InsertHistory(service.db, &note3)

	req := protoHistory.RecentHistoryRequest{
		ObjectID: "bf24b117-1c0f-4c4f-82bc-7586c99b8d40",
		Count:    2,
	}
	res := protoHistory.HistoryListResponse{}
	service.FindMostRecentHistory(context.Background(), &req, &res)

	assert.Equal(t, "success", res.Status, fmt.Sprintf("%s", res.Message))
	assert.Equal(t, 2, len(res.Data.History), "must be 2 history")
	assert.Equal(t, "2018-08-18T05:54:00Z", res.Data.History[0].Timestamp, "first timestamp should be most recent")
	assert.Equal(t, "2018-08-18T05:44:00Z", res.Data.History[1].Timestamp, "second timestamp incorrect")

	repoUser.DeleteUserHard(service.db, user.ID)
}
