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
	router.GET("/items", JsonResponse(ItemsList))
	router.GET("/items/:itemId", JsonResponse(ItemsGet))
	router.PUT("/items", JsonResponse(ItemsPut))
	router.DELETE("/items/:itemId", JsonResponse(ItemsDelete))
	router.DELETE("/items", JsonResponse(ItemsFlush))

	fmt.Printf("Listening on port %v...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
