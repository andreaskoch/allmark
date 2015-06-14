// Copyright 2014 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package orchestrator

import (
	"fmt"
	"strings"
	"time"

	"allmark.io/modules/common/config"
	"allmark.io/modules/common/logger"
	"allmark.io/modules/common/paths"
	"allmark.io/modules/common/route"
	"allmark.io/modules/dataaccess"
	"allmark.io/modules/model"
	"allmark.io/modules/services/converter"
	"allmark.io/modules/services/parser"
	"allmark.io/modules/web/orchestrator/index"
	"allmark.io/modules/web/orchestrator/search"
	"allmark.io/modules/web/view/viewmodel"
	"allmark.io/modules/web/webpaths"
)

type UpdateType int

const (
	UpdateTypeUnchanged UpdateType = iota
	UpdateTypeNew
	UpdateTypeModified
	UpdateTypeDeleted
)

func (updateType UpdateType) String() string {
	switch updateType {
	case UpdateTypeUnchanged:
		return "Unchanged"

	case UpdateTypeNew:
		return "new"

	case UpdateTypeModified:
		return "modified"

	case UpdateTypeDeleted:
		return "deleted"
	}

	panic("Unknown update type")
}

func NewUpdate(updateType UpdateType, route route.Route) Update {
	return Update{updateType, route}
}

type Update struct {
	updateType UpdateType
	route      route.Route
}

func (update *Update) String() string {
	return fmt.Sprintf("%s (%s)", update.route.String(), update.updateType.String())
}

func (update *Update) Route() route.Route {
	return update.route
}

func (update *Update) Type() UpdateType {
	return update.updateType
}

func cacheUpdate(name string, updateType UpdateType, callback func(updatedRoute route.Route)) CacheUpdateCallback {
	return CacheUpdateCallback{
		name,
		updateType,
		callback,
	}
}

type CacheUpdateCallback struct {
	Name   string
	Type   UpdateType
	Update func(updatedRoute route.Route)
}

func newBaseOrchestrator(logger logger.Logger, config config.Config, repository dataaccess.Repository, parser parser.Parser, converter converter.Converter, webPathProvider webpaths.WebPathProvider) *Orchestrator {

	startTime := time.Now()

	orchestrator := &Orchestrator{
		logger: logger,

		config:     config,
		repository: repository,
		parser:     parser,
		converter:  converter,

		webPathProvider: webPathProvider,

		updateSubscribers: make([]chan Update, 0),
		updateCallbacks:   make(map[UpdateType][]CacheUpdateCallback),
	}

	// warm up caches
	orchestrator.blockingCacheWarmup()

	stopTime := time.Now()
	duration := stopTime.Sub(startTime)
	logger.Statistics("Priming the base orchestrator cache took %f seconds.", duration.Seconds())

	return orchestrator
}

type Orchestrator struct {
	logger logger.Logger

	config     config.Config
	repository dataaccess.Repository
	parser     parser.Parser
	converter  converter.Converter

	webPathProvider webpaths.WebPathProvider

	// caches and indizes (do not initialize!)
	fulltextIndex   *search.ItemSearch
	repositoryIndex *index.Index
	items           []*model.Item
	itemsByAlias    map[string]*model.Item

	// update handling
	updateCallbacks   map[UpdateType][]CacheUpdateCallback
	updateSubscribers []chan Update
}

// Get the full-page title for a given headline.
func (orchestrator *Orchestrator) GetPageTitle(headline string) string {
	rootItem := orchestrator.rootItem()
	return fmt.Sprintf("%s - %s", headline, rootItem.Title)
}

// blockingCacheWarmup triggers a cache-warmup.
func (orchestrator *Orchestrator) blockingCacheWarmup() {
	orchestrator.index()
	orchestrator.getAllItems()
	orchestrator.search("", 5)
	orchestrator.getItemByAlias("")
}

func (orchestrator *Orchestrator) Subscribe(update chan Update) {
	orchestrator.updateSubscribers = append(orchestrator.updateSubscribers, update)
}

