// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package renderer

import (
	"fmt"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/view"
)

func attachNavigation(item *repository.Item) {
	item.Model.Navigation = &view.Navigation{
		Entries: getNavigationEntries(item),
	}
}

func getNavigationEntries(item *repository.Item) []*view.NavigationEntry {
	navigationEntries := make([]*view.NavigationEntry, 0)

	// abort if item or model is nil
	if item == nil || item.Model == nil {
		fmt.Println("model is nil")
		return navigationEntries
	}

	// recurse
	if item.Parent != nil {
		navigationEntries = append(navigationEntries, getNavigationEntries(item.Parent)...)
	}

	// route := item.RootPathProvider().GetWebRoute(item)
	model := item.Model

	// append a new navigation entry and return it
	return append(navigationEntries, &view.NavigationEntry{
		Level: item.Level,
		Title: model.Title,
		Path:  "/" + model.AbsoluteRoute,
	})
}
