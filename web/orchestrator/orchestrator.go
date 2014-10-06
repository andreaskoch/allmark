// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"github.com/andreaskoch/allmark2/common/config"
	"github.com/andreaskoch/allmark2/common/logger"
	"github.com/andreaskoch/allmark2/common/paths"
	"github.com/andreaskoch/allmark2/common/route"
	"github.com/andreaskoch/allmark2/dataaccess"
	"github.com/andreaskoch/allmark2/model"
	"github.com/andreaskoch/allmark2/services/converter"
	"github.com/andreaskoch/allmark2/services/parser"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
	"github.com/andreaskoch/allmark2/web/webpaths"
	"sort"
	"strings"
	"time"
)

func newBaseOrchestrator(logger logger.Logger, config config.Config, repository dataaccess.Repository, parser parser.Parser, converter converter.Converter, webPathProvider webpaths.WebPathProvider) *Orchestrator {
	return &Orchestrator{
		logger: logger,

		config:     config,
		repository: repository,
		parser:     parser,
		converter:  converter,

		webPathProvider: webPathProvider,
	}
}

type Orchestrator struct {
	logger logger.Logger

	config     config.Config
	repository dataaccess.Repository
	parser     parser.Parser
	converter  converter.Converter

	webPathProvider webpaths.WebPathProvider

	// caches
	itemsByAlias  map[string]*model.Item
	leafesByRoute map[string][]route.Route
}

func (orchestrator *Orchestrator) ResetCache() {
	orchestrator.itemsByAlias = make(map[string]*model.Item)
	orchestrator.leafesByRoute = make(map[string][]route.Route)
}

func (orchestrator *Orchestrator) ItemExists(route route.Route) bool {
	_, exists := orchestrator.repository.Item(route)
	return exists
}

func (orchestrator *Orchestrator) absolutePather(prefix string) paths.Pather {
	return orchestrator.webPathProvider.AbsolutePather(prefix)
}

func (orchestrator *Orchestrator) itemPather() paths.Pather {
	return orchestrator.webPathProvider.ItemPather()
}

func (orchestrator *Orchestrator) tagPather() paths.Pather {
	return orchestrator.webPathProvider.TagPather()
}

func (orchestrator *Orchestrator) relativePather(baseRoute route.Route) paths.Pather {
	return orchestrator.webPathProvider.RelativePather(baseRoute)
}

func (orchestrator *Orchestrator) parseItem(item *dataaccess.Item) *model.Item {
	parsedItem, err := orchestrator.parser.ParseItem(item)
	if err != nil {
		orchestrator.logger.Warn(err.Error())
		return nil
	}

	return parsedItem
}

func (orchestrator *Orchestrator) parseFile(file *dataaccess.File) *model.File {
	parsedFile, err := orchestrator.parser.ParseFile(file)
	if err != nil {
		orchestrator.logger.Warn(err.Error())
		return nil
	}

	return parsedFile
}

func (orchestrator *Orchestrator) rootItem() *model.Item {
	return orchestrator.parseItem(orchestrator.repository.Root())
}

func (orchestrator *Orchestrator) getItem(route route.Route) *model.Item {
	item, exists := orchestrator.repository.Item(route)
	if !exists {
		return nil
	}

	return orchestrator.parseItem(item)
}

func (orchestrator *Orchestrator) getLatestRoutesByPage(parentRoute route.Route, pageSize, page int) (routes []route.Route, found bool) {

	latestRoutes, found := orchestrator.getLatestRoutes(parentRoute)
	if !found {
		return []route.Route{}, false
	}

	// determine the start index
	startIndex := pageSize * (page - 1)
	if startIndex >= len(latestRoutes) {
		return []route.Route{}, false
	}

	// determine the end index
	endIndex := startIndex + pageSize
	if endIndex > len(latestRoutes) {
		endIndex = len(latestRoutes)
	}

	return latestRoutes[startIndex:endIndex], true
}

