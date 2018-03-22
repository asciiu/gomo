package sql_test

import (
	"log"
	"testing"

	"github.com/asciiu/gomo/common/database"
	"github.com/asciiu/gomo/common/database/sql"
	"github.com/asciiu/gomo/common/models"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func TestInsertUser(t *testing.T) {
	db, err := database.NewDB("postgres://postgres@localhost/gomo_test?&sslmode=disable")
	checkErr(err)
	defer db.Close()

	user := models.NewUser("test@email", "hash")
	_, error := sql.InsertUser(db, user)
	if error != nil {
		t.Errorf("%s", error)
	}
	//fmt.Printf("%#v", *user)
}

func TestFindUser(t *testing.T) {
	db, err := database.NewDB("postgres://postgres@localhost/gomo_test?&sslmode=disable")
	checkErr(err)
	defer db.Close()

	email := "test@email"
	user, err := sql.FindUser(db, email)
	if err != nil {
		t.Errorf("%s", err)
	}
	if user == nil {
		t.Errorf("user not found %s", email)
	}

	sqlStatement := `delete from users where email = $1`
	db.Exec(sqlStatement, email)
}
