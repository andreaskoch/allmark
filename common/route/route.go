// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	regexpWhitespacePattern          = regexp.MustCompile(`\s+`)
	regexpBackSlashPattern           = regexp.MustCompile(`\\+`)
	regexpdoubleForwardSlashPattern  = regexp.MustCompile(`/+`)
	regexpForbiddenCharactersPattern = regexp.MustCompile(`[&%§\)\(}{\]\["|]`)
)

type Route struct {
	value         string
	originalValue string
	isFileRoute   bool
}

// Creates a new route from the given base directory (e.g. "/home/user/repository") and item path (e.g. "/home/user/repository/documents/sample/document.md") and strips the file name from the route.
func NewFromItemPath(baseDirectory, itemPath string) Route {

	// normalize the base path
	normalizedBasePath := normalize(baseDirectory)

	// normalize the item path
	normalizedItemPath := normalize(itemPath)

	// return a root if both paths are the same
	if normalizedItemPath == normalizedBasePath {
		return New()
	}

	// prepare the route value:
	// strip the repository path from the item path
	routeValue := strings.Replace(normalizedItemPath, normalizedBasePath, "", 1)

	// strip the file name
	routeValue = routeValue[:strings.LastIndex(routeValue, "/")]

	// trim leading slashes
	routeValue = strings.TrimLeft(routeValue, "/")

	return Route{
		value:         toURL(routeValue),
		originalValue: routeValue,
		isFileRoute:   true,
	}
}

// Creates a new route from the given base directory (e.g. "/home/user/repository") and item directory (e.g. "/home/user/repository/documents/sample").
func NewFromItemDirectory(baseDirectory, itemDirectory string) Route {

	// normalize the base path
	normalizedBasePath := normalize(baseDirectory)

	// normalize the item path
	normalizedItemPath := normalize(itemDirectory)

	// return a root if both paths are the same
	if normalizedItemPath == normalizedBasePath {
		return New()
	}

	// prepare the route value:
	// strip the repository path from the item path
	routeValue := strings.Replace(normalizedItemPath, normalizedBasePath, "", 1)

	// trim leading slashes
	routeValue = strings.TrimLeft(routeValue, "/")

	return Route{
		value:         toURL(routeValue),
		originalValue: routeValue,
		isFileRoute:   false,
	}
}

// Creates a new route from the given base directory (e.g. "/home/user/repository") and item path (e.g. "/home/user/repository/documents/sample/files/image.jpg").
func NewFromFilePath(baseDirectory, itemPath string) Route {

	// normalize the base path
	normalizedBasePath := normalize(baseDirectory)

	// normalize the item path
	normalizedItemPath := normalize(itemPath)

	// return a root if both paths are the same
	if normalizedItemPath == normalizedBasePath {
		return New()
	}

	// prepare the route value:
	// strip the repository path from the item path
	routeValue := strings.Replace(normalizedItemPath, normalizedBasePath, "", 1)

	// trim leading slashes
	routeValue = strings.TrimLeft(routeValue, "/")

	return Route{
		value:         toURL(routeValue),
		originalValue: routeValue,
		isFileRoute:   true,
	}
}

// Create a new route from the given request path.
func NewFromRequest(requestPath string) Route {

	// normalize the request path
	routeValue := normalize(requestPath)

	return Route{
		value:         toURL(routeValue),
		originalValue: routeValue,
	}
}

// Create an empty route.
func New() Route {

	// normalize the request path
	routeValue := normalize("")

	return Route{
		value:         toURL(routeValue),
		originalValue: routeValue,
		isFileRoute:   false,
	}
}

// combines two routes
func Combine(route1, route2 Route) Route {
	return NewFromRequest(route1.OriginalValue() + "/" + route2.OriginalValue())
}

func Intersect(baseRoute, subRoute Route) Route {

	// abort if the base route is empty
	if baseRoute.IsEmpty() {
		return baseRoute
	}

	// abort if the subroute is empty
	if subRoute.IsEmpty() {
		return baseRoute
	}

	baseDirectory := baseRoute.OriginalValue()
	itemDirectory := subRoute.OriginalValue()

	// normalize the base path
	normalizedBasePath := normalize(baseDirectory)

	// normalize the item path
	normalizedItemPath := normalize(itemDirectory)

	// return the base if the sub-route is the same
	if normalizedItemPath == normalizedBasePath {
		return baseRoute
	}

	// prepare the route value:
	// strip the repository path from the item path
	routeValue := strings.Replace(normalizedItemPath, normalizedBasePath, "", 1)

	// trim leading slashes
	routeValue = strings.TrimLeft(routeValue, "/")

	return Route{
		value:         toURL(routeValue),
		originalValue: routeValue,
	}
}

func (route Route) String() string {
	return strings.Join(route.Components(), " > ")
}

