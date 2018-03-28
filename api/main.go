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
	// Create a new service. Optionally include some options here.
	// service := micro.NewService(micro.Name("greeter.client"))
	// service.Init()

	// client := pb.NewSessionServiceClient("go.micro.srv.session", service.Client())
	// session := pb.SessionRequest{
	// 	Id:     "what",
	// 	UserId: "userid",
	// }

	// r, err := client.CreateSession(context.Background(), &session)
	// if err != nil {
	// 	log.Fatalf("Could not greet: %v", err)
	// }
	// log.Printf("Created: #%v", r.Response)

	dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	fmt.Println(dbUrl)

	gomoDB, err := db.NewDB(dbUrl)
	checkErr(err)
	defer gomoDB.Close()

	e := routes.New(gomoDB)
	e.Logger.Fatal(e.Start(":5000"))
}
