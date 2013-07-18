// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"
	"strings"
)

type Route struct {
	originalRoute   string
	normalizedRoute string
	filepath        string
}

func newRoute(route, filepath string) (*Route, error) {
	if strings.TrimSpace(route) == "" {
		return nil, fmt.Errorf("A route cannot be empty.")
	}

	if strings.TrimSpace(filepath) == "" {
		return nil, fmt.Errorf("The filepath of a route cannot be empty (route: %s)", route)
	}

	return &Route{
		originalRoute:   route,
		normalizedRoute: normalizeRoute(route),
		filepath:        filepath,
	}, nil
}

func (route *Route) String() string {
	return route.Normalized()
}

func (route *Route) Original() string {
	return route.originalRoute
}

func (route *Route) Normalized() string {
	return route.normalizedRoute
}

func (route *Route) Filepath() string {
	return route.filepath
}
