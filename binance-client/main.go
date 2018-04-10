package main

import (
	"fmt"
	"time"
)

func main() {

	msg := "Hello"
	go func(m string) {
		fmt.Println(m)
	}(msg)

	msg = "Goodbye"
	time.Sleep(1000 * time.Millisecond)
}
