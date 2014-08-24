// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"sort"
)

type SortItemBy func(item1, item2 *Item) bool

func (by SortItemBy) Sort(items []*Item) {
	sorter := &itemSorter{
		items: items,
		by:    by,
	}

	sort.Sort(sorter)
}

type itemSorter struct {
	items []*Item
	by    SortItemBy
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
