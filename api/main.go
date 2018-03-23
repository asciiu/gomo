package main

import (
	"log"

	"github.com/asciiu/gomo/api/routes"
	"github.com/asciiu/gomo/common/database"
	_ "github.com/lib/pq"
)

const (
	privateKeyPath = "keys/gomo-key-ecdsa"
	publicKeyPath  = "keys/gomo-key-ecdsa.pub"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func main() {
	db, err := database.NewDB("postgres://postgres@localhost/gomo_dev?&sslmode=disable")
	checkErr(err)
	defer db.Close()

	e := routes.New(db)
	e.Logger.Fatal(e.Start(":5000"))
}
