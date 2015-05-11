package main

import "time"

func (items Items) GetItems() Items {
	lock.RLock()
	defer lock.RUnlock()

	for k, v := range items {
		v.Age = time.Now().Unix() - v.Created
		items[k] = v
	}

	return items
}

func (items Items) CreateOneItem(i Item) Result {
	lock.Lock()
	defer lock.Unlock()

	id := i.Id
	items[id] = i

	result.Status = "OK"
	return result
}

func (items Items) FindOneItemById(id string) Item {
	lock.RLock()
	defer lock.RUnlock()

	i, present := items[id]
	if present == false {
		return Item{}
	}

	i.Age = time.Now().Unix() - i.Created
	return i
}

func (items Items) DeleteOneItemById(id string) Result {
	lock.Lock()
	defer lock.Unlock()

	delete(items, id)

	result.Status = "OK"
	return result
}

func (items Items) FlushItems() Result {
	lock.Lock()
	defer lock.Unlock()

	for k := range items {
		delete(items, k)
	}

	result.Status = "OK"
	return result
}
