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

func (itemlist ItemList) Add(item *Item) {

	for _, existingItem := range itemlist {
		if existingItem == item {
			return // abort, item already exists
		}
	}

	itemlist = append(itemlist, item)

}

func (itemlist ItemList) Remove(item *Item) {

	newlist := make(ItemList, 0)

	itemRemoved := false
	for _, existingItem := range itemlist {
		if existingItem != item {
			itemRemoved = true
			newlist = append(newlist, existingItem)
		}
	}

	if itemRemoved {
		itemlist = newlist
	}

}
