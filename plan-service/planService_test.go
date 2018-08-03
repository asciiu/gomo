package main

import (
	"context"
	"log"
	"testing"

	"github.com/asciiu/gomo/common/db"
	repoKey "github.com/asciiu/gomo/key-service/db/sql"
	protoKey "github.com/asciiu/gomo/key-service/proto/key"
	protoOrder "github.com/asciiu/gomo/plan-service/proto/order"
	protoPlan "github.com/asciiu/gomo/plan-service/proto/plan"
	repoUser "github.com/asciiu/gomo/user-service/db/sql"
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
	keys := make([]*protoKey.Key, 0)
	keys = append(keys, &protoKey.Key{
		KeyID:    "examplekey",
		UserID:   "testuser",
		Exchange: "testex",
	})

	return &protoKey.KeyListResponse{
		Status: "success",
		Data: &protoKey.UserKeysData{
			Keys: keys,
		},
	}, nil
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

func setupService() (*PlanService, *user.User, *protoKey.Key) {
	dbUrl := "postgres://postgres@localhost:5432/gomo_test?&sslmode=disable"
	db, _ := db.NewDB(dbUrl)

	planService := PlanService{
		DB:        db,
		KeyClient: MokieKeyService(),
	}

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := repoUser.InsertUser(db, user)
	keyRequest := protoKey.KeyRequest{
		KeyID:    "92512e72-b7e8-49f4-bab5-271f4ba450d9",
		UserID:   user.ID,
		Exchange: "test",
	}

	key, error := repoKey.InsertKey(db, &keyRequest)
	checkErr(error)

	return &planService, user, key
}

// You should not be able to insert a plan with no orders. A new plan requires at least a single order.
func TestEmptyOrderPlan(t *testing.T) {
	service, user, _ := setupService()

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

	repoUser.DeleteUserHard(service.DB, user.ID)
}

// Successfully inserting a new plan
func TestSuccessfulOrderPlan(t *testing.T) {
	service, user, key := setupService()

	defer service.DB.Close()

	orders := make([]*protoOrder.NewOrderRequest, 0)
	order := protoOrder.NewOrderRequest{
		OrderID:         "4d671984-d7dd-4dce-a20f-23f25d6daf7f",
		KeyID:           key.KeyID,
		OrderType:       "paper",
		OrderTemplateID: "mokie",
		ParentOrderID:   "00000000-0000-0000-0000-000000000000",
		MarketName:      "ADA-BTC",
		Side:            "buy",
		ActiveCurrencyBalance: 100,
	}
	orders = append(orders, &order)
	req := protoPlan.NewPlanRequest{
		UserID:          user.ID,
		Status:          "active",
		CloseOnComplete: false,
		PlanTemplateID:  "thing",
		Orders:          orders,
	}
	res := protoPlan.PlanResponse{}
	service.NewPlan(context.Background(), &req, &res)

	assert.Equal(t, res.Status, "success", "return status of inserting plan should be success")

	repoUser.DeleteUserHard(service.DB, user.ID)
}

//func TestInsertPlan(t *testing.T) {}
