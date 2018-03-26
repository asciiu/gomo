package main

import (
	"fmt"
	"log"
	"os"

	"github.com/asciiu/gomo/api/routes"
	"github.com/asciiu/gomo/common/database"
	_ "github.com/lib/pq"
)

func checkErr(err error) {
	if err != nil {
		log.Printf("ERROR: %s", err)
		panic(err)
	}
}

func main() {
	dbUrl := fmt.Sprintf("postgres://postgres@%s:%s/gomo_dev?&sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	fmt.Println(dbUrl)

	db, err := database.NewDB(dbUrl)
	checkErr(err)
	defer db.Close()

	e := routes.New(db)
	e.Logger.Fatal(e.Start(":5000"))
}