// Update all caches
func (orchestrator *Orchestrator) UpdateCache(dataaccessLayerUpdate dataaccess.Update) {

	// inform subscribers ...
	// ... about new items
	for _, newItemRoute := range dataaccessLayerUpdate.New() {

		// execute cache update callbacks
		for _, callbackDefinition := range orchestrator.updateCallbacks[UpdateTypeNew] {
			callbackDefinition.Update(newItemRoute)
		}

		// notify subscribers
		for _, subscriber := range orchestrator.updateSubscribers {
			subscriber <- NewUpdate(UpdateTypeNew, newItemRoute)
		}
	}

	// ... about modified items
	for _, modifiedItemRoute := range dataaccessLayerUpdate.Modified() {

		// execute cache update callbacks
		for _, callbackDefinition := range orchestrator.updateCallbacks[UpdateTypeModified] {
			callbackDefinition.Update(modifiedItemRoute)
		}

		// notify subscribers
		for _, subscriber := range orchestrator.updateSubscribers {
			subscriber <- NewUpdate(UpdateTypeModified, modifiedItemRoute)
		}
	}

	// ... about deleted items
	for _, deletedItemRoute := range dataaccessLayerUpdate.Deleted() {

		// execute cache update callbacks
		for _, callbackDefinition := range orchestrator.updateCallbacks[UpdateTypeDeleted] {
			callbackDefinition.Update(deletedItemRoute)
		}

		// notify subscribers
		for _, subscriber := range orchestrator.updateSubscribers {
			subscriber <- NewUpdate(UpdateTypeDeleted, deletedItemRoute)
		}
	}
}

