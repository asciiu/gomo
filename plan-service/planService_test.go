package main

import (
	"context"
	"log"
	"testing"

	"github.com/asciiu/gomo/common/db"
	protoKey "github.com/asciiu/gomo/key-service/proto/key"
	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	protoPlan "github.com/asciiu/gomo/plan-service/proto/plan"
	userRepo "github.com/asciiu/gomo/user-service/db/sql"
	user "github.com/asciiu/gomo/user-service/models"
	"github.com/micro/go-micro/client"
	"github.com/stretchr/testify/assert"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

type mockKeyService struct {
}

func (m *mockKeyService) GetKeys(ctx context.Context, req *protoKey.GetKeysRequest, opts ...client.CallOption) (*protoKey.KeyListResponse, error) {
	return &protoKey.KeyListResponse{}, nil
}
func (m *mockKeyService) AddKey(ctx context.Context, in *protoKey.KeyRequest, opts ...client.CallOption) (*protoKey.KeyResponse, error) {
	return &protoKey.KeyResponse{}, nil
}
func (m *mockKeyService) GetUserKey(ctx context.Context, in *protoKey.GetUserKeyRequest, opts ...client.CallOption) (*protoKey.KeyResponse, error) {
	return &protoKey.KeyResponse{}, nil
}
func (m *mockKeyService) GetUserKeys(ctx context.Context, in *protoKey.GetUserKeysRequest, opts ...client.CallOption) (*protoKey.KeyListResponse, error) {
	return &protoKey.KeyListResponse{}, nil
}
func (m *mockKeyService) RemoveKey(ctx context.Context, in *protoKey.RemoveKeyRequest, opts ...client.CallOption) (*protoKey.KeyResponse, error) {
	return &protoKey.KeyResponse{}, nil
}
func (m *mockKeyService) UpdateKeyDescription(ctx context.Context, in *protoKey.KeyRequest, opts ...client.CallOption) (*protoKey.KeyResponse, error) {
	return &protoKey.KeyResponse{}, nil
}

func MokieKeyService() protoKey.KeyServiceClient {
	return new(mockKeyService)
}

func setupService() (*PlanService, *user.User) {
	dbUrl := "postgres://postgres@localhost:5432/gomo_test?&sslmode=disable"
	db, _ := db.NewDB(dbUrl)

	planService := PlanService{
		DB:        db,
		KeyClient: MokieKeyService(),
	}

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := userRepo.InsertUser(db, user)
	checkErr(error)

	return &planService, user
}

func TestEmptyOrderPlan(t *testing.T) {
	service, user := setupService()

	defer service.DB.Close()

	orders := make([]*protoOrder.NewOrderRequest, 0)
	req := protoPlan.NewPlanRequest{
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "thing",
		Orders:          orders,
	}
	res := protoPlan.PlanResponse{}
	service.NewPlan(context.Background(), &req, &res)

	assert.Equal(t, res.Status, "fail", "return status of inserting an empty order plan should be fail")

	userRepo.DeleteUserHard(service.DB, user.ID)
}

//func TestInsertPlan(t *testing.T) {}
