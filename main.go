package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/dimfeld/httptreemux"
)

// Instantiate a global RWMutex. This is necessary because
// map operations are not atomic.
var gMu = &sync.RWMutex{}

func main() {
	// Read in the supplied flags from the command line to determine the port.
	// By default, listen to port :9532 (:ykdb)
	var portFlag int
	flag.IntVar(&portFlag, "port", 9532, "Port to listen on for http requests.")
	flag.Parse()
	port := strconv.Itoa(portFlag)

	router := httptreemux.New()
	router.GET("/", Index)
	router.GET("/collections", JsonResponse(CollectionsList))
	router.GET("/:collectionId", JsonResponse(CollectionsGet))
	router.PUT("/:collectionId", JsonResponse(CollectionsPut))
	router.DELETE("/:collectionId", JsonResponse(CollectionsDelete))

	router.GET("/:collectionId/items", JsonResponse(ItemsList))
	router.GET("/:collectionId/:itemId", JsonResponse(ItemsGet))
	router.PUT("/:collectionId/:itemId", JsonResponse(ItemsPut))
	router.DELETE("/:collectionId/:itemId", JsonResponse(ItemsDelete))
	router.DELETE("/:collectionId/items", JsonResponse(ItemsFlush))

	fmt.Printf("Listening on port %v...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
