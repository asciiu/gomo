package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	repoActivity "github.com/asciiu/gomo/activity-bulletin/db/sql"
	protoActivity "github.com/asciiu/gomo/activity-bulletin/proto"
	"github.com/asciiu/gomo/common/db"
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

func setupService() (*Bulletin, *user.User) {
	dbUrl := "postgres://postgres@localhost:5432/gomo_test?&sslmode=disable"
	db, _ := db.NewDB(dbUrl)

	hs := Bulletin{
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
func TestActivity(t *testing.T) {
	service, user := setupService()

	defer service.db.Close()

	note := protoActivity.Activity{
		UserID:      user.ID,
		Type:        "order",
		ObjectID:    "bf24b117-1c0f-4c4f-82bc-7586c99b8d40",
		Title:       "Test",
		Subtitle:    "test",
		Description: "this is a test",
		Timestamp:   string(pq.FormatTimestamp(time.Now().UTC())),
	}
	repoActivity.InsertActivity(service.db, &note)

	req := protoActivity.ActivityRequest{
		UserID:   user.ID,
		ObjectID: "bf24b117-1c0f-4c4f-82bc-7586c99b8d40",
		Page:     0,
		PageSize: 10,
	}
	res := protoActivity.ActivityPagedResponse{}
	service.FindUserActivity(context.Background(), &req, &res)

	assert.Equal(t, "success", res.Status, fmt.Sprintf("%s", res.Message))
	assert.Equal(t, uint32(1), res.Data.Total, "must be 1 history entry")

	repoUser.DeleteUserHard(service.db, user.ID)
}

func TestRecentActivity(t *testing.T) {
	service, user := setupService()

	defer service.db.Close()

	note1 := protoActivity.Activity{
		UserID:      user.ID,
		Type:        "order",
		ObjectID:    "bf24b117-1c0f-4c4f-82bc-7586c99b8d40",
		Title:       "Test",
		Subtitle:    "test",
		Description: "this is a test",
		Timestamp:   "2018-08-18 05:34:27.462218561Z",
	}
	repoActivity.InsertActivity(service.db, &note1)
	note2 := protoActivity.Activity{
		UserID:      user.ID,
		Type:        "order",
		ObjectID:    "bf24b117-1c0f-4c4f-82bc-7586c99b8d40",
		Title:       "Test",
		Subtitle:    "test",
		Description: "this is a test",
		Timestamp:   "2018-08-18 05:44:00.000000000Z",
	}
	repoActivity.InsertActivity(service.db, &note2)
	note3 := protoActivity.Activity{
		UserID:      user.ID,
		Type:        "order",
		ObjectID:    "bf24b117-1c0f-4c4f-82bc-7586c99b8d40",
		Title:       "Test",
		Subtitle:    "test",
		Description: "this is a test",
		Timestamp:   "2018-08-18 05:54:00.000000000Z",
	}
	repoActivity.InsertActivity(service.db, &note3)

	req := protoActivity.RecentActivityRequest{
		ObjectID: "bf24b117-1c0f-4c4f-82bc-7586c99b8d40",
		Count:    2,
	}
	res := protoActivity.ActivityListResponse{}
	service.FindMostRecentActivity(context.Background(), &req, &res)

	assert.Equal(t, "success", res.Status, fmt.Sprintf("%s", res.Message))
	assert.Equal(t, 2, len(res.Data.Activity), "must be 2 history")
	assert.Equal(t, "2018-08-18T05:54:00Z", res.Data.Activity[0].Timestamp, "first timestamp should be most recent")
	assert.Equal(t, "2018-08-18T05:44:00Z", res.Data.Activity[1].Timestamp, "second timestamp incorrect")

	repoUser.DeleteUserHard(service.db, user.ID)
}

func TestActivityCount(t *testing.T) {
	service, user := setupService()

	defer service.db.Close()

	note1 := protoActivity.Activity{
		UserID:      user.ID,
		Type:        "order",
		ObjectID:    "bf24b117-1c0f-4c4f-82bc-7586c99b8d40",
		Title:       "Test",
		Subtitle:    "test",
		Description: "this is a test",
		Timestamp:   "2018-08-18 05:34:27.462218561Z",
	}
	repoActivity.InsertActivity(service.db, &note1)
	note2 := protoActivity.Activity{
		UserID:      user.ID,
		Type:        "order",
		ObjectID:    "bf24b117-1c0f-4c4f-82bc-7586c99b8d40",
		Title:       "Test",
		Subtitle:    "test",
		Description: "this is a test",
		Timestamp:   "2018-08-18 05:44:00.000000000Z",
	}
	repoActivity.InsertActivity(service.db, &note2)

	req := protoActivity.ActivityCountRequest{
		ObjectID: "bf24b117-1c0f-4c4f-82bc-7586c99b8d40",
	}
	res := protoActivity.ActivityCountResponse{}
	service.FindActivityCount(context.Background(), &req, &res)

	assert.Equal(t, "success", res.Status, fmt.Sprintf("%s", res.Message))
	assert.Equal(t, uint32(2), res.Data.Count, "must be 2 history")

	repoUser.DeleteUserHard(service.db, user.ID)
}

func TestUpdateActivity(t *testing.T) {
	service, user := setupService()

	defer service.db.Close()

	note1 := protoActivity.Activity{
		UserID:      user.ID,
		Type:        "order",
		ObjectID:    "bf24b117-1c0f-4c4f-82bc-7586c99b8d40",
		Title:       "Test",
		Subtitle:    "test",
		Description: "this is a test",
		Timestamp:   "2018-08-18 05:34:27.462218561Z",
	}
	history, _ := repoActivity.InsertActivity(service.db, &note1)

	req := protoActivity.UpdateActivityRequest{
		ActivityID: history.ActivityID,
		SeenAt:     "2018-08-18T05:34:00Z",
		ClickedAt:  "2018-08-18T05:54:00Z",
	}
	res := protoActivity.ActivityResponse{}
	service.UpdateActivity(context.Background(), &req, &res)

	assert.Equal(t, "success", res.Status, fmt.Sprintf("%s", res.Message))
	assert.Equal(t, "2018-08-18T05:34:00Z", res.Data.Activity.SeenAt, "clicked did not match")
	assert.Equal(t, "2018-08-18T05:54:00Z", res.Data.Activity.ClickedAt, "seen did not match")

	repoUser.DeleteUserHard(service.db, user.ID)
}
