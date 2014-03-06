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

func NewFromItemPath(basePath, itemPath string) (*Route, error) {

	// normalize the base path
	normalizedBasePath := normalize(basePath)

	// normalize the item path
	normalizedItemPath := normalize(itemPath)

	// prepare the route value:
	// strip the repository path from the item path
	routeValue := strings.Replace(normalizedItemPath, normalizedBasePath, "", 1)

	// strip the file name
	routeValue = routeValue[:strings.LastIndex(routeValue, "/")]

	// trim leading slashes
	routeValue = strings.TrimLeft(routeValue, "/")

	return &Route{routeValue}, nil
}

func NewFromFilePath(basePath, itemPath string) (*Route, error) {

	// normalize the base path
	normalizedBasePath := normalize(basePath)

	// normalize the item path
	normalizedItemPath := normalize(itemPath)

	// prepare the route value:
	// strip the repository path from the item path
	routeValue := strings.Replace(normalizedItemPath, normalizedBasePath, "", 1)

	// trim leading slashes
	routeValue = strings.TrimLeft(routeValue, "/")

	return &Route{routeValue}, nil
}

func NewFromRequest(requestPath string) (*Route, error) {

	// normalize the request path
	normalizedRequestPath := normalize(requestPath)

	return &Route{
		normalizedRequestPath,
	}, nil
}

func New() (*Route, error) {

	// normalize the request path
	normalizedRequestPath := normalize("")

	return &Route{
		normalizedRequestPath,
	}, nil
}

func Combine(route1, route2 *Route) (*Route, error) {
	return NewFromRequest(route1.Value() + "/" + route2.Value())
}

func (route *Route) String() string {
	return route.value
}

func (route *Route) Value() string {
	return route.value
}

func (route *Route) FolderName() string {
	lastSlashPosition := strings.LastIndex(route.value, "/")
	if lastSlashPosition == -1 {
		return route.value
	}

	return route.value[lastSlashPosition:]
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

func (route *Route) SubRoute(level int) (*Route, error) {

	// root level
	if level == 0 {
		return NewFromRequest("")
	}

	// same level
	if level == route.Level() {
		return route, nil
	}

	// split path into components
	components := strings.Split(route.value, "/")

	// abort if the requested level is out of range
	if level > len(components)-1 {
		return nil, fmt.Errorf("The route %q does nof a have a sub-route with the level %d.", route, level)
	}

	// assemble the sub route
	subset := components[0:level]
	subRoutePath := strings.Join(subset, "/")

	subRoute, err := NewFromRequest(subRoutePath)
	if err != nil {
		return nil, fmt.Errorf("Unable to create a route from the path %q. Error: %s", subRoutePath, err)
	}

	return subRoute, nil
}

func (route *Route) IsMatch(path string) bool {
	cleanedRoute := strings.ToLower(route.Value())
	cleanedPath := strings.TrimSpace(strings.ToLower(path))

	// check if the current route ends with the supplied path
	routeEndsWithSpecifiedPath := strings.HasSuffix(cleanedRoute, cleanedPath)

	return routeEndsWithSpecifiedPath
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

// Check if the current route is a child of the supplied (parent) route.
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

	return true
}

// Normalize the supplied path to be used for an Item or File
func normalize(path string) string {

	// trim spaces
	path = strings.TrimSpace(path)

	// check if the path is empty
	if path == "" {
		return ""
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

	// replace brackets
	path = strings.Replace(path, "(", "", -1)
	path = strings.Replace(path, ")", "", -1)

	return path
}
