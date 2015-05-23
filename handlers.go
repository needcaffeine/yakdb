package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	//"time"

	"github.com/dimfeld/httptreemux"
)

// This is our JSON response renderer that just sets the right content type.
func JsonResponse(h httptreemux.HandlerFunc) httptreemux.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, p map[string]string) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		h(w, r, p)
	}
}

// List out all the routes.
func Index(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	fmt.Fprintln(
		w,
		"This is yakdb, a highly performant key-value store written in Go.\n\n"+
			"Usage:\n"+
			"------\n"+
			"List all items: GET /items\n"+
			"Get an item: GET /items/{itemid}\n"+
			"Put an item: PUT /items\n"+
			"Delete an item: DELETE /items/{itemid}\n"+
			"Delete all items: DELETE /items\n\n"+
			"More documentation: https://github.com/needcaffeine/yakdb",
	)
}

// List all the collections.
func CollectionsList(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	w.WriteHeader(http.StatusOK)

	data := collections.ListCollections()
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

// Get the contents of a collection.
func CollectionsGet(w http.ResponseWriter, r *http.Request, p map[string]string) {
	w.WriteHeader(http.StatusOK)

	collectionId := p["collectionId"]

	data := collections[collectionId]

	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

// Add a collection.
func CollectionsPut(w http.ResponseWriter, r *http.Request, p map[string]string) {
	w.WriteHeader(http.StatusCreated)

	collectionId := p["collectionId"]

	data := collections.CreateEmptyCollection(collectionId)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

// Delete a collection.
func CollectionsDelete(w http.ResponseWriter, r *http.Request, p map[string]string) {
	w.WriteHeader(http.StatusOK)

	data := collections.DeleteOneCollection(p["collectionId"])
	json.NewEncoder(w).Encode(data)
}

// Add an item to a collection.
func ItemsPut(w http.ResponseWriter, r *http.Request, p map[string]string) {
	// Read in the json document supplied to us.
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// Decode the json document.
	var document interface{}
	if err := json.Unmarshal(body, &document); err != nil {
		if err := json.NewEncoder(w).Encode(err); err != nil {
			w.WriteHeader(422)
			panic(err)
		}
	}

	if len(body) > 0 {
		collectionId := p["collectionId"]
		itemId := p["itemId"]

		//collection := collections.GetCollection(collectionId)
		data := collections.AddItem(collectionId, itemId, document)

		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(data); err != nil {
			panic(err)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)

		fmt.Fprintln(w, "Empty body.")
	}
}

// List all the items in the collection.
func ItemsList(w http.ResponseWriter, r *http.Request, p map[string]string) {
	w.WriteHeader(http.StatusOK)

	collectionId := p["collectionId"]
	collection := collections[collectionId]

	data := collection.GetItems()
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

// Get an item.
func ItemsGet(w http.ResponseWriter, r *http.Request, p map[string]string) {
	w.WriteHeader(http.StatusOK)

	collectionId := p["collectionId"]
	collection := collections[collectionId]
	itemId := p["itemId"]

	data := collection.FindOneItemById(itemId)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

// Delete an item.
func ItemsDelete(w http.ResponseWriter, r *http.Request, p map[string]string) {
	w.WriteHeader(http.StatusOK)

	collectionId := p["collectionId"]
	collection := collections[collectionId]
	itemId := p["itemId"]

	data := collection.DeleteOneItemById(itemId)
	json.NewEncoder(w).Encode(data)
}

// Delete all items.
func ItemsFlush(w http.ResponseWriter, r *http.Request, p map[string]string) {
	w.WriteHeader(http.StatusOK)

	collectionId := p["collectionId"]
	collection := collections[collectionId]

	data := collection.FlushItems()
	json.NewEncoder(w).Encode(data)
}
