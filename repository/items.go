// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

import (
	"sort"
)

type Items []*Item

type By func(item1, item2 *Item) bool

func (by By) Sort(items Items) {
	sorter := &itemSorter{
		items: items,
		by:    by,
	}

	sort.Sort(sorter)
}

type itemSorter struct {
	items Items
	by    By
}

func (sorter *itemSorter) Len() int {
	return len(sorter.items)
}

func (sorter *itemSorter) Swap(i, j int) {
	sorter.items[i], sorter.items[j] = sorter.items[j], sorter.items[i]
}

func (sorter *itemSorter) Less(i, j int) bool {
	return sorter.by(sorter.items[i], sorter.items[j])
}
