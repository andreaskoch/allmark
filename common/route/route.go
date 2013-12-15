// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	regexpWhitespacePattern         = regexp.MustCompile(`\s+`)
	regexpBackSlashPattern          = regexp.MustCompile(`\\+`)
	regexpdoubleForwardSlashPattern = regexp.MustCompile(`/+`)
)

type Route struct {
	value string
}

func NewFromPath(repositoryPath, itemPath string) (*Route, error) {

	// normalize the repository path
	normalizedRepositoryPath, err := normalize(repositoryPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to normalize the supplied repository path %q. Error: %s", repositoryPath, err)
	}

	// normalize the item path
	normalizedItemPath, err := normalize(itemPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to normalize the supplied item path %q. Error: %s", itemPath, err)
	}

	// prepare the route value:
	// strip the repository path from the item path
	routeValue := strings.Replace(normalizedItemPath, normalizedRepositoryPath, "", 1)

	// strip the file name
	routeValue = routeValue[:strings.LastIndex(routeValue, "/")]

	// trim leading slashes
	routeValue = strings.TrimLeft(routeValue, "/")

	return &Route{routeValue}, nil
}

func NewFromRequest(requestPath string) (*Route, error) {

	// normalize the request path
	normalizedRequestPath, err := normalize(requestPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to normalize the supplied request path %q. Error: %s", requestPath, err)
	}

	return &Route{normalizedRequestPath}, nil
}

func (route *Route) String() string {
	return route.value
}

func (route *Route) Value() string {
	return route.value
}

func (route *Route) Level() int {

	// empty routes have the level 0
	if route.value == "" {
		return 0
	}

	// routes without a slash are 1st level
	if !strings.Contains(route.value, "/") {
		return 1
	}

	// routes with slashes have a level equal to the number of slashes
	return strings.Count(route.value, "/") + 1
}

func (route *Route) Parent() *Route {
	routeValue := route.Value()

	// check if the route contains a slash, if not there is no parent
	if !strings.Contains(routeValue, "/") {
		return nil // no parent available
	}

	positionOfLastSlash := strings.LastIndex(routeValue, "/")
	parentRouteValue := routeValue[:positionOfLastSlash]

	return &Route{parentRouteValue}
}

// Check if the the current route is direct parent for the supplied (child) route.
func (parent *Route) IsParentOf(child *Route) bool {
	parentRoute := parent.Value()
	childRoute := child.Value()

	// the current route cannot be a parent for the supplied (child) route if the parent route length greater or equal than the child route length.
	if len(parentRoute) >= len(childRoute) {
		return false
	}

	// if the parent route is not the start of the child route it cannot be its parent
	if !strings.HasPrefix(childRoute, parentRoute) {
		return false
	}

	// if there is more than one slash in the relative child route, the child is not a direct descendant of the parent route
	relativeChildRoute := strings.TrimLeft(strings.Replace(childRoute, parentRoute, "", 1), "/")
	if strings.Count(relativeChildRoute, "/") > 0 {
		return false
	}

	// the child is a direct desecendant of the parent
	return true

}

// Check if the current route is a direct child of the supplied (parent) route.
func (child *Route) IsChildOf(parent *Route) bool {
	childRoute := child.Value()
	parentRoute := parent.Value()

	// the current route cannot be a child of the supplied (parent) route if the child route length less or equal than the parent route length.
	if len(childRoute) <= len(parentRoute) {
		return false
	}

	// if the child route does not start with the parent route it cannot be a child
	if !strings.HasPrefix(childRoute, parentRoute) {
		return false
	}

	// if there is more than one slash in the relative child route, the child is not a direct descendant of the parent route
	relativeChildRoute := strings.TrimLeft(strings.Replace(childRoute, parentRoute, "", 1), "/")
	if strings.Count(relativeChildRoute, "/") > 0 {
		return false
	}

	// the child is a direct desecendant of the parent
	return true
}

// Normalize the supplied path to be used for an Item or File
func normalize(path string) (string, error) {

	// trim spaces
	path = strings.TrimSpace(path)

	// check if the path is empty
	if path == "" {
		return path, fmt.Errorf("A path cannot be empty.")
	}

	// replace all backslashes with a (single) forward slash
	path = regexpBackSlashPattern.ReplaceAllString(path, "/")

	// replace multiple forward slashes with a single forward slash
	path = regexpdoubleForwardSlashPattern.ReplaceAllString(path, "/")

	// remove leading slashes
	path = strings.TrimLeft(path, "/")

	// remove trailing slashes
	path = strings.TrimRight(path, "/")

	// replace duplicate spaces with a (single) url safe character
	path = regexpWhitespacePattern.ReplaceAllString(path, "+")

	return path, nil
}
