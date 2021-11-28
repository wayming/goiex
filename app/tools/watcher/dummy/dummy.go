package main

import (
	"log"
	"time"
)

func main() {
	for {
		log.Println("Delay 1 second")
		time.Sleep(time.Second)
	}
}
