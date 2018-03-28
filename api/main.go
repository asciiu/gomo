package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/asciiu/gomo/api/routes"
	"github.com/asciiu/gomo/common/db"
	pb "github.com/asciiu/gomo/session-service/proto/session"
	_ "github.com/lib/pq"
	"github.com/micro/go-micro/client"
)

func checkErr(err error) {
	if err != nil {
		log.Printf("ERROR: %s", err)
		panic(err)
	}
}

func main() {
	client := pb.NewSessionServiceClient("go.micro.srv.sessions", client.DefaultClient)
	session := pb.Session{
		Id:     "what",
		UserId: "userid",
	}

	r, err := client.CreateSession(context.Background(), &session)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Created: #%v", r.Session)

	dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	fmt.Println(dbUrl)

	gomoDB, err := db.NewDB(dbUrl)
	checkErr(err)
	defer gomoDB.Close()

	e := routes.New(gomoDB)
	e.Logger.Fatal(e.Start(":5000"))
}
