// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"allmark.io/modules/common/route"
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

type RepositoryUpdater interface {
	AfterReindex(notificationChannel chan bool)
	OnUpdate(callback func(route.Route))
	StartWatching(route route.Route)
	StopWatching(route route.Route)
}

type Repository interface {
	PathProvider
	ItemsProvider
	RoutesProvider

	// update handling
	RepositoryUpdater
}
