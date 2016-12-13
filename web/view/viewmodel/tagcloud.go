// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewmodel

import (
	"sort"
)

type TagCloud []TagCloudEntry

type TagCloudEntry struct {
	Name           string `json:"name"`
	Anchor         string `json:"anchor"`
	Route          string `json:"route"`
	Level          int    `json:"level"`
	NumberOfChildren int    `json:"numberofchildren"`
}

type SortTagCloudBy func(tagCloudEntry1, tagCloudEntry2 TagCloudEntry) bool

func (by SortTagCloudBy) Sort(tagCloud TagCloud) {
	sorter := &tagCloudSorter{
		tagCloud: tagCloud,
		by:       by,
	}

	sort.Sort(sorter)
}

type tagCloudSorter struct {
	tagCloud TagCloud
	by       SortTagCloudBy
}

func (sorter *tagCloudSorter) Len() int {
	return len(sorter.tagCloud)
}

func (sorter *tagCloudSorter) Swap(i, j int) {
	sorter.tagCloud[i], sorter.tagCloud[j] = sorter.tagCloud[j], sorter.tagCloud[i]
}

func (sorter *tagCloudSorter) Less(i, j int) bool {
	return sorter.by(sorter.tagCloud[i], sorter.tagCloud[j])
}