// Equals compares the current route with the supplied route. Returns true if the routes are alike; otherwise false.
func (route Route) Equals(otherRoute Route) bool {
	return route.String() == otherRoute.String()
}

func (route Route) OriginalValue() string {
	return route.originalValue
}

func (route Route) Components() []string {
	return strings.Split(route.OriginalValue(), "/")
}

func (route Route) Value() string {
	return route.value
}

func (route Route) IsEmpty() bool {
	return len(route.value) == 0
}

func (route Route) IsFileRoute() bool {
	return route.isFileRoute
}

func (route Route) Path() string {
	lastSlashPosition := strings.LastIndex(route.originalValue, "/")
	if lastSlashPosition == -1 {
		return route.originalValue
	}

	return strings.TrimSuffix(route.originalValue[:lastSlashPosition], "/")
}

func (route Route) FirstComponentName() string {

	if route.Level() == 0 {
		return ""
	}

	// handle root-level routes
	if route.Level() == 1 {

		if route.IsFileRoute() {
			return ""
		}

		return route.originalValue
	}

	firstSlashPosition := strings.Index(route.originalValue, "/")
	return strings.TrimSuffix(route.originalValue[:firstSlashPosition], "/")
}

func (route Route) LastComponentName() string {
	lastSlashPosition := strings.LastIndex(route.originalValue, "/")
	if lastSlashPosition == -1 {
		return route.originalValue
	}

	return strings.TrimPrefix(route.originalValue[lastSlashPosition:], "/")
}

func (route Route) Level() int {

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

func (route Route) SubRoute(level int) (Route, error) {

	// root level
	if level == 0 {
		return NewFromRequest(""), nil
	}

	// same level
	if level == route.Level() {
		return route, nil
	}

	// split path into components
	components := strings.Split(route.value, "/")

	// abort if the requested level is out of range
	if level > len(components)-1 {
		return Route{}, fmt.Errorf("The route %q does nof a have a sub-route with the level %d.", route, level)
	}

	// assemble the sub route
	subset := components[0:level]
	subRoutePath := strings.Join(subset, "/")

	subRoute := NewFromRequest(subRoutePath)
	return subRoute, nil
}

func (route Route) IsMatch(path string) bool {
	cleanedRoute := strings.ToLower(route.Value())
	normalizedPath := strings.ToLower(toURL(normalize(path)))

	// check if the current route ends with the supplied path
	routeEndsWithSpecifiedPath := strings.HasSuffix(cleanedRoute, normalizedPath)

	return routeEndsWithSpecifiedPath
}

func (route Route) Parent() (parent Route, exists bool) {

	if route.IsEmpty() {
		return parent, false
	}

	routeValue := route.Value()

	// if there is no slash, the parent must be the root
	if !strings.Contains(routeValue, "/") {
		return New(), true
	}

	positionOfLastSlash := strings.LastIndex(routeValue, "/")
	parentRouteValue := routeValue[:positionOfLastSlash]

	parentRoute := NewFromRequest(parentRouteValue)
	return parentRoute, true
}

// Check if the the current route is direct parent for the supplied (child) route.
func (parent Route) IsParentOf(child Route) bool {
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
func (child Route) IsChildOf(parent Route) bool {
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

// Returns a normalized version of the supplied path
func normalize(path string) string {

	// trim spaces
	path = strings.TrimSpace(path)

	// check if the path is empty
	if path == "" {
		return ""
	}

	path = fromURL(path)

	// replace all forbidden characters
	path = regexpForbiddenCharactersPattern.ReplaceAllString(path, "")

	// replace all backslashes with a (single) forward slash
	path = regexpBackSlashPattern.ReplaceAllString(path, "/")

	// replace multiple forward slashes with a single forward slash
	path = regexpdoubleForwardSlashPattern.ReplaceAllString(path, "/")

	// remove leading slashes
	path = strings.TrimLeft(path, "/")

	// remove trailing slashes
	path = strings.TrimRight(path, "/")

	// replace duplicate spaces
	path = regexpWhitespacePattern.ReplaceAllString(path, " ")

	return path
}

// Returns an "url-safe" version of the supplied path
func toURL(path string) string {

	// replace duplicate spaces with a (single) url safe character
	path = strings.Replace(path, " ", "+", -1)

	// replace brackets
	path = strings.Replace(path, "(", "%28", -1)
	path = strings.Replace(path, ")", "%29", -1)

	return path
}

// Returns an "url-safe" version of the supplied path
func fromURL(path string) string {

	// replace duplicate spaces with a (single) url safe character
	path = strings.Replace(path, "+", " ", -1)

	// replace brackets
	path = strings.Replace(path, "%28", "(", -1)
	path = strings.Replace(path, "%29", ")", -1)

	return path
}
