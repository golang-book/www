// +build !appengine

package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("starting server on 127.0.0.1:8002")
	err := http.ListenAndServe("127.0.0.1:8002", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
