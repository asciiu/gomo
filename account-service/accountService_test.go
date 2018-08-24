package main

import (
	"log"

	"github.com/asciiu/gomo/common/db"
	repoUser "github.com/asciiu/gomo/user-service/db/sql"
	user "github.com/asciiu/gomo/user-service/models"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func setupService() (*AccountService, *user.User) {
	dbUrl := "postgres://postgres@localhost:5432/gomo_test?&sslmode=disable"
	db, _ := db.NewDB(dbUrl)
	service := AccountService{DB: db}

	user := user.NewUser("first", "last", "test@email", "hash")
	_, error := repoUser.InsertUser(db, user)
	checkErr(error)

	return &service, user
}

// func TestUnsupportedKey(t *testing.T) {
// 	service, user := setupService()

// 	defer service.DB.Close()

// 	key := keyProto.KeyRequest{
// 		UserID:      user.ID,
// 		Exchange:    "monex",
// 		Key:         "public",
// 		Secret:      "secret",
// 		Description: "shit test!",
// 	}

// 	response := keyProto.KeyResponse{}
// 	service.AddKey(context.Background(), &key, &response)

// 	assert.Equal(t, "fail", response.Status, "return status is expected to be fail due to unsupported exchange")

// 	repoUser.DeleteUserHard(service.DB, user.ID)
// }
