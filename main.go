package main

import (
	"log"

	"github.com/asciiu/gomo/database"
	"github.com/asciiu/gomo/router"
	_ "github.com/lib/pq"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func main() {
	db, err := database.NewDB("postgres://postgres@localhost/fomo_dev?&sslmode=disable")
	checkErr(err)
	defer db.Close()

	e := router.New(db)
	e.Logger.Fatal(e.Start(":5000"))
}
