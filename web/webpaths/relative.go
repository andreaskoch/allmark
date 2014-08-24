// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webpaths

import (
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/dataaccess"
	"strings"
)

// Create a new relative web path provider
func newRelativeWebPathProvider(logger logger.Logger, repository dataaccess.Repository, baseRoute route.Route) *RelativeWebPathProvider {
	return &RelativeWebPathProvider{
		logger:     logger,
		repository: repository,
		baseRoute:  baseRoute,
	}
}

type RelativeWebPathProvider struct {
	logger     logger.Logger
	repository dataaccess.Repository
	baseRoute  route.Route
}

// Get the path relative for the supplied item
func (webPathProvider *RelativeWebPathProvider) Path(itemPath string) string {
	baseRouteString := webPathProvider.baseRoute.Value()

	// filter expression which includes only childs with matching routes
	onlyMatchingItemsExpression := func(child *dataaccess.Item) bool {
		return child.Route().IsMatch(itemPath)
	}

	// get all childs which have a matching route
	baseRouteChilds := webPathProvider.repository.AllMatchingChilds(webPathProvider.baseRoute, onlyMatchingItemsExpression)

	// abort if no matching routes have been found
	if noMatchingChildsFound := len(baseRouteChilds) == 0; noMatchingChildsFound {
		// path could not be resolved, try to trim the path
		return strings.TrimPrefix(strings.TrimPrefix(itemPath, baseRouteString), "/")
	}

	// use only the first child
	child := baseRouteChilds[0]

	// intersect the child route with the base route to get full path
	childRouteString := child.Route().Value()
	path := strings.TrimPrefix(strings.TrimPrefix(childRouteString, baseRouteString), "/")

	return path
}

func (webPathProvider *RelativeWebPathProvider) Base() route.Route {
	return webPathProvider.baseRoute
}
