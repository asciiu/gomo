package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	dbUrl := fmt.Sprintf("%s", os.Getenv("DB_URL"))
	srv := NewKeyService("go.srv.apikey-service", dbUrl)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
