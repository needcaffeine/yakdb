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

// List all the items in the system.
func ItemsList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)

	itemsCopy := GetItems()
	if err := json.NewEncoder(w).Encode(itemsCopy); err != nil {
		panic(err)
	}
}

// Get an item.
func ItemsGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusOK)

	item := FindOneItemById(ps.ByName("itemId"))

	if err := json.NewEncoder(w).Encode(item); err != nil {
		panic(err)
	}
}

func ItemsPut(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

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
		i := CreateOneItem(item)
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(i); err != nil {
			panic(err)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)

		fmt.Fprintln(w, "`id` and `value` are required fields.")
	}
}

func ItemsDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusOK)

	result := DeleteOneItemById(ps.ByName("itemId"))
	json.NewEncoder(w).Encode(result)
}

func ItemsFlush(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)

	result := FlushItems()
	json.NewEncoder(w).Encode(result)
}
