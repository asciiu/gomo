package main

import (
	"fmt"
	"log"
	"os"

	"github.com/asciiu/gomo/api/routes"
	"github.com/asciiu/gomo/common/db"
	_ "github.com/lib/pq"
)

func checkErr(err error) {
	if err != nil {
		log.Printf("ERROR: %s", err)
		panic(err)
	}
}

func main() {

	dbUrl := fmt.Sprintf("postgres://postgres@%s:%s/%s?&sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	fmt.Println(dbUrl)

	gomoDB, err := db.NewDB(dbUrl)
	checkErr(err)
	defer gomoDB.Close()

	e := routes.New(gomoDB)
	e.Logger.Fatal(e.Start(":5000"))
}
