package sql_test

import (
	"database/sql"
	"log"
	"testing"

	keyRepo "github.com/asciiu/gomo/apikey-service/db/sql"
	keyProto "github.com/asciiu/gomo/apikey-service/proto/apikey"
	"github.com/asciiu/gomo/common/db"
	orderRepo "github.com/asciiu/gomo/order-service/db/sql"
	orderProto "github.com/asciiu/gomo/order-service/proto/order"
	userRepo "github.com/asciiu/gomo/user-service/db/sql"
	user "github.com/asciiu/gomo/user-service/models"
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
	_, error := userRepo.InsertUser(db, user)
	checkErr(error)

	return db, user
}

func TestInsertOrder(t *testing.T) {
	db, user := setupService()
	defer db.Close()

	key := keyProto.ApiKeyRequest{
		UserId:      user.Id,
		Exchange:    "test",
		Key:         "key",
		Secret:      "secret",
		Description: "Hey this worked!",
	}
	apikey, error := keyRepo.InsertApiKey(db, &key)
	checkErr(error)

	orderReq := orderProto.OrderRequest{
		UserId:     user.Id,
		ApiKeyId:   apikey.ApiKeyId,
		Exchange:   apikey.Exchange,
		MarketName: "ShitCoin!",
		Price:      1.1,
		Qty:        500.10,
		Conditions: "{price <= 0.0004}",
		Side:       "buy",
		OrderType:  "limit",
	}

	order, error := orderRepo.InsertOrder(db, &orderReq)
	checkErr(error)

	if order.UserId != user.Id {
		t.Errorf("user IDs do not match")
	}
	if order.Status != "pending" {
		t.Errorf("Should be pending")
	}

	// cleanup
	error = userRepo.DeleteUserHard(db, user.Id)
	checkErr(error)
}
