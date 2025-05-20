package main

import (
	"log"

	"github.com/birabittoh/go-lift/src"
)

func main() {
	err := src.Run()
	if err != nil {
		log.Fatal(err)
	}
}
