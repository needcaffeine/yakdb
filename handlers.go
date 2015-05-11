package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// This is our JSON response renderer that just sets the right content type.
func JsonResponse(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		h(w, r, p)
	}
}

// List out all the routes.
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

// List all the items in the database.
func ItemsList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)

	data := items.GetItems()
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

// Get an item.
func ItemsGet(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.WriteHeader(http.StatusOK)

	data := items.FindOneItemById(p.ByName("itemId"))

	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

// Add an item to the database.
func ItemsPut(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Read in the json document supplied to us.
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// Decode the json document into an item.
	var item Item
	if err := json.Unmarshal(body, &item); err != nil {
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	// Ensure that we are passed an id and value.
	if len(item.Id) > 0 && len(item.Value) > 0 {
		item.Created = time.Now().Unix()

		data := items.CreateOneItem(item)
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(data); err != nil {
			panic(err)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)

		fmt.Fprintln(w, "`id` and `value` are required fields.")
	}
}

// Delete an item.
func ItemsDelete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.WriteHeader(http.StatusOK)

	data := items.DeleteOneItemById(p.ByName("itemId"))
	json.NewEncoder(w).Encode(data)
}

// Delete all items.
func ItemsFlush(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)

	result := items.FlushItems()
	json.NewEncoder(w).Encode(result)
}
