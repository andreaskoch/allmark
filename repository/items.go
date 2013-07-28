// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

type Items []*Item

func (items Items) Len() int {
	return len(items)
}

func (items Items) Less(i, j int) bool {
	return items[i].Less(items[j])
}

func (items Items) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}
