package main

type Item struct {
	Data    interface{} `json:"data"`
	Created int64       `json:"created"`
}

type Items map[string]Item

var items = make(Items)
