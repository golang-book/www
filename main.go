//go:generate go run scripts/fileversions/main.go

package main

import (
	"log"
)

func main() {
	err := build()
	if err != nil {
		log.Fatalln(err)
	}
}
