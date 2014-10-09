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

	var matchingRouteHasBeenFound bool
	var matchingRoute route.Route
	for _, route := range webPathProvider.repository.Routes() {

		// ignore all routes which are not a child of the base route
		if !route.IsChildOf(webPathProvider.baseRoute) {
			continue
		}

		// ignore all non-matching routes
		if !route.IsMatch(itemPath) {
			continue
		}

		// a matching route has been found
		matchingRouteHasBeenFound = true
		matchingRoute = route
		break
	}

	if !matchingRouteHasBeenFound {
		// path could not be resolved, return fallback
		return "/" + itemPath
	}

	// intersect the child route with the base route to get full path
	path := strings.TrimPrefix(matchingRoute.Value(), webPathProvider.baseRoute.Value())
	return strings.TrimPrefix(path, "/")
}

func (webPathProvider *RelativeWebPathProvider) Base() route.Route {
	return webPathProvider.baseRoute
}
