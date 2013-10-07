// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"strings"
)

type ItemMap map[string]*Item

func NewItemMap() ItemMap {
	return make(ItemMap)
}

func (itemmap ItemMap) Register(item *Item) {

	key := itemmap.getMapKey(item)
	itemmap[key] = item

}

func (itemmap ItemMap) Remove(item *Item) {

	key := itemmap.getMapKey(item)
	delete(itemmap, key)

}

func (itemmap ItemMap) LookupByAlias(alias string) *Item {
	key := itemmap.normalizeKey(alias)
	if item, exists := itemmap[key]; exists {
		return item
	}

	return nil
}

func (itemmap ItemMap) getMapKey(item *Item) string {
	return itemmap.normalizeKey(item.MetaData.Alias)
}

func (itemmap ItemMap) normalizeKey(key string) string {
	return strings.ToLower(key)
}
