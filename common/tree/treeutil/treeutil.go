// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package treeutil

import (
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/common/tree"
)

func RouteToPath(route *route.Route) tree.Path {
	if route == nil {
		return tree.NewPath()
	}

	return tree.NewPath(route.Components()...)
}