func (orchestrator *Orchestrator) getLatestRoutes(parentRoute route.Route) (routes []route.Route, found bool) {

	leafes := orchestrator.getAllLeafes(parentRoute)

	// collect the creation dates for all leafes
	routesAndDates := make([]routeAndDate, 0, len(leafes))
	for _, leaf := range leafes {
		creationDate, found := orchestrator.getCreationDate(leaf)
		if !found {
			// todo: log info
			continue
		}

		routesAndDates = append(routesAndDates, routeAndDate{leaf, creationDate})
	}

	// sort the leafes by date
	SortItemRoutesAndDatesBy(sortRoutesAndDatesDescending).Sort(routesAndDates)

	routes = make([]route.Route, 0)
	for _, routeAndDate := range routesAndDates {
		routes = append(routes, routeAndDate.route)
	}

	return routes, true
}

func (orchestrator *Orchestrator) getAllLeafes(parentRoute route.Route) []route.Route {

	// cache lookup
	key := parentRoute.Value()
	if leafes, isset := orchestrator.leafesByRoute[key]; isset {
		return leafes
	}

	childRoutes := make([]route.Route, 0)

	childItems := orchestrator.getChilds(parentRoute)
	if hasNoMoreChilds := len(childItems) == 0; hasNoMoreChilds {
		return []route.Route{parentRoute}
	}

	// recurse
	for _, childItem := range childItems {

		// skip locations
		if childItem.Type == model.TypeLocation {
			continue
		}

		childRoutes = append(childRoutes, orchestrator.getAllLeafes(childItem.Route())...)
	}

	// store the value
	orchestrator.leafesByRoute[key] = childRoutes

	return childRoutes

}

func (orchestrator *Orchestrator) getAllItems() []*model.Item {

	allItems := make([]*model.Item, 0)

	for _, repositoryItem := range orchestrator.repository.Items() {
		item := orchestrator.parseItem(repositoryItem)
		if item == nil {
			continue
		}

		allItems = append(allItems, item)
	}

	model.SortItemsBy(sortItemsByDate).Sort(allItems)

	return allItems
}

func (orchestrator *Orchestrator) getItems(pageSize, page int) []*model.Item {

	allItems := orchestrator.getAllItems()

	// determine the start index
	startIndex := pageSize * (page - 1)
	if startIndex >= len(allItems) {
		return []*model.Item{}
	}

	// determine the end index
	endIndex := startIndex + pageSize
	if endIndex > len(allItems) {
		endIndex = len(allItems)
	}

	return allItems[startIndex:endIndex]
}

func (orchestrator *Orchestrator) getCreationDate(itemRoute route.Route) (creationDate time.Time, found bool) {

	item := orchestrator.getItem(itemRoute)
	if item == nil {
		return time.Time{}, false
	}

	return item.MetaData.CreationDate, true
}

func (orchestrator *Orchestrator) getFile(route route.Route) *model.File {
	file, exists := orchestrator.repository.File(route)
	if !exists {
		return nil
	}

	return orchestrator.parseFile(file)
}

func (orchestrator *Orchestrator) getParent(route route.Route) *model.Item {
	parent := orchestrator.repository.Parent(route)
	if parent == nil {
		return nil
	}

	return orchestrator.parseItem(parent)
}

func (orchestrator *Orchestrator) getPrevious(currentRoute route.Route) *model.Item {

	latestRoutes, found := orchestrator.getLatestRoutes(route.New())
	if !found {
		return nil
	}

	// determine the position of the supplied route
	matchingIndex := -1
	for index, route := range latestRoutes {
		if route.Value() == currentRoute.Value() {
			matchingIndex = index
			break
		}
	}

	// abort if the route was not found at all
	if noMatchFound := (matchingIndex == -1); noMatchFound {
		return nil
	}

	// abort if there is no next item
	nextIndex := matchingIndex + 1
	if nextIndex >= len(latestRoutes) {
		return nil
	}

	// determine the next route
	nextRoute := latestRoutes[nextIndex]

	return orchestrator.getItem(nextRoute)
}

