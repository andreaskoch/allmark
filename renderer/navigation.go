// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"github.com/andreaskoch/allmark/mapper"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
)

func attachBreadcrumbNavigation(item *repository.Item) {
	item.Model.BreadcrumbNavigation = &view.BreadcrumbNavigation{
		Entries: mapper.MapBreadcrumbNavigationEntries(item),
	}
}

func attachToplevelNavigation(root, item *repository.Item) {
	if item == nil {
		return
	}

	toplevelNavigation := mapper.MapToplevelNavigation(root)
	item.ToplevelNavigation = toplevelNavigation
}
