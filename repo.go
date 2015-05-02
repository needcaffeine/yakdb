package main

import "time"

func FindOneItemById(id string) Item {
	i, present := items[id]
	if present == false {
		return Item{}
	}

	i.Age = time.Now().Unix() - i.Created
	return i
}

func GetItems() map[string]Item {
	for k, v := range items {
		v.Age = time.Now().Unix() - v.Created
		items[k] = v
	}

	return items
}

func CreateOneItem(i Item) Result {
	id := i.Id
	items[id] = i

	result.Status = "OK"
	return result
}

func DeleteOneItemById(id string) Result {
	delete(items, id)

	result.Status = "OK"
	return result
}

func FlushItems() Result {
	for k := range items {
		delete(items, k)
	}

	result.Status = "OK"
	return result
}
