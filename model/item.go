// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"sort"
	"strings"

	"github.com/elWyatt/allmark/common/route"
	"github.com/elWyatt/allmark/dataaccess"
)

type ItemType int

func (itemType ItemType) String() string {
	switch itemType {
	case TypeDocument:
		return "document"

	case TypePresentation:
		return "presentation"

	case TypeRepository:
		return "repository"

	default:
		return "unknown"

	}

	panic("Unreachable")
}

const (
	TypeDocument ItemType = iota
	TypePresentation
	TypeRepository
	TypeUnknown
)

// An Item represents a single document.
type Item struct {
	route      route.Route
	files      []*File
	sourceType dataaccess.ItemType

	Type ItemType

	Title       string
	Description string
	Content     string
	Markdown    string

	Hash string

	MetaData MetaData
}

func NewItem(route route.Route, files []*File, sourceType dataaccess.ItemType) *Item {

	return &Item{
		route:      route,
		files:      files,
		sourceType: sourceType,
	}

}

func (item *Item) String() string {
	return fmt.Sprintf("%s", item.route)
}

func (item *Item) FolderName() string {
	return item.route.LastComponentName()
}

func (item *Item) Route() route.Route {
	return item.route
}

func (item *Item) Files() []*File {
	return item.files
}

// Get the file which matches the supplied route. Returns nil if there is no matching file.
func (item *Item) GetFile(fileRoute route.Route) *File {
	for _, file := range item.Files() {
		if !strings.HasSuffix(fileRoute.Value(), file.Route().Value()) {
			continue
		}

		return file
	}

	return nil
}

func (item *Item) IsPhysical() bool {
	return item.sourceType == dataaccess.TypePhysical
}

func (item *Item) IsVirtual() bool {
	return item.sourceType == dataaccess.TypeVirtual
}

func (item *Item) IsFileCollection() bool {
	return item.sourceType == dataaccess.TypeFileCollection
}

type SortItemsBy func(item1, item2 *Item) bool

func (by SortItemsBy) Sort(items []*Item) {
	sorter := &modelSorter{
		items: items,
		by:    by,
	}

	sort.Sort(sorter)
}

type modelSorter struct {
	items []*Item
	by    SortItemsBy
}

func (sorter *modelSorter) Len() int {
	return len(sorter.items)
}

func (sorter *modelSorter) Swap(i, j int) {
	sorter.items[i], sorter.items[j] = sorter.items[j], sorter.items[i]
}

func (sorter *modelSorter) Less(i, j int) bool {
	return sorter.by(sorter.items[i], sorter.items[j])
}
