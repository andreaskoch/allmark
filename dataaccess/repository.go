// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"github.com/andreaskoch/allmark2/common/route"
)

type Repository interface {
	Path() string

	Items() []*Item
	Routes() []route.Route

	// update handling
	AfterReindex() chan bool
	OnUpdate(callback func(route.Route))
	StartWatching(route route.Route)
	StopWatching(route route.Route)
}
