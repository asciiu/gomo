package sql_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/asciiu/gomo/common/database"
	"github.com/asciiu/gomo/common/database/sql"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func TestSum(t *testing.T) {
	db, err := database.NewDB("postgres://postgres@localhost/gomo_dev?&sslmode=disable")
	checkErr(err)
	defer db.Close()

	email := "test@email"
	user, err := sql.GetUser(db, email)
	if err != nil {
		fmt.Println(err)
	}
	if user == nil {
		t.Errorf("user not found %s", email)
	}

	fmt.Printf("%+v\n", user)
}
