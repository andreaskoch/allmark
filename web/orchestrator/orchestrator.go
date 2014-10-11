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
	"github.com/andreaskoch/allmark2/web/orchestrator/index"
	"github.com/andreaskoch/allmark2/web/orchestrator/search"
	"github.com/andreaskoch/allmark2/web/view/viewmodel"
	"github.com/andreaskoch/allmark2/web/webpaths"
	"strings"
	"time"
)

type CacheState int

const (
	CacheStateStale CacheState = iota
	CacheStatePrimed
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

	// cache control
	cacheStatusMap    map[string]CacheState
	cachePrimerMap    map[string]func()
	cachePrimerStatus map[string]bool

	// caches and indizes
	fulltextIndex   *search.ItemSearch
	repositoryIndex *index.Index
	items           []*model.Item
	itemsByAlias    map[string]*model.Item
	leafesByRoute   map[string][]route.Route
}

// Reset all Caches
func (orchestrator *Orchestrator) ResetCache() {

	// mark all caches as stale
	for cacheType := range orchestrator.cacheStatusMap {
		orchestrator.cacheStatusMap[cacheType] = CacheStateStale
	}

	// prime all caches asynchronously
	go func() {
		orchestrator.primeCaches()
	}()

}

func (orchestrator *Orchestrator) setCache(cacheType string, primer func()) {

	// initialize the primer map on first use
	if orchestrator.cachePrimerMap == nil {
		orchestrator.cachePrimerMap = make(map[string]func())
	}

	// initialize the status map on first use
	if orchestrator.cacheStatusMap == nil {
		orchestrator.cacheStatusMap = make(map[string]CacheState)
	}

	// store the primer
	orchestrator.cachePrimerMap[cacheType] = primer

	// fill the cache
	primer()

	// mark the cache type as primed
	orchestrator.cacheStatusMap[cacheType] = CacheStatePrimed
}

func (orchestrator *Orchestrator) isCacheStale(cacheType string) bool {
	if status, exists := orchestrator.cacheStatusMap[cacheType]; exists {
		return status == CacheStateStale
	}

	// if there is no status it is definitly stale
	return true
}

// Prime all caches
func (orchestrator *Orchestrator) primeCaches() {
	for cacheType := range orchestrator.cacheStatusMap {
		orchestrator.primeCache(cacheType)
	}
}

// Prime a particular cache
func (orchestrator *Orchestrator) primeCache(cacheType string) {

	// initialize the mutex map
	if orchestrator.cachePrimerStatus == nil {
		orchestrator.cachePrimerStatus = make(map[string]bool)
	}

	// check if there is a mutex
	if exists, _ := orchestrator.cachePrimerStatus[cacheType]; exists {

		// abort. There is already a primer running for the supplied cache type
		return
	}

	// set a mutex for the supplied cache type
	orchestrator.cachePrimerStatus[cacheType] = true

	// execute the primer func
	primerFunc := orchestrator.cachePrimerMap[cacheType]
	primerFunc()

	// set the cache status to "primed"
	orchestrator.cacheStatusMap[cacheType] = CacheStatePrimed

	// release the mutex
	defer delete(orchestrator.cachePrimerStatus, cacheType)
}

func (orchestrator *Orchestrator) ItemExists(route route.Route) bool {
	_, exists := orchestrator.index().IsMatch(route)
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
	return orchestrator.index().Root()
}

