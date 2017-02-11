// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

import (
	"sort"
)

type Tags struct {
	Model
	Tags []Tag
}

type Tag struct {
	Name     string  `json:"name"`
	Anchor   string  `json:"anchor"`
	Route    string  `json:"route"`
	Children []Model `json:"children"`
}

type SortTagBy func(tag1, tag2 Tag) bool

func (by SortTagBy) Sort(tags []Tag) {
	sorter := &tagSorter{
		tags: tags,
		by:   by,
	}

	sort.Sort(sorter)
}

type tagSorter struct {
	tags []Tag
	by   SortTagBy
}

func (sorter *tagSorter) Len() int {
	return len(sorter.tags)
}

func (sorter *tagSorter) Swap(i, j int) {
	sorter.tags[i], sorter.tags[j] = sorter.tags[j], sorter.tags[i]
}

func (sorter *tagSorter) Less(i, j int) bool {
	return sorter.by(sorter.tags[i], sorter.tags[j])
}
