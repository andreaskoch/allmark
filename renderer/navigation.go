// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"fmt"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
)

func attachBreadcrumbNavigation(item *repository.Item) {
	item.Model.BreadcrumbNavigation = &view.BreadcrumbNavigation{
		Entries: getBreadcrumbNavigationEntries(item),
	}
}

func getBreadcrumbNavigationEntries(item *repository.Item) []*view.Breadcrumb {
	navigationEntries := make([]*view.Breadcrumb, 0)

	// abort if item or model is nil
	if item == nil || item.Model == nil {
		fmt.Println("model is nil")
		return navigationEntries
	}

	// recurse
	if item.Parent != nil {
		navigationEntries = append(navigationEntries, getBreadcrumbNavigationEntries(item.Parent)...)
	}

	// route := item.RootPathProvider().GetWebRoute(item)
	model := item.Model

	// append a new navigation entry and return it
	return append(navigationEntries, &view.Breadcrumb{
		Level: item.Level,
		Title: model.Title,
		Path:  "/" + model.AbsoluteRoute,
	})
}
