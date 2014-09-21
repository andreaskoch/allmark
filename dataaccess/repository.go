// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"github.com/andreaskoch/allmark2/common/route"
)

type Repository interface {
	String() string

	// Meta Data
	Id() string
	Path() string
	Size() int

	// Content Access
	Root() *Item
	Items() []*Item
	Item(route route.Route) (*Item, bool)
	File(route route.Route) (*File, bool)
	Parent(route route.Route) *Item
	Childs(route route.Route) []*Item
	AllChilds(route route.Route) []*Item
	AllMatchingChilds(route route.Route, matchExpression func(item *Item) bool) []*Item

	AfterReindex() chan bool

	// update handling
	OnUpdate(callback func(route.Route))
	StartWatching(route route.Route)
	StopWatching(route route.Route)

	// Fulltext Search
	Search(keywords string, maxiumNumberOfResults int) []SearchResult
}
