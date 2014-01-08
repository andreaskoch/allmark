// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webpaths

import (
	"github.com/andreaskoch/allmark2/common/index"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/route"
	"strings"
)

// Create a new relative web path provider
func newRelativeWebPathProvider(logger logger.Logger, index *index.Index, baseRoute *route.Route) *RelativeWebPathProvider {
	return &RelativeWebPathProvider{
		logger:    logger,
		index:     index,
		baseRoute: baseRoute,
	}
}

type RelativeWebPathProvider struct {
	logger    logger.Logger
	index     *index.Index
	baseRoute *route.Route
}

// Get the path relative for the supplied item
func (webPathProvider *RelativeWebPathProvider) Path(itemPath string) string {
	baseRouteString := webPathProvider.baseRoute.Value()
	baseRouteChilds := webPathProvider.index.GetChilds(webPathProvider.baseRoute)

	for _, child := range baseRouteChilds {

		// ignore childs which don't match
		if !child.Route().IsMatch(itemPath) {
			continue
		}

		// intersect the child route with the base route to get full path
		childRouteString := child.Route().Value()
		path := "dsdsd/" + strings.TrimPrefix(childRouteString, baseRouteString)

		return path
	}

	// path could not be resolved
	webPathProvider.logger.Warn("The path %q could not be resolved", itemPath)
	return itemPath
}
