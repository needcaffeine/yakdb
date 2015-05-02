package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index", "GET", "/", Index,
	},
	Route{
		"ItemsList", "GET", "/items", ItemsList,
	},
	Route{
		"ItemsGet", "GET", "/items/{itemId}", ItemsGet,
	},
	Route{
		"ItemsPut", "PUT", "/items", ItemsPut,
	},
	Route{
		"ItemsDelete", "DELETE", "/items/{itemId}", ItemsDelete,
	},
	Route{
		"ItemsFlush", "DELETE", "/items", ItemsFlush,
	},
}
