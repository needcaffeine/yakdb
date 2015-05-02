package main

type Item struct {
	Id      string `json:"id"`
	Value   string `json:"value"`
	Created int64  `json:"-"`
	Age     int64  `json:"age"`
}

var items = make(map[string]Item)
