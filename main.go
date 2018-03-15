package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/asciiu/gomo/router"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
)

const (
	DB_URL      = "postgresql://postgres@localhost:5432/fomo_dev"
	DB_HOST     = "localhost"
	DB_USER     = "postgres"
	DB_PASSWORD = ""
	DB_NAME     = "fomo_dev"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func main() {
	dbinfo := fmt.Sprintf("host=%s dbname=%s sslmode=disable",
		DB_HOST, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error: Could not establish a connection with the database")
	}

	// panic on error
	u1 := uuid.Must(uuid.NewV4())

	stmt, err := db.Prepare("INSERT INTO users(id, first_name, last_name, email, password, salt) VALUES($1,$2,$3,$4,$5,$6)")
	if err != nil {
		log.Print("HERE")
		log.Fatal(err)
	}
	res, err := stmt.Exec(u1, "test", "name", "test@email", "password", "salt")
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("affected = %d\n", rowCnt)

	e := router.New()
	e.Logger.Fatal(e.Start(":5000"))
}