// registerUpdateCallback registers callbacks for new, modified and deleted items.
func (orchestrator *Orchestrator) registerUpdateCallback(name string, updateType UpdateType, callback func(updatedRoute route.Route)) {

	if orchestrator.updateCallbacks[updateType] == nil {
		orchestrator.updateCallbacks[updateType] = make([]CacheUpdateCallback, 0)
	}

	orchestrator.updateCallbacks[updateType] = append(orchestrator.updateCallbacks[updateType], cacheUpdate(name, updateType, callback))
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

func (orchestrator *Orchestrator) parseItem(item dataaccess.Item) *model.Item {
	parsedItem, err := orchestrator.parser.ParseItem(item)
	if err != nil {
		orchestrator.logger.Warn(err.Error())
		return nil
	}

	return parsedItem
}

func (orchestrator *Orchestrator) parseFile(file dataaccess.File) *model.File {
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

func (orchestrator *Orchestrator) getLatestItems(parentRoute route.Route) []*model.Item {

	leafes := orchestrator.index().GetLeafes(parentRoute)

	// sort the leafes by date
	model.SortItemsBy(sortItemsByDate).Sort(leafes)

	return leafes
}

func (orchestrator *Orchestrator) index() *index.Index {

	if orchestrator.repositoryIndex != nil {
		return orchestrator.repositoryIndex
	}

	// newItem fetches the item with the given route and adds it to the index.
	updateItem := func(updatedRoute route.Route) {

		// get the updated item from the repository
		updatedRepositoryItem := orchestrator.repository.Item(updatedRoute)

		// parse the item
		parsedItem := orchestrator.parseItem(updatedRepositoryItem)
		if parsedItem == nil {
			orchestrator.logger.Warn("Unable to parse item %q", parsedItem.String())
			return
		}

		// update the index
		if orchestrator.repositoryIndex == nil {
			orchestrator.logger.Warn("Cannot add item %q, the index has not been initialized yet.", parsedItem.String())
			return
		}

		orchestrator.repositoryIndex.Add(parsedItem)
	}

	// deleteItem deletes the item with the given route from the index.
	deleteItem := func(deletedRoute route.Route) {
		orchestrator.repositoryIndex.Remove(deletedRoute)
	}

	// create a new index
	orchestrator.repositoryIndex = index.New(orchestrator.logger)

	// parse all items
	repositoryItems := orchestrator.repository.Items()
	for _, repositoryItem := range repositoryItems {
		parsedItem := orchestrator.parseItem(repositoryItem)
		if parsedItem == nil {
			orchestrator.logger.Warn("Unable to parse item %q", repositoryItem.String())
			continue
		}

		orchestrator.repositoryIndex.Add(parsedItem)
	}

	// register update callbacks
	orchestrator.registerUpdateCallback("update index", UpdateTypeNew, updateItem)
	orchestrator.registerUpdateCallback("update index", UpdateTypeModified, updateItem)
	orchestrator.registerUpdateCallback("update index", UpdateTypeDeleted, deleteItem)

	return orchestrator.repositoryIndex
}

func (orchestrator *Orchestrator) search(keywords string, maxiumNumberOfResults int) []search.Result {

	if orchestrator.fulltextIndex != nil {
		return orchestrator.fulltextIndex.Search(keywords, maxiumNumberOfResults)
	}

	newFullTextIndex := search.NewItemSearch(orchestrator.logger, orchestrator.getAllItems())

	// destroy the old index
	if orchestrator.fulltextIndex != nil {
		oldIndex := orchestrator.fulltextIndex
		go oldIndex.Destroy()
	}

	orchestrator.fulltextIndex = newFullTextIndex

	return orchestrator.fulltextIndex.Search(keywords, maxiumNumberOfResults)
}

func (orchestrator *Orchestrator) getAllItems() []*model.Item {

	allItems := orchestrator.index().GetAllItems()
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

	latestItems := orchestrator.getLatestItems(route.New())
	if len(latestItems) == 0 {
		return nil
	}

	// determine the position of the supplied route
	matchingIndex := -1
	for index, item := range latestItems {
		if item.Route().Value() == currentRoute.Value() {
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
	if nextIndex >= len(latestItems) {
		return nil
	}

	return latestItems[nextIndex]
}

func (orchestrator *Orchestrator) getNext(currentRoute route.Route) *model.Item {

	latestItems := orchestrator.getLatestItems(route.New())
	if len(latestItems) == 0 {
		return nil
	}

	// determine the position of the supplied route
	matchingIndex := -1
	for index, item := range latestItems {
		if item.Route().Value() == currentRoute.Value() {
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

	return latestItems[previousIndex]
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

	// return from cache
	alias = normalizeAlias(alias)
	if orchestrator.itemsByAlias != nil {
		return orchestrator.itemsByAlias[alias]
	}

	// updateAliasMap updates the alias map for the given route.
	updateAliasMap := func(route route.Route) {
		item := orchestrator.getItem(route)
		alias := normalizeAlias(item.MetaData.Alias)
		if alias == "" {
			return
		}

		// add item to map
		orchestrator.itemsByAlias[alias] = item
	}

	// removeItemFromAliasMap removed the item with the given route from the alias map.
	removeItemFromAliasMap := func(route route.Route) {

		aliasToRemove := ""
		for alias, item := range orchestrator.itemsByAlias {
			if item.Route().Equals(route) {
				aliasToRemove = alias
				break
			}
		}

		// abort if no alias was found
		if aliasToRemove == "" {
			return
		}

		// remove item from map
		delete(orchestrator.itemsByAlias, aliasToRemove)
	}

	// build cache
	itemsByAlias := make(map[string]*model.Item)

	for _, item := range orchestrator.getAllItems() {

		// ignore items without an alias
		if item.MetaData.Alias == "" {
			continue
		}

		itemAlias := normalizeAlias(item.MetaData.Alias)
		itemsByAlias[itemAlias] = item
	}

	orchestrator.itemsByAlias = itemsByAlias

	// register update callbacks
	orchestrator.registerUpdateCallback("update alias map", UpdateTypeNew, updateAliasMap)
	orchestrator.registerUpdateCallback("update alias map", UpdateTypeModified, updateAliasMap)
	orchestrator.registerUpdateCallback("update alias map", UpdateTypeDeleted, removeItemFromAliasMap)

	return orchestrator.itemsByAlias[alias]
}

// Get the publisher information view model.
func (orchestrator *Orchestrator) getPublisherInformation() viewmodel.Publisher {
	return viewmodel.Publisher{
		Name:  orchestrator.config.Web.Publisher.Name,
		Email: orchestrator.config.Web.Publisher.Email,
		Url:   orchestrator.config.Web.Publisher.Url,

		GooglePlusHandle: orchestrator.config.Web.Publisher.GooglePlusHandle,
		TwitterHandle:    orchestrator.config.Web.Publisher.TwitterHandle,
		FacebookHandle:   orchestrator.config.Web.Publisher.FacebookHandle,
	}
}

// Get the publisher information view model.
func (orchestrator *Orchestrator) getAuthorInformation(authorName string) viewmodel.Author {

	if authorName == "" {
		return viewmodel.Author{}
	}

	author, exists := orchestrator.config.Web.Authors[authorName]
	if !exists {
		return viewmodel.Author{
			Name: authorName,
		}
	}

	return viewmodel.Author{
		Name:  author.Name,
		Email: author.Email,
		Url:   author.Url,

		GooglePlusHandle: author.GooglePlusHandle,
		TwitterHandle:    author.TwitterHandle,
		FacebookHandle:   author.FacebookHandle,
	}
}

// Get the analytics view model.
func (orchestrator *Orchestrator) getAnalyticsSettings() viewmodel.Analytics {
	return viewmodel.Analytics{
		Enabled: orchestrator.config.Analytics.Enabled,
		GoogleAnalytics: viewmodel.GoogleAnalytics{
			Enabled:    orchestrator.config.Analytics.GoogleAnalytics.Enabled,
			TrackingId: orchestrator.config.Analytics.GoogleAnalytics.TrackingId,
		},
	}
}

// normalizeAlias takes the supplied alias and normalizes it.
func normalizeAlias(alias string) string {
	return strings.TrimSpace(strings.ToLower(alias))
}
