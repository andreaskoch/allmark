// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mapper

import (
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
)

func MapToplevelNavigation(root *repository.Item) *view.ToplevelNavigation {

	if root == nil || root.Childs == nil {
		return nil
	}

	toplevelEntries := make([]*view.ToplevelEntry, 0, len(root.Childs))
	for _, child := range root.Childs {
		toplevelEntries = append(toplevelEntries, &view.ToplevelEntry{
			Title: child.Title,
			Path:  "/" + child.AbsoluteRoute,
		})
	}

	return &view.ToplevelNavigation{
		Entries: toplevelEntries,
	}
}

func MapBreadcrumbNavigationEntries(item *repository.Item) []*view.Breadcrumb {
	navigationEntries := make([]*view.Breadcrumb, 0)

	// abort if item or model is nil
	if item == nil || item.Model == nil {
		return navigationEntries
	}

	// recurse
	if item.Parent != nil {
		navigationEntries = append(navigationEntries, MapBreadcrumbNavigationEntries(item.Parent)...)
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
