// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"strings"
)

type ResolverExpression func(item *Item) bool

type ItemResolver func(itemName string, expression ResolverExpression) *Item

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

func (itemmap ItemMap) Lookup(alias string, expression ResolverExpression) *Item {
	results := itemmap.lookupByAlias(alias)
	if results == nil || len(results) == 0 {
		return nil
	}

	for _, item := range results {
		if expression(item) {
			return item
		}
	}

	return nil
}

func (itemmap ItemMap) lookupByAlias(alias string) ItemList {
	key := itemmap.normalizeKey(alias)
	return itemmap[key]
}

func (itemmap ItemMap) getMapKey(item *Item) string {
	return itemmap.normalizeKey(item.MetaData.Alias)
}

func (itemmap ItemMap) normalizeKey(key string) string {
	return strings.ToLower(key)
}
