package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func main() {
	// Read in the supplied flags from the command line to determine the port.
	// By default, listen to port :9532 (:ykdb)
	var portFlag int
	flag.IntVar(&portFlag, "port", 9532, "Port to listen on for http requests.")
	flag.Parse()
	port := strconv.Itoa(portFlag)

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/items", ItemsList)
	router.GET("/items/:itemId", ItemsGet)
	router.PUT("/items", ItemsPut)
	router.DELETE("/items/:itemId", ItemsDelete)
	router.DELETE("/items", ItemsFlush)

	fmt.Printf("Listening on port %v...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
