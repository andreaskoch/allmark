// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dataaccess

import (
	"github.com/elWyatt/allmark/common/content"
	"github.com/elWyatt/allmark/common/route"
)

// A File represents a file ressource that is associated with an Item.
type File interface {
	content.ContentProviderInterface

	String() string
	Id() string
	Name() string
	Parent() route.Route
	Route() route.Route
}
