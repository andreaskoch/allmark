// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"github.com/andreaskoch/allmark2/common/route"
	"sort"
	"strings"
)

type ItemType int

func (itemType ItemType) String() string {
	switch itemType {
	case TypeDocument:
		return "document"

	case TypePresentation:
		return "presentation"

	case TypeMessage:
		return "message"

	case TypeLocation:
		return "location"

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
	TypeMessage
	TypeLocation
	TypeRepository
	TypeUnknown
)

// An Item represents a single document.
type Item struct {
	isVirtual bool
	route     *route.Route
	files     []*File

	Type        ItemType
	Title       string
	Description string
	Content     string

	MetaData *MetaData
}

func NewVirtualItem(route *route.Route, itemType ItemType) (*Item, error) {
	return &Item{
		isVirtual: true,
		route:     route,
		Type:      itemType,
	}, nil
}

func NewItem(route *route.Route, files []*File) (*Item, error) {
	return &Item{
		isVirtual: false,
		route:     route,
		files:     files,
	}, nil
}

func (item *Item) String() string {
	return fmt.Sprintf("%s", item.route)
}

func (item *Item) FolderName() string {
	return item.route.LastComponentName()
}

func (item *Item) IsVirtual() bool {
	return item.isVirtual
}

func (item *Item) Route() *route.Route {
	return item.route
}

func (item *Item) Files() []*File {
	return item.files
}

func (item *Item) GetFile(fileRoute route.Route) *File {
	for _, file := range item.Files() {
		if !strings.HasSuffix(fileRoute.Value(), file.Route().Value()) {
			continue
		}

		return file
	}

	return nil
}

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
