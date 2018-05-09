// FOMO API
//
// Endpoints labeled open do not require authentication. The protected endpoints on the other hand, do require
// authentication. I'm not sure how long the jwt token should last. I'm thinking we should set the expire on
// that token to be super short - like 5 minutes (upto an hour maybe?) to minimize the amount of time an
// attacker can use that token. The refresh token will last longer - currently 7 days. If you make a request
// to a protected endpoint using a "Refresh" token in your request headers, you will receive a new
// authorization token (set-authorization) and refresh token (set-refresh) in the response headers when
// you make a request with an expired authorization token. You MUST replace both tokens in your request headers
// to stay authenticated. The old refresh token gets replaced on the backend therefore, you need to use the
// new refresh token to remain actively logged in.
//
//     Schemes: https
//     BasePath: /api
//     Version: 0.0.1
//     Author: The flo
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
	dbURL := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	gomoDB, err := db.NewDB(dbURL)
	checkErr(err)
	defer gomoDB.Close()

	e := NewRouter(gomoDB)

	e.Logger.Fatal(e.Start(":5000"))
}
