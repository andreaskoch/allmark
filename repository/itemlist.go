// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

type ItemList []*Item

func NewItemList(items ...*Item) ItemList {
	itemlist := make(ItemList, 0)

	for _, item := range items {
		itemlist = append(itemlist, item)
	}

	return itemlist
}

func (itemlist ItemList) IsEmpty() bool {
	return len(itemlist) == 0
}

func (itemlist ItemList) Add(item *Item) ItemList {

	for _, existingItem := range itemlist {
		if existingItem == item {
			return itemlist // abort, item already exists
		}
	}

	return append(itemlist, item)

}

func (itemlist ItemList) Remove(item *Item) ItemList {

	newlist := make(ItemList, 0)

	itemRemoved := false
	for _, existingItem := range itemlist {
		if existingItem != item {
			itemRemoved = true
			newlist = append(newlist, existingItem)
		}
	}

	if !itemRemoved {
		return itemlist
	}

	return newlist
}
