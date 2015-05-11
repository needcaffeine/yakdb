package main

import "sync"

type Item struct {
	Id      string `json:"id"`
	Value   string `json:"value"`
	Created int64  `json:"-"`
	Age     int64  `json:"age"`
}

type Items map[string]Item

var items = make(Items)

// Instantiate a global RWMutex. This is necessary because
// map operations are not atomic. @TODO: Maybe this should
// be a lock on an individual item.
var lock = &sync.RWMutex{}
