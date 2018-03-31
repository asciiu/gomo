// FOMO API
//
//     Schemes: https
//     BasePath: /api
//     Version: 0.0.1
//     Contact: Flowy <ellyssin.gimhae@gmail.com>
//     Host: stage.fomo.exchange
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearer
//
//     SecurityDefinitions:
//     Bearer:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package main

import (
	"fmt"
	"log"
	"os"

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
	dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	fmt.Println(dbUrl)

	gomoDB, err := db.NewDB(dbUrl)
	checkErr(err)
	defer gomoDB.Close()

	e := NewRouter(gomoDB)
	e.Logger.Fatal(e.Start(":5000"))
}
