package main

import (
	"sync"
	"time"
)

// A collection of Items.
type Collection struct {
	Items map[string]Item `json:"items"`
	cMu   sync.RWMutex
}

type Collections map[string]Collection

var collections = make(Collections)

func (collections Collections) ListCollections() Collections {
	gMu.RLock()
	defer gMu.RUnlock()

	return collections
}

func (collections Collections) GetCollection(id string) Collection {
	gMu.RLock()
	defer gMu.RUnlock()

	return collections[id]
}

func (collections Collections) CreateEmptyCollection(id string) Collection {
	var c Collection
	c.Items = make(Items)
	collections[id] = c

	return c
}

func (collections Collections) AddItem(collectionId string, itemId string, data interface{}) interface{} {
	gMu.Lock()
	defer gMu.Unlock()

	collection, present := collections[collectionId]
	if present == false {
		collection = collections.CreateEmptyCollection(collectionId)
	}

	var item Item
	item.Data = data
	item.Created = time.Now().Unix()

	collection.Items[itemId] = item

	result.Status = "OK"
	return result
}

func (collections Collections) DeleteOneCollection(id string) Result {
	gMu.Lock()
	defer gMu.Unlock()

	delete(collections, id)

	result.Status = "OK"
	return result
}

func (collections Collections) FlushCollections() Result {
	gMu.Lock()
	defer gMu.Unlock()

	for k := range collections {
		delete(collections, k)
	}

	result.Status = "OK"
	return result
}

func (collection Collection) GetItems() Items {
	collection.cMu.RLock()
	defer collection.cMu.RUnlock()

	return collection.Items
}

func (collection Collection) FindOneItemById(itemId string) Item {
	collection.cMu.RLock()
	defer collection.cMu.RUnlock()

	i, present := collection.Items[itemId]
	if present == false {
		return Item{}
	}

	return i
}

func (collection Collection) DeleteOneItemById(itemId string) Result {
	collection.cMu.Lock()
	defer collection.cMu.Unlock()

	items := collection.Items
	delete(items, itemId)

	result.Status = "OK"
	return result
}

func (collection Collection) FlushItems() Result {
	collection.cMu.Lock()
	defer collection.cMu.Unlock()

	items := collection.Items
	for k := range items {
		delete(items, k)
	}

	result.Status = "OK"
	return result
}