func (orchestrator *Orchestrator) getItem(route route.Route) *model.Item {

	if item, exists := orchestrator.index().IsMatch(route); exists {
		return item
	}

	return nil
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

	// initialize the leafes map on first use
	if orchestrator.leafesByRoute == nil {
		orchestrator.leafesByRoute = make(map[string][]route.Route)
	}

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

func (orchestrator *Orchestrator) index() *index.Index {

	cacheType := "index"

	// load from cache
	if orchestrator.repositoryIndex != nil {

		// re-prime the cache if it is stale
		if orchestrator.isCacheStale(cacheType) {
			go orchestrator.primeCache(cacheType)
		}

		return orchestrator.repositoryIndex
	}

	orchestrator.setCache(cacheType, func() {

		// parse all items
		repositoryItems := orchestrator.repository.Items()
		parsedItems := make([]*model.Item, 0, len(repositoryItems))
		for _, repositoryItem := range repositoryItems {
			parsedItem := orchestrator.parseItem(repositoryItem)
			if parsedItem == nil {
				continue
			}

			parsedItems = append(parsedItems, parsedItem)
		}

		// create a new index
		newIndex := index.New(orchestrator.logger, parsedItems)

		// store to cache
		orchestrator.repositoryIndex = newIndex
	})

	return orchestrator.repositoryIndex
}

func (orchestrator *Orchestrator) search(keywords string, maxiumNumberOfResults int) []search.Result {
	cacheType := "fulltextIndex"

	// load from cache
	if orchestrator.fulltextIndex != nil {

		// re-prime the cache if it is stale
		if orchestrator.isCacheStale(cacheType) {
			go orchestrator.primeCache(cacheType)
		}

		return orchestrator.fulltextIndex.Search(keywords, maxiumNumberOfResults)
	}

	orchestrator.setCache(cacheType, func() {

		newFullTextIndex := search.NewItemSearch(orchestrator.logger, orchestrator.getAllItems())

		// destroy the old index
		if orchestrator.fulltextIndex != nil {
			oldIndex := orchestrator.fulltextIndex
			go oldIndex.Destroy()
		}

		// store to cache
		orchestrator.fulltextIndex = newFullTextIndex
	})

	return orchestrator.fulltextIndex.Search(keywords, maxiumNumberOfResults)
}

func (orchestrator *Orchestrator) getAllItems() []*model.Item {

	cacheType := "allItems"

	// load from cache
	if orchestrator.items != nil {

		// re-prime the cache if it is stale
		if orchestrator.isCacheStale(cacheType) {
			go orchestrator.primeCache(cacheType)
		}

		return orchestrator.items
	}

	orchestrator.setCache(cacheType, func() {

		// get all items
		allItems := orchestrator.index().GetAllItems()

		// sort the items by date
		model.SortItemsBy(sortItemsByDate).Sort(allItems)

		// store to cache
		orchestrator.items = allItems
	})

	return orchestrator.items
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
	file, exists := orchestrator.index().IsFileMatch(route)
	if !exists {
		return nil
	}

	return file
}

func (orchestrator *Orchestrator) getParent(route route.Route) *model.Item {
	parent := orchestrator.index().GetParent(route)
	if parent == nil {
		return nil
	}

	return parent
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

func (orchestrator *Orchestrator) getChilds(route route.Route) []*model.Item {

	// get all childs
	childs := orchestrator.index().GetDirectChilds(route)

	// sort the childs by date
	model.SortItemsBy(sortItemsByDate).Sort(childs)

	return childs
}

// Get the item that has the specified alias. Returns nil if there is no matching item.
func (orchestrator *Orchestrator) getItemByAlias(alias string) *model.Item {

	cacheType := "itembyalias"
	alias = strings.TrimSpace(strings.ToLower(alias))

	// load from cache
	if orchestrator.itemsByAlias != nil {

		// re-prime the cache if it is stale
		if orchestrator.isCacheStale(cacheType) {
			go orchestrator.primeCache(cacheType)
		}

		return orchestrator.itemsByAlias[alias]
	}

	orchestrator.setCache(cacheType, func() {

		itemsByAlias := make(map[string]*model.Item)

		for _, item := range orchestrator.getAllItems() {

			// continue items without an alias
			if item.MetaData.Alias == "" {
				continue
			}

			itemAlias := strings.TrimSpace(strings.ToLower(item.MetaData.Alias))
			itemsByAlias[itemAlias] = item
		}

		orchestrator.itemsByAlias = itemsByAlias
	})

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
