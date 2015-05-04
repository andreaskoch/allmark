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
	Items() []Item
	Item(route route.Route) Item
}

type Subscriber interface {
	Subscribe(updates chan Update)
}

type LiveReload interface {
	StartWatching(route route.Route)
	StopWatching(route route.Route)
}

type Repository interface {
	PathProvider
	ItemsProvider
	RoutesProvider
	Subscriber
	LiveReload
}

func NewModifiedItemUpdate(modifiedRoute route.Route) Update {
	return NewUpdate([]route.Route{}, []route.Route{modifiedRoute}, []route.Route{})
}

func NewUpdate(newItemRoutes, modifiedItemRoutes, deletedItemRoutes []route.Route) Update {
	return Update{newItemRoutes, modifiedItemRoutes, deletedItemRoutes}
}

type Update struct {
	newItemRoutes      []route.Route
	modifiedItemRoutes []route.Route
	deletedItemRoutes  []route.Route
}

func (update *Update) String() string {
	return fmt.Sprintf("Update (New: %v, Modified: %v, Deleted: %v)",
		len(update.newItemRoutes), len(update.modifiedItemRoutes), len(update.deletedItemRoutes))
}

func (update *Update) IsEmpty() bool {
	return len(update.newItemRoutes) == 00 && len(update.modifiedItemRoutes) == 0 && len(update.deletedItemRoutes) == 0
}

func (update *Update) New() []route.Route {
	return update.newItemRoutes
}

func (update *Update) Modified() []route.Route {
	return update.modifiedItemRoutes
}

func (update *Update) Deleted() []route.Route {
	return update.deletedItemRoutes
}
