package main

import "time"

// A collection of Items.
type Collection struct {
	Items map[string]Item `json:"items"`
}

type Collections map[string]Collection

var collections = make(Collections)

func (collections Collections) ListCollections() Collections {
	lock.RLock()
	defer lock.RUnlock()

	return collections
}

func (collections Collections) GetCollection(id string) Collection {
	lock.RLock()
	defer lock.RUnlock()

	return collections[id]
}

func (collections Collections) CreateEmptyCollection(id string) Collection {
	var c Collection
	c.Items = make(Items)
	collections[id] = c

	return c
}

func (collections Collections) AddItem(collectionId string, itemId string, data interface{}) interface{} {
	lock.Lock()
	defer lock.Unlock()

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
	lock.Lock()
	defer lock.Unlock()

	delete(collections, id)

	result.Status = "OK"
	return result
}

func (collections Collections) FlushCollections() Result {
	lock.Lock()
	defer lock.Unlock()

	for k := range collections {
		delete(collections, k)
	}

	result.Status = "OK"
	return result
}

func (collection Collection) GetItems() Items {
	lock.RLock()
	defer lock.RUnlock()

	return collection.Items
}

func (collection Collection) FindOneItemById(itemId string) Item {
	lock.RLock()
	defer lock.RUnlock()

	i, present := collection.Items[itemId]
	if present == false {
		return Item{}
	}

	return i
}

func (collection Collection) DeleteOneItemById(itemId string) Result {
	lock.Lock()
	defer lock.Unlock()

	items := collection.Items
	delete(items, itemId)

	result.Status = "OK"
	return result
}

func (collection Collection) FlushItems() Result {
	lock.Lock()
	defer lock.Unlock()

	items := collection.Items
	for k := range items {
		delete(items, k)
	}

	result.Status = "OK"
	return result
}
