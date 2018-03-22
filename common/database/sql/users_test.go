package sql_test

import (
	"log"
	"testing"

	"github.com/satori/go.uuid"

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

	newId, _ := uuid.NewV4()
	user := models.User{
		Id:            newId.String(),
		First:         "test",
		Last:          "one",
		Email:         "test@email",
		EmailVerified: true,
		PasswordHash:  "hash",
		Salt:          "salt",
	}
	_, error := sql.InsertUser(db, &user)
	if error != nil {
		t.Errorf("%s", error)
	}
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
