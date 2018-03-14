package main

import (
	"github.com/asciiu/gomo/router"
)

func main() {
	e := router.New()

	e.Logger.Fatal(e.Start(":5000"))
}
