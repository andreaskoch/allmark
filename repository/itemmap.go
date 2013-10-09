// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"fmt"
	"strings"
)

type ItemMap map[string]ItemList

func NewItemMap() ItemMap {
	return make(ItemMap)
}

func (itemmap ItemMap) Register(item *Item) {

	key := itemmap.getMapKey(item)
	if itemlist, exists := itemmap[key]; exists {

		// add the item to the item list
		itemmap[key] = itemlist.Add(item)

	} else {

		// create a new item list
		itemmap[key] = NewItemList(item)
	}

}

func (itemmap ItemMap) Remove(item *Item) {

	key := itemmap.getMapKey(item)
	if itemlist, exists := itemmap[key]; exists {

		newItemList := itemlist.Remove(item)
		if len(newItemList) > 0 {
			itemmap[key] = newItemList
		} else {
			delete(itemmap, key)
		}

	}

}

func (itemmap ItemMap) LookupByAlias(alias string) *Item {
	key := itemmap.normalizeKey(alias)
	if itemList, exists := itemmap[key]; exists {
		if len(itemList) == 0 {
			return nil
		}

		// debug
		if len(itemList) > 1 {
			fmt.Printf("There is more than one document for the alias %q.\n", alias)
		}

		return itemList[0]
	}

	return nil
}

func (itemmap ItemMap) getMapKey(item *Item) string {
	return itemmap.normalizeKey(item.MetaData.Alias)
}

func (itemmap ItemMap) normalizeKey(key string) string {
	return strings.ToLower(key)
}
