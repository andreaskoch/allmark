// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tree

import (
	"github.com/elWyatt/allmark/common/route"
)

func RouteToPath(route route.Route) Path {
	if route.IsEmpty() {
		return NewPath()
	}

	return NewPath(route.Components()...)
}