func (orchestrator *Orchestrator) getNext(currentRoute route.Route) *model.Item {

	latestRoutes, found := orchestrator.getLatestRoutes(route.New())
	if !found {
		return nil
	}

	// determine the position of the supplied route
	matchingIndex := -1
	for index, route := range latestRoutes {
		if route.Value() == currentRoute.Value() {
			matchingIndex = index
			break
		}
	}

	// abort if the route was not found at all
	if noMatchFound := (matchingIndex == -1); noMatchFound {
		return nil
	}

	// abort if there is no previous item
	previousIndex := matchingIndex - 1
	if noPreviousItem := (previousIndex < 0); noPreviousItem {
		return nil
	}

	// determine the previous route
	previousRoute := latestRoutes[previousIndex]

	return orchestrator.getItem(previousRoute)
}

func (orchestrator *Orchestrator) getChilds(route route.Route) (childs []*model.Item) {

	childs = make([]*model.Item, 0)

	for _, child := range orchestrator.repository.Childs(route) {
		parsed := orchestrator.parseItem(child)
		if parsed == nil {
			continue
		}

		childs = append(childs, parsed)
	}

	return childs
}

// Get the item that has the specified alias. Returns nil if there is no matching item.
func (orchestrator *Orchestrator) getItemByAlias(alias string) *model.Item {

	alias = strings.TrimSpace(strings.ToLower(alias))

	if orchestrator.itemsByAlias == nil {

		orchestrator.logger.Info("Initializing alias list")
		itemsByAlias := make(map[string]*model.Item)

		for _, repositoryItem := range orchestrator.repository.Items() {

			item := orchestrator.parseItem(repositoryItem)
			if item == nil {
				orchestrator.logger.Warn("Cannot parse repository item %q.", repositoryItem.String())
				continue
			}

			// continue items without an alias
			if item.MetaData.Alias == "" {
				continue
			}

			itemAlias := strings.TrimSpace(strings.ToLower(item.MetaData.Alias))
			itemsByAlias[itemAlias] = item
		}

		// refresh control
		go func() {
			for {
				select {
				case <-orchestrator.repository.AfterReindex():
					// reset the list
					orchestrator.logger.Info("Resetting the alias list")
					orchestrator.itemsByAlias = nil
				}
			}
		}()

		orchestrator.itemsByAlias = itemsByAlias
	}

	return orchestrator.itemsByAlias[alias]
}

func (orchestrator *Orchestrator) getAnalyticsSettings() viewmodel.Analytics {
	return viewmodel.Analytics{
		Enabled: orchestrator.config.Analytics.Enabled,
		GoogleAnalytics: viewmodel.GoogleAnalytics{
			Enabled:    orchestrator.config.Analytics.GoogleAnalytics.Enabled,
			TrackingId: orchestrator.config.Analytics.GoogleAnalytics.TrackingId,
		},
	}
}

// sort the models by date and name
func sortItemsByDate(model1, model2 *model.Item) bool {

	return model1.MetaData.CreationDate.Before(model2.MetaData.CreationDate)

}

func sortRoutesAndDatesDescending(itemRoute1, itemRoute2 routeAndDate) bool {
	return itemRoute1.date.After(itemRoute2.date)
}

type routeAndDate struct {
	route route.Route
	date  time.Time
}

type SortItemRoutesAndDatesBy func(itemRoute1, itemRoute2 routeAndDate) bool

func (by SortItemRoutesAndDatesBy) Sort(routesAndDates []routeAndDate) {
	sorter := &routeAndDateSorter{
		routesAndDates: routesAndDates,
		by:             by,
	}

	sort.Sort(sorter)
}

type routeAndDateSorter struct {
	routesAndDates []routeAndDate
	by             SortItemRoutesAndDatesBy
}

func (sorter *routeAndDateSorter) Len() int {
	return len(sorter.routesAndDates)
}

func (sorter *routeAndDateSorter) Swap(i, j int) {
	sorter.routesAndDates[i], sorter.routesAndDates[j] = sorter.routesAndDates[j], sorter.routesAndDates[i]
}

func (sorter *routeAndDateSorter) Less(i, j int) bool {
	return sorter.by(sorter.routesAndDates[i], sorter.routesAndDates[j])
}
