// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webpaths

import (
	"testing"

	"allmark.io/modules/common/route"
)

func Test_RelativeWebPathProvider_RootAsBasePath_Path_ReturnsPathWithLeadingSlash(t *testing.T) {
	// arrange
	baseRoute := route.New()
	routes := getRoutesFromStrings([]string{
		"",
	})
	routesProvider := dummyRoutesProvider{routes}
	pathProvider := newRelativeWebPathProvider(routesProvider, baseRoute)
	inputPath := "ya/da/ya/da"
	expected := "/ya/da/ya/da"

	// act
	result := pathProvider.Path(inputPath)

	// assert
	if result != expected {
		t.Errorf("The result for pathProvider.Path(%q) should be %q but was %q.", inputPath, expected, result)
	}
}

func Test_RelativeWebPathProvider_RootAsBasePath_Path_ParameterIsAbsolute_ReturnsPathAsSpecified(t *testing.T) {
	// arrange
	baseRoute := route.New()
	routes := getRoutesFromStrings([]string{
		"",
	})
	routesProvider := dummyRoutesProvider{routes}
	pathProvider := newRelativeWebPathProvider(routesProvider, baseRoute)
	inputPath := "ftp://example.com/ya/da/ya/da"
	expected := "ftp://example.com/ya/da/ya/da"

	// act
	result := pathProvider.Path(inputPath)

	// assert
	if result != expected {
		t.Errorf("The result for pathProvider.Path(%q) should be %q but was %q.", inputPath, expected, result)
	}
}

// Get an array of route.Route objects from a string array of Uris.
func getRoutesFromStrings(uris []string) []route.Route {

	routes := []route.Route{}

	for _, uri := range uris {
		route, err := route.NewFromRequest(uri)
		if err != nil {
			continue
		}
		routes = append(routes, route)
	}

	return routes
}

// A dummy route provider
type dummyRoutesProvider struct {
	routes []route.Route
}

func (provider dummyRoutesProvider) Routes() []route.Route {
	return provider.routes
}
