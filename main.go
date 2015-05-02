package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := NewRouter()

	// By default, listen to port :9532 (:ykdb)
	fmt.Println("Listening on port 9532...")
	log.Fatal(http.ListenAndServe(":9532", router))
}
