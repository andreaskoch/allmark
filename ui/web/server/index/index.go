// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/model"
)

func New(logger logger.Logger) *Index {
	return &Index{
		logger: logger,
		items:  make(map[route.Route]*model.Item),
	}
}

type Index struct {
	logger logger.Logger
	items  map[route.Route]*model.Item
}

func (index *Index) IsMatch(route route.Route) (item *model.Item, isMatch bool) {
	item, isMatch = index.items[route]
	return
}

func (index *Index) GetParent(item *model.Item) *model.Item {
	childRoute := item.Route()

	for parentRoute, parentItem := range index.items {
		if !parentRoute.IsParentOf(childRoute) {
			continue
		}

		return parentItem
	}

	// no parent found
	return nil
}

func (index *Index) GetChilds(item *model.Item) []*model.Item {
	route := item.Route()
	childs := make([]*model.Item, 0)
	for itemRoute, item := range index.items {
		if !itemRoute.IsChildOf(route) {
			continue
		}

		childs = append(childs, item)
	}

	return childs
}

func (index *Index) Routes() []route.Route {
	routes := make([]route.Route, 0)
	for route, _ := range index.items {
		routes = append(routes, route)
	}
	return routes
}

func (index *Index) Add(item *model.Item) {
	index.logger.Debug("Adding item %q to index", item)
	index.items[*item.Route()] = item
}
