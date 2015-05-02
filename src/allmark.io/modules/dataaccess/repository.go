// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"allmark.io/modules/common/route"
	"fmt"
)

type PathProvider interface {
	Path() string
}

type RoutesProvider interface {
	Routes() []route.Route
}

type ItemsProvider interface {
	Items() []*Item
	Item(route route.Route) *Item
}

type Subscriber interface {
	Subscribe(updates chan Update)
}

type Repository interface {
	PathProvider
	ItemsProvider
	RoutesProvider
	Subscriber
}

func NewUpdate(newItems, modifiedItems, deletedItems []*Item) Update {
	return Update{newItems, modifiedItems, deletedItems}
}

type Update struct {
	newItems      []*Item
	modifiedItems []*Item
	deletedItems  []*Item
}

func (update *Update) String() string {
	return fmt.Sprintf("Update (New: %v, Modified: %v, Deleted: %v)",
		len(update.newItems), len(update.modifiedItems), len(update.deletedItems))
}

func (update *Update) IsEmpty() bool {
	return len(update.newItems) == 00 && len(update.modifiedItems) == 0 && len(update.deletedItems) == 0
}

func (update *Update) New() []*Item {
	return update.newItems
}

func (update *Update) Modified() []*Item {
	return update.modifiedItems
}

func (update *Update) Deleted() []*Item {
	return update.deletedItems
}
