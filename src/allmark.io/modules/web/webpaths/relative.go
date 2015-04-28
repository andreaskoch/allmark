// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webpaths

import (
	"strings"

	"allmark.io/modules/common/route"
	"allmark.io/modules/dataaccess"
)

// Create a new relative web path provider
func newRelativeWebPathProvider(routesProvider dataaccess.RoutesProvider, baseRoute route.Route) *RelativeWebPathProvider {
	return &RelativeWebPathProvider{
		routesProvider: routesProvider,
		baseRoute:      baseRoute,
	}
}

type RelativeWebPathProvider struct {
	routesProvider dataaccess.RoutesProvider
	baseRoute      route.Route
}

// Get the path relative for the supplied item
func (webPathProvider *RelativeWebPathProvider) Path(itemPath string) string {

	// return the supplied item path if it is already absolute
	if isAbsoluteUri(itemPath) {
		return itemPath
	}

	var matchingRouteHasBeenFound bool
	var matchingRoute route.Route
	for _, route := range webPathProvider.routesProvider.Routes() {

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
